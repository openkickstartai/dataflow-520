class PipelineBuilder {
    constructor() {
        this.canvas = document.getElementById('pipeline-canvas');
        this.ctx = this.canvas.getContext('2d');
        this.nodes = [];
        this.connections = [];
        this.selectedNode = null;
        this.dragOffset = { x: 0, y: 0 };
        this.isDragging = false;
        
        this.setupEventListeners();
        this.loadPipelines();
    }

    setupEventListeners() {
        this.canvas.addEventListener('mousedown', this.handleMouseDown.bind(this));
        this.canvas.addEventListener('mousemove', this.handleMouseMove.bind(this));
        this.canvas.addEventListener('mouseup', this.handleMouseUp.bind(this));
        
        document.getElementById('add-source-node').addEventListener('click', () => {
            this.addNode('source', 'Data Source', { x: 100, y: 100 });
        });
        
        document.getElementById('add-transform-node').addEventListener('click', () => {
            this.addNode('transform', 'Transform', { x: 300, y: 100 });
        });
        
        document.getElementById('add-sink-node').addEventListener('click', () => {
            this.addNode('sink', 'Data Sink', { x: 500, y: 100 });
        });
        
        document.getElementById('save-pipeline').addEventListener('click', this.savePipeline.bind(this));
        document.getElementById('execute-pipeline').addEventListener('click', this.executePipeline.bind(this));
    }

    addNode(type, name, position) {
        const node = {
            id: 'node_' + Date.now(),
            type: type,
            name: name,
            position: position,
            config: {}
        };
        
        this.nodes.push(node);
        this.render();
    }

    handleMouseDown(event) {
        const rect = this.canvas.getBoundingClientRect();
        const x = event.clientX - rect.left;
        const y = event.clientY - rect.top;
        
        const clickedNode = this.getNodeAt(x, y);
        if (clickedNode) {
            this.selectedNode = clickedNode;
            this.isDragging = true;
            this.dragOffset = {
                x: x - clickedNode.position.x,
                y: y - clickedNode.position.y
            };
        }
    }

    handleMouseMove(event) {
        if (this.isDragging && this.selectedNode) {
            const rect = this.canvas.getBoundingClientRect();
            const x = event.clientX - rect.left;
            const y = event.clientY - rect.top;
            
            this.selectedNode.position.x = x - this.dragOffset.x;
            this.selectedNode.position.y = y - this.dragOffset.y;
            
            this.render();
        }
    }

    handleMouseUp() {
        this.isDragging = false;
        this.selectedNode = null;
    }

    getNodeAt(x, y) {
        return this.nodes.find(node => {
            const dx = x - node.position.x;
            const dy = y - node.position.y;
            return dx >= 0 && dx <= 120 && dy >= 0 && dy <= 60;
        });
    }

    render() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Draw connections
        this.connections.forEach(conn => {
            const sourceNode = this.nodes.find(n => n.id === conn.source);
            const targetNode = this.nodes.find(n => n.id === conn.target);
            
            if (sourceNode && targetNode) {
                this.ctx.beginPath();
                this.ctx.moveTo(sourceNode.position.x + 120, sourceNode.position.y + 30);
                this.ctx.lineTo(targetNode.position.x, targetNode.position.y + 30);
                this.ctx.strokeStyle = '#007bff';
                this.ctx.lineWidth = 2;
                this.ctx.stroke();
            }
        });
        
        // Draw nodes
        this.nodes.forEach(node => {
            this.drawNode(node);
        });
    }

    drawNode(node) {
        const { x, y } = node.position;
        
        // Node background
        this.ctx.fillStyle = this.getNodeColor(node.type);
        this.ctx.fillRect(x, y, 120, 60);
        
        // Node border
        this.ctx.strokeStyle = '#333';
        this.ctx.lineWidth = 2;
        this.ctx.strokeRect(x, y, 120, 60);
        
        // Node text
        this.ctx.fillStyle = '#fff';
        this.ctx.font = '12px Arial';
        this.ctx.textAlign = 'center';
        this.ctx.fillText(node.name, x + 60, y + 35);
    }

    getNodeColor(type) {
        switch (type) {
            case 'source': return '#28a745';
            case 'transform': return '#007bff';
            case 'sink': return '#dc3545';
            default: return '#6c757d';
        }
    }

    async loadPipelines() {
        try {
            const response = await fetch('/api/v1/pipelines');
            const data = await response.json();
            
            const pipelineList = document.getElementById('pipeline-list');
            pipelineList.innerHTML = '';
            
            data.pipelines.forEach(pipeline => {
                const item = document.createElement('div');
                item.className = 'pipeline-item';
                item.innerHTML = `
                    <h4>${pipeline.name}</h4>
                    <p>Status: <span class="status-${pipeline.status}">${pipeline.status}</span></p>
                    <p>Nodes: ${pipeline.nodes.length}</p>
                `;
                pipelineList.appendChild(item);
            });
        } catch (error) {
            console.error('Failed to load pipelines:', error);
        }
    }

    async savePipeline() {
        const pipelineName = prompt('Enter pipeline name:');
        if (!pipelineName) return;
        
        const pipeline = {
            name: pipelineName,
            nodes: this.nodes,
            connections: this.connections,
            config: {}
        };
        
        try {
            const response = await fetch('/api/v1/pipelines', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(pipeline)
            });
            
            if (response.ok) {
                alert('Pipeline saved successfully!');
                this.loadPipelines();
            } else {
                alert('Failed to save pipeline');
            }
        } catch (error) {
            console.error('Error saving pipeline:', error);
            alert('Error saving pipeline');
        }
    }

    async executePipeline() {
        if (this.nodes.length === 0) {
            alert('Please add nodes to the pipeline first');
            return;
        }
        
        // For demo, execute the first pipeline
        try {
            const response = await fetch('/api/v1/pipelines');
            const data = await response.json();
            
            if (data.pipelines.length > 0) {
                const pipelineId = data.pipelines[0].id;
                const execResponse = await fetch(`/api/v1/pipelines/${pipelineId}/execute`, {
                    method: 'POST'
                });
                
                if (execResponse.ok) {
                    alert('Pipeline execution started!');
                    setTimeout(() => this.loadPipelines(), 1000);
                } else {
                    alert('Failed to execute pipeline');
                }
            }
        } catch (error) {
            console.error('Error executing pipeline:', error);
        }
    }
}

// Initialize the pipeline builder when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new PipelineBuilder();
});