#version 330 core

uniform vec4 uColor;
uniform sampler2D textSampler;

in vec2 texCoord;
out vec4 FragColor; 

void main() { 
    FragColor = uColor*texture(textSampler, texCoord); 
} 