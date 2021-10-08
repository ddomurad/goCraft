package simple2d

import (
	"unsafe"

	"github.com/ddomurad/goCraft/core"
	"github.com/ddomurad/goCraft/resource"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DRI_MESH_QUAD             = "default_quad_mesh"
	DRI_MESH_CIRCLE           = "default_circle_mesh"
	DRI_MESH_QUAD_BORDER      = "default_quad_border_mesh"
	DRI_MESH_CIRCLE_BORDER    = "default_circle_border_mesh"
	DRI_SHADER_SIMPLE         = "default_simple_shader_program"
	DRI_SHADER_SIMPLE_TEXTURE = "default_simple_texture_shader_program"
)

type Scene2d interface {
	Render(dt float64, renderere *Renderer2d, app *core.App)
}

type Renderer2d struct {
	scene Scene2d

	clearColor          core.Color
	activeShaderProgram resource.ShaderData
	quadMesh            resource.MeshData
	circleMesh          resource.MeshData
	quadBorderMesh      resource.MeshData
	circleBorderMesh    resource.MeshData
	projectionMatrix    mgl32.Mat4
	activeViewMatrix    mgl32.Mat4
	alphaEnabled        bool
	updateNeeded        bool
	app                 *core.App
}

func (r *Renderer2d) Render(dt float64, app *core.App) {
	if r.updateNeeded {
		gl.ClearColor(r.clearColor[0], r.clearColor[1], r.clearColor[2], r.clearColor[3])

		if r.activeShaderProgram.ProgramId == 0 {
			r.activeShaderProgram = app.ResourceManager.GetResource(DRI_SHADER_SIMPLE).Data.(resource.ShaderData)
			gl.UseProgram(r.activeShaderProgram.ProgramId)
		}

		r.quadMesh = app.ResourceManager.GetResource(DRI_MESH_QUAD).Data.(resource.MeshData)
		r.circleMesh = app.ResourceManager.GetResource(DRI_MESH_CIRCLE).Data.(resource.MeshData)
		r.quadBorderMesh = app.ResourceManager.GetResource(DRI_MESH_QUAD_BORDER).Data.(resource.MeshData)
		r.circleBorderMesh = app.ResourceManager.GetResource(DRI_MESH_CIRCLE_BORDER).Data.(resource.MeshData)

		if r.alphaEnabled {
			gl.Enable(gl.BLEND)
			gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		} else {
			gl.Disable(gl.BLEND)
		}

		wh := float32(app.Window.Width) / float32(app.Window.Height)
		r.projectionMatrix = mgl32.Ortho2D(-wh, wh, -1, 1)

		r.activeShaderProgram.SetProjectionMat(r.projectionMatrix)
		gl.Viewport(0, 0, int32(app.Window.Width), int32(app.Window.Height))
		r.updateNeeded = false
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r.scene.Render(dt, r, app)
}

func NewRenderer2d(app *core.App, scene Scene2d) *Renderer2d {
	return &Renderer2d{
		scene:        scene,
		updateNeeded: true,
		app:          app,
	}
}

func (r *Renderer2d) HandleEvent(e core.Event) bool {
	switch e.(type) {
	case core.ResizeEvent:
		r.updateNeeded = true
	}
	return false
}

func (r *Renderer2d) Init() {
	r.app.ResourceManager.
		AddLoader(resource.NewFileTextureLoader()).
		AddLoader(resource.NewShaderLoader()).
		AddLoader(resource.NewProceduralMesh2dLoader()).
		PreloadReource(resource.RT_MESH, DRI_MESH_QUAD, resource.PMT_QUAD).
		PreloadReource(resource.RT_MESH, DRI_MESH_CIRCLE, resource.PMT_CIRCLE).
		PreloadReource(resource.RT_MESH, DRI_MESH_QUAD_BORDER, resource.PMT_QUAD_BORDER).
		PreloadReource(resource.RT_MESH, DRI_MESH_CIRCLE_BORDER, resource.PMT_CIRCLE_BORDER).
		PreloadReource(resource.RT_SHADER, DRI_SHADER_SIMPLE, resource.EmbededShaderSource{
			ShaderName: "simple",
		}).
		PreloadReource(resource.RT_SHADER, DRI_SHADER_SIMPLE_TEXTURE, resource.EmbededShaderSource{
			ShaderName: "simple_texture",
		})

	// PreloadReource(resource.RT_SHADER, DRI_SHADER_PROGRAM, resource.ShaderFileSource{
	// 	VertexShaderPath:   "/home/work/Projects/goCraftProject/goCraftTestApp/res/shader.vs",
	// 	FragmentShaderPath: "/home/work/Projects/goCraftProject/goCraftTestApp/res/shader.fs",
	// })

	r.app.EventManager.RegisterHandler(r)
}

func (r *Renderer2d) SetWirframe(enable bool) {
	if enable {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func (r *Renderer2d) SetCrictleSegments(segments uint) {
	resourceManager := r.app.ResourceManager
	mesh := resourceManager.GetResource(DRI_MESH_CIRCLE)
	mesh.Unload()

	borderMesh := resourceManager.GetResource(DRI_MESH_CIRCLE_BORDER)
	borderMesh.Unload()

	vertexData, indexData, drawingType := resource.GetCircleVertices(segments)
	meshReshource, _ := resource.CreateMesh2dResource(DRI_MESH_CIRCLE, vertexData, indexData, drawingType)
	resourceManager.AddResource(meshReshource)

	vertexBorderData, indexBorderData, drawingBorderType := resource.GetCircleBorderVertices(segments)
	meshBorderReshource, _ := resource.CreateMesh2dResource(DRI_MESH_CIRCLE_BORDER, vertexBorderData, indexBorderData, drawingBorderType)
	resourceManager.AddResource(meshBorderReshource)
}

func (r *Renderer2d) SetClearColor(color core.Color) {
	r.clearColor = color
	r.updateNeeded = true
}

func (r *Renderer2d) SetShader(uri string) {
	r.activeShaderProgram = r.app.ResourceManager.GetResource(uri).Data.(resource.ShaderData)
	gl.UseProgram(r.activeShaderProgram.ProgramId)
	r.activeShaderProgram.SetProjectionMat(r.projectionMatrix)
	r.activeShaderProgram.SetViewMat(r.activeViewMatrix)
}

func (r *Renderer2d) GetShader() resource.ShaderData {
	return r.activeShaderProgram
}

func (r *Renderer2d) SetTexture(uri string) {
	textureData := r.app.ResourceManager.GetResource(uri).Data.(resource.TextureData)
	gl.BindTexture(gl.TEXTURE_2D, textureData.Id)
}

func (r *Renderer2d) SetAlpha(enabled bool) {
	r.alphaEnabled = enabled
	r.updateNeeded = true
}

func (r *Renderer2d) ApplyCamera(camera *Camera2d) {
	r.activeViewMatrix = camera.GetViewMatrix()
	r.activeShaderProgram.SetViewMat(r.activeViewMatrix)
}

func (r *Renderer2d) DrawRectV(pos, size mgl32.Vec2, rot float32, color core.Color) {
	r.DrawRect(pos.X(), pos.Y(), size.X(), size.Y(), rot, color)
}

func (r *Renderer2d) DrawRect(x, y, w, h, rot float32, color core.Color) {
	var transformMat = getTransformMattrix(x, y, w, h, rot)

	r.activeShaderProgram.SetColor(color)
	r.activeShaderProgram.SetTransformationMat(transformMat)

	gl.BindVertexArray(r.quadMesh.VAO)
	gl.DrawElements(r.quadMesh.Drawing, int32(r.quadMesh.VCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (r *Renderer2d) DrawRectBorderV(pos, size mgl32.Vec2, rot, width float32, color core.Color) {
	r.DrawRectBorder(pos.X(), pos.Y(), size.X(), size.Y(), rot, width, color)
}

func (r *Renderer2d) DrawRectBorder(x, y, w, h, rot, width float32, color core.Color) {
	var transformMat = getTransformMattrix(x, y, w, h, rot)

	r.activeShaderProgram.SetColor(color)
	r.activeShaderProgram.SetTransformationMat(transformMat)

	gl.LineWidth(width)
	gl.BindVertexArray(r.quadBorderMesh.VAO)
	gl.DrawElements(r.quadBorderMesh.Drawing, int32(r.quadBorderMesh.VCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (r *Renderer2d) DrawElipseV(pos, size mgl32.Vec2, rot float32, color core.Color) {
	r.DrawElipse(pos[0], pos[1], size[0], size[1], rot, color)
}

func (r *Renderer2d) DrawElipse(x, y, w, h, rot float32, color core.Color) {
	var transformMat = getTransformMattrix(x, y, w, h, rot)

	r.activeShaderProgram.SetColor(color)
	r.activeShaderProgram.SetTransformationMat(transformMat)

	gl.BindVertexArray(r.circleMesh.VAO)
	gl.DrawElements(r.circleMesh.Drawing, int32(r.circleMesh.VCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (r *Renderer2d) DrawElipseBorder(x, y, w, h, rot, width float32, color core.Color) {
	var transformMat = getTransformMattrix(x, y, w, h, rot)

	r.activeShaderProgram.SetColor(color)
	r.activeShaderProgram.SetTransformationMat(transformMat)

	gl.LineWidth(width)
	gl.BindVertexArray(r.circleBorderMesh.VAO)
	gl.DrawElements(r.circleBorderMesh.Drawing, int32(r.circleBorderMesh.VCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
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
