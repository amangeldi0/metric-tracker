package main

import (
	"github.com/amangeldi0/metric-tracker/internal/server/config"
	filestorage "github.com/amangeldi0/metric-tracker/internal/server/filestorage"
	"github.com/amangeldi0/metric-tracker/internal/server/handlers"
	"github.com/amangeldi0/metric-tracker/internal/server/middlewares"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	sugarLogger := logger.Sugar()

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
	}

	memStorage := storage.NewMem()

	var fileStorage *filestorage.Storage
	if config.Config.FileStoragePath != "" {
		fileStorage, err = filestorage.New(&memStorage, sugarLogger)
		if err != nil {
			sugarLogger.Panicf("Failed loading file storage: %s", err)
		}

		if err = fileStorage.Restore(); err != nil {
			sugarLogger.Panicf("Failed to recover data from file: %s", err)
		}
		fileStorage.Start()
	}

	defer func(logger *zap.Logger, fileStorage *filestorage.Storage) {
		if err = logger.Sync(); err != nil {
			panic(err)
		}
		if fileStorage != nil {
			if err = fileStorage.Close(); err != nil {
				panic(err)
			}
		}
	}(logger, fileStorage)

	r := setupRouter(&memStorage, fileStorage, sugarLogger)
	if err = r.Run(config.Config.Address); err != nil {
		sugarLogger.Panicf("Failed start server: %s", err)
	}
}

func setupRouter(storage *storage.Mem, fileStorage *filestorage.Storage, logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	baseHandler := handlers.NewBase(storage, logger)
	baseMiddleware := middlewares.NewBase(logger)

	r.Use(baseMiddleware.Compress)
	r.Use(baseMiddleware.Logger)

	if fileStorage != nil {
		r.Use(fileStorage.GetMiddleware())
	}

	r.GET("/", baseHandler.Values())

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
