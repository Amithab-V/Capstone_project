package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	nodeHealthGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_node_health_status",
			Help: "Node health (1=Ready, 0=NotReady)",
		},
		[]string{"node"},
	)
	podReadyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_pod_ready_status",
			Help: "Pod readiness (1=Ready, 0=NotReady)",
		},
		[]string{"namespace", "pod"},
	)
)

func getClient() (*kubernetes.Clientset, error) {
	kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig")
	flag.Parse()
	if *kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
		return kubernetes.NewForConfig(cfg)
	}
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func recordMetrics(client *kubernetes.Clientset) {
	for {
		// --- Node health ---
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			for _, n := range nodes.Items {
				ready := 0.0
				for _, cond := range n.Status.Conditions {
					if cond.Type == "Ready" && cond.Status == "True" {
						ready = 1.0
					}
				}
				nodeHealthGauge.WithLabelValues(n.Name).Set(ready)
			}
		}

		// --- Pod readiness ---
		pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			for _, p := range pods.Items {
				ready := 0.0
				for _, cond := range p.Status.Conditions {
					if cond.Type == "Ready" && cond.Status == "True" {
						ready = 1.0
					}
				}
				podReadyGauge.WithLabelValues(p.Namespace, p.Name).Set(ready)
			}
		}

		time.Sleep(30 * time.Second) // scrape interval
	}
}

func main() {
	client, err := getClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to connect: %v\n", err)
		os.Exit(1)
	}

	prometheus.MustRegister(nodeHealthGauge, podReadyGauge)

	go recordMetrics(client)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("📡 Serving metrics on :9090/metrics")
	http.ListenAndServe(":9090", nil)
}


package main

import (
	"log"
	"k8s-selfhealer/pkg/healer"
)

func main() {
	log.Println("🚀 Starting Kubernetes Self-Healer...")
	go healer.HealPods() // Run the self-healing loop

	select {} // keep running forever
}
