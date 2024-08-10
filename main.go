package main

import (
	"log"
	"os"

	"github.com/HarshitNagpal29/go-ecommerce-cart/controllers"
	"github.com/HarshitNagpal29/go-ecommerce-cart/database"
	"github.com/HarshitNagpal29/go-ecommerce-cart/middleware"
	"github.com/HarshitNagpal29/go-ecommerce-cart/routes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests.",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	log.Fatal(router.Run(":" + port))
}
