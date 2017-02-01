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
		if m[str] == 0 {
			m[str] = n+1
		} else {
			fmt.Println("Duplicate word:", str);
			os.Exit(1)
		}
	}
	for str := range m {
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
	}
	output := strings.Join(os.Args[1:], " ")
	expect := 1
	firstCheck := time.Now()
	var numOfCompletes int64 = 0
	for str := range ch {
		num := m[str]
		if num == expect {
			expect = expect+1
			if (expect == noOfWords+1) {
				numOfCompletes = numOfCompletes + 1
				avgDelay := int64(time.Now().Sub(firstCheck)) / numOfCompletes
				fmt.Println(output, time.Duration(avgDelay))
				expect = 1
			}
		} else if num == 1 {
			// we have not received what we expected, so drop counter to beginning
			// but we just received word 1, so next should be 2
			// (if we have just 1 word, we never get in this branch)
			expect = 2
		} else if num != 0 {
			// we have not received what we expected, but ignore words not from list
			expect = 1
		}
	}
}
