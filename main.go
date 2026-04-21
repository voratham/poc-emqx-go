package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var testHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("🦄 topic=%s, payload=%s", msg.Topic(), string(msg.Payload()))
}

var count = 0

func main() {
	log.Println("🟢 start program")

	stopPublisher := false

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-mqtt-client")

	// how to singla term
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	opts.SetDefaultPublishHandler(testHandler)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic("🔴 client.Connect Error : ", token.Error())
	}

	go func() {
		if token := client.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
			log.Panic("🔴 client.Subscribe Error : ", token.Error())
		}
	}()

	go func() {
		<-signalChan
		log.Println("🔴 Received an interrupt, stopping publisher...")
		stopPublisher = true
		if token := client.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
			log.Panic("🔴 client.Unsubscribe Error : ", token.Error())
		}
	}()

	for !stopPublisher {
		token := client.Publish("testtopic/1", 0, false, fmt.Sprintf("Hello world from Go! %d", count))
		token.Wait()
		count++
		time.Sleep(time.Second * 1)

	}

	client.Disconnect(250) // 250 = 250ms to wait for pending work to complete before disconnecting
	log.Println("🟢 end program")

}
