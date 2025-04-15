#version 460 core
flat in int orientation;
in vec2 uv;

uniform sampler2D atlas;

out vec4 FragColor;
void main()
{		
	vec4 col = texture(atlas, uv);
	
	switch (orientation) {
	case 0:
		FragColor = col;
		break;
	case 1:
		FragColor = col;
		break;
	case 2:
		FragColor = vec4(col.xyz * 0.8, col.w);
		break;
	case 3:
		FragColor = vec4(col.xyz * 0.8, col.w);
		break;
	case 4:
		FragColor = vec4(col.xyz * 0.6, col.w);
		break;
	case 5:
		FragColor = vec4(col.xyz * 0.6, col.w);
		break;
	default:
		FragColor = col;
		break;
	}
}
