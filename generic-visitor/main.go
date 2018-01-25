package main

import "fmt"


type event interface {
	visit(v eventVisitor)
}

// The handlers for each event type are defined in instances of this struct
type eventVisitor struct {
	// Notice these can have different function signatures
	visitSubscribe func(subscribeEvent)
	visitPublish   func(publishEvent)
}

type subscribeEvent struct {
	messageChan chan<- string
}

func (s subscribeEvent) visit(v eventVisitor) {
	v.visitSubscribe(s)
}

type publishEvent struct {
	message string
}

func (p publishEvent) visit(v eventVisitor) {
	v.visitPublish(p)
}

type pubsubBus struct {
	subs      []chan<- string
	eventChan chan event
}

func (p *pubsubBus) Run() {
	go func() {
		for event := range p.eventChan {
			p.handleEvent(event)
		}
	}()
}

// Handler implementations are passed in here.
// Alternative handler implementations could be defined by creating an
// alternative version of this method.
func (p *pubsubBus) handleEvent(e event) {
	e.visit(eventVisitor{
		visitSubscribe: func(sE subscribeEvent) {
			p.subs = append(p.subs, sE.messageChan)
		},
		visitPublish: func(pE publishEvent) {
			for _, sub := range p.subs {
				sub <- pE.message
			}
		},
	})
}

func main() {
	bus := pubsubBus{make([]chan<- string, 0), make(chan event)}
	bus.Run()
	messageChan := make(chan string, 1)
	bus.eventChan<-subscribeEvent{messageChan}
	bus.eventChan<-publishEvent{"message"}
	msg := <-messageChan
	fmt.Println("Received:", msg)
}
