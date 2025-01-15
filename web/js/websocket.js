export class WebSocketClient extends EventTarget {
    constructor() {
        super();
        this.connect();
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws`;

        this.ws = new WebSocket(wsUrl);
        this.ws.onopen = () => this.onOpen();
        this.ws.onclose = () => this.onClose();
        this.ws.onmessage = (e) => this.onMessage(e);
        this.ws.onerror = (e) => this.onError(e);
    }

    onOpen() {
        console.log('WebSocket connection established');
        this.startHeartbeat();
    }

    onClose() {
        console.log('WebSocket connection closed');
        this.stopHeartbeat();
        // Attempt to reconnect after 5 seconds
        setTimeout(() => this.connect(), 5000);
    }

    onMessage(event) {
        try {
            const data = JSON.parse(event.data);
            
            switch (data.type) {
                case 'update':
                    this.dispatchEvent(new CustomEvent('update', {
                        detail: data.payload
                    }));
                    break;

                case 'pong':
                    // Handle heartbeat response
                    break;

                default:
                    console.warn('Unknown message type:', data.type);
            }
        } catch (error) {
            console.error('Error processing message:', error);
        }
    }

    onError(error) {
        console.error('WebSocket error:', error);
    }

    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            if (this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000); // Send heartbeat every 30 seconds
    }

    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

    send(data) {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('WebSocket is not open');
        }
    }
} 