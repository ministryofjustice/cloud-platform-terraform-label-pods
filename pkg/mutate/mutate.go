// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it, it also makes testing Mutate() kind of easy w/o need for a fake http server, etc.
package mutate

import (
	"encoding/json"
	"fmt"

	"github.com/ministryofjustice/cloud-platform-label-pods/pkg/namespace"
	"github.com/ministryofjustice/cloud-platform-label-pods/utils"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var systemNamespaces = []string{} // TODO maybe we could get this list from environments?

func Mutate(body []byte) ([]byte, error) {
	admReview := v1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}
	admReq := admReview.Request
	resp := v1.AdmissionResponse{}

	if admReq != nil {
		// get the Pod object and unmarshal it into its struct
		if err := json.Unmarshal(admReq.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}

		isSystemNs := utils.Contains(systemNamespaces, pod.GetNamespace())
		if isSystemNs {
			return nil, fmt.Errorf("do not need to label this namespace as it is not a user namespace")
		}

		githubTeamName, nsErr := namespace.GetTeamName(pod.GetNamespace())
		if nsErr != nil {
			return nil, fmt.Errorf("unable to get pod namespace %v", nsErr)
		}

		// set response options
		resp.Allowed = true
		resp.UID = admReq.UID
		pT := v1.PatchTypeJSONPatch
		resp.PatchType = &pT

		// add some audit annotations, helpful to know why a object was modified
		resp.AuditAnnotations = map[string]string{
			"metadata.label.github_teams": "mutation added for identification",
		}

		// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
		// tell K8S how it should modifiy it
		p := []map[string]string{}
		patch := map[string]string{
			"op":    "add",
			"path":  "/metadata/labels/github_teams",
			"value": githubTeamName,
		}
		p = append(p, patch)

		// parse the []map into JSON
		resp.Patch, err = json.Marshal(p)

		// Success, of course ;)
		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		// back into JSON so we can return the finished AdmissionReview w/ Response directly
		// w/o needing to convert things in the http handler
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err // untested section
		}
	}

	return responseBody, nil
}
