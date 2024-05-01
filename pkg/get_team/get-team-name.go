package get_team

import (
	"context"
	"log"
	"strings"

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
			return "all-org-members"
		}

		isSystemNs := utils.Contains(systemNamespaces, ns)
		if isSystemNs {
			return "all-org-members"
		}

		return githubTeamName
	}
}

func GetTeamName(ns string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	teamNameFromNS := getTeamNameFromNS(clientset, ns)
	teamNamesFromRBAC := getTeamNameFromRBAC(clientset, ns)

	teamName := teamNameFromNS + " " + teamNamesFromRBAC

	return teamName, nil
}

func getTeamNameFromNS(client *kubernetes.Clientset, ns string) string {
	api, err := client.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return api.Annotations["cloud-platform.justice.gov.uk/team-name"]
}

func getTeamNameFromRBAC(client *kubernetes.Clientset, ns string) string {
	teamNames := ""

	rolebindings, err := client.RbacV1().RoleBindings(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, rb := range rolebindings.Items {
		for _, subj := range rb.Subjects {
			if strings.Contains(subj.Name, "github:") {
				teamName := subj.Name[6 : len(subj.Name)-1]

				teamNames = teamNames + " " + teamName
			}
		}
	}

	return teamNames
}
