#version 460 core
layout (location = 0) in int ver;

uniform mat4 view;
uniform mat4 projection;

out vec2 uv;
flat out int orientation;
void main()
{	
	int x = (ver>>28) & 0xF;
	int y = (ver>>20) & 0xFF;
	int z = (ver>>16) & 0xF;
	orientation = (ver>>12) & 0xF;
	int textIndex = (ver>>4) & 0xFF;

	// level is only used for vertices of height y+1 
	int level = ver & 0xF;
	gl_Position = projection * view * vec4(float(x), float(y) - float(15 - level)/16.0, float(z), 1.0);

	uv = vec2(textIndex % 16, textIndex / 16) / 16.0;

	// adjust the side texture to match the level
	if (orientation != 0) {
		uv.y += float(15 - level)/256.0;
	}
}
