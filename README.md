# Capstone_project
Capstone Project

Kubernetes Cluster Health Checker and Auto-Healing



Project Goals
1. Develop an automated health monitoring system for Kubernetes clusters, focusing on key metrics like node health, pod statuses, and resource utilization.
2. Implement self-healing actions that restart failed pods, reschedule workloads, and, if necessary, trigger scaling events to balance workloads.
3. Provide real-time alerts and notifications to inform the team of critical issues that may require manual intervention.
4. Create a web dashboard to display real-time health status, historical data, and auto-healing logs for transparency and traceability.


Tools Used
- Programming Languages: Go (for performance and Kubernetes API compatibility), Python (for scripting and initial data processing)
- Kubernetes API: For interacting with and managing Kubernetes resources
- Prometheus: For monitoring and collecting metrics from the Kubernetes cluster
- Grafana: For visualizing metrics and health status on a real-time dashboard
- Alertmanager (Prometheus): For sending alerts to Slack or other communication tools
- Slack API: For notifications to DevOps teams
- Docker: To containerize the application and deploy it as a microservice



 Sprint 1: Project Setup and Kubernetes Cluster Access
  - Define project structure and initialize the repository.
  - Set up access to the Kubernetes cluster.
  - Configure basic API access to interact with Kubernetes resources (nodes, pods, services).
  - Install and configure Prometheus for Kubernetes to monitor cluster metrics.
  - Goal: Establish a foundation for the project by setting up the necessary environment, tools, and access to Kubernetes resources.

Sprint 1 Deliverables
•	Repo initialized with project structure
•	Access to Kubernetes cluster (via kubeconfig or in-cluster SA)
•	Minimal Go program that lists nodes & pods
•	Prometheus + Grafana installed in monitoring namespace
	k8s-selfhealer
	k8s-selfhealer
	   git init
	    cmd/selfhealer pkg/ deploy/ docs/
	touch README.md
Initialize Go module
	cd k8s-selfhealer
	go mod init github.com/you/k8s-selfhealer
	go get k8s.io/client-go@v0.29.4
	go get k8s.io/apimachinery@v0.29.4

Run locally - go run cmd/selfhealer/main.go --kubeconfig=$HOME/.kube/config

Apply- 
o	kubectl create ns monitoring
o	kubectl apply -f deploy/rbac.yaml

Install Prometheus + Grafana
o	helm repo add prometheus-community https://prometheus community.github.io/helm-charts
o	helm repo update
o	helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring
Verify:
o	kubectl get pods -n monitoring


Sprint 2: Health Monitoring Module (Node & Pod Checks)
  - Develop a module to monitor the health of nodes and pods using the Kubernetes API.
  - Set up Prometheus metrics collection for node and pod statuses (e.g., CPU, memory usage, pod readiness).
  - Begin integrating Prometheus with Grafana to visualize key metrics.
  - Implement initial threshold-based alerts for critical health issues using Alertmanager.
  - Goal: Enable basic health checks on nodes and pods and set up alerts for critical issues.


Sprint 2 Deliverables
•	Go module that monitors node & pod health and exposes metrics (/metrics)
•	Prometheus scraping your service (via ServiceMonitor)
•	Grafana dashboard panels showing node/pod health
•	Alertmanager rules for “critical” states (e.g., node not ready, pod failures)

Extend Go app with metrics
o	go get github.com/prometheus/client_golang@v1.17.0
o	Update cmd/selfhealer/main.go

• Checks node readiness
• Checks pod readiness
• Exposes Prometheus metrics at :9090/metrics

Expose service to Prometheus
	Create a Service for the app –
deploy/service.yaml
deploy/service-monitor.yaml
	Grafana dashboards
•	Import a new dashboard (create or use templates).
•	Panels to add:
•	Node Health Gauge: query avg(k8s_node_health_status)
•	Pod Readiness: count by (namespace) (k8s_pod_ready_status==0)
•	Failed Pods Over Time: use PromQL with rate(kube_pod_container_status_restarts_total[5m])

