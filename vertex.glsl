#version 330

layout (location = 0) in vec2 vert;
layout (location = 1) in vec2 tex;
layout (location = 2) in vec4 mask;

uniform mat3 camera;
uniform vec2 viewportSize;
uniform vec2 textureSize;

out vec2 fragTex;
out vec4 fragMask;
void main() {
    fragMask = mask;
    fragTex = tex/textureSize;
    gl_Position = vec4(camera * vec3(vert/viewportSize, 0), 1);
}