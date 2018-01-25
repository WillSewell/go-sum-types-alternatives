package main

import "fmt"

type event interface {
	isEvent()
}

type subscribeEvent struct {
	messageChan chan<- string
}

func (subscribeEvent) isEvent() {}

type publishEvent struct {
	message string
}

func (publishEvent) isEvent() {}

type pubsubBus struct {
	subs      []chan<- string
	eventChan chan event // Now only types which implement `event` can be sent
}

func (p *pubsubBus) Run() {
	go func() {
		for event := range p.eventChan {
			// This is not type safe because we might remove a handler, but forget to
			// remove the function which sends the now unhandled event on the channel.
			switch e := event.(type) {
			case subscribeEvent:
				p.subs = append(p.subs, e.messageChan)
			case publishEvent:
				for _, sub := range p.subs {
					sub <- e.message
				}
			default:
				panic(fmt.Sprint("Unknown event type"))
			}
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
