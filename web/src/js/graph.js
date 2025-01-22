import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import * as d3 from 'd3';

export class YOLOGraph {
    constructor() {
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        this.renderer = new THREE.WebGLRenderer({ antialias: true });
        this.controls = new OrbitControls(this.camera, this.renderer.domElement);
        
        this.init();
        this.setupMaterials();
        this.setupSimulation();
        this.setupEvents();
        this.animate();
    }

    init() {
        const container = document.getElementById('graph');
        this.renderer.setSize(container.clientWidth, container.clientHeight);
        this.renderer.setClearColor(0x111827); // Match Tailwind's gray-900
        container.appendChild(this.renderer.domElement);
        
        this.camera.position.z = 500;
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.05;
    }

    setupMaterials() {
        this.materials = {
            epic: new THREE.MeshPhongMaterial({ color: 0x3B82F6 }), // blue-500
            feature: new THREE.MeshPhongMaterial({ color: 0x10B981 }), // emerald-500
            task: new THREE.MeshPhongMaterial({ color: 0xF59E0B }), // amber-500
            line: new THREE.LineBasicMaterial({ color: 0x4B5563 }) // gray-600
        };

        // Add lighting
        const light = new THREE.DirectionalLight(0xffffff, 1);
        light.position.set(1, 1, 1);
        this.scene.add(light);
        
        const ambientLight = new THREE.AmbientLight(0x404040);
        this.scene.add(ambientLight);
    }

    setupSimulation() {
        this.simulation = d3.forceSimulation()
            .force('link', d3.forceLink().id(d => d.id).distance(100))
            .force('charge', d3.forceManyBody().strength(-300))
            .force('center', d3.forceCenter(0, 0))
            .on('tick', () => this.updatePositions());
    }

    setupEvents() {
        window.addEventListener('resize', () => this.onWindowResize());
        
        // Raycaster for node selection
        this.raycaster = new THREE.Raycaster();
        this.mouse = new THREE.Vector2();
        
        this.renderer.domElement.addEventListener('click', (event) => this.onNodeClick(event));
        this.renderer.domElement.addEventListener('mousemove', (event) => this.onMouseMove(event));
    }

    updateData(nodes, links) {
        // Clear existing nodes and links
        this.nodes?.forEach(node => this.scene.remove(node.mesh));
        this.links?.forEach(link => this.scene.remove(link.line));
        
        this.nodes = nodes.map(node => {
            let geometry;
            switch(node.type) {
                case 'epic':
                    geometry = new THREE.SphereGeometry(10);
                    break;
                case 'feature':
                    geometry = new THREE.BoxGeometry(15, 15, 15);
                    break;
                case 'task':
                    geometry = new THREE.ConeGeometry(8, 20, 16);
                    break;
                default:
                    geometry = new THREE.SphereGeometry(8);
            }
            
            const mesh = new THREE.Mesh(geometry, this.materials[node.type]);
            mesh.userData = node;
            this.scene.add(mesh);
            return { ...node, mesh };
        });

        this.links = links.map(link => {
            const geometry = new THREE.BufferGeometry();
            const line = new THREE.Line(geometry, this.materials.line);
            this.scene.add(line);
            return { ...link, line };
        });

        this.simulation.nodes(this.nodes);
        this.simulation.force('link').links(this.links);
        this.simulation.alpha(1).restart();
    }

    updatePositions() {
        this.nodes?.forEach(node => {
            node.mesh.position.set(node.x, node.y, 0);
        });

        this.links?.forEach(link => {
            const positions = new Float32Array([
                link.source.x, link.source.y, 0,
                link.target.x, link.target.y, 0
            ]);
            link.line.geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
            link.line.geometry.attributes.position.needsUpdate = true;
        });
    }

    onWindowResize() {
        const container = document.getElementById('graph');
        this.camera.aspect = container.clientWidth / container.clientHeight;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(container.clientWidth, container.clientHeight);
    }

    onNodeClick(event) {
        event.preventDefault();
        this.mouse.x = (event.clientX / this.renderer.domElement.clientWidth) * 2 - 1;
        this.mouse.y = -(event.clientY / this.renderer.domElement.clientHeight) * 2 + 1;
        
        this.raycaster.setFromCamera(this.mouse, this.camera);
        const intersects = this.raycaster.intersectObjects(this.scene.children);
        
        if (intersects.length > 0) {
            const node = intersects[0].object.userData;
            if (node) {
                const event = new CustomEvent('nodeSelected', { detail: node });
                window.dispatchEvent(event);
            }
        }
    }

    onMouseMove(event) {
        event.preventDefault();
        this.mouse.x = (event.clientX / this.renderer.domElement.clientWidth) * 2 - 1;
        this.mouse.y = -(event.clientY / this.renderer.domElement.clientHeight) * 2 + 1;
    }

    animate() {
        requestAnimationFrame(() => this.animate());
        this.controls.update();
        this.renderer.render(this.scene, this.camera);
    }
} 