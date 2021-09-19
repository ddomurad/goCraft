package resources

import (
	"errors"
	"io/ioutil"

	"github.com/ddomurad/goCraft/core"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderStringSource struct {
	VertexShader   string
	FragmentShader string
}

type ShaderFileSource struct {
	VertexShaderPath   string
	FragmentShaderPath string
}

const (
	RT_SHADER core.ResourceType = "shader"
)

type ShaderData struct {
	ProgramId        uint32
	uniformLocations map[string]int32
}

func (s *ShaderData) GetUniformLocation(name string) int32 {
	if s.uniformLocations == nil {
		return gl.GetUniformLocation(s.ProgramId, gl.Str(name+"\x00"))
	}

	loc, ok := s.uniformLocations[name]
	if ok {
		return loc
	}

	loc = gl.GetUniformLocation(s.ProgramId, gl.Str(name+"\x00"))
	s.uniformLocations[name] = loc
	return loc
}

func GetEmptyShader(uri string) core.Resource {
	return core.Resource{
		Type:   RT_SHADER,
		Uri:    uri,
		Empty:  true,
		Data:   ShaderData{},
		Unload: nil,
	}
}

type ShaderLoader struct{}

func (l ShaderLoader) CanLoad(resourceType core.ResourceType, uri string, param core.LoaderParam) bool {
	if resourceType != RT_SHADER {
		return false
	}

	switch param.(type) {
	case ShaderStringSource:
		return true
	case ShaderFileSource:
		return true
	default:
		return false
	}
}

func (l ShaderLoader) Load(uri string, param core.LoaderParam) (core.Resource, error) {
	var shaderData ShaderData
	var loadError error

	switch source := param.(type) {
	case ShaderStringSource:
		shaderData, loadError = loadShadersFromString(source.FragmentShader, source.VertexShader)
	case ShaderFileSource:
		shaderData, loadError = loadShadersFromFiles(source.FragmentShaderPath, source.VertexShaderPath)
	default:
		return GetEmptyShader(uri), errors.New("unsuported shader source")
	}

	if loadError != nil {
		return GetEmptyShader(uri), loadError
	}

	return core.Resource{
		Type:  RT_SHADER,
		Uri:   uri,
		Empty: false,
		Data:  shaderData,
		Unload: func() {
			gl.DeleteProgram(shaderData.ProgramId)
		},
	}, nil
}

func GetDefaultShaderSource() ShaderStringSource {
	return ShaderStringSource{
		VertexShader:   "#version 330 core\n layout (location = 0) in vec3 aPos; out vec4 vertexColor; void main(){ gl_Position = vec4(aPos, 1.0); vertexColor = vec4(0.5, 0.0, 0.0, 1.0); }",
		FragmentShader: "#version 330 core\n out vec4 FragColor; in vec4 vertexColor;  void main() { FragColor = vertexColor; } ",
	}
}

func NewShaderLoader() ShaderLoader {
	return ShaderLoader{}
}

func loadShadersFromFiles(fsPath string, vsPath string) (ShaderData, error) {
	fs_text, err := ioutil.ReadFile(fsPath)
	if err != nil {
		return ShaderData{}, err
	}

	vs_text, err := ioutil.ReadFile(vsPath)

	if err != nil {
		return ShaderData{}, err
	}

	return loadShadersFromString(string(fs_text), string(vs_text))
}

func loadShadersFromString(fs string, vs string) (ShaderData, error) {
	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)

	fs_src, fs_free_fnc := gl.Strs(fs + "\x00")
	vs_src, vs_free_fnc := gl.Strs(vs + "\x00")

	defer fs_free_fnc()
	defer vs_free_fnc()

	gl.ShaderSource(fragment_shader, 1, fs_src, nil)
	gl.ShaderSource(vertex_shader, 1, vs_src, nil)

	defer gl.DeleteShader(fragment_shader)
	defer gl.DeleteShader(vertex_shader)

	gl.CompileShader(fragment_shader)

	var status int32
	gl.GetShaderiv(fragment_shader, gl.COMPILE_STATUS, &status)

	if status == 0 {
		return ShaderData{}, errors.New("failed to compile fragment shader")
	}

	gl.CompileShader(vertex_shader)
	gl.GetShaderiv(vertex_shader, gl.COMPILE_STATUS, &status)

	if status == 0 {
		return ShaderData{}, errors.New("failed to compile vertex shader")
	}

	shader_program := gl.CreateProgram()
	gl.AttachShader(shader_program, vertex_shader)
	gl.AttachShader(shader_program, fragment_shader)
	gl.LinkProgram(shader_program)

	gl.GetProgramiv(shader_program, gl.LINK_STATUS, &status)

	if status == 0 {
		gl.DeleteProgram(shader_program)
		return ShaderData{}, errors.New("failed to link shader program")
	}

	return ShaderData{
		ProgramId:        shader_program,
		uniformLocations: make(map[string]int32),
	}, nil
}
