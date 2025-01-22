export class WebSocketClient extends EventTarget {
    constructor(wsUrl = `ws://${window.location.host}/ws`) {
        super();
        this.connect(wsUrl);
        this.startHeartbeat();
    }

    connect(wsUrl) {
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            console.log('WebSocket connection established');
            this.dispatchEvent(new Event('connected'));
        };

        this.ws.onclose = () => {
            console.log('WebSocket connection closed');
            this.dispatchEvent(new Event('disconnected'));
            // Try to reconnect after 5 seconds
            setTimeout(() => this.connect(wsUrl), 5000);
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.dispatchEvent(new Event('error'));
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type === 'pong') return; // Ignore heartbeat responses
                
                this.dispatchEvent(new CustomEvent('message', {
                    detail: data
                }));
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };
    }

    startHeartbeat() {
        // Send a ping every 30 seconds to keep the connection alive
        setInterval(() => {
            if (this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000);
    }

    send(data) {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('WebSocket is not connected');
        }
    }
} 