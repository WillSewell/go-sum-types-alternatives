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
	bus := pubsubBus{make([]chan<- string, 0), make(chan interface{})}
	bus.Run()
	messageChan := make(chan string, 1)
	bus.eventChan<-subscribeEvent{messageChan}
	bus.eventChan<-publishEvent{"message"}
	msg := <-messageChan
	fmt.Println("Received:", msg)
	bus.eventChan <- "This will panic!"
	<-messageChan
}
