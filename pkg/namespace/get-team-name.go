package namespace

import (
	"context"
	"log"

	"github.com/ministryofjustice/cloud-platform-label-pods/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var systemNamespaces = []string{
	"cloud-platform-label-pods",
	"calico-apiserver",
	"calico-system",
	"cert-manager",
	"concourse",
	"gatekeeper-system",
	"ingress-controllers",
	"kube-system",
	"kuberos",
	"logging",
	"monitoring",
	"tigera-operator",
	"trivy-system",
	"velero",
	"cloud-platform-canary-app-eks",
} // TODO maybe we could get this list from environments (anything that's not in env)?

func InitGetGithubTeamName(getTeamName func(string) (string, error)) func(string) string {
	return func(ns string) string {
		var githubTeamName string
		var err error

		githubTeamName, err = getTeamName(ns)
		if err != nil {
			return "webops"
		}

		isSystemNs := utils.Contains(systemNamespaces, ns)
		if isSystemNs {
			githubTeamName = "webops"
		}

		return githubTeamName
	}
}

func GetTeamNameFromNs(ns string) (string, error) {
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
