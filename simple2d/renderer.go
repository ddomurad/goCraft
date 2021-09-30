package simple2d

import (
	"unsafe"

	"github.com/ddomurad/goCraft/core"
	"github.com/ddomurad/goCraft/resources"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DRI_SHADER_PROGRAM   = "default_shader_program"
	DRI_QUAD_MESH        = "default_quad_mesh"
	DRI_QUAD_BORDER_MESH = "default_quad_border_mesh"
)

type Scene2d interface {
	Render(dt float64, renderere *Renderer2d, app *core.App)
}

type Renderer2d struct {
	scene Scene2d

	clearColor     core.Color
	shaderProgram  resources.ShaderData
	quadMesh       resources.MeshData
	quadBorderMesh resources.MeshData

	projectionMatrix mgl32.Mat4
	alphaEnabled     bool
	updateNeeded     bool
}

func (r *Renderer2d) Render(dt float64, app *core.App) {
	if r.updateNeeded {
		gl.ClearColor(r.clearColor[0], r.clearColor[1], r.clearColor[2], r.clearColor[3])

		r.shaderProgram = app.ResourceManager.GetResource(DRI_SHADER_PROGRAM).Data.(resources.ShaderData)
		r.quadMesh = app.ResourceManager.GetResource(DRI_QUAD_MESH).Data.(resources.MeshData)
		r.quadBorderMesh = app.ResourceManager.GetResource(DRI_QUAD_BORDER_MESH).Data.(resources.MeshData)

		if r.alphaEnabled {
			gl.Enable(gl.BLEND)
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		} else {
			gl.Disable(gl.BLEND)
		}

		hw := float32(app.Window.Height) / float32(app.Window.Width)
		r.projectionMatrix = mgl32.Ortho2D(-1.0, 1.0, -hw, hw)

		gl.UseProgram(r.shaderProgram.ProgramId)
		r.shaderProgram.SetProjectionMat(r.projectionMatrix)

		r.updateNeeded = false
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.scene.Render(dt, r, app)
}

func NewRenderer2d(scene Scene2d) *Renderer2d {
	return &Renderer2d{
		scene: scene,
	}
}

func (r *Renderer2d) HandleEvent(e core.Event) bool {
	return false
}

func (r *Renderer2d) Init(app *core.App) {
	app.ResourceManager.
		AddLoader(resources.NewFileTextureLoader()).
		AddLoader(resources.NewShaderLoader()).
		AddLoader(resources.NewProceduralMesh2dLoader()).
		PreloadReource(resources.RT_MESH, DRI_QUAD_MESH, resources.PMT_QUAD).
		PreloadReource(resources.RT_MESH, DRI_QUAD_BORDER_MESH, resources.PMT_QUAD_BORDER).
		// PreloadReource(resources.RT_SHADER, DRI_SHADER_PROGRAM, resources.GetDefaultShaderSource())
		PreloadReource(resources.RT_SHADER, DRI_SHADER_PROGRAM, resources.ShaderFileSource{
			VertexShaderPath:   "/home/work/Projects/goCraftProject/goCraftTestApp/res/shader.vs",
			FragmentShaderPath: "/home/work/Projects/goCraftProject/goCraftTestApp/res/shader.fs",
		})

	app.EventManager.RegisterHandler(r)
}

func (r *Renderer2d) SetClearColor(color core.Color) {
	r.clearColor = color
	r.updateNeeded = true
}

func (r *Renderer2d) SetAlpha(enabled bool) {
	r.alphaEnabled = enabled
	r.updateNeeded = true
}

func (r *Renderer2d) DrawRectV(pos, size mgl32.Vec2, rot float32, color core.Color) {
	r.DrawRect(pos.X(), pos.Y(), size.X(), size.Y(), rot, color)
}

func (r *Renderer2d) DrawRect(x, y, w, h, rot float32, color core.Color) {
	var transformMat = getTransformMattrix(x, y, w, h, rot)

	r.shaderProgram.SetColor(color)
	r.shaderProgram.SetTransformationMat(transformMat)

	gl.BindVertexArray(r.quadMesh.VAO)
	gl.DrawElements(r.quadMesh.Drawing, int32(r.quadMesh.VCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (r *Renderer2d) DrawRectBorderV(pos, size mgl32.Vec2, rot, width float32, color core.Color) {
	r.DrawRectBorder(pos.X(), pos.Y(), size.X(), size.Y(), rot, width, color)
}

func (r *Renderer2d) DrawRectBorder(x, y, w, h, rot, width float32, color core.Color) {
	var transformMat = getTransformMattrix(x, y, w, h, rot)

	r.shaderProgram.SetColor(color)
	r.shaderProgram.SetTransformationMat(transformMat)

	gl.LineWidth(width)
	gl.BindVertexArray(r.quadBorderMesh.VAO)
	gl.DrawElements(r.quadBorderMesh.Drawing, int32(r.quadBorderMesh.VCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func getTransformMattrix(x, y, sx, sy, rot float32) mgl32.Mat4 {
	var transformMat mgl32.Mat4

	if rot != 0 {
		transformMat = mgl32.Translate3D(x, y, 0).
			Mul4(mgl32.HomogRotate3DZ(rot)).
			Mul4(mgl32.Scale3D(sx, sy, 1))
	} else {
		transformMat = mgl32.Translate3D(x, y, 0).
			Mul4(mgl32.Scale3D(sx, sy, 1))
	}

	return transformMat
}
