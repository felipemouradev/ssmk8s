package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"ssmk8s/utils"
	"time"
)

func CreateCMWithStoreParameters(clientSet *kubernetes.Clientset, name string, namespace string, path string) {
	var ssmValues map[string]string
	json.Unmarshal([]byte(utils.GetPathStoreParameters(path)), &ssmValues)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: ssmValues,
	}

	resultCreate, errCreate := clientSet.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if errCreate != nil {
		fmt.Println("Failed to create configmap: ", errCreate, "... Trying update...")
		resultUpdate, errUpdate := clientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
		if errUpdate != nil {
			fmt.Println("Failed to update configmap: ", errUpdate)
			return
		}
		fmt.Printf("Updated configmap: %q\n", resultUpdate.GetObjectMeta().GetName())
		return
	}
	fmt.Printf("Created configmap: %q\n", resultCreate.GetObjectMeta().GetName())

}

func UpdateDeployment(clientSet *kubernetes.Clientset, deployment *v1.Deployment) {
	deploymentClient := clientSet.AppsV1().Deployments(deployment.Namespace)
	deployment.CreationTimestamp = metav1.Time{time.Now()}
	deployment.ResourceVersion = ""
	deployment.SelfLink = ""
	deployment.UID = ""
	result, err := deploymentClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error on update deployment: ", err)
		return
	}
	fmt.Println("Updated deployment: ", result.GetObjectMeta().GetName())
}
