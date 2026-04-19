# UTM System

Hệ thống quản lý không phận drone (Unmanned Traffic Management).

## Kiến trúc tổng quan

```
                         utm-alert :8081
                    ┌────────────────────────┐
                    │  /ws/backend           │
utm-backend ───WS──►│  BackendConn           │
  (WS client)       │     │  bridge.ToFE     │
                    │     ▼                  │
                    │  FeHub                 │◄──WS── FE-1
                    │     │  broadcast       │◄──WS── FE-2
                    │     │                  │◄──WS── FE-N
                    │  bridge.ToBackend      │
                    └────────────────────────┘

Telemetry
  Service ───WS──►  /ws/backend  (cùng endpoint, khác connection)
```

## Services

| Service | Port | Vai trò |
|---|---|---|
| `utm-backend` | 8080 | Business logic: deviation check, xử lý lệnh điều khiển |
| `utm-alert` | 8081 | Gateway WebSocket: bridge giữa backend và FE |
| `utm-frontend` | 3000 | Giao diện hiển thị bản đồ và cảnh báo |
| `postgres` | 5432 | Lưu flight plan, user, aircraft |
| `mongo` | 27017 | Lưu lịch sử telemetry, command logs |

---

## Luồng hoạt động

### 1. Khởi động hệ thống

```
utm-alert khởi động trước
    │  lắng nghe :8081/ws/backend  (cho backend + telemetry service)
    │  lắng nghe :8081/ws/monitor  (cho FE)
    │
utm-backend khởi động
    │  AlertClient.Run() → dial ws://utm-alert:8081/ws/backend
    │  Nếu utm-alert chưa sẵn sàng → retry mỗi 5s
    │
FE mở tab
    │  ws = new WebSocket("ws://utm-alert:8081/ws/monitor")
```

---

### 2. Luồng cảnh báo lệch làn (backend → FE)

```
Drone
 │  POST /api/v1/telemetry/push  {lat, lng, altitude_m, auth_request_id, ...}
 ▼
utm-backend
 ├─[goroutine]─► MongoDB Insert  (lưu lịch sử bay)
 │
 └─[goroutine]─► checkAndBroadcastAlert()
                     │
                     ├─ isMuted? → bỏ qua nếu đang tắt
                     ├─ Postgres: lấy AuthorizationRequest (FlightAreaJSON, PlannedAltM)
                     ├─ pointInPolygon(lat, lng)  ← ray casting
                     ├─ haversineM()              ← tính khoảng cách ra biên
                     │
                     ├─ Trong vùng + altitude OK → không làm gì
                     │
                     └─ VI PHẠM → tạo LaneAlert
                                      │
                              sendToAlert(alert)
                              → AlertClient.Send()
                              → WS frame → utm-alert
                              → BackendConn.readPump()
                              → FeHub broadcast
                              → writePump() × N FE clients
                                      │
                                      ▼
                                     FE nhận:
                             {
                               "type": "lane_alert",
                               "severity": "warning",
                               "deviation_m": 45.2,
                               "aircraft_id": "drone-001",
                               "message": "Drone bay ngoài vùng 45m",
                               "timestamp": "..."
                             }
```

**Mức độ cảnh báo:**

| Severity | Điều kiện |
|---|---|
| `normal` | Trong vùng, altitude OK |
| `warning` | Lệch ra ngoài 0–50m |
| `danger` | Lệch > 50m hoặc vượt altitude |

---

### 3. Luồng telemetry real-time (service khác → FE)

Telemetry không đi qua utm-backend. Service riêng kết nối thẳng vào utm-alert:

```
Telemetry Service
 │  WS connect ws://utm-alert:8081/ws/backend
 │  gửi: {"type":"telemetry","payload":{lat, lng, speed_kmh, ...}}
 ▼
utm-alert → BackendConn.readPump() → FeHub broadcast → FE
```

FE dùng để vẽ vị trí drone lên bản đồ real-time.

---

### 4. Luồng điều khiển (FE → backend)

```
FE
 │  ws.send({
 │    "type": "control",
 │    "action": "mute_alert",
 │    "params": {"duration_ms": 5000}
 │  })
 ▼
utm-alert → feClient.readPump() → br.ToBackend
         → BackendConn.writePump() → WS → utm-backend
         → AlertClient.readPump() → handleRawCommand()
 │
 ├─ cập nhật ControlConfig (muteUntil / sensitivity)
 │
 └─[goroutine]─► publishAndLog()
                     ├─ MQTT Publish "utm/control/commands"
                     └─ MongoDB Insert ControlCommandLog
                            {action, params, received_at, published_at}
```

**Các lệnh điều khiển:**

| Action | Params | Tác dụng |
|---|---|---|
| `mute_alert` | `{"duration_ms": 5000}` | Tắt cảnh báo trong N ms |
| `set_sensitivity` | `{"value": 0.75}` | Điều chỉnh ngưỡng phát hiện (0.0–1.0) |
| `acknowledge` | — | Xác nhận đã đọc cảnh báo |

---

### 5. Message format WebSocket

**Backend → FE (qua utm-alert):**
```json
{
  "type": "lane_alert",
  "aircraft_id": "drone-001",
  "session_id": "sess-abc",
  "severity": "warning",
  "deviation_m": 45.2,
  "lat": 10.7769,
  "lng": 106.7009,
  "altitude_m": 120.5,
  "message": "Drone bay ngoài vùng cấp phép 45m",
  "timestamp": "2026-04-20T10:00:00Z"
}
```

**FE → Backend (qua utm-alert):**
```json
{
  "type": "control",
  "action": "mute_alert",
  "params": { "duration_ms": 5000 }
}
```

---

## Cấu trúc thư mục

```
utm/
├── backend/          # Go service — business logic
│   ├── cmd/server/
│   ├── internal/
│   │   ├── hapi/         # HTTP handlers
│   │   ├── service/      # Business logic (deviation check, command handling)
│   │   ├── ws/           # WebSocket alert client
│   │   ├── mq/           # MQTT client
│   │   ├── model/        # Data models
│   │   ├── repository/   # MongoDB + Postgres repositories
│   │   └── config/
│   └── Dockerfile
├── alert/            # Go service — WebSocket gateway
│   ├── cmd/server/
│   ├── internal/
│   │   ├── ws/           # FeHub + BackendConn
│   │   ├── bridge/       # Channel bridge
│   │   └── config/
│   └── Dockerfile
├── frontend/         # Next.js
├── docker-compose.yml
└── README.md
```

## Chạy local

```bash
cp .env.example .env
docker compose up --build
```

Biến môi trường cần thiết (`.env`):

```env
POSTGRES_USER=utm_user
POSTGRES_PASSWORD=utm_pass
POSTGRES_DB=utm_db
MONGO_USER=utm_user
MONGO_PASSWORD=utm_pass
MONGO_DB=utm_db
JWT_SECRET=your-secret
ALERT_SERVICE_URL=ws://utm-alert:8081/ws/backend

# Tuỳ chọn
MQTT_BROKER=tcp://localhost:1883
MQTT_CLIENT_ID=utm-backend
```
