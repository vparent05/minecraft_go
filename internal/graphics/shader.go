package graphics

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var programs = map[Program]uint32{}

type Program = int

const (
	BLOCK  Program = iota
	SKYBOX Program = iota
)

type Shader struct {
	sourcePath string
	xtype      uint32
}

func NewShader(path string, xtype uint32) Shader {
	return Shader{path, xtype}
}

func readSource(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("os.ReadFile(): %w", err)
	}
	return string(data), nil
}

func (s *Shader) compile() (uint32, error) {
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

func NewProgram(id int, shaders ...Shader) (uint32, error) {
	program := gl.CreateProgram()
	for _, s := range shaders {
		sh, err := s.compile()
		if err != nil {
			return 0, fmt.Errorf("newShader(): %w", err)
		}
		gl.AttachShader(program, sh)
		gl.DeleteShader(sh)
	}
	gl.LinkProgram(program)

	programs[id] = program
	return program, nil
}

func SelectProgram(id Program) (uint32, error) {
	program, ok := programs[id]
	if !ok {
		return 0, fmt.Errorf("program %v doesn't exist", id)
	}
	gl.UseProgram(program)
	return program, nil
}
