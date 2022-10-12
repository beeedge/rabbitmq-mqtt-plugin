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
	msgs, channel, err := connect(cfg)
	defer channel.Close()
	if err != nil {
		klog.Errorf("Failed to connect rabbitmq: %s\n", err.Error())
		return
	}
	if err != nil {
		klog.Errorf("Failed to create msg: %s\n", err.Error())
		return
	}
	klog.Info("Connected to rabbitmq server")

	// goroutine detect connection close
	closeChan := make(chan bool, 1)
	go connectCloseNotice(channel, closeChan)
	for {
		select {
		case msg := <-msgs:
			var body model.RabbitMQMsg
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				klog.Errorf("Unmarshal json error: %s", err.Error())
				continue
			}
			message <- body
		case <-closeChan:
			// Connect again
			msgs, channel, err = connect(cfg)
			if err != nil {
				klog.Errorf("Connect again error: %s", err.Error())
				klog.Info("Connect rabbitmq fail,sleep 5 minute")
				time.Sleep(5 * time.Minute)
				go connectCloseNotice(channel, closeChan)
				continue
			}
		}
	}
}

func connectCloseNotice(channel *amqp.Channel, closeChan chan<- bool) {
	cc := make(chan *amqp.Error)
	err := <-channel.NotifyClose(cc)
	klog.Errorf("Rabbitmq connection close: %s\n", err.Error())
	closeChan <- true
}

func connect(cfg *config.RabbitMQConfig) (<-chan amqp.Delivery, *amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Addr, cfg.Port))
	if err != nil {
		klog.Errorf("Failed to connect to rabbitmq: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Channel{}, err
	}
	channel, err := conn.Channel()
	if err != nil {
		klog.Errorf("Failed to create rabbitmq channel: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Channel{}, err
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
		return make(<-chan amqp.Delivery), &amqp.Channel{}, err
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
		return make(<-chan amqp.Delivery), &amqp.Channel{}, err
	}
	if err = channel.QueueBind(queue.Name, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
		klog.Errorf("Failed to bind queue: %s\n", err.Error())
		return make(<-chan amqp.Delivery), &amqp.Channel{}, err
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
	return msgs, channel, err
}
