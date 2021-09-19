package core

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderable interface {
	Initialize(app *App) error
	Render(app *App)
	Move(dt float64, app *App)
	LateRender() bool
}

type Scene struct {
	ClearColor Color
	Actors     []Renderable
	IsValid    bool

	lateRender []Renderable
}

func (scene *Scene) HandleEvent(e Event) bool {
	return false
}

func (scene *Scene) Initialize(app *App) error {
	gl.ClearColor(scene.ClearColor[0], scene.ClearColor[1], scene.ClearColor[2], scene.ClearColor[3])

	for _, actor := range scene.Actors {
		err := actor.Initialize(app)
		if err != nil {
			return err
		}
	}

	scene.lateRender = make([]Renderable, 0, len(scene.Actors))

	return nil
}

func (scene *Scene) Render(app *App) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, act := range scene.Actors {
		if act.LateRender() {
			scene.lateRender = append(scene.lateRender, act)
		} else {
			act.Render(app)
		}
	}

	for i, act := range scene.lateRender {
		act.Render(app)
		scene.lateRender[i] = nil
	}

	scene.lateRender = scene.lateRender[:0]
}

func (scene *Scene) Move(dt float64, app *App) {
	for _, act := range scene.Actors {
		act.Move(dt, app)
	}
}
func (s *Scene) LateRender() bool {
	return false
}
