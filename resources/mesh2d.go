package resources

import (
	"errors"

	"github.com/ddomurad/goCraft/core"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ProceduralMesh2dLoader struct{}

const (
	PMT_QUAD        ProceduralMeshType = "pmt_quad"
	PMT_QUAD_BORDER ProceduralMeshType = "pmt_quad_border"
)

func (l ProceduralMesh2dLoader) CanLoad(resourceType core.ResourceType, uri string, param core.LoaderParam) bool {
	if resourceType != RT_MESH {
		return false
	}

	switch meshType := param.(type) {
	case ProceduralMeshType:
		return meshType == PMT_QUAD || meshType == PMT_QUAD_BORDER
	default:
		return false
	}
}

func (l ProceduralMesh2dLoader) Load(uri string, param core.LoaderParam) (core.Resource, error) {
	var verticesData []float32
	var indices []uint32
	var drawingType uint32

	switch param.(ProceduralMeshType) {
	case PMT_QUAD:
		verticesData, indices = getQuadMesh()
		drawingType = gl.TRIANGLE_FAN
	case PMT_QUAD_BORDER:
		verticesData, indices = getQuadBorderMesh()
		drawingType = gl.LINE_STRIP
	default:
		return GetEmptyMesh(uri), errors.New("unsuported procedural mesh type")
	}

	var vao uint32
	var vbo uint32
	var ibo uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verticesData)*4, gl.Ptr(verticesData), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return core.Resource{
		Type:  RT_MESH,
		Uri:   uri,
		Empty: false,
		Data: MeshData{
			VAO:     vao,
			VBO:     vbo,
			IBO:     ibo,
			VCount:  int32(len(indices)),
			Drawing: drawingType,
		},
		Unload: func() {
			gl.DeleteBuffers(1, &vbo)
			gl.DeleteBuffers(1, &ibo)
			gl.DeleteVertexArrays(1, &vao)
		},
	}, nil
}

func NewProceduralMesh2dLoader() ProceduralMesh2dLoader {
	return ProceduralMesh2dLoader{}
}

func getQuadMesh() ([]float32, []uint32) {
	return []float32{
		-0.5, -0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.0, 0.0, 0.0,
	}, []uint32{0, 1, 2, 3}
}

func getQuadBorderMesh() ([]float32, []uint32) {
	return []float32{
		-0.5, -0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.0, 0.0, 0.0,
	}, []uint32{0, 1, 2, 3, 0}
}

// func getCircleMesh() ([]float32, []uint32) {
// 	return []float32{
// 		-0.5, -0.5, 0.0, 0.0, 1.0,
// 		0.5, -0.5, 0.0, 1.0, 1.0,
// 		0.5, 0.5, 0.0, 1.0, 0.0,
// 		-0.5, 0.5, 0.0, 0.0, 0.0,
// 	}, []uint32{0, 1, 2, 3}
// }
