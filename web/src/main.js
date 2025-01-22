import '@/css/style.css'
import { App } from '@/js/app.js'

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    // Setup UI components
    new App();
    
    // Hide loading screen
    const loading = document.getElementById('loading')
    loading.classList.add('hidden')
})
