#version 400 core

layout (location = 0) in Vec3 position;
layout (location = 1) in Vec2 uv_v;

uniform mat4 projection;

out Vec2 uv_f;
void main() {	
	gl_Position = projection * vec4(position, 1.0);
	uv_f = uv_v;
}
