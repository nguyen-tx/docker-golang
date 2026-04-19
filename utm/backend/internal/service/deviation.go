package service

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/utm/backend/internal/model"
)

type geoJSONPolygon struct {
	Type        string          `json:"type"`
	Coordinates [][][2]float64  `json:"coordinates"` // [ring][point][lng,lat]
}

// checkDeviation trả về LaneAlert nếu drone bay ngoài vùng cấp phép hoặc vượt altitude.
// Trả về nil nếu drone vẫn trong giới hạn.
func checkDeviation(t *model.Telemetry, flightAreaJSON string, plannedAltM int) *model.LaneAlert {
	inside, err := pointInPolygon(t.Latitude, t.Longitude, flightAreaJSON)
	if err != nil || inside {
		if inside && t.AltitudeM <= float64(plannedAltM) {
			return nil
		}
	}

	devM := distToPolygonBoundaryM(t.Latitude, t.Longitude, flightAreaJSON)

	severity := model.AlertSeverityWarning
	msg := fmt.Sprintf("Drone bay ngoài vùng cấp phép %.0fm", devM)

	if !inside && devM > 50 {
		severity = model.AlertSeverityDanger
	}
	if t.AltitudeM > float64(plannedAltM) {
		severity = model.AlertSeverityDanger
		msg = fmt.Sprintf("Vượt độ cao cho phép: %.0fm (giới hạn %dm)", t.AltitudeM, plannedAltM)
	}

	return &model.LaneAlert{
		Type:       "lane_alert",
		AircraftID: t.AircraftID,
		SessionID:  t.SessionID,
		Severity:   severity,
		DeviationM: devM,
		Lat:        t.Latitude,
		Lng:        t.Longitude,
		AltitudeM:  t.AltitudeM,
		Message:    msg,
		Timestamp:  time.Now().UTC(),
	}
}

// pointInPolygon dùng ray casting, GeoJSON coordinates là [lng, lat].
func pointInPolygon(lat, lng float64, geoJSON string) (bool, error) {
	var poly geoJSONPolygon
	if err := json.Unmarshal([]byte(geoJSON), &poly); err != nil {
		return false, err
	}
	if len(poly.Coordinates) == 0 {
		return false, nil
	}
	ring := poly.Coordinates[0]
	inside := false
	j := len(ring) - 1
	for i := 0; i < len(ring); i++ {
		xi, yi := ring[i][0], ring[i][1]
		xj, yj := ring[j][0], ring[j][1]
		if ((yi > lat) != (yj > lat)) && (lng < (xj-xi)*(lat-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside, nil
}

// distToPolygonBoundaryM trả về khoảng cách xấp xỉ tới đỉnh gần nhất của polygon (đơn vị mét).
func distToPolygonBoundaryM(lat, lng float64, geoJSON string) float64 {
	var poly geoJSONPolygon
	if err := json.Unmarshal([]byte(geoJSON), &poly); err != nil {
		return 0
	}
	if len(poly.Coordinates) == 0 {
		return 0
	}
	minDist := math.MaxFloat64
	for _, pt := range poly.Coordinates[0] {
		if d := haversineM(lat, lng, pt[1], pt[0]); d < minDist {
			minDist = d
		}
	}
	return minDist
}

func haversineM(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6_371_000.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
