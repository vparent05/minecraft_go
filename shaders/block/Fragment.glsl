#version 460 core
flat in int orientation;
in vec2 uv;

uniform sampler2D atlas;

out vec4 FragColor;
void main()
{		
	vec3 col = texture(atlas, uv).xyz;
	
	switch (orientation) {
	case 0:
		FragColor = vec4(col, 1);
		break;
	case 1:
		FragColor = vec4(col, 1);
		break;
	case 2:
		FragColor = vec4(col * 0.8, 1);
		break;
	case 3:
		FragColor = vec4(col * 0.8, 1);
		break;
	case 4:
		FragColor = vec4(col * 0.6, 1);
		break;
	case 5:
		FragColor = vec4(col * 0.6, 1);
		break;
	default:
		FragColor = vec4(col, 1);
		break;
	}
}
