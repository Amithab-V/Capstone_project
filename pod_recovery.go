import "k8s-selfhealer/pkg/notifier"

const slackWebhook = "https://hooks.slack.com/services/REPLACE_THIS"

func deletePod(clientset *kubernetes.Clientset, pod *v1.Pod) {
	err := clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Error deleting pod %s: %v", pod.Name, err)
	} else {
		msg := fmt.Sprintf("🤖 Self-Healer deleted pod *%s* in namespace *%s*", pod.Name, pod.Namespace)
		notifier.SendSlackAlert(slackWebhook, msg)
		log.Printf("✅ Deleted pod %s in namespace %s", pod.Name, pod.Namespace)
	}
}
