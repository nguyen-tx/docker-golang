package hapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/utm/backend/internal/config"
	"github.com/utm/backend/internal/middleware"
	"github.com/utm/backend/internal/mq"
	mongoRepo "github.com/utm/backend/internal/repository/mongo"
	pgRepo "github.com/utm/backend/internal/repository/postgres"
	scrpclient "github.com/utm/backend/internal/scrp"
	"github.com/utm/backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(cfg *config.Config, pgPool *pgxpool.Pool, mongoDB *mongo.Database, scrpClient *scrpclient.Client) *gin.Engine {
	h := buildHandlers(cfg, pgPool, mongoDB, scrpClient)
	return setupEngine(cfg, h)
}

func buildHandlers(cfg *config.Config, pgPool *pgxpool.Pool, mongoDB *mongo.Database, scrpClient *scrpclient.Client) *handlers {
	userRepo := pgRepo.NewUserRepository(pgPool)
	aircraftRepo := pgRepo.NewAircraftRepository(pgPool)
	authzRepo := pgRepo.NewAuthorizationRepository(pgPool)
	airspaceRepo := pgRepo.NewAirspaceRepository(pgPool)
	telemetryRepo := mongoRepo.NewTelemetryRepository(mongoDB)
	cmdLogRepo := mongoRepo.NewCommandLogRepository(mongoDB)
	notifRepo := mongoRepo.NewNotificationRepository(mongoDB)
	fnRepo := pgRepo.NewFlightNotificationRepository(pgPool)
	laneRepo := pgRepo.NewLaneRepository(pgPool)
	waypointRepo := pgRepo.NewWaypointRepository(pgPool)
	vertiportRepo := pgRepo.NewVertiportRepository(pgPool)

	var mqttClient *mq.MQTTClient
	if cfg.MQTT.Broker != "" {
		var err error
		mqttClient, err = mq.NewMQTTClient(cfg.MQTT)
		if err != nil {
			log.Warn().Err(err).Msg("mqtt init failed, running without mqtt")
		}
	}

	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userSvc := service.NewUserService(userRepo)
	aircraftSvc := service.NewAircraftService(aircraftRepo)
	authzSvc := service.NewAuthorizationService(authzRepo)
	airspaceSvc := service.NewAirspaceService(airspaceRepo)
	telemetrySvc := service.NewTelemetryService(
		telemetryRepo, authzRepo, cmdLogRepo, mqttClient,
		cfg.AlertService.URL,
	)
	go telemetrySvc.AlertClient().Run()
	notifSvc := service.NewNotificationService(notifRepo)
	flightNotifSvc := service.NewFlightNotificationService(fnRepo, laneRepo, waypointRepo, vertiportRepo, scrpClient)

	return &handlers{
		auth:               NewAuthHandler(authSvc),
		user:               NewUserHandler(userSvc),
		aircraft:           NewAircraftHandler(aircraftSvc),
		airspace:           NewAirspaceHandler(airspaceSvc),
		authorization:      NewAuthorizationHandler(authzSvc),
		telemetry:          NewTelemetryHandler(telemetrySvc),
		notification:       NewNotificationHandler(notifSvc),
		flightNotification: NewFlightNotificationHandler(flightNotifSvc),
	}
}

func setupEngine(cfg *config.Config, h *handlers) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "utm-backend"})
	})
	public := r.Group("/api/v1")
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth(cfg.JWT.Secret))
	h.auth.RegisterRoutes(public, protected)
	h.user.RegisterRoutes(protected)
	h.aircraft.RegisterRoutes(protected)
	h.airspace.RegisterRoutes(public, protected)
	h.authorization.RegisterRoutes(protected)
	h.telemetry.RegisterRoutes(protected)
	h.notification.RegisterRoutes(protected)
	h.flightNotification.RegisterRoutes(protected)
	return r
}

type handlers struct {
	auth               *AuthHandler
	user               *UserHandler
	aircraft           *AircraftHandler
	airspace           *AirspaceHandler
	authorization      *AuthorizationHandler
	telemetry          *TelemetryHandler
	notification       *NotificationHandler
	flightNotification *FlightNotificationHandler
}
