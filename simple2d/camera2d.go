package simple2d

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera2d struct {
	position  mgl32.Vec2
	rootation float32
	zoom      float32

	viewMatrix        mgl32.Mat4
	shouldRecalculate bool
}

func NewCamera() *Camera2d {
	return &Camera2d{
		position:          mgl32.Vec2{0, 0},
		rootation:         0,
		zoom:              1,
		viewMatrix:        mgl32.Ident4(),
		shouldRecalculate: true,
	}
}

func (c *Camera2d) SetPositionV(v mgl32.Vec2) {
	c.position = v
	c.shouldRecalculate = true
}

func (c *Camera2d) MovePositionV(v mgl32.Vec2) {
	c.position[0] += v[0]
	c.position[1] += v[1]
	c.shouldRecalculate = true
}

func (c *Camera2d) SetPosition(x, y float32) {
	c.position[0] = x
	c.position[1] = y
	c.shouldRecalculate = true
}

func (c *Camera2d) MovePosition(x, y float32) {
	c.position[0] += x
	c.position[1] += y
	c.shouldRecalculate = true
}

func (c *Camera2d) SetZoom(z float32) {
	c.zoom = z
	c.shouldRecalculate = true
}

func (c *Camera2d) Zoom(z float32) {
	c.zoom += z
	c.shouldRecalculate = true
}

func (c *Camera2d) SetRotation(rot float32) {
	c.rootation = rot
	c.shouldRecalculate = true
}

func (c *Camera2d) Rotate(rot float32) {
	c.rootation += rot
	c.shouldRecalculate = true
}

func (c *Camera2d) GetViewMatrix() mgl32.Mat4 {
	if c.shouldRecalculate {
		c.updateViewMatrix()
		c.shouldRecalculate = false
	}

	return c.viewMatrix
}

func (c *Camera2d) updateViewMatrix() {
	c.viewMatrix = mgl32.Translate3D(-c.position[0], -c.position[1], 0).
		Mul4(mgl32.HomogRotate3DZ(-c.rootation)).
		Mul4(mgl32.Scale3D(c.zoom, c.zoom, 1))
}
