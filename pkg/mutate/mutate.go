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

func createAdmReviewSucc(admReview v1.AdmissionReview, resp v1.AdmissionResponse, githubTeamName string, createAnnoField bool) ([]byte, error) {
	resp.AuditAnnotations = map[string]string{
		"metadata.annotations.github_teams": "mutation added for identification",
	}

	// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
	// tell K8S how it should modifiy it
	pT := v1.PatchTypeJSONPatch
	resp.PatchType = &pT

	type Annotations struct {
		Value string `json:"value,omitempty"`
	}

	type Patch struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value any    `json:"value"`
	}

	p := []Patch{}

	if createAnnoField {
		patch := Patch{
			"add",
			"/metadata/annotations",
			new(Annotations),
		}
		p = append(p, patch)
	}

	patch := Patch{
		"add",
		"/metadata/annotations/github_teams",
		githubTeamName,
	}
	p = append(p, patch)

	log.Printf("patch: %s\n", p)

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

	annoFieldExists := true

	var pod *corev1.Pod

	admReview := v1.AdmissionReview{}
	resp := v1.AdmissionResponse{
		Allowed: true,
	}

	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, err
	}

	admReq := admReview.Request

	if admReq == nil {
		responseBody, failErr := createAdmReviewFail(admReview, resp, "AdmissionReview request body is nil")
		if failErr != nil {
			return nil, failErr
		}
		log.Printf("resp adm review fail admReq nil: %s\n", string(responseBody))
		return responseBody, nil
	}

	resp.UID = admReq.UID

	if err := json.Unmarshal(admReq.Object.Raw, &pod); err != nil {
		responseBody, failErr := createAdmReviewFail(admReview, resp, fmt.Sprintf("unable unmarshal pod json object %v", err.Error()))
		if failErr != nil {
			return nil, failErr
		}
		log.Printf("resp unable to unmarshall pod json: %s\n", string(responseBody))
		return responseBody, nil
	}

	if pod.Annotations == nil {
		annoFieldExists = false
	}

	githubTeamName := getGithubTeamName(pod.GetNamespace())

	responseBody, succErr := createAdmReviewSucc(admReview, resp, githubTeamName, !annoFieldExists)
	if succErr != nil {
		return nil, succErr
	}

	log.Printf("resp success: %s\n", string(responseBody))

	return responseBody, nil
}
