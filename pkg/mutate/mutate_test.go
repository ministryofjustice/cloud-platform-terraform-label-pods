// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it
package mutate

import (
	"reflect"
	"testing"
)

// need a mock for get github team name fn
// need a test for happy path
// need a test for adm review fail -- cannot unmarshal admin review
// need a test for adm review fail -- adm review request body is empty
// need a test for adm review fail -- unable to unmarshal pod json
// need a test for create adm review fail fn -- test marshalling err

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
		// TODO: Add test cases.
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
