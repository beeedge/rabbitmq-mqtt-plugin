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

package options

import (
	"time"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/component-base/config"
)

const (
	RabbitmqMqttPluginUserAgent = "rabbitmq-mqtt-plugin"
)

type Options struct {
	Master     string
	Kubeconfig string
	QPS        float32
	Burst      int
	SyncPeriod int
	config.LeaderElectionConfiguration

	// TODO: add credentials
	BaseCfgPath string
}

func NewRabbitmqMqttPluginOptions() *Options {
	return &Options{
		QPS:        float32(1000),
		Burst:      1000,
		SyncPeriod: 30,
		LeaderElectionConfiguration: config.LeaderElectionConfiguration{
			ResourceName:      RabbitmqMqttPluginUserAgent,
			ResourceNamespace: metav1.NamespaceSystem,
			ResourceLock:      "lease",
			LeaseDuration:     metav1.Duration{Duration: time.Second * time.Duration(15)},
			RenewDeadline:     metav1.Duration{Duration: time.Second * time.Duration(10)},
			RetryPeriod:       metav1.Duration{Duration: time.Second * time.Duration(2)},
		},
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Master, "master", o.Master, "apiserver master address")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "kubeconfig path, empty path means in cluster mode")
	fs.Float32Var(&o.QPS, "kube-qps", o.QPS, "kubeconfig qps setting")
	fs.IntVar(&o.Burst, "kube-burst", o.Burst, "kubeconfig burst setting")
	fs.IntVar(&o.SyncPeriod, "sync-period", o.SyncPeriod, "Period for syncing the objects")

	fs.BoolVar(&o.LeaderElect, "leader-elect", o.LeaderElect, ""+
		"Start a leader election client and gain leadership before "+
		"executing the main loop. Enable this when running replicated "+
		"components for high availability.")
	fs.DurationVar(&o.LeaseDuration.Duration, "leader-elect-lease-duration", o.LeaseDuration.Duration, ""+
		"The duration that non-leader candidates will wait after observing a leadership "+
		"renewal until attempting to acquire leadership of a led but unrenewed leader "+
		"slot. This is effectively the maximum duration that a leader can be stopped "+
		"before it is replaced by another candidate. This is only applicable if leader "+
		"election is enabled.")
	fs.DurationVar(&o.RenewDeadline.Duration, "leader-elect-renew-deadline", o.RenewDeadline.Duration, ""+
		"The interval between attempts by the acting master to renew a leadership slot "+
		"before it stops leading. This must be less than or equal to the lease duration. "+
		"This is only applicable if leader election is enabled.")
	fs.DurationVar(&o.RetryPeriod.Duration, "leader-elect-retry-period", o.RetryPeriod.Duration, ""+
		"The duration the clients should wait between attempting acquisition and renewal "+
		"of a leadership. This is only applicable if leader election is enabled.")
	fs.StringVar(&o.ResourceLock, "leader-elect-resource-lock", o.ResourceLock, ""+
		"The type of resource object that is used for locking during "+
		"leader election. Supported options are `endpoints` (default) and `configmaps`.")
	fs.StringVar(&o.ResourceName, "rabbitmq-mqtt-plugin", o.ResourceName, ""+
		"The name of resource object that is used for locking during "+
		"leader election.")
	fs.StringVar(&o.ResourceNamespace, "leader-elect-resource-namespace", o.ResourceNamespace, ""+
		"The namespace of resource object that is used for locking during "+
		"leader election.")
	fs.StringVar(&o.BaseCfgPath, "base-config", "/etc/rabbitmq-mqtt-plugin/config.yaml", "base configuration path.")
}
