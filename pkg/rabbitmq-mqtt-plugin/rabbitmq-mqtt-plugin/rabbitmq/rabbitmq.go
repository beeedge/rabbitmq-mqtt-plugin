package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/config"
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/model"
	"github.com/streadway/amqp"
	"k8s.io/klog"
)

/*
	Step1: Connect to rabbitmq server
	Step2: Subscription queue message,then deliver message to mqtt_process goroutine
	Step3: Create a goroutine in order to detect connection close
*/
func RabbitMQProcess(cfg *config.RabbitMQConfig, message chan<- model.RabbitMQMsg) {
	// Connect rabbitmq
	msgs, conn, err := connect(cfg)
	defer conn.Close()
	if err != nil {
		klog.Errorf("Failed to connect rabbitmq: %s\n", err.Error())
		return
	}
	if err != nil {
		klog.Errorf("Failed to create msg: %s\n", err.Error())
		return
	}
	klog.Info("Connected to rabbitmq server")

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case msg := <-msgs:
			var body model.RabbitMQMsg
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				klog.Errorf("Unmarshal json error: %s", err.Error())
				continue
			}
			message <- body
		case <-ticker.C:
			if conn.IsClosed() {
				// Connect again
				msgs, conn, err = connect(cfg)
				if err != nil {
					klog.Errorf("Connect again error: %s", err.Error())
					continue
				}
			}
		}
	}
}

func connect(cfg *config.RabbitMQConfig) (<-chan amqp.Delivery, *amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Addr, cfg.Port))
	if err != nil {
		klog.Errorf("Failed to connect to rabbitmq: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Connection{}, err
	}
	channel, err := conn.Channel()
	if err != nil {
		klog.Errorf("Failed to create rabbitmq channel: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Connection{}, err
	}
	if err = channel.ExchangeDeclarePassive(
		cfg.Exchange,
		cfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		klog.Errorf("Failed to declare rabbitmq exchange: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Connection{}, err
	}
	queue, err := channel.QueueDeclare(
		cfg.Queue,
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("Failed to declare queue: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Connection{}, err
	}
	if err = channel.QueueBind(queue.Name, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
		klog.Errorf("Failed to bind queue: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Connection{}, err
	}
	msgs, err := channel.Consume(
		cfg.Queue,
		cfg.ConsumerTag,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("Failed to create msg: %s\n", err.Error())
	}
	return msgs, conn, err
}
