package main

import (
	"fmt"

	"strings"

	"github.com/pubnub/go/messaging"
)

func main() {
	pubnub := messaging.NewPubnub("demo", "demo", "", "", false, "")
	channel := "t18538"
	msg := &email{
		Sender:    "me@example.com",
		Receivers: []string{"other@example.com", "person@example.com"},
		Title:     "Hello",
		Text:      "Hi",
	}

	successPublish := make(chan []byte)
	errorPublish := make(chan []byte)
	successSubscribe := make(chan []byte)
	errorSubscribe := make(chan []byte)
	awaitConnected := make(chan bool)
	awaitMessage := make(chan bool)

	go pubnub.Subscribe(channel, "", successSubscribe, false, errorSubscribe)
	go func() {
		for {
			select {
			case msg := <-successSubscribe:
				if strings.Contains(string(msg), "connected") {
					fmt.Println("connected...")
					awaitConnected <- true
				} else {
					fmt.Printf("%s\n", msg)
					awaitMessage <- true
				}
			case er := <-errorSubscribe:
				panic(string(er))
			case <-messaging.Timeout():
				panic("Subscribe timeout")
			}
		}
	}()

	<-awaitConnected
	go pubnub.Publish(channel, msg, successPublish, errorPublish)
	select {
	case <-successPublish:
	case er := <-errorPublish:
		panic(string(er))
	case <-messaging.Timeout():
		panic("Publish timeout")
	}

	<-awaitMessage
}

type email struct {
	Sender    string
	Receivers []string
	Title     string
	Text      string
}
