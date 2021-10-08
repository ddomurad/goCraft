#version 330 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTex;

uniform mat4 uTrans;
uniform mat4 uProj;
uniform mat4 uView;

out vec2 texCoord;

void main(){
    gl_Position = uProj * uView * uTrans * vec4(aPos, 1.0);
    texCoord = aTex;
}