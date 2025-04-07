#version 460 core

in vec3 textureCoordinates;

uniform samplerCube skybox;

out vec4 FragColor;
void main() {	
	FragColor = texture(skybox, textureCoordinates);
}
