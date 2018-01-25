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

func (p *pubsubBus) run() {
	for event := range p.eventChan {
		// Handler implementations are passed in here.
		// Alternative handler implementations could be defined by
		// creating an alternative version of this method.
		event.visit(eventVisitor{
			visitSubscribe: p.handleSubscribe,
			visitPublish: p.handlePublish,
		})
	}
}

func (p *pubsubBus) handleSubscribe(subscribeEvent subscribeEvent) {
	p.subs = append(p.subs, subscribeEvent.messageChan)
}

func (p *pubsubBus) handlePublish(publishEvent publishEvent) {
	for _, sub := range p.subs {
		sub <- publishEvent.message
	}
}

func main() {
	bus := pubsubBus{make([]chan<- string, 0), make(chan event)}
	go bus.run()
	messageChan := make(chan string, 1)
	bus.eventChan<-subscribeEvent{messageChan}
	bus.eventChan<-publishEvent{"message"}
	msg := <-messageChan
	fmt.Println("Received:", msg)
}
