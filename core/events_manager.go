package core

type Event interface {
}

type EventHandler interface {
	HandleEvent(e Event) bool
}

type CustomEventHandler struct {
	handler func(e Event) bool
}

type EventManager struct {
	queue         []Event
	eventHandlers []EventHandler
}

type MouseMoveEvent struct {
	Pos  [2]float64
	NPos [2]float64
}

type ResizeEvent struct {
	Size [2]int
}

type FpsEvent struct {
	Fps int
}

func NewCustomEventHandler(handler func(e Event) bool) CustomEventHandler {
	return CustomEventHandler{
		handler: handler,
	}
}

func (eh CustomEventHandler) HandleEvent(e Event) bool {
	return eh.handler(e)
}

func NewEventManager(cap int) *EventManager {
	return &EventManager{
		queue: make([]Event, 0, cap),
	}
}

func (em *EventManager) Push(e Event) {
	em.queue = append(em.queue, e)
}

func (em *EventManager) Fulsh() {
	for i := len(em.queue) - 1; i >= 0; i-- {
		e := em.queue[i]

		for _, eh := range em.eventHandlers {
			if eh.HandleEvent(e) {
				break
			}
		}
		em.queue[i] = nil
	}

	em.queue = em.queue[:0]
}

func (em *EventManager) RegisterHandler(handler EventHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}

func (em *EventManager) RegisterFncHandler(handler func(e Event) bool) {
	em.eventHandlers = append(em.eventHandlers, NewCustomEventHandler(handler))
}

func (em *EventManager) UnregisterHandler(handler EventHandler) {
	panic("Not implemented !")
}
