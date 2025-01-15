import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';

export class YOLOGraph {
    constructor(container) {
        this.container = container;
        this.nodes = new Map();
        this.links = new Map();
        this.nodeObjects = new Map();
        this.linkObjects = new Map();
        this.selectedNode = null;
        this.hoveredNode = null;
        this.raycaster = new THREE.Raycaster();
        this.mouse = new THREE.Vector2();

        this.initScene();
        this.initEvents();
        this.animate();
    }

    initScene() {
        this.scene = new THREE.Scene();
        this.scene.background = new THREE.Color(0x1a1a1a);
        
        this.camera = new THREE.PerspectiveCamera(
            75,
            this.container.clientWidth / this.container.clientHeight,
            0.1,
            1000
        );
        this.camera.position.z = 200;

        this.renderer = new THREE.WebGLRenderer({ antialias: true });
        this.renderer.setSize(this.container.clientWidth, this.container.clientHeight);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        this.container.appendChild(this.renderer.domElement);

        this.controls = new OrbitControls(this.camera, this.renderer.domElement);
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.05;

        // Create node geometry and materials
        this.nodeGeometry = new THREE.CircleGeometry(5, 32);
        this.nodeMaterials = {
            epic: {
                default: new THREE.MeshBasicMaterial({ color: 0x3b82f6 }),
                hover: new THREE.MeshBasicMaterial({ color: 0x60a5fa }),
                selected: new THREE.MeshBasicMaterial({ color: 0x93c5fd })
            },
            feature: {
                default: new THREE.MeshBasicMaterial({ color: 0x22c55e }),
                hover: new THREE.MeshBasicMaterial({ color: 0x4ade80 }),
                selected: new THREE.MeshBasicMaterial({ color: 0x86efac })
            },
            task: {
                default: new THREE.MeshBasicMaterial({ color: 0xeab308 }),
                hover: new THREE.MeshBasicMaterial({ color: 0xfacc15 }),
                selected: new THREE.MeshBasicMaterial({ color: 0xfde047 })
            }
        };

        // Create link materials
        this.linkMaterials = {
            default: new THREE.LineBasicMaterial({
                color: 0x666666,
                transparent: true,
                opacity: 0.3
            }),
            highlighted: new THREE.LineBasicMaterial({
                color: 0x9ca3af,
                transparent: true,
                opacity: 0.6
            })
        };
    }

    initEvents() {
        window.addEventListener('resize', () => this.onWindowResize());
        this.renderer.domElement.addEventListener('mousemove', (e) => this.onMouseMove(e));
        this.renderer.domElement.addEventListener('click', (e) => this.onClick(e));
    }

    updateData(nodes, links) {
        console.log('Updating data:', nodes, links);
        
        // Clear existing objects
        this.clearScene();

        // Create nodes
        nodes.forEach(node => {
            const mesh = new THREE.Mesh(
                this.nodeGeometry,
                this.nodeMaterials[node.type]?.default || this.nodeMaterials.task.default
            );
            
            // Random position in a sphere
            const theta = Math.random() * Math.PI * 2;
            const phi = Math.acos(Math.random() * 2 - 1);
            const r = 100;
            
            mesh.position.x = r * Math.sin(phi) * Math.cos(theta);
            mesh.position.y = r * Math.sin(phi) * Math.sin(theta);
            mesh.position.z = r * Math.cos(phi);

            // Scale based on type
            const scale = node.type === 'epic' ? 1.5 : node.type === 'feature' ? 1.2 : 1;
            mesh.scale.setScalar(scale);

            this.scene.add(mesh);
            this.nodeObjects.set(node.id, { node, mesh });
        });

        // Create links
        links.forEach(link => {
            const sourceNode = this.nodeObjects.get(link.source);
            const targetNode = this.nodeObjects.get(link.target);

            if (sourceNode && targetNode) {
                const geometry = new THREE.BufferGeometry();
                const positions = new Float32Array([
                    sourceNode.mesh.position.x, sourceNode.mesh.position.y, sourceNode.mesh.position.z,
                    targetNode.mesh.position.x, targetNode.mesh.position.y, targetNode.mesh.position.z
                ]);
                geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
                
                const line = new THREE.Line(geometry, this.linkMaterials.default);
                this.scene.add(line);
                this.linkObjects.set(`${link.source}-${link.target}`, {
                    line,
                    source: link.source,
                    target: link.target
                });
            }
        });

        // Ensure camera is positioned correctly
        this.camera.position.z = 200;
        this.camera.lookAt(0, 0, 0);
        this.controls.update();
    }

