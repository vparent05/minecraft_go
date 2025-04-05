#version 400 core
layout (location = 0) in int ver;

uniform mat4 view;
uniform mat4 projection;

flat out int orientation;
void main()
{	
	int x = (ver>>28) & 0xF;
	int y = (ver>>20) & 0xFF;
	int z = (ver>>16) & 0xF;
	gl_Position = projection * view * vec4(float(x), float(y), float(z), 1.0);

	orientation = (ver>>12) & 0xF;
}
