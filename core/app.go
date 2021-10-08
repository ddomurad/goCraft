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
	Width  int
	Height int

	glfwWindow *glfw.Window
}

type App struct {
	Window          Window
	EventManager    *EventManager
	ResourceManager *ResourceManager
	ShouldRun       bool
	fpsTime         float64
	fpsCount        int
	lastRnderTime   float64
}

type Renderable interface {
	Render(dt float64, app *App)
}

func InitApp(title string, width, height int, resizable bool, syncSwap bool) (app *App) {
	var err error

	if err = glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	app = &App{}
	app.EventManager = NewEventManager(100)
	app.EventManager.RegisterHandler(app)
	app.ResourceManager = NewResourceManager()

	glfw.WindowHint(glfw.Resizable, IfThenElse(resizable, glfw.True, glfw.False).(int))
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	app.Window = Window{
		Width:  width,
		Height: height,
	}

	app.Window.Height = height
	app.Window.glfwWindow, err = glfw.CreateWindow(app.Window.Width, app.Window.Height, title, nil, nil)

	if err != nil {
		log.Fatalln("failed to create glfw window:", err)
	}

	app.Window.glfwWindow.MakeContextCurrent()

	app.Window.glfwWindow.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		app.EventManager.Push(ResizeEvent{
			Size: [2]int{width, height},
		})
	})

	app.Window.glfwWindow.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		app.EventManager.Push(MouseButtonEvent{
			Button: button,
			Action: action,
			Mods:   mods,
		})
	})

	app.Window.glfwWindow.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		app.EventManager.Push(MouseScrollEvent{
			X: xoff,
			Y: yoff,
		})
	})

	app.Window.glfwWindow.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		app.EventManager.Push(MouseMoveEvent{
			Pos:  [2]float64{xpos, ypos},
			NPos: [2]float64{xpos / float64(app.Window.Width), ypos / float64(app.Window.Height)},
		})
	})

	if err = gl.Init(); err != nil {
		log.Fatalln("failed to initialize GL:", err)
	}

	glfw.SwapInterval(IfThenElse(syncSwap, 1, 0).(int))

	app.ShouldRun = true

	return
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

func (a *App) Render(renderable Renderable) {
	now := glfw.GetTime()

	if now-a.fpsTime > 1.0 {
		a.EventManager.Push(FpsEvent{
			Fps: a.fpsCount,
		})
		a.fpsTime = now
		a.fpsCount = 0
	}

	dt := now - a.lastRnderTime

	renderable.Render(dt, a)
	a.Window.glfwWindow.SwapBuffers()
	a.fpsCount++
	a.lastRnderTime = now
}

func (a *App) HandleEvent(e Event) bool {
	switch te := e.(type) {
	case ResizeEvent:
		a.Window.Width = te.Size[0]
		a.Window.Height = te.Size[1]
	}
	return false
}