    clearScene() {
        // Remove all nodes and links
        this.nodeObjects.forEach(({ mesh }) => this.scene.remove(mesh));
        this.linkObjects.forEach(({ line }) => this.scene.remove(line));
        
        this.nodeObjects.clear();
        this.linkObjects.clear();
        this.selectedNode = null;
        this.hoveredNode = null;
    }

    onWindowResize() {
        this.camera.aspect = this.container.clientWidth / this.container.clientHeight;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(this.container.clientWidth, this.container.clientHeight);
    }

    onMouseMove(event) {
        const rect = this.renderer.domElement.getBoundingClientRect();
        this.mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
        this.mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;

        this.checkIntersection();
    }

    onClick(event) {
        if (this.hoveredNode) {
            this.selectNode(this.hoveredNode);
        } else {
            this.selectNode(null);
        }
    }

    checkIntersection() {
        this.raycaster.setFromCamera(this.mouse, this.camera);
        const intersects = [];
        
        this.nodeObjects.forEach(({ node, mesh }) => {
            const intersect = this.raycaster.intersectObject(mesh);
            if (intersect.length > 0) {
                intersects.push({ node, mesh, distance: intersect[0].distance });
            }
        });

        // Sort by distance and get the closest
        intersects.sort((a, b) => a.distance - b.distance);
        const closest = intersects[0];

        if (closest && (!this.hoveredNode || this.hoveredNode.node.id !== closest.node.id)) {
            this.setHoveredNode(closest);
        } else if (!closest && this.hoveredNode) {
            this.setHoveredNode(null);
        }
    }

    setHoveredNode(nodeData) {
        // Reset previous hover state
        if (this.hoveredNode) {
            const materials = this.nodeMaterials[this.hoveredNode.node.type];
            this.hoveredNode.mesh.material = this.selectedNode?.node.id === this.hoveredNode.node.id 
                ? materials.selected 
                : materials.default;
        }

        this.hoveredNode = nodeData;

        // Set new hover state
        if (this.hoveredNode) {
            const materials = this.nodeMaterials[this.hoveredNode.node.type];
            this.hoveredNode.mesh.material = materials.hover;
            this.container.style.cursor = 'pointer';
        } else {
            this.container.style.cursor = 'default';
        }

        this.updateLinkHighlights();
    }

    selectNode(nodeData) {
        // Reset previous selection
        if (this.selectedNode) {
            const materials = this.nodeMaterials[this.selectedNode.node.type];
            this.selectedNode.mesh.material = materials.default;
        }

        this.selectedNode = nodeData;

        // Set new selection
        if (this.selectedNode) {
            const materials = this.nodeMaterials[this.selectedNode.node.type];
            this.selectedNode.mesh.material = materials.selected;
            this.container.dispatchEvent(new CustomEvent('nodeSelect', { 
                detail: this.selectedNode.node 
            }));
        }

        this.updateLinkHighlights();
    }

    updateLinkHighlights() {
        const highlightedNodeId = this.selectedNode?.node.id || this.hoveredNode?.node.id;

        this.linkObjects.forEach(({ line, source, target }) => {
            if (highlightedNodeId && (source === highlightedNodeId || target === highlightedNodeId)) {
                line.material = this.linkMaterials.highlighted;
            } else {
                line.material = this.linkMaterials.default;
            }
        });
    }

    animate() {
        requestAnimationFrame(() => this.animate());
        
        if (this.controls) {
            this.controls.update();
        }

        if (this.renderer && this.scene && this.camera) {
            this.renderer.render(this.scene, this.camera);
        }
    }
} 