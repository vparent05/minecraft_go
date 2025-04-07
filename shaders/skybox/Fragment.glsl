#version 460 core

in vec3 textureCoordinates;

uniform samplerCube constellation;
uniform samplerCube skybox;

out vec4 FragColor;
void main() {	
	vec4 col1 = texture(constellation, textureCoordinates);
	vec4 col2 = texture(skybox, textureCoordinates);

	FragColor = col1 * 0.2 + col2;
}
