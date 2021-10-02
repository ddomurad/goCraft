package simple2d

import (
	"github.com/ddomurad/goCraft/core"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type MouseDragMonitor struct {
	button       glfw.MouseButton
	dragActive   bool
	lastPos      mgl32.Vec2
	currentPos   mgl32.Vec2
	deltaPos     mgl32.Vec2
	deltaApplied bool
	multiplier   float32
}

func NewMouseDragMonitor(button glfw.MouseButton, multiplier float32) *MouseDragMonitor {
	return &MouseDragMonitor{
		button:     button,
		multiplier: multiplier,
	}
}

func (m *MouseDragMonitor) HandleEvent(e core.Event) bool {
	switch te := e.(type) {
	case core.MouseMoveEvent:
		m.lastPos[0] = float32(te.Pos[0])
		m.lastPos[1] = float32(te.Pos[1])

		if m.dragActive {
			m.deltaApplied = false
		}
	case core.MouseButtonEvent:
		if te.Button == m.button {
			if te.Action == glfw.Press {
				m.dragActive = true
				m.deltaPos = mgl32.Vec2{0, 0}
				m.currentPos = m.lastPos
				m.deltaApplied = false
			} else if te.Action == glfw.Release {
				m.dragActive = false
			}
		}
	}

	return false
}

func (m *MouseDragMonitor) applyDelta() {
	if m.deltaApplied {
		return
	}

	m.deltaPos[0] = (m.currentPos[0] - m.lastPos[0]) * m.multiplier
	m.currentPos[0] = m.lastPos[0]

	m.deltaPos[1] = (m.lastPos[1] - m.currentPos[1]) * m.multiplier
	m.currentPos[1] = m.lastPos[1]

	m.deltaApplied = true
}

func (m *MouseDragMonitor) IsActive() bool {
	return m.dragActive
}

func (m *MouseDragMonitor) GetDeltaV() (ov mgl32.Vec2) {
	m.applyDelta()

	ov = m.deltaPos

	m.deltaPos[0] = 0
	m.deltaPos[1] = 0
	return
}

func (m *MouseDragMonitor) GetDelta() (x, y float32) {
	m.applyDelta()
	x = m.deltaPos[0]
	y = m.deltaPos[1]

	m.deltaPos[0] = 0
	m.deltaPos[1] = 0
	return
}
