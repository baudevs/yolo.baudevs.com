import { YOLOGraph } from './graph.js';
import { Sidebar } from './sidebar.js';
import { Controls } from './controls.js';
import { WebSocketClient } from './websocket.js';

export class App {
    constructor() {
        this.graph = new YOLOGraph();
        this.sidebar = new Sidebar();
        this.controls = new Controls(this.graph, this.sidebar);
        this.ws = new WebSocketClient();
        
        this.initEvents();
        this.loadData();
    }

    initEvents() {
        // Handle WebSocket messages
        this.ws.addEventListener('message', (event) => {
            const data = event.detail;
            if (data.type === 'update') {
                this.updateData(data.nodes, data.links);
            }
        });

        // Handle search events
        window.addEventListener('search', (event) => {
            const query = event.detail.query.toLowerCase();
            const filteredNodes = this.nodes.filter(node => 
                node.title.toLowerCase().includes(query) ||
                node.description?.toLowerCase().includes(query) ||
                node.tags?.some(tag => tag.toLowerCase().includes(query))
            );
            this.graph.updateData(filteredNodes, this.links);
        });
    }

    async loadData() {
        try {
            const response = await fetch('/api/nodes');
            if (!response.ok) throw new Error('Failed to fetch nodes');
            
            const data = await response.json();
            this.nodes = data.nodes;
            this.links = data.links;
            
            this.updateData(this.nodes, this.links);
            document.getElementById('loading').classList.add('hidden');
        } catch (error) {
            console.error('Error loading data:', error);
            document.getElementById('loading').innerHTML = `
                <div class="text-center">
                    <div class="text-red-500 text-xl mb-4">Failed to load project data</div>
                    <button onclick="window.location.reload()" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded">
                        Retry
                    </button>
                </div>
            `;
        }
    }

    updateData(nodes, links) {
        this.graph.updateData(nodes, links);
    }
}