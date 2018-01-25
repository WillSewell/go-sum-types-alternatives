package main

import "fmt"

type event interface {
	// Instances now implement the handler in this method
	visit(*pubsubBus)
}

type subscribeEvent struct {
	messageChan chan<- string
}

func (sE subscribeEvent) visit(p *pubsubBus) {
	p.handleSubscribe(sE)
}

type publishEvent struct {
	message string
}

func (pE publishEvent) visit(p *pubsubBus) {
	p.handlePublish(pE)
}

type pubsubBus struct {
	subs      []chan<- string
	eventChan chan event
}

func (p *pubsubBus) run() {
	for event := range p.eventChan {
		// Type switch is not required, so it's type-safe
		event.visit(p)
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
	bus.run()
	messageChan := make(chan string, 1)
	bus.eventChan<-subscribeEvent{messageChan}
	bus.eventChan<-publishEvent{"message"}
	msg := <-messageChan
	fmt.Println("Received:", msg)
}
