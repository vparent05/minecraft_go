#version 460 core
layout (location = 0) in int vertex;

uniform vec2 chunkCoordinates;
uniform mat4 view;
uniform mat4 projection;

out vec2 uv;
flat out int orientation;
void main()
{	
	float x = ((vertex>>28) & 0xF) + chunkCoordinates.x * 15;
	float y = (vertex>>20) & 0xFF;
	float z = ((vertex>>16) & 0xF) + chunkCoordinates.y * 15;
	orientation = (vertex>>12) & 0xF;
	int textIndex = (vertex>>4) & 0xFF;

	// level is only used for vertices of height y+1 
	int level = vertex & 0xF;
	gl_Position = projection * view * vec4(x, y - (15 - level)/16.0, float(z), 1.0);

	uv = vec2(textIndex % 16, textIndex / 16) / 16.0;

	// adjust the side texture to match the level
	if (orientation != 0) {
		uv.y += float(15 - level)/256.0;
	}
}
