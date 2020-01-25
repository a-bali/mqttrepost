package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

type Topic struct {
	Source string
	Target string
	JSON   []struct {
		Source string
		Target string
	}
}
type Config struct {
	Debug  bool
	Server struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	Topics []Topic
}

var config Config
var topics = make(map[string]Topic)
var client mqtt.Client

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " <configfile>")
		os.Exit(1)
	}

	filename, _ := filepath.Abs(os.Args[1])
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	for _, v := range config.Topics {
		topics[v.Source] = v
	}
	debug("Loaded config: " + fmt.Sprintf("%+v", config))

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", "tcp", config.Server.Host, config.Server.Port))
	opts.SetUsername(config.Server.User)
	opts.SetPassword(config.Server.Password)
	opts.OnConnect = func(c mqtt.Client) {
		topics := make(map[string]byte)
		for _, t := range config.Topics {
			topics[t.Source] = byte(0)
		}
		if token := c.SubscribeMultiple(topics, onMessageReceived); token.Wait() && token.Error() != nil {
			log("Unable to subscribe to MQTT: " + token.Error().Error())
		} else {
			for k := range topics {
				log("Subscribed to " + k)
			}
		}
	}

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log("Unable to connect to MQTT: " + token.Error().Error())
	} else {
		log("Connected to MQTT server at " + opts.Servers[0].String())
	}

	for {
		select {}
	}

}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {

	debug("MQTT receive: " + message.Topic() + ": " + string(message.Payload()))

	if topics[message.Topic()].Target != "" {
		postMQTT(topics[message.Topic()].Target, string(message.Payload()))
	}

	if len(topics[message.Topic()].JSON) > 0 {
		for _, m := range topics[message.Topic()].JSON {
			v := gjson.Get(string(message.Payload()), m.Source)
			if v.Exists() {
				postMQTT(m.Target, v.String())
			}
		}

	}

}

func postMQTT(topic string, value string) {
	debug("MQTT send: " + topic + ": " + value)
	client.Publish(topic, 0, false, value)
}

func log(s string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), s)
}

func debug(s string) {
	if config.Debug {
		log("(" + s + ")")
	}
}
