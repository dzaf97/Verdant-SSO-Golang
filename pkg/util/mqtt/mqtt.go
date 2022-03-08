package mqtt

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	ClientStruct struct {
		mq1     mqtt.Client
		mq2     mqtt.Client
		servers []string
	}
)

var loraBroker mqtt.Client
var nbiotBroker mqtt.Client

func NewMQTTClient(server string) *ClientStruct {

	uname := os.Getenv("MQ_UNAME")
	passwd := os.Getenv("MQ_PASSWD")

	// tls := NewTLSConfig()
	rand.Seed(time.Now().UnixNano())
	log.Println(server)

	client1 := mqtt.NewClientOptions().AddBroker(server).
		SetClientID(fmt.Sprintf("vecto-%d", rand.Int())).
		SetUsername(uname).
		SetPassword(passwd).
		// SetTLSConfig(tls).
		SetDefaultPublishHandler(MsgHandlerBrokerNBIOT).
		SetCleanSession(false)

	mqttObj := mqtt.NewClient(client1)

	if token := mqttObj.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	nbiotBroker = mqttObj
	return &ClientStruct{
		mq1: mqttObj,
	}

}

func (c *ClientStruct) GetClient1() mqtt.Client {
	return c.mq1
}

func GetMqttClient() mqtt.Client {

	return nbiotBroker

}

func MsgHandlerBrokerNBIOT(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Payload())
	payload := msg.Payload()
	// log.Println(payload)
	ParseCanbusData(payload)

}
