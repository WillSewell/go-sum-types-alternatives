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
	visit(v eventVisitor)
}

// The handlers for each event type are defined in instances of this struct
type eventVisitor struct {
	// Notice these can have different function signatures
	visitSubscribe func(subscribeEvent)
	visitPublish   func(publishEvent)
}

type subscribeEvent struct {
	c channel
	s subscriber
}

func (s subscribeEvent) visit(v eventVisitor) {
	v.visitSubscribe(s)
}

type publishEvent struct {
	c       channel
	message interface{}
}

func (p publishEvent) visit(v eventVisitor) {
	v.visitPublish(p)
}

// Handler implementations are passed in here.
// Alternative handler implementations could be defined by creating an
// alternative version of this method.
func (p *pubsubBus) handleEvent(e event) {
	e.visit(eventVisitor{
		visitSubscribe: func(sE subscribeEvent) {
			p.subs[sE.c] = append(p.subs[sE.c], sE.s)
		},
		visitPublish: func(pE publishEvent) {
			for _, sub := range p.subs[pE.c] {
				sub <- pE.message
			}
		},
	})
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
			p.handleEvent(event)
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
