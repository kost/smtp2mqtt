package main

import (
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttc mqtt.Client

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func initmqtt() (err error) {
	if *mqttKeep {
		mqttc, err = connect2mqtt(true)
		return err
	}
	return nil
}

func connect2mqtt(init bool) (c mqtt.Client, err error) {
	if !init {
		if *mqttKeep {
			return mqttc, nil
		}
	}
	if *enableDebug {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stderr, "", 0)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(*mqttServer)
	opts.SetUsername(*mqttUser)
	opts.SetPassword(*mqttPassword)
	opts.SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		errt := token.Error()
		log.Printf("Error connecting to MQTT server %s: %v", *mqttServer, errt)
		return c, errt
	}
	return c, nil
}

func disconnectmqtt(c mqtt.Client) {
	if !*mqttKeep {
		c.Disconnect(250)
	}
}

func sanitizeTopic(topic string) (retstr string) {
	retstr = strings.ReplaceAll(topic, "/", "_")
	return
}

func closemqtt() {
	if *mqttKeep {
		mqttc.Disconnect(250)
	}
}

func send2mqtt(topic string, text string) {
	c, _ := connect2mqtt(false)
	token := c.Publish(topic, 0, false, text)
	token.Wait()
	disconnectmqtt(c)
	// time.Sleep(1 * time.Second)
}
