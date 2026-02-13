package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	router := gin.Default()

	// Limit request body size to 1MB to prevent resource exhaustion
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})

	api.SetupRoutes(router, engine)
	api.SetupRoutes(router, engine)

	// Serve static files for web dashboard
	router.Static("/static", "./web/static")
	log.Printf("DataFlow server starting on port %s", port)
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "DataFlow Pipeline Dashboard",
		})
	})

	log.Printf("DataFlow server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}