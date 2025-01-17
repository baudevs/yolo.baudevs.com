export class Controls {
    constructor(graph, sidebar) {
        this.graph = graph;
        this.sidebar = sidebar;
        this.initEvents();
    }

    initEvents() {
        // Sidebar toggle
        const toggleButton = document.getElementById('toggle-sidebar');
        toggleButton.addEventListener('click', () => {
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
            this.graph.camera.position.set(0, 0, 500);
            this.graph.camera.lookAt(0, 0, 0);
            this.graph.controls.reset();
        });

        // Keyboard shortcuts
        window.addEventListener('keydown', (event) => {
            // Escape to hide sidebar
            if (event.key === 'Escape') {
                this.sidebar.hide();
            }
            
            // Ctrl/Cmd + F to focus search
            if ((event.ctrlKey || event.metaKey) && event.key === 'f') {
                event.preventDefault();
                document.getElementById('search').focus();
            }
            
            // Ctrl/Cmd + R to reset view
            if ((event.ctrlKey || event.metaKey) && event.key === 'r') {
                event.preventDefault();
                document.getElementById('reset-view').click();
            }
        });
    }
} 