package resource

import (
	"errors"
	"math"
	"strings"

	"github.com/ddomurad/goCraft/core"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ProceduralMesh2dLoader struct{}

const (
	PMT_QUAD          ProceduralMeshType = "pmt_2d_quad"
	PMT_QUAD_BORDER   ProceduralMeshType = "pmt_2d_quad_border"
	PMT_CIRCLE        ProceduralMeshType = "pmt_2d_circle"
	PMT_CIRCLE_BORDER ProceduralMeshType = "pmt_2d_circle_border"
)

func (l ProceduralMesh2dLoader) CanLoad(resourceType core.ResourceType, uri string, param core.LoaderParam) bool {
	if resourceType != RT_MESH {
		return false
	}

	switch meshType := param.(type) {
	case ProceduralMeshType:
		return strings.HasPrefix(string(meshType), "pmt_2d_")
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
		verticesData, indices, drawingType = GetQuadVertices()
	case PMT_QUAD_BORDER:
		verticesData, indices, drawingType = GetQuadBorderVertices()
	case PMT_CIRCLE:
		verticesData, indices, drawingType = GetCircleVertices(24)
	case PMT_CIRCLE_BORDER:
		verticesData, indices, drawingType = GetCircleBorderVertices(24)
	default:
		return GetEmptyMesh(uri), errors.New("unsuported procedural mesh type")
	}

	return CreateMesh2dResource(uri, verticesData, indices, drawingType)
}

func NewProceduralMesh2dLoader() ProceduralMesh2dLoader {
	return ProceduralMesh2dLoader{}
}

func CreateMesh2dResource(uri string, verticesData []float32, indexData []uint32, drawingType uint32) (core.Resource, error) {
	meshData := CreateMesh2d(verticesData, indexData, drawingType)

	return core.Resource{
		Type:  RT_MESH,
		Uri:   uri,
		Empty: false,
		Data:  meshData,
		Unload: func() {
			meshData.Unload()
		},
	}, nil
}

func CreateMesh2d(verticesData []float32, indexData []uint32, drawingType uint32) MeshData {
	vao, vbo, ibo := CreateMeshBuffers(verticesData, indexData)

	return MeshData{
		VAO:     vao,
		VBO:     vbo,
		IBO:     ibo,
		VCount:  int32(len(indexData)),
		Drawing: drawingType,
	}
}

func CreateMeshBuffers(verticesData []float32, indexData []uint32) (vao, vbo, ibo uint32) {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verticesData)*4, gl.Ptr(verticesData), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indexData)*4, gl.Ptr(indexData), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)
	return
}

func GetQuadVertices() ([]float32, []uint32, uint32) {
	return []float32{
			-0.5, -0.5, 0.0, 0.0, 1.0,
			0.5, -0.5, 0.0, 1.0, 1.0,
			0.5, 0.5, 0.0, 1.0, 0.0,
			-0.5, 0.5, 0.0, 0.0, 0.0,
		}, []uint32{0, 1, 2, 3},
		gl.TRIANGLE_FAN
}

func GetQuadBorderVertices() ([]float32, []uint32, uint32) {
	return []float32{
			-0.5, -0.5, 0.0, 0.0, 1.0,
			0.5, -0.5, 0.0, 1.0, 1.0,
			0.5, 0.5, 0.0, 1.0, 0.0,
			-0.5, 0.5, 0.0, 0.0, 0.0,
		}, []uint32{0, 1, 2, 3, 0},
		gl.LINE_STRIP
}

func GetCircleVertices(segments uint) (vertices []float32, indices []uint32, drawingType uint32) {
	vertices = make([]float32, (segments+1)*5)
	indices = make([]uint32, segments+2)

	dAnlge := math.Pi * 2 / float64(segments)

	vertices[3] = 0.5
	vertices[4] = 0.5

	for i := uint(1); i < segments+1; i++ {
		vertex := vertices[i*5 : i*5+5]
		vertexPos := vertex[0:3]
		vertexUv := vertex[3:5]

		angle := dAnlge * float64(i-1)

		vertexPos[0] = float32(math.Cos(angle)) / 2.0
		vertexPos[1] = float32(math.Sin(angle)) / 2.0
		vertexPos[2] = 0

		vertexUv[0] = 0.5 + float32(math.Cos(angle))/2.0
		vertexUv[1] = 0.5 + float32(math.Sin(angle))/2.0

		indices[i] = uint32(i)
	}

	indices[segments+1] = 1
	drawingType = gl.TRIANGLE_FAN

	return
}

func GetCircleBorderVertices(segments uint) (vertices []float32, indices []uint32, drawingType uint32) {
	vertices = make([]float32, segments*5)
	indices = make([]uint32, segments+1)

	dAnlge := math.Pi * 2 / float64(segments)

	for i := uint(0); i < segments; i++ {
		vertex := vertices[i*5 : i*5+5]
		vertexPos := vertex[0:3]
		vertexUv := vertex[3:5]

		angle := dAnlge * float64(i)

		vertexPos[0] = float32(math.Cos(angle)) / 2.0
		vertexPos[1] = float32(math.Sin(angle)) / 2.0
		vertexPos[2] = 0

		vertexUv[0] = 0.5 + float32(math.Cos(angle))/2.0
		vertexUv[1] = 0.5 + float32(math.Sin(angle))/2.0

		indices[i] = uint32(i)
	}

	drawingType = gl.LINE_STRIP

	return
}
