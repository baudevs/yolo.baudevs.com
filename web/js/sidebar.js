export class Sidebar {
    constructor(container) {
        this.container = container;
        this.visible = false;
        this.currentNode = null;
        this.searchInput = container.querySelector('#search');
        this.detailsContainer = container.querySelector('#node-details');

        this.initEvents();
    }

    initEvents() {
        this.searchInput.addEventListener('input', (e) => this.onSearch(e));
    }

    show() {
        this.container.classList.remove('hidden');
        this.visible = true;
    }

    hide() {
        this.container.classList.add('hidden');
        this.visible = false;
    }

    toggle() {
        if (this.visible) {
            this.hide();
        } else {
            this.show();
        }
    }

    showNode(node) {
        this.currentNode = node;
        this.show();
        this.render();
    }

    render() {
        if (!this.currentNode) {
            this.detailsContainer.innerHTML = '<p class="text-gray-400">No node selected</p>';
            return;
        }

        const statusColor = this.getStatusColor(this.currentNode.status);
        const typeColor = this.getTypeColor(this.currentNode.type);

        this.detailsContainer.innerHTML = `
            <div class="mb-6">
                <h2 class="text-xl font-bold mb-2">${this.currentNode.title}</h2>
                <div class="flex gap-2 mb-4">
                    <span class="px-2 py-1 rounded text-sm ${typeColor}">
                        ${this.currentNode.type}
                    </span>
                    <span class="px-2 py-1 rounded text-sm ${statusColor}">
                        ${this.currentNode.status}
                    </span>
                </div>
                <div class="text-sm text-gray-400">
                    Version: ${this.currentNode.version}
                </div>
                <div class="text-sm text-gray-400">
                    Last modified: ${this.formatDate(this.currentNode.modified)}
                </div>
            </div>

            <div class="prose prose-invert max-w-none">
                ${this.renderMarkdown(this.currentNode.content)}
            </div>

            ${this.renderRelationships()}
        `;
    }

    renderRelationships() {
        if (!this.currentNode.links || this.currentNode.links.length === 0) {
            return '';
        }

        return `
            <div class="mt-6">
                <h3 class="text-lg font-semibold mb-2">Relationships</h3>
                <ul class="space-y-2">
                    ${this.currentNode.links.map(link => `
                        <li class="flex items-center gap-2">
                            <span class="w-2 h-2 rounded-full bg-blue-500"></span>
                            <a href="#" class="text-blue-400 hover:text-blue-300">${link}</a>
                        </li>
                    `).join('')}
                </ul>
            </div>
        `;
    }

    renderMarkdown(content) {
        // TODO: Implement markdown rendering
        // For now, just return the content with basic formatting
        return content
            .replace(/\n\n/g, '</p><p>')
            .replace(/\n/g, '<br>')
            .replace(/^(.+)$/gm, '<p>$1</p>');
    }

    getStatusColor(status) {
        switch (status.toLowerCase()) {
            case 'active': return 'bg-green-900 text-green-300';
            case 'planned': return 'bg-yellow-900 text-yellow-300';
            case 'completed': return 'bg-blue-900 text-blue-300';
            case 'deprecated': return 'bg-red-900 text-red-300';
            default: return 'bg-gray-900 text-gray-300';
        }
    }

    getTypeColor(type) {
        switch (type.toLowerCase()) {
            case 'epic': return 'bg-blue-900 text-blue-300';
            case 'feature': return 'bg-green-900 text-green-300';
            case 'task': return 'bg-yellow-900 text-yellow-300';
            default: return 'bg-gray-900 text-gray-300';
        }
    }

    formatDate(date) {
        return new Date(date).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    onSearch(event) {
        const searchTerm = event.target.value.toLowerCase();
        // TODO: Implement search functionality
        // This should emit an event that the main app can listen to
        this.container.dispatchEvent(new CustomEvent('search', {
            detail: { searchTerm }
        }));
    }
} 