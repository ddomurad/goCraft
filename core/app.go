package core

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type Window struct {
	Width     int
	Hieght    int
	Title     string
	Resizable bool

	glfwWindow *glfw.Window
}

type App struct {
	Window          Window
	EventManager    *EventManager
	ResourceManager *ResourceManager
	ShouldRun       bool
	fpsTime         float64
	fpsCount        int
	lastMoveTime    float64
}

func InitApp(app App) *App {
	var err error

	if err = glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	if app.Window.Width == 0 {
		app.Window.Width = 800
	}

	if app.Window.Hieght == 0 {
		app.Window.Hieght = 600
	}

	app.EventManager = NewEventManager(100)
	app.EventManager.RegisterHandler(&app)
	app.ResourceManager = NewResourceManager()

	glfw.WindowHint(glfw.Resizable, IfThenElse(app.Window.Resizable, glfw.True, glfw.False).(int))
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	app.Window.glfwWindow, err = glfw.CreateWindow(app.Window.Width, app.Window.Hieght, app.Window.Title, nil, nil)

	if err != nil {
		log.Fatalln("failed to create glfw window:", err)
	}

	app.Window.glfwWindow.MakeContextCurrent()

	app.Window.glfwWindow.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		app.EventManager.Push(ResizeEvent{
			Size: [2]int{width, height},
		})
	})

	app.Window.glfwWindow.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		app.EventManager.Push(MouseMoveEvent{
			Pos:  [2]float64{xpos, ypos},
			NPos: [2]float64{xpos / float64(app.Window.Width), ypos / float64(app.Window.Hieght)},
		})
	})

	if err = gl.Init(); err != nil {
		log.Fatalln("failed to initialize GL:", err)
	}

	app.ShouldRun = true
	return &app
}

func (a *App) Close() {
	glfw.Terminate()
}

func (a *App) Run() bool {
	if !a.ShouldRun || a.Window.glfwWindow.ShouldClose() {
		return false
	}

	glfw.PollEvents()

	a.EventManager.Fulsh()
	return true
}

func (a *App) Move(movable Renderable) {
	now := glfw.GetTime()
	dt := now - a.lastMoveTime
	movable.Move(dt, a)

	a.lastMoveTime = now
}

func (a *App) Initialize(initializable Renderable) {
	initializable.Initialize(a)
}

func (a *App) Render(renderable Renderable) {
	if now := glfw.GetTime(); now-a.fpsTime > 1.0 {
		a.EventManager.Push(FpsEvent{
			Fps: a.fpsCount,
		})
		a.fpsTime = now
		a.fpsCount = 0
	}

	renderable.Render(a)
	a.Window.glfwWindow.SwapBuffers()
	a.fpsCount++
}

func (a *App) HandleEvent(e Event) bool {
	switch te := e.(type) {
	case ResizeEvent:
		a.Window.Width = te.Size[0]
		a.Window.Hieght = te.Size[1]
	}
	return false
}
