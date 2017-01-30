package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
	"os"
	"strings"
)

var noOfWords int = len(os.Args)-1
var ch chan string = make(chan string, noOfWords)

var receiveCallback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	ch <- string(msg.Payload())
}

func main() {
	m := make(map[string]int)
	for n,str := range os.Args[1:] {
		topic := fmt.Sprintf("topic_%s", str)
		opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
		c := MQTT.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := c.Subscribe(topic, 0, receiveCallback); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("Listening for", topic)
		m[str] = n+1
	}
	output := strings.Join(os.Args[1:], " ")
	expect := 1
	lastCheck := time.Now()
	for str := range ch {
		num := m[str]
		if num == expect {
			expect = expect+1
			if (expect == noOfWords) {
				fmt.Println(output, time.Now().Sub(lastCheck))
				expect = 1
				lastCheck = time.Now()
			}
		} else if num != 0 {
			expect = 1
		}
	}
}
