export const nodeVertexShader = `
    uniform float time;
    varying vec2 vUv;
    varying vec3 vColor;

    void main() {
        vUv = uv;
        vColor = instanceColor;
        
        // Add subtle animation
        vec3 pos = position;
        float scale = 1.0 + 0.1 * sin(time * 2.0);
        pos *= scale;

        vec4 mvPosition = modelViewMatrix * instanceMatrix * vec4(pos, 1.0);
        gl_Position = projectionMatrix * mvPosition;
    }
`;

export const nodeFragmentShader = `
    uniform float time;
    uniform float selected;
    uniform float hovered;
    varying vec2 vUv;
    varying vec3 vColor;

    void main() {
        // Calculate distance from center
        vec2 center = vec2(0.5, 0.5);
        float dist = length(vUv - center);

        // Base color with gradient
        vec3 color = vColor;
        float alpha = 1.0 - smoothstep(0.45, 0.5, dist);

        // Add glow effect
        float glow = 0.0;
        if (selected > 0.5) {
            glow = 0.5 * (1.0 + sin(time * 3.0));
            color += vec3(0.3, 0.3, 0.3) * glow;
            alpha += 0.2 * glow;
        } else if (hovered > 0.5) {
            glow = 0.3;
            color += vec3(0.2, 0.2, 0.2) * glow;
            alpha += 0.1;
        }

        // Add subtle pulsing
        float pulse = 0.05 * sin(time * 2.0);
        color += vec3(pulse);

        // Add edge highlight
        float edge = smoothstep(0.48, 0.5, dist);
        color += vec3(0.2) * edge;

        gl_FragColor = vec4(color, alpha);
    }
`;

export const createNodeGeometry = () => {
    const geometry = new THREE.CircleGeometry(1, 32);
    return geometry;
};

export const createLinkGeometry = () => {
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute([], 3));
    return geometry;
}; 