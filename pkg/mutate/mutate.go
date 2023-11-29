// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it
package mutate

import (
	"encoding/json"
	"fmt"
	"log"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createAdmReviewFail(admReview v1.AdmissionReview, resp v1.AdmissionResponse, failureMsg string) ([]byte, error) {
	resp.Allowed = false
	resp.Result = &metav1.Status{
		Status:  "Failure",
		Message: failureMsg,
	}

	admReview.Response = &resp
	responseBody, err := json.Marshal(admReview)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func createAdmReviewSucc(admReview v1.AdmissionReview, resp v1.AdmissionResponse, githubTeamName string) ([]byte, error) {
	resp.AuditAnnotations = map[string]string{
		"metadata.label.github_teams": "mutation added for identification",
	}

	// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
	// tell K8S how it should modifiy it
	pT := v1.PatchTypeJSONPatch
	resp.PatchType = &pT
	p := []map[string]string{}
	patch := map[string]string{
		"op":    "add",
		"path":  "/metadata/labels/github_teams",
		"value": githubTeamName,
	}
	p = append(p, patch)

	respPatch, patchErr := json.Marshal(p)
	if patchErr != nil {
		return nil, patchErr
	}

	resp.Patch = respPatch
	resp.Result = &metav1.Status{
		Status: "Success",
	}

	admReview.Response = &resp

	// back into JSON so we can return the finished AdmissionReview w/ Response directly
	// w/o needing to convert things in the http handler
	responseBody, err := json.Marshal(admReview)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func Mutate(body []byte, getGithubTeamName func(string) string) ([]byte, error) {
	log.Printf("recv: %s\n", string(body))

	var pod *corev1.Pod

	admReview := v1.AdmissionReview{}
	resp := v1.AdmissionResponse{
		Allowed: true,
	}

	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, err
	}

	admReq := admReview.Request
	resp.UID = admReq.UID

	if admReq == nil {
		responseBody, failErr := createAdmReviewFail(admReview, resp, "AdmissionReview request body is nil")
		if failErr != nil {
			return nil, failErr
		}
		return responseBody, nil
	}

	if err := json.Unmarshal(admReq.Object.Raw, &pod); err != nil {
		responseBody, failErr := createAdmReviewFail(admReview, resp, fmt.Sprintf("unable unmarshal pod json object %v", err.Error()))
		if failErr != nil {
			return nil, failErr
		}
		return responseBody, err
	}

	githubTeamName := getGithubTeamName(pod.GetNamespace())

	responseBody, succErr := createAdmReviewSucc(admReview, resp, githubTeamName)
	if succErr != nil {
		return nil, succErr
	}

	log.Printf("resp: %s\n", string(responseBody))

	return responseBody, nil
}
