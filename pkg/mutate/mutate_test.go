// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it
package mutate

import (
	"reflect"
	"testing"
)

var validInput string = `{ "kind": "AdmissionReview", "apiVersion": "admission.k8s.io/v1beta1", "request": { "uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c", "kind": { "group": "", "version": "v1", "kind": "Pod" }, "resource": { "group": "", "version": "v1", "resource": "pods" }, "requestKind": { "group": "", "version": "v1", "kind": "Pod" }, "requestResource": { "group": "", "version": "v1", "resource": "pods" }, "namespace": "yolo", "operation": "CREATE", "userInfo": { "username": "kubernetes-admin", "groups": [ "system:masters", "system:authenticated" ] }, "object": { "kind": "Pod", "apiVersion": "v1", "metadata": { "name": "c7m", "namespace": "yolo", "creationTimestamp": null, "labels": { "name": "c7m" }, "annotations": { "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"c7m\"},\"name\":\"c7m\",\"namespace\":\"yolo\"},\"spec\":{\"containers\":[{\"args\":[\"-c\",\"trap \\\"killall sleep\\\" TERM; trap \\\"kill -9 sleep\\\" KILL; sleep infinity\"],\"command\":[\"/bin/bash\"],\"image\":\"centos:7\",\"name\":\"c7m\"}]}}\n" } }, "spec": { "volumes": [ { "name": "default-token-5z7xl", "secret": { "secretName": "default-token-5z7xl" } } ], "containers": [ { "name": "c7m", "image": "centos:7", "command": [ "/bin/bash" ], "args": [ "-c", "trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity" ], "resources": {}, "volumeMounts": [ { "name": "default-token-5z7xl", "readOnly": true, "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount" } ], "terminationMessagePath": "/dev/termination-log", "terminationMessagePolicy": "File", "imagePullPolicy": "IfNotPresent" } ], "restartPolicy": "Always", "terminationGracePeriodSeconds": 30, "dnsPolicy": "ClusterFirst", "serviceAccountName": "default", "serviceAccount": "default", "securityContext": {}, "schedulerName": "default-scheduler", "tolerations": [ { "key": "node.kubernetes.io/not-ready", "operator": "Exists", "effect": "NoExecute", "tolerationSeconds": 300 }, { "key": "node.kubernetes.io/unreachable", "operator": "Exists", "effect": "NoExecute", "tolerationSeconds": 300 } ], "priority": 0, "enableServiceLinks": true }, "status": {} }, "oldObject": null, "dryRun": false, "options": { "kind": "CreateOptions", "apiVersion": "meta.k8s.io/v1" } } }`

var emptyReqJSON string = `{ "kind": "AdmissionReview", "apiVersion": "admission.k8s.io/v1beta1", "request": null }`

var validResp string = `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"7f0b2891-916f-4ed6-b7cd-27bff1815a8c","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"namespace":"yolo","operation":"CREATE","userInfo":{"username":"kubernetes-admin","groups":["system:masters","system:authenticated"]},"object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"c7m","namespace":"yolo","creationTimestamp":null,"labels":{"name":"c7m"},"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"c7m\"},\"name\":\"c7m\",\"namespace\":\"yolo\"},\"spec\":{\"containers\":[{\"args\":[\"-c\",\"trap \\\"killall sleep\\\" TERM; trap \\\"kill -9 sleep\\\" KILL; sleep infinity\"],\"command\":[\"/bin/bash\"],\"image\":\"centos:7\",\"name\":\"c7m\"}]}}\n"}},"spec":{"volumes":[{"name":"default-token-5z7xl","secret":{"secretName":"default-token-5z7xl"}}],"containers":[{"name":"c7m","image":"centos:7","command":["/bin/bash"],"args":["-c","trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity"],"resources":{},"volumeMounts":[{"name":"default-token-5z7xl","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true},"status":{}},"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}},"response":{"uid":"7f0b2891-916f-4ed6-b7cd-27bff1815a8c","allowed":true,"status":{"metadata":{},"status":"Success"},"patch":"W3sib3AiOiJhZGQiLCJwYXRoIjoiL21ldGFkYXRhL2xhYmVscy9naXRodWJfdGVhbXMiLCJ2YWx1ZSI6InlvbG8ifV0=","patchType":"JSONPatch","auditAnnotations":{"metadata.annotations.github_teams":"mutation added for identification"}}}`

var invalidNullReqBody string = `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","response":{"uid":"","allowed":false,"status":{"metadata":{},"status":"Failure","message":"AdmissionReview request body is nil"}}}`

var invalidPodObject string = `{ "kind": "AdmissionReview", "apiVersion": "admission.k8s.io/v1beta1", "request": { "uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c", "kind": { "group": "", "version": "v1", "kind": "Pod" }, "resource": { "group": "", "version": "v1", "resource": "pods" }, "requestKind": { "group": "", "version": "v1", "kind": "Pod" }, "requestResource": { "group": "", "version": "v1", "resource": "pods" }, "namespace": "yolo", "operation": "CREATE", "userInfo": { "username": "kubernetes-admin", "groups": [ "system:masters", "system:authenticated" ] }, "object": null, "oldObject": null, "dryRun": false, "options": { "kind": "CreateOptions", "apiVersion": "meta.k8s.io/v1" } } }`

var invalidPodObjResp string = `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"7f0b2891-916f-4ed6-b7cd-27bff1815a8c","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"namespace":"yolo","operation":"CREATE","userInfo":{"username":"kubernetes-admin","groups":["system:masters","system:authenticated"]},"object":null,"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}},"response":{"uid":"7f0b2891-916f-4ed6-b7cd-27bff1815a8c","allowed":false,"status":{"metadata":{},"status":"Failure","message":"unable unmarshal pod json object unexpected end of JSON input"}}}`

func mockGetGithubTeamName(ns string) string {
	return ns
}

func TestMutate(t *testing.T) {
	type args struct {
		body              []byte
		getGithubTeamName func(string) string
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"mutate incoming admission review and append the correct github team name", args{[]byte(validInput), mockGetGithubTeamName}, []byte(validResp), false},
		{"cannot unmarshal incoming admission review request w/ no body", args{[]byte(emptyReqJSON), mockGetGithubTeamName}, []byte(invalidNullReqBody), false},
		{"cannot unmarshal incoming admission review pod object", args{[]byte(invalidPodObject), mockGetGithubTeamName}, []byte(invalidPodObjResp), false},
		{"cannot unmarshal incoming admission review", args{[]byte(`foo`), mockGetGithubTeamName}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mutate(tt.args.body, tt.args.getGithubTeamName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mutate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mutate() = %v, want %v", got, tt.want)
			}
		})
	}
}
