import { YOLOGraph } from './graph.js';
import { Sidebar } from './sidebar.js';
import { WebSocketClient } from './websocket.js';

class App {
    constructor() {
        this.initComponents();
        this.initEvents();
        this.loadData();
    }

    initComponents() {
        // Initialize graph
        const graphContainer = document.getElementById('graph-container');
        this.graph = new YOLOGraph(graphContainer);

        // Initialize sidebar
        const sidebarContainer = document.getElementById('sidebar');
        this.sidebar = new Sidebar(sidebarContainer);

        // Initialize WebSocket
        this.ws = new WebSocketClient();

        // Initialize controls
        this.initControls();
    }

    initControls() {
        // Sidebar toggle
        document.getElementById('toggle-sidebar').addEventListener('click', () => {
            this.sidebar.toggle();
        });

        // Zoom controls
        document.getElementById('zoom-in').addEventListener('click', () => {
            this.graph.controls.dollyIn(1.5);
        });

        document.getElementById('zoom-out').addEventListener('click', () => {
            this.graph.controls.dollyOut(1.5);
        });

        document.getElementById('reset-view').addEventListener('click', () => {
            this.graph.camera.position.set(0, 0, 200);
            this.graph.camera.lookAt(0, 0, 0);
            this.graph.controls.reset();
        });

        // Filter buttons
        const filterButtons = document.querySelectorAll('[data-filter]');
        filterButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                filterButtons.forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
                this.applyFilter(button.dataset.filter);
            });
        });

        // Keyboard shortcuts
        window.addEventListener('keydown', (e) => {
            switch (e.key) {
                case 'Escape':
                    this.sidebar.hide();
                    this.graph.selectNode(null);
                    break;
                case 'f':
                    if (e.ctrlKey || e.metaKey) {
                        e.preventDefault();
                        document.getElementById('search').focus();
                    }
                    break;
            }
        });
    }

    initEvents() {
        // Graph events
        this.graph.container.addEventListener('nodeSelect', (e) => {
            if (e.detail) {
                this.sidebar.showNode(e.detail);
                this.sidebar.show();
            } else {
                this.sidebar.hide();
            }
        });

        // Sidebar events
        this.sidebar.container.addEventListener('search', (e) => {
            this.handleSearch(e.detail.searchTerm);
        });

        // WebSocket events
        this.ws.addEventListener('update', (e) => {
            this.handleDataUpdate(e.detail);
        });
    }

    async loadData() {
        try {
            // Show loading overlay
            const loadingOverlay = document.getElementById('loading');
            loadingOverlay.classList.remove('hidden');

            const [nodesResponse, linksResponse] = await Promise.all([
                fetch('/api/nodes'),
                fetch('/api/links')
            ]);

            if (!nodesResponse.ok || !linksResponse.ok) {
                throw new Error('Failed to fetch data');
            }

            const nodes = await nodesResponse.json();
            const links = await linksResponse.json();

            // Update graph with data
            this.graph.updateData(nodes, links);

            // Hide loading overlay after graph is updated
            loadingOverlay.classList.add('hidden');
        } catch (error) {
            console.error('Error loading data:', error);
            // Show error message
            const loadingOverlay = document.getElementById('loading');
            loadingOverlay.innerHTML = `
                <div class="text-center">
                    <p class="text-red-500">Error loading data. Please try again.</p>
                    <button onclick="window.location.reload()" class="mt-4 px-4 py-2 bg-red-600 rounded hover:bg-red-500">
                        Reload
                    </button>
                </div>
            `;
        }
    }

    applyFilter(filter) {
        if (filter === 'all') {
            this.graph.nodeObjects.forEach(({ mesh }) => {
                mesh.visible = true;
            });
        } else {
            this.graph.nodeObjects.forEach(({ node, mesh }) => {
                mesh.visible = node.type === filter;
            });
        }

        // Update link visibility
        this.graph.linkObjects.forEach(({ line, source, target }) => {
            const sourceNode = this.graph.nodeObjects.get(source);
            const targetNode = this.graph.nodeObjects.get(target);
            line.visible = filter === 'all' || 
                (sourceNode?.mesh.visible && targetNode?.mesh.visible);
        });
    }

    handleSearch(searchTerm) {
        const term = searchTerm.toLowerCase();
        this.graph.nodeObjects.forEach(({ node, mesh }) => {
            const matches = node.title.toLowerCase().includes(term) ||
                          node.content.toLowerCase().includes(term) ||
                          node.id.toLowerCase().includes(term);
            mesh.visible = matches;
        });

        // Update link visibility based on visible nodes
        this.graph.linkObjects.forEach(({ line, source, target }) => {
            const sourceNode = this.graph.nodeObjects.get(source);
            const targetNode = this.graph.nodeObjects.get(target);
            line.visible = sourceNode?.mesh.visible && targetNode?.mesh.visible;
        });
    }

    handleDataUpdate(data) {
        // Handle real-time updates from WebSocket
        if (data.nodes || data.links) {
            this.graph.updateData(data.nodes || [], data.links || []);
        }
    }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.app = new App();
}); 