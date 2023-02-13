package main

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"ssmk8s/k8s"
	"time"
)

func main() {
	path := os.Args[1]
	deployName := os.Args[2]
	log.Println("Listener " + deployName)
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientSet := kubernetes.NewForConfigOrDie(config)
	informer := informers.NewSharedInformerFactory(clientSet, time.Second*30).Apps().V1().Deployments().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//deployment := obj.(*v1.Deployment)
			//fmt.Println(time.Now(), "Deployment Added", deployment.Name, "Namespace: ", deployment.Namespace)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldDeployment := oldObj.(*v1.Deployment)
			newDeployment := newObj.(*v1.Deployment)
			if deployName == oldDeployment.Name && (oldDeployment.Spec.Template.Spec.Containers[0].Image != newDeployment.Spec.Template.Spec.Containers[0].Image) {
				fmt.Println(time.Now(), "Deployment Image Changed:", oldDeployment.Name)
				k8s.CreateCMWithStoreParameters(clientSet, oldDeployment.Name, oldDeployment.Namespace, path)

				newEnvFrom := []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: newDeployment.Name,
							},
						},
					},
				}

				newDeployment.Spec.Template.Spec.Containers[0].EnvFrom = newEnvFrom
				k8s.UpdateDeployment(clientSet, newDeployment)
			}
			/*else {
				fmt.Println(time.Now(), "Deployment Changed But Not Interest", oldDeployment.Name, "Namespace:", oldDeployment.Namespace)
			}*/
		},
		DeleteFunc: func(obj interface{}) {
			deployment := obj.(*v1.Deployment)
			fmt.Println("Deployment Deleted:", deployment.Name)
		},
	})
	informer.Run(stopper)
}
