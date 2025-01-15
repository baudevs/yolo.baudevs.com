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
            if (e.key === 'Escape') {
                this.sidebar.hide();
            }
        });
    }

    applyFilter(filter) {
        // Emit filter change event
        this.graph.container.dispatchEvent(new CustomEvent('filterChange', {
            detail: { filter }
        }));
    }
} 