Sprint 3: Self-Healing Mechanisms (Pod Recovery)
  - Implement pod restart and rescheduling based on health checks (e.g., if a pod is unresponsive or has failed).
  - Add automated resource cleanup for pods in "CrashLoopBackOff" or "Evicted" states.
  - Test self-healing actions in a staging environment.
  - Log each action for auditing and reporting.
  - Goal: Automate pod recovery processes and ensure that self-healing actions are logged for transparency.

Step 1: Create a “healer” module
Create a new Go file for self-healing logic.
o	New-Item -ItemType File -Path ".\pkg\healer\pod_recovery.go"
Step 2: Update the main file
o	cmd/selfhealer/main.go
Step 3: Build and Dockerize
o	Dockerfile/docker

	docker build -t yourdockerhubusername/k8s-selfhealer:v3 -f .\docker\Dockerfile .
	docker push yourdockerhubusername/k8s-selfhealer:v3


Sprint 4: Advanced Self-Healing (Node Scaling & Resource Balancing)
  - Implement automatic node scaling (up/down) based on resource utilization.
  - Configure horizontal pod autoscaling for workloads with high resource demands.
  - Set up resource balancing by redistributing pods across nodes to optimize cluster usage.
  - Test these features under simulated load to ensure reliability.
  - Goal: Ensure the system can scale and balance resources automatically, optimizing cluster performance and cost.


	Add Prometheus metrics to your self-healer service

	Instrument your code
o	internal/monitor/metrics.go
o	go run main.go


Deploy Prometheus

	kubectl create namespace monitoring
	kubectl apply -f .\deploy\prometheus\prometheus-deploy.yaml


Add - prometheus-deploy.yaml and run

o	kubectl port-forward svc/prometheus -n monitoring 9090:9090


Sprint 5: Alerting and Notification System Integration
  - Integrate Slack or Teams API for real-time notifications on critical issues or auto-healing actions.
  - Configure Alertmanager with customizable alerting rules for different severity levels.
  - Test alerts and notifications to ensure that DevOps teams receive timely updates.
  - Create documentation for setting up and managing alerts.
  - Goal: Implement a comprehensive alerting and notification system to keep the team informed of critical events and actions.
•	Go to 👉 https://api.slack.com/apps
•	Click Create New App → From Scratch
•	Give it a name (e.g. K8sSelfHealerAlerts)
•	Select your Slack workspace.
•	Go to Incoming Webhooks → Activate Incoming Webhooks
•	Click Add New Webhook to Workspace
•	Select a channel (e.g. #devops-alerts)
•	Copy the Webhook URL, e.g.:

Channel - #k8s-self-healer-alerts
URL-https://hooks.slack.com/services/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

 Sprint 6: Web Dashboard and Project Documentation
  - Develop a web dashboard (using Grafana or a custom interface) to display health metrics, self-healing actions, and historical data.
  - Integrate Prometheus metrics and alert logs into the dashboard.
  - Create user documentation, including setup, usage, and troubleshooting guides.
  - Conduct final testing and gather feedback for potential improvements.
  - Goal: Deliver a user-friendly dashboard for monitoring cluster health and document the project for deployment in real-world environments.


	✅ Create a web dashboard
•	Node / Pod health
•	Self-healing actions & logs
•	Historical metrics
	✅ Connect Prometheus metrics and alerts
✅ Write clear documentation (setup, usage, troubleshooting)
Decide the Dashboard Approach
•	Go to http://localhost:3000 (after port-forwarding).
•	Click “+ Create → Dashboard”.
•	Add a Stat Panel:
	Query:
	sum(selfhealer_pods_healed_total)

	Add another Stat Panel:
	sum(selfhealer_pods_failed_total)
	Add another for cluster load:
	Add another for cluster load	
