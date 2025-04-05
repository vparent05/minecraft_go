#version 400 core
in Vec2 uv_f;

uniform sampler2D img2DSampler;

out vec4 FragColor;
void main() {		
	FragColor = texture(img2DSampler, uv_f);
}
