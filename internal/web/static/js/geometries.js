import * as THREE from 'three';

export function createNodeGeometry() {
    const geometry = new THREE.CircleGeometry(1, 32);
    return geometry;
}

export function createLinkGeometry() {
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute([], 3));
    return geometry;
}

export function createNodeMaterial(color) {
    return new THREE.MeshBasicMaterial({
        color: color,
        transparent: true,
        opacity: 0.8,
        side: THREE.DoubleSide
    });
}

export function createLinkMaterial() {
    return new THREE.LineBasicMaterial({
        color: 0x666666,
        transparent: true,
        opacity: 0.3
    });
}

export function createNodeInstance(geometry, material, position, scale = 1) {
    const mesh = new THREE.Mesh(geometry, material);
    mesh.position.copy(position);
    mesh.scale.setScalar(scale);
    return mesh;
}

export function createLinkInstance(startPos, endPos) {
    const geometry = new THREE.BufferGeometry();
    const positions = new Float32Array([
        startPos.x, startPos.y, startPos.z,
        endPos.x, endPos.y, endPos.z
    ]);
    geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    return new THREE.Line(geometry, createLinkMaterial());
} 