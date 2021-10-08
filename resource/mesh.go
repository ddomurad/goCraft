package resource

import (
	"github.com/ddomurad/goCraft/core"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ProceduralMeshType string

const (
	RT_MESH core.ResourceType = "mesh"
)

type MeshData struct {
	VAO     uint32
	VBO     uint32
	IBO     uint32
	VCount  int32
	Drawing uint32
}

func GetEmptyMesh(uri string) core.Resource {
	return core.Resource{
		Type:   RT_MESH,
		Uri:    uri,
		Empty:  true,
		Data:   MeshData{},
		Unload: nil,
	}
}

func (m MeshData) Unload() {
	gl.DeleteBuffers(1, &m.VBO)
	gl.DeleteBuffers(1, &m.IBO)
	gl.DeleteVertexArrays(1, &m.VAO)
}
