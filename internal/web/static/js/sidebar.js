export class Sidebar {
    constructor() {
        this.sidebar = document.getElementById('sidebar');
        this.searchInput = document.getElementById('search');
        this.nodeDetails = document.getElementById('node-details');
        this.isOpen = false;
        
        this.setupEvents();
    }

    setupEvents() {
        window.addEventListener('nodeSelected', (event) => {
            this.showNode(event.detail);
            if (!this.isOpen) this.show();
        });

        this.searchInput.addEventListener('input', (event) => this.onSearch(event));
    }

    show() {
        this.sidebar.classList.remove('-translate-x-full');
        this.isOpen = true;
    }

    hide() {
        this.sidebar.classList.add('-translate-x-full');
        this.isOpen = false;
    }

    toggle() {
        if (this.isOpen) {
            this.hide();
        } else {
            this.show();
        }
    }

    showNode(node) {
        if (!node) {
            this.nodeDetails.innerHTML = '<p class="text-gray-400">No node selected</p>';
            return;
        }

        const statusColors = {
            'todo': 'bg-yellow-500',
            'in_progress': 'bg-blue-500',
            'done': 'bg-green-500',
            'blocked': 'bg-red-500'
        };

        const typeColors = {
            'epic': 'bg-blue-500',
            'feature': 'bg-emerald-500',
            'task': 'bg-amber-500'
        };

        this.nodeDetails.innerHTML = `
            <div class="bg-gray-800 rounded-lg p-4 mb-4">
                <div class="flex items-center justify-between mb-2">
                    <span class="px-2 py-1 text-xs rounded ${typeColors[node.type] || 'bg-gray-500'}">${node.type}</span>
                    <span class="px-2 py-1 text-xs rounded ${statusColors[node.status] || 'bg-gray-500'}">${node.status}</span>
                </div>
                <h2 class="text-xl font-bold mb-2">${node.title}</h2>
                <p class="text-gray-300 mb-4">${node.description || 'No description available'}</p>
                ${this.renderRelationships(node)}
            </div>
        `;
    }

    renderRelationships(node) {
        if (!node.relationships || node.relationships.length === 0) {
            return '<p class="text-gray-400">No relationships</p>';
        }

        return `
            <div class="space-y-2">
                <h3 class="text-sm font-semibold text-gray-400">Relationships</h3>
                <ul class="space-y-1">
                    ${node.relationships.map(rel => `
                        <li class="flex items-center text-sm">
                            <span class="w-16 text-gray-400">${rel.type}:</span>
                            <span class="text-blue-400">${rel.target}</span>
                        </li>
                    `).join('')}
                </ul>
            </div>
        `;
    }

    onSearch(event) {
        const query = event.target.value.toLowerCase();
        const searchEvent = new CustomEvent('search', { 
            detail: { query } 
        });
        window.dispatchEvent(searchEvent);
    }
} 