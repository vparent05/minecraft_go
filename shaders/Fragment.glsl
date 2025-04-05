#version 400 core
flat in int orientation;

out vec4 FragColor;
void main()
{		
		switch (orientation) {
		case 0:
			FragColor = vec4(1, 0, 0, 1);
			break;
		case 1:
			FragColor = vec4(1, 0, 0, 1);
			break;
		case 2:
			FragColor = vec4(0, 0, 1, 1);
			break;
		case 3:
			FragColor = vec4(0, 0, 1, 1);
			break;
		case 4:
			FragColor = vec4(0, 1, 0, 1);
			break;
		case 5:
			FragColor = vec4(0, 1, 0, 1);
			break;
		default:
			FragColor = vec4(1, 1, 1, 1);
			break;
		}
}
