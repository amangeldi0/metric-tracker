package app

import (
	"github.com/amangeldi0/metric-tracker/internal/server/config"
	"github.com/amangeldi0/metric-tracker/internal/server/handlers"
	"github.com/amangeldi0/metric-tracker/internal/server/middlewares"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	sugarLogger := logger.Sugar()

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
	}

	str, err := storage.SetupStorage(sugarLogger)

	if err != nil {
		sugarLogger.Panicf("Failed setup storage: %s", err)
	}

	defer func() {
		if err = sugarLogger.Sync(); err != nil {
			panic(err)
		}

		if err = str.Close(); err != nil {
			panic(err)
		}
	}()

	r := setupRouter(str, sugarLogger)
	if err = r.Run(config.Config.Address); err != nil {
		sugarLogger.Panicf("Failed start server: %s", err)
	}
}

func setupRouter(storage storage.Storage, logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	baseHandler := handlers.NewBase(storage, logger)
	baseMiddleware := middlewares.NewBase(logger)

	r.Use(baseMiddleware.Compress)
	r.Use(baseMiddleware.Logger)

	r.GET("/", baseHandler.Values())

	r.GET("/ping", baseHandler.Ping())

	r.POST("/value", baseHandler.ValueByBody())
	r.POST("/value/", baseHandler.ValueByBody())

	r.GET("/value/:type/:name", baseHandler.ValueByURI())
	r.GET("/value/:type/:name/", baseHandler.ValueByURI())

	r.POST("/update", baseHandler.UpdateByBody())
	r.POST("/update/", baseHandler.UpdateByBody())

	r.POST("/update/:type", baseHandler.UpdateByURI())
	r.POST("/update/:type/", baseHandler.UpdateByURI())

	r.POST("/update/:type/:name/:value", baseHandler.UpdateByURI())
	r.POST("/update/:type/:name/:value/", baseHandler.UpdateByURI())

	r.NoRoute(baseHandler.BadRequest)

	return r
}
