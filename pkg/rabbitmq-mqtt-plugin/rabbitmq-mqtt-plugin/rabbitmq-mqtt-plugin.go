/*
Copyright 2022 The BeeThings Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rabbitmqMqttPlugin

import (
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/config"
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/model"
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/rabbitmq-mqtt-plugin/mqtt"
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/rabbitmq-mqtt-plugin/rabbitmq"
	"k8s.io/klog"
)

type RabbitmqMqttPlugin struct {
	config *config.Config
}

func NewRabbitmqMqttPlugin(cfg *config.Config) (*RabbitmqMqttPlugin, error) {
	return &RabbitmqMqttPlugin{config: cfg}, nil
}

func (rmp *RabbitmqMqttPlugin) Run(stopCh <-chan struct{}) {
	klog.Info("Run rabbitmq mqtt plugin")
	stopch := make(chan interface{})
	messageChan := make(chan model.RabbitMQMsg, 10000)
	// Handle rabbitmq message
	go rabbitmq.RabbitMQProcess(rmp.config.RabbitMQ, messageChan)
	// Send to mqtt broker
	go mqtt.MqttProcess(rmp.config.Mqtt, messageChan)
	<-stopch
}
