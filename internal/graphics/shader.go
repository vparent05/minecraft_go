package graphics

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type program struct {
	id              uint32
	uniformLocation map[string]int32
}

type shader struct {
	sourcePath string
	xtype      uint32
}

func NewShader(path string, xtype uint32) shader {
	return shader{path, xtype}
}

func readSource(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("os.ReadFile(): %w", err)
	}
	return string(data), nil
}

func (s *shader) compile() (uint32, error) {
	shader := gl.CreateShader(s.xtype)
	source, err := readSource(s.sourcePath)
	if err != nil {
		return 0, fmt.Errorf("readSource(): %w", err)
	}
	cSource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, cSource, nil)
	free()
	gl.CompileShader(shader)

	var isCompiled int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &isCompiled)
	if isCompiled == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)

		log := make([]uint8, length)
		gl.GetShaderInfoLog(shader, length, &length, &log[0])

		return 0, fmt.Errorf("gl.CompileShader(): %s", string(log))
	}

	return shader, nil
}

func NewProgram(shaders ...shader) (*program, error) {
	p := gl.CreateProgram()
	for _, s := range shaders {
		sh, err := s.compile()
		if err != nil {
			return nil, fmt.Errorf("newShader(): %w", err)
		}
		gl.AttachShader(p, sh)
		gl.DeleteShader(sh)
	}
	gl.LinkProgram(p)

	return &program{
		p,
		make(map[string]int32),
	}, nil
}

func (p *program) use() {
	gl.UseProgram(p.id)
}

func (p *program) getUniformLocation(name string) (int32, error) {
	if location, ok := p.uniformLocation[name]; ok {
		return location, nil
	}
	location := gl.GetUniformLocation(p.id, gl.Str(name+"\x00"))
	if location == -1 {
		return -1, fmt.Errorf("gl.GetUniformLocation(): uniform \"%s\" doesn't exist", name)
	}
	return location, nil
}
