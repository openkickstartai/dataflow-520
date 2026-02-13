package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/dataflow/internal/pipeline"
)

// SetupRoutes configures API routes
func SetupRoutes(router *gin.Engine, engine *pipeline.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/pipelines", listPipelines(engine))
		api.POST("/pipelines", createPipeline(engine))
		api.GET("/pipelines/:id", getPipeline(engine))
		api.POST("/pipelines/:id/execute", executePipeline(engine))
		api.GET("/health", healthCheck)
	}
}

func listPipelines(engine *pipeline.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		pipelines := engine.ListPipelines()
		c.JSON(http.StatusOK, gin.H{
			"pipelines": pipelines,
			"count":     len(pipelines),
		})
	}
}

func createPipeline(engine *pipeline.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pipeline.Pipeline
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate UUID if not provided; validate format if user-supplied
		if req.ID == "" {
			req.ID = uuid.New().String()
		} else if _, err := uuid.Parse(req.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pipeline ID must be a valid UUID"})
			return
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pipeline name is required"})
			return
		}

		if len(req.Name) > 256 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pipeline name exceeds maximum length of 256 characters"})
			return
		}
			c.JSON(http.StatusBadRequest, gin.H{"error": "pipeline name is required"})
			return
		}

		if err := engine.CreatePipeline(&req); err != nil {
func getPipeline(engine *pipeline.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pipeline ID format"})
			return
		}
		p, err := engine.GetPipeline(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "pipeline not found"})
			return
		}

func executePipeline(engine *pipeline.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pipeline ID format"})
			return
		}
		if err := engine.ExecutePipeline(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "pipeline not found or execution failed"})
			return
		}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "dataflow",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
		id := c.Param("id")
		if err := engine.ExecutePipeline(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pipeline execution started"})
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "dataflow",
		"timestamp": strconv.FormatInt(c.Request.Context().Value("timestamp").(int64), 10),
	})
}