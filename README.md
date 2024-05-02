# Cloud Platform Label Pods

The service is a mutating admission controller webhook.

## Purpose

> Label all user pods with their github team name.

## Why? 

> We need to identify who owns which logs, so that we can prevent users seeing logs they shouldn't

## How?

Mutating admission webhooks allow you to “modify” a (e.g.) Pod (or any kubernetes resource) request. The webhook can live anywhere, k8s just needs to know where that is.

Webhooks can _only_ be called over a HTTPS, so we need sign signed certs and a https server.

What a webhook has to do is relatively simple, it receives a AdmissionReview from the k8s api server with a "request" field, and responds with an AdmissionReview containing an additional "response" field.

Although we are creating a mutating webhook, we actually aren't mutating anything. We are telling the k8s api _how_ to mutate the object (via the admission review object). We do this via a JSONPatch in the "response".

# Development

The easiest way to make changes to the api is to deploy to a new namespace in a test cluster (you'll also need to apply a new [cluster role](https://github.com/ministryofjustice/cloud-platform-environments/blob/9ac868069ab515ffaca66685d5cb59ee6cf1717e/namespaces/live.cloud-platform.service.justice.gov.uk/user-roles.yaml#L86), a [cluster role binding](https://github.com/ministryofjustice/cloud-platform-environments/blob/9ac868069ab515ffaca66685d5cb59ee6cf1717e/namespaces/live.cloud-platform.service.justice.gov.uk/cloud-platform-label-pods/01-rbac.yaml#L17), a [self signed cert](https://github.com/ministryofjustice/cloud-platform-environments/blob/main/namespaces/live.cloud-platform.service.justice.gov.uk/cloud-platform-label-pods/05-self-signed-cert.yaml) and a [mutating web hook](https://github.com/ministryofjustice/cloud-platform-environments/blob/main/namespaces/live.cloud-platform.service.justice.gov.uk/cloud-platform-label-pods/06-mutating-webhook.yaml))

Once deployed you can watch the logs and inspect the admission requests and responses as you spin up pods.

### Steps to production

1. on deploy increment helm chart app version (stretch)
1. pull list of system namespaces from somewhere (stretch)
