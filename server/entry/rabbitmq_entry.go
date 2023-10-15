package entry

import (
	"encoding/json"
	"fmt"
	"log"

	data_json "github.com/Fonzeca/UserHub/server/entry/json"
	"github.com/Fonzeca/UserHub/server/entry/manager"
	"github.com/Fonzeca/UserHub/server/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqDataEntry struct {
	inputs <-chan amqp.Delivery
}

var DataEntryManager *manager.DataEntryManager

func NewRabbitMqDataEntry() RabbitMqDataEntry {
	channel := services.GlobalChannel

	q, err := channel.QueueDeclare("userhub", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(q.Name, "notification.zone.userhub.preparing", "carmind", false, nil)
	if err != nil {
		panic(err)
	}

	// Subscribing to QueueService1 for getting messages.
	messages, err := channel.Consume(
		q.Name,    // queue name
		"userhub", // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		log.Println(err)
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatal(err)
	}

	instance := RabbitMqDataEntry{inputs: messages}

	go instance.Run()
	return instance
}

func (m *RabbitMqDataEntry) Run() {
	print("runnning")
	for message := range m.inputs {
		print(message.Body)
		switch message.RoutingKey {
		case "notification.zone.userhub.preparing":
			pojo := data_json.ZoneNotification{}
			err := json.Unmarshal(message.Body, &pojo)
			if err != nil {
				fmt.Println(err)
				break
			}

			err = DataEntryManager.ProcessData(pojo)
			if err != nil {
				fmt.Println(err)
				break
			}
			message.Ack(false)
		}
	}
}
