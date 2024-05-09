package get_team

import (
	"context"
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
	"overprovision",
} // TODO maybe we could get this list from environments (anything that's not in env)?

func InitGetGithubTeamName(getTeamName func(string) (string, error)) func(string) string {
	return func(ns string) string {
		var githubTeamName string
		var err error

		isSystemNs := utils.Contains(systemNamespaces, ns)
		if isSystemNs {
			return "all-org-members"
		}

		githubTeamName, err = getTeamName(ns)
		if err != nil {
			return "all-org-members"
		}

		if githubTeamName == "" {
			return "all-org-members"
		}

		return githubTeamName
	}
}

func GetTeamName(ns string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	teamNamesFromRBAC, rbacErr := getTeamNameFromRBAC(clientset, ns)
	if rbacErr != nil {
		return "", rbacErr
	}

	return teamNamesFromRBAC, nil
}

func getTeamNameFromRBAC(client *kubernetes.Clientset, ns string) (string, error) {
	teamNames := ""

	rolebindings, err := client.RbacV1().RoleBindings(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, rb := range rolebindings.Items {
		for i, subj := range rb.Subjects {
			if strings.Contains(subj.Name, "github:") {
				teamName := subj.Name[7:len(subj.Name)]

				if i == 0 {
					teamNames = teamName
					continue
				}
				teamNames = teamNames + "_" + teamName
			}
		}
	}

	return teamNames, nil
}
