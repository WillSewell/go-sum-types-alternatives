package main

import "fmt"

type subscribeEvent struct {
	messageChan chan<- string
}

type publishEvent struct {
	message string
}

type pubsubBus struct {
	subs      []chan<- string
	eventChan chan interface{} // Note the interface{} type for events
}

func (p *pubsubBus) run() {
	for event := range p.eventChan {
		// This is not type safe because we might remove a handler, but forget to
		// remove the function which sends the now unhandled event on the channel.
		switch e := event.(type) {
		case subscribeEvent:
			p.handleSubscribe(e)
		case publishEvent:
			p.handlePublish(e)
		default:
			panic(fmt.Sprint("Unknown event type"))
		}
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
	bus := pubsubBus{make([]chan<- string, 0), make(chan interface{})}
	go bus.run()
	messageChan := make(chan string, 1)
	bus.eventChan<-subscribeEvent{messageChan}
	bus.eventChan<-publishEvent{"message"}
	msg := <-messageChan
	fmt.Println("Received:", msg)
	bus.eventChan <- "This will panic!"
	<-messageChan
}
