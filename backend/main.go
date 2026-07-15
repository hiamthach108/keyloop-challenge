package main

import (
	"github.com/hiamthach108/dreon-sdk/logger"
	"github.com/hiamthach108/keyloop-challenge/backend/config"
	_ "github.com/hiamthach108/keyloop-challenge/backend/docs"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/repository"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/service"
	"github.com/hiamthach108/keyloop-challenge/backend/pkg/database"
	"github.com/hiamthach108/keyloop-challenge/backend/presentation/http"
	"github.com/hiamthach108/keyloop-challenge/backend/presentation/http/handler"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// @title Keyloop Intelligent Inventory API
// @version 1.0
// @description Scenario B API for dealership inventory, aging-stock actions, and stock movement history.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	app := fx.New(
		fx.WithLogger(func(appLogger logger.ILogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: appLogger.GetZapLogger()}
		}),
		fx.Provide(
			// Core
			config.NewAppConfig,
			newAppLogger,
			database.NewDbClient,
			http.NewHttpServer,

			// Repositories
			repository.NewDealershipRepository,
			repository.NewVehicleRepository,
			repository.NewInventoryStockRepository,
			repository.NewInventoryActionRepository,
			repository.NewStockMovementRepository,

			// Services
			service.NewInventorySvc,

			// Handlers
			handler.NewInventoryHandler,
		),
		fx.Invoke(http.RegisterHooks),
	)

	app.Run()
}

func newAppLogger(config *config.AppConfig) (logger.ILogger, error) {
	return logger.NewLogger(logger.Config{
		Service: config.App.Name,
		Level:   config.Logger.Level,
	})
}
