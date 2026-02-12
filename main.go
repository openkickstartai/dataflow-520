package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/dataflow/internal/api"
	"github.com/dataflow/internal/pipeline"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize pipeline engine
	engine := pipeline.NewEngine()
	go engine.Start()

	// Setup HTTP router
	router := gin.Default()
	api.SetupRoutes(router, engine)

	// Serve static files for web dashboard
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "DataFlow Pipeline Dashboard",
		})
	})

	log.Printf("DataFlow server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}