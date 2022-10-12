package mqtt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/config"
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog"
)

/*
	Step 1: Connect to mqtt server
	Step 2: Send message to mqtt server when receive message from rabbitmq server
	Step 3: Detect connection status every 5 minutes; If connection closed,will try to connect again
*/
func MqttProcess(cfg *config.MqttConfig, messageChan <-chan model.RabbitMQMsg) {
	// Connect mqtt client
	brokerAddr := fmt.Sprintf("tcp://%s:%d", cfg.Addr, cfg.Port)
	opts := mqtt.NewClientOptions().AddBroker(brokerAddr).SetClientID("rabbitmq_mqtt_plugin")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetProtocolVersion(4)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		klog.Infof("Connect to mqtt server error:", token.Error())
	}
	klog.Info("Connected mqtt server")

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case msg := <-messageChan:
			topic := strings.ReplaceAll(msg.Topic, ".", "/")
			bytes, _ := json.Marshal(msg)
			token := client.Publish(topic, byte(cfg.Qos), false, bytes)
			token.Wait()
		case <-ticker.C:
			// 5 minute detect connection
			if !client.IsConnectionOpen() {
				// Connect again
				if token := client.Connect(); token.Wait() && token.Error() != nil {
					klog.Info("False to connect,connect again")
					continue
				}
			}
		}
	}
}
