package namespace

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetTeamName(ns string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	api, err := clientset.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	teamName := api.Annotations["cloud-platform.justice.gov.uk/team-name"]

	log.Println("teamName", teamName)

	return teamName, nil
}
