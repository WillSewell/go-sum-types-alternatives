package main

import (
	"fmt"
)

type (
	channel    string
	message    interface{}
	subscriber chan<- message
)

type event interface {
	// Instances now implement the handler in this method
	visit(*pubsubBus)
}

type subscribeEvent struct {
	c channel
	s subscriber
}

func (sE subscribeEvent) visit(p *pubsubBus) {
	p.subs[sE.c] = append(p.subs[sE.c], sE.s)
}

type publishEvent struct {
	c       channel
	message interface{}
}

func (pE publishEvent) visit(p *pubsubBus) {
	for _, sub := range p.subs[pE.c] {
		sub <- pE.message
	}
}

type pubsubBus struct {
	subs      map[channel][]subscriber
	eventChan chan event
}

// Public interface

func (p *pubsubBus) Subscribe(c channel, s subscriber) {
	p.eventChan <- subscribeEvent{c, s}
}

func (p *pubsubBus) Publish(c channel, message interface{}) {
	p.eventChan <- publishEvent{c, message}
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
	bus := pubsubBus{make(map[channel][]subscriber), make(chan event)}
	bus.Run()
	subChan := make(chan message)
	bus.Subscribe("testChan", subChan)
	bus.Publish("testChan", "message")
	msg := <-subChan
	fmt.Println("Received:", msg)
}
