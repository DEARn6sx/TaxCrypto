package routes

import (
	"extaxcrypto/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/transactions", controllers.AddTransaction)
	e.GET("/profit-loss", controllers.CalculateProfit)
	e.GET("/averages", controllers.CalculateAveragePrice)
}
