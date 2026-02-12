package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Engine manages pipeline execution and monitoring
type Engine struct {
	pipelines map[string]*Pipeline
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// Pipeline represents a data processing pipeline
type Pipeline struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Nodes       []Node                 `json:"nodes"`
	Connections []Connection           `json:"connections"`
	Status      string                 `json:"status"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Node represents a processing step in the pipeline
type Node struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Config   map[string]interface{} `json:"config"`
	Position Position               `json:"position"`
}

// Connection represents data flow between nodes
type Connection struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Position represents node position in visual editor
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// NewEngine creates a new pipeline engine
func NewEngine() *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		pipelines: make(map[string]*Pipeline),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start begins the engine's monitoring loop
func (e *Engine) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.monitorPipelines()
		}
	}
}

// CreatePipeline creates a new pipeline
func (e *Engine) CreatePipeline(pipeline *Pipeline) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	pipeline.CreatedAt = time.Now()
	pipeline.UpdatedAt = time.Now()
	pipeline.Status = "created"

	e.pipelines[pipeline.ID] = pipeline
	log.Printf("Created pipeline: %s", pipeline.Name)
	return nil
}

// GetPipeline retrieves a pipeline by ID
func (e *Engine) GetPipeline(id string) (*Pipeline, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	pipeline, exists := e.pipelines[id]
	if !exists {
		return nil, fmt.Errorf("pipeline not found: %s", id)
	}
	return pipeline, nil
}

// ListPipelines returns all pipelines
func (e *Engine) ListPipelines() []*Pipeline {
	e.mu.RLock()
	defer e.mu.RUnlock()

	pipelines := make([]*Pipeline, 0, len(e.pipelines))
	for _, pipeline := range e.pipelines {
		pipelines = append(pipelines, pipeline)
	}
	return pipelines
}

// ExecutePipeline starts pipeline execution
func (e *Engine) ExecutePipeline(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	pipeline, exists := e.pipelines[id]
	if !exists {
		return fmt.Errorf("pipeline not found: %s", id)
	}

	pipeline.Status = "running"
	pipeline.UpdatedAt = time.Now()

	// Start pipeline execution in goroutine
	go e.runPipeline(pipeline)

	return nil
}

func (e *Engine) runPipeline(pipeline *Pipeline) {
	log.Printf("Executing pipeline: %s", pipeline.Name)
	
	// Simulate pipeline execution
	time.Sleep(2 * time.Second)
	
	e.mu.Lock()
	pipeline.Status = "completed"
	pipeline.UpdatedAt = time.Now()
	e.mu.Unlock()
	
	log.Printf("Pipeline completed: %s", pipeline.Name)
}

func (e *Engine) monitorPipelines() {
	e.mu.RLock()
	runningCount := 0
	for _, pipeline := range e.pipelines {
		if pipeline.Status == "running" {
			runningCount++
		}
	}
	e.mu.RUnlock()
	
	if runningCount > 0 {
		log.Printf("Monitoring: %d pipelines running", runningCount)
	}
}

// Stop gracefully shuts down the engine
func (e *Engine) Stop() {
	e.cancel()
}