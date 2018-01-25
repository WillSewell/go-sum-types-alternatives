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
	p.subs = append(p.subs, sE.messageChan)
}

type publishEvent struct {
	message string
}

func (pE publishEvent) visit(p *pubsubBus) {
	for _, sub := range p.subs {
		sub <- pE.message
	}
}

type pubsubBus struct {
	subs      []chan<- string
	eventChan chan event
}

func (p *pubsubBus) Run() {
	go func() {
		for event := range p.eventChan {
			// Type switch is not required, so it's type-safe
			event.visit(p)
		}
	}()
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
