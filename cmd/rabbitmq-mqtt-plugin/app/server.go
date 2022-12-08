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

package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/beeedge/beethings/cmd/rabbitmq-mqtt-plugin/app/options"
	"github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/config"
	rabbitmqMqttPlugin "github.com/beeedge/beethings/pkg/rabbitmq-mqtt-plugin/rabbitmq-mqtt-plugin"
	"github.com/beeedge/beethings/pkg/util"
	"github.com/beeedge/beethings/pkg/version"
	"github.com/beeedge/beethings/pkg/version/verflag"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	clientset "k8s.io/client-go/kubernetes"
	clientgokubescheme "k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
)

func NewRabbitmqMqttPluginCommand() *cobra.Command {
	o := options.NewRabbitmqMqttPluginOptions()

	cmd := &cobra.Command{
		Use: "rabbitmq-mqtt-plugin",
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()

			klog.Infof("Versions: %#v\n", version.Get())
			util.PrintFlags(cmd.Flags())

			kubeconfig, err := clientcmd.BuildConfigFromFlags(o.Master, o.Kubeconfig)
			if err != nil {
				klog.Fatalf("failed to create kubeconfig: %v", err)
			}

			kubeconfig.QPS = o.QPS
			kubeconfig.Burst = o.Burst

			copyConfig := *kubeconfig
			copyConfig.Timeout = time.Second * o.RenewDeadline.Duration

			leaderElectionClient := clientset.NewForConfigOrDie(restclient.AddUserAgent(&copyConfig, "leader-election"))
			kubeClient := clientset.NewForConfigOrDie(kubeconfig)

			var electionChecker *leaderelection.HealthzAdaptor
			if !o.LeaderElect {
				runRabbitmqMqttPlugin(context.TODO(), o)
				panic("unreachable")
			}

			electionChecker = leaderelection.NewLeaderHealthzAdaptor(time.Second * 20)

			id, err := os.Hostname()
			if err != nil {
				klog.Fatalf("failed to get hostname %v", err)
			}
			id = id + "_" + string(uuid.NewUUID())
			eventRecorder := createRecorder(kubeClient, options.RabbitmqMqttPluginUserAgent)
			rl, err := resourcelock.New(o.ResourceLock, o.ResourceNamespace, o.ResourceName,
				leaderElectionClient.CoreV1(),
				leaderElectionClient.CoordinationV1(),
				resourcelock.ResourceLockConfig{
					Identity:      id,
					EventRecorder: eventRecorder,
				})
			if err != nil {
				klog.Fatalf("error creating lock: %v", err)
			}

			leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
				Lock:          rl,
				LeaseDuration: o.LeaseDuration.Duration,
				RenewDeadline: o.RenewDeadline.Duration,
				RetryPeriod:   o.RetryPeriod.Duration,
				Callbacks: leaderelection.LeaderCallbacks{
					OnStartedLeading: func(ctx context.Context) {
						/*
						 */
						runRabbitmqMqttPlugin(ctx, o)
					},
					OnStoppedLeading: func() {
						klog.Fatalf("leaderelection lost")
					},
				},
				WatchDog: electionChecker,
				Name:     options.RabbitmqMqttPluginUserAgent,
			})
			panic("unreachable")
		},
	}

	fs := cmd.Flags()
	o.AddFlags(fs)

	return cmd
}

func runRabbitmqMqttPlugin(parent context.Context, o *options.Options) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ssConfig, err := config.NewRabbitmqMqttPluginConfig(o)
	fmt.Println(ssConfig)
	if err != nil {
		klog.V(5).Info("Fail to create new server subscription config: %v", err)
		fmt.Println("err:", err)
		return
	}
	ss, err := rabbitmqMqttPlugin.NewRabbitmqMqttPlugin(ssConfig)
	if err != nil {
		fmt.Println("err2:", err)
		klog.V(5).Info("Fail to start new server subscription: %v", err)
		return
	}
	go ss.Run(ctx.Done())
	<-ctx.Done()
}

func createRecorder(kubeClient clientset.Interface, userAgent string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(clientgokubescheme.Scheme, v1.EventSource{Component: userAgent})
}
