#version 330

uniform sampler2D tex;
uniform int useTexture;

in vec2 fragTex;
in vec4 fragMask;

out vec4 outputColor;

void main() {
    if (useTexture == 1) {
        outputColor = texture(tex, fragTex) * fragMask;
    } else {
        outputColor = fragMask;
    }
}