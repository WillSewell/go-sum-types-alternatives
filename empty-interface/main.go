package main

import (
	"fmt"
	"reflect"
)

type (
	channel    string
	message    interface{}
	subscriber chan<- message
)

type subscribeEvent struct {
	c channel
	s subscriber
}

type publishEvent struct {
	c       channel
	message interface{}
}

type pubsubBus struct {
	subs      map[channel][]subscriber
	eventChan chan interface{} // Note the interface{} type for events
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
			// This is not type safe because we might remove a handler, but forget to
			// remove the function which sends the now unhandled event on the channel.
			switch e := event.(type) {
			case subscribeEvent:
				p.handleSubscribe(e)
			case publishEvent:
				p.handlePublish(e)
			default:
				panic(fmt.Sprint("Unknown event type: ", reflect.TypeOf(event)))
			}
		}
	}()
}

func (p *pubsubBus) handleSubscribe(subscribeEvent subscribeEvent) {
	p.subs[subscribeEvent.c] = append(p.subs[subscribeEvent.c], subscribeEvent.s)
}

func (p *pubsubBus) handlePublish(publishEvent publishEvent) {
	for _, sub := range p.subs[publishEvent.c] {
		sub <- publishEvent.message
	}
}

func main() {
	bus := pubsubBus{make(map[channel][]subscriber), make(chan interface{})}
	bus.Run()
	subChan := make(chan message)
	bus.Subscribe("testChan", subChan)
	bus.Publish("testChan", "message")
	msg := <-subChan
	fmt.Println("Received:", msg)
	bus.eventChan <- "This will panic!"
	<-subChan
}
