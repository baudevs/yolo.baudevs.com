class ShortcutsApp {
    constructor() {
        this.shortcuts = [];
        this.recording = false;
        this.currentKeys = [];
        this.editingId = null;
        this.ws = null;

        // Key display mapping for special keys
        this.keyDisplayMap = {
            'ArrowUp': '↑',
            'ArrowDown': '↓',
            'ArrowLeft': '←',
            'ArrowRight': '→',
            'Enter': '↵',
            'Escape': 'Esc',
            'Backspace': '⌫',
            'Delete': '⌦',
            'Tab': '⇥',
            'CapsLock': '⇪',
            'Space': 'Space'
        };

        this.initElements();
        this.initEvents();
        this.initWebSocket();
        this.loadShortcuts();
    }

    initElements() {
        // Buttons
        this.addButton = document.getElementById('addShortcut');
        this.cancelButton = document.getElementById('cancelShortcut');
        
        // Modal
        this.modal = document.getElementById('shortcutModal');
        this.form = document.getElementById('shortcutForm');
        
        // Form inputs
        this.nameInput = document.getElementById('shortcutName');
        this.keysInput = document.getElementById('shortcutKeys');
        this.keysDisplay = document.getElementById('keysDisplay');
        this.commandInput = document.getElementById('shortcutCommand');
        this.descriptionInput = document.getElementById('shortcutDescription');
        
        // Lists
        this.shortcutsList = document.getElementById('shortcutsList');
        
        // Status
        this.statusMessage = document.getElementById('statusMessage');
    }

    initEvents() {
        // Add shortcut button
        this.addButton.addEventListener('click', () => this.showModal());
        
        // Cancel button
        this.cancelButton.addEventListener('click', () => this.hideModal());
        
        // Form submission
        this.form.addEventListener('submit', (e) => {
            e.preventDefault();
            this.saveShortcut();
        });
        
        // Keys recording
        this.keysDisplay.addEventListener('click', () => this.startRecording());
        
        // Global key events for recording
        document.addEventListener('keydown', (e) => {
            if (this.recording) {
                e.preventDefault();
                this.handleKeyDown(e);
            }
        });
        
        // Global key events for stopping recording
        document.addEventListener('keyup', (e) => {
            if (this.recording) {
                e.preventDefault();
                this.handleKeyUp(e);
            }
        });
    }

    initWebSocket() {
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);
        
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.showStatus('Connected to server', 'green');
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket disconnected');
            this.showStatus('Disconnected from server', 'red');
            setTimeout(() => this.initWebSocket(), 5000);
        };
        
        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleWebSocketMessage(data);
        };
    }

    async loadShortcuts() {
        try {
            const response = await fetch('/api/list');
            if (!response.ok) throw new Error('Failed to load shortcuts');
            const shortcuts = await response.json();
            this.shortcuts = shortcuts;
            this.renderShortcuts();
        } catch (error) {
            console.error('Failed to load shortcuts:', error);
            this.showStatus('Failed to load shortcuts', 'red');
        }
    }

    renderShortcuts() {
        this.shortcutsList.innerHTML = '';
        
        if (this.shortcuts.length === 0) {
            this.shortcutsList.innerHTML = `
                <div class="text-gray-500 text-center py-8">
                    No shortcuts configured yet. Click "Add Shortcut" to create one.
                </div>
            `;
            return;
        }
        
        this.shortcuts.forEach(shortcut => {
            const item = document.createElement('div');
            item.className = `shortcut-item bg-gray-700 rounded-lg p-4 ${shortcut.enabled ? '' : 'disabled'}`;
            item.innerHTML = `
                <div class="flex justify-between items-start">
                    <div>
                        <h3 class="font-semibold mb-1">${shortcut.name}</h3>
                        <div class="text-sm text-gray-400 mb-2">${shortcut.description || 'No description'}</div>
                        <div class="flex flex-wrap gap-2">
                            ${shortcut.keys.map(key => `
                                <span class="key-badge">${key}</span>
                            `).join('')}
                        </div>
                    </div>
                    <div class="flex space-x-2">
                        <button class="toggle-btn bg-gray-600 hover:bg-gray-500 px-3 py-1 rounded-lg text-sm">
                            ${shortcut.enabled ? 'Disable' : 'Enable'}
                        </button>
                        <button class="edit-btn bg-blue-600 hover:bg-blue-700 px-3 py-1 rounded-lg text-sm">
                            Edit
                        </button>
                        <button class="delete-btn bg-red-600 hover:bg-red-700 px-3 py-1 rounded-lg text-sm">
                            Delete
                        </button>
                    </div>
                </div>
            `;
            
            // Add event listeners
            item.querySelector('.toggle-btn').addEventListener('click', () => 
                this.toggleShortcut(shortcut.id));
            item.querySelector('.edit-btn').addEventListener('click', () => 
                this.editShortcut(shortcut));
            item.querySelector('.delete-btn').addEventListener('click', () => 
                this.deleteShortcut(shortcut.id));
            
            this.shortcutsList.appendChild(item);
        });
    }

    showModal(shortcut = null) {
        this.editingId = shortcut ? shortcut.id : null;
        this.modal.classList.remove('hidden');
        this.modal.classList.add('flex');
        
        if (shortcut) {
            this.nameInput.value = shortcut.name;
            this.keysInput.value = shortcut.keys.join('+');
            this.keysDisplay.textContent = shortcut.keys.join(' + ');
            this.commandInput.value = shortcut.command;
            this.descriptionInput.value = shortcut.description || '';
        } else {
            this.form.reset();
            this.keysDisplay.textContent = 'Click to record...';
        }
    }

    hideModal() {
        this.modal.classList.remove('flex');
        this.modal.classList.add('hidden');
        this.editingId = null;
        this.form.reset();
        this.stopRecording();
    }

    startRecording() {
        this.recording = true;
        this.currentKeys = [];
        this.keysInput.value = '';
        this.keysDisplay.textContent = 'Recording...';
        this.keysDisplay.classList.add('recording');
    }

    stopRecording() {
        this.recording = false;
        this.keysDisplay.classList.remove('recording');
    }

    handleKeyDown(e) {
        // Prevent recording modifier keys alone
        if (['Alt', 'Shift', 'Control', 'Meta'].includes(e.key)) {
            return;
        }

        const keys = [];
        if (e.metaKey) keys.push('⌘');
        if (e.ctrlKey) keys.push('⌃');
        if (e.altKey) keys.push('⌥');
        if (e.shiftKey) keys.push('⇧');

        // Get the main key
        let mainKey = e.key;
        
        // Handle dead keys and Option key combinations
        if (mainKey === 'Dead' || mainKey === '∆' || mainKey === 'ˆ' || mainKey === '©') {
            // Use the physical key code instead
            mainKey = e.code.replace('Key', '');
        } else if (mainKey.length === 1) {
            // For letter keys, always use uppercase
            mainKey = mainKey.toUpperCase();
        } else {
            // For special keys, use the mapping if available
            mainKey = this.keyDisplayMap[mainKey] || mainKey;
        }
        
        keys.push(mainKey);

        this.currentKeys = keys;
        this.keysInput.value = keys.join('+');
        this.keysDisplay.textContent = keys.join(' ');
        this.recording = false;
        this.keysDisplay.classList.remove('recording');
    }

    handleKeyUp(e) {
        // Only stop recording if all modifier keys are released
        if (!e.metaKey && !e.ctrlKey && !e.altKey && !e.shiftKey) {
            this.recording = false;
            this.keysDisplay.classList.remove('recording');
        }
    }

    async saveShortcut() {
        const shortcut = {
            id: this.editingId || crypto.randomUUID(),
            name: this.nameInput.value,
            keys: this.keysInput.value.split('+'),
            command: this.commandInput.value,
            description: this.descriptionInput.value,
            enabled: true
        };
        
        try {
            const response = await fetch('/api/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(shortcut)
            });
            
            if (!response.ok) throw new Error('Failed to save shortcut');
            
            this.showStatus('Shortcut saved successfully', 'green');
            this.hideModal();
            await this.loadShortcuts();
        } catch (error) {
            console.error('Failed to save shortcut:', error);
            this.showStatus('Failed to save shortcut', 'red');
        }
    }

    async toggleShortcut(id) {
        const shortcut = this.shortcuts.find(s => s.id === id);
        if (!shortcut) return;
        
        try {
            const response = await fetch('/api/toggle', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id, enabled: !shortcut.enabled })
            });
            
            if (!response.ok) throw new Error('Failed to toggle shortcut');
            
            this.showStatus(`Shortcut ${shortcut.enabled ? 'disabled' : 'enabled'}`, 'green');
            await this.loadShortcuts();
        } catch (error) {
            console.error('Failed to toggle shortcut:', error);
            this.showStatus('Failed to toggle shortcut', 'red');
        }
    }

    async deleteShortcut(id) {
        if (!confirm('Are you sure you want to delete this shortcut?')) return;
        
        try {
            const response = await fetch('/api/delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id })
            });
            
            if (!response.ok) throw new Error('Failed to delete shortcut');
            
            this.showStatus('Shortcut deleted successfully', 'green');
            await this.loadShortcuts();
        } catch (error) {
            console.error('Failed to delete shortcut:', error);
            this.showStatus('Failed to delete shortcut', 'red');
        }
    }

    editShortcut(shortcut) {
        this.showModal(shortcut);
    }

    handleWebSocketMessage(data) {
        switch (data.type) {
            case 'shortcut_triggered':
                this.showStatus(`Shortcut triggered: ${data.shortcut.name}`, 'green');
                break;
            case 'shortcuts_updated':
                this.loadShortcuts();
                break;
            case 'error':
                this.showStatus(data.message, 'red');
                break;
        }
    }

    showStatus(message, type = 'green') {
        this.statusMessage.querySelector('div').textContent = message;
        this.statusMessage.querySelector('div').className = 
            `bg-${type}-600 text-white px-6 py-3 rounded-lg shadow-lg`;
        
        this.statusMessage.classList.remove('hidden');
        this.statusMessage.classList.add('show');
        
        setTimeout(() => {
            this.statusMessage.classList.remove('show');
            this.statusMessage.classList.add('hidden');
        }, 3000);
    }
}

// Initialize the app
document.addEventListener('DOMContentLoaded', () => {
    window.app = new ShortcutsApp();
}); 