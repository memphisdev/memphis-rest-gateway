package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"rest-gateway/conf"
	"rest-gateway/handlers"
)

func InitilizeMonitoringRoutes(app *fiber.App) {
	configuration := conf.GetConfig()

	monitoringHandlerHandler := handlers.MonitoringHandler{}
	api := app.Group("/monitoring", logger.New())
	api.Get("/status", monitoringHandlerHandler.Status)
	if configuration.DEV_ENV == "true" {
		api.Get("/getResourcesUtilization", monitoringHandlerHandler.GetResourcesUtilization)
	}
}
