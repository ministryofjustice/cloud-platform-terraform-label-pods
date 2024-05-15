# Cloud Platform Label Pods

The service is a mutating admission controller webhook. Given the [hard char limit of 63](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set) for pod labels, we are forced to use annotations to pass users' github team data up to fluent-bit. The location of the "label" can be found on the pod at `.metadata.annotations.github_teams`. A strict hard limit does not exist on annotations, instead annotations values are limited on size, in this case [256kb](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/api/validation/objectmeta.go#L44-L67) is more than enough.

## Purpose

> Label all user pods with their github team name.

## Why? 

> We need to identify who owns which logs, so that we can prevent users seeing logs they shouldn't. You can find the lua responsible for pulling the data from the label out and attaching it to the log field [here](https://github.com/ministryofjustice/cloud-platform-terraform-logging/blob/main/templates/fluent-bit.yaml.tpl)

## How?

Mutating admission webhooks allow you to “modify” a (e.g.) Pod (or any kubernetes resource) request. The webhook can live anywhere, k8s just needs to know where that is.

Webhooks can _only_ be called over a HTTPS, so we need sign signed certs and a https server.

What a webhook has to do is relatively simple, it receives a AdmissionReview from the k8s api server with a "request" field, and responds with an AdmissionReview containing an additional "response" field.

Although we are creating a mutating webhook, we actually aren't mutating anything. We are telling the k8s api _how_ to mutate the object (via the admission review object). We do this via a JSONPatch in the "response".

# Development

The easist way to make changes to the api is to deploy the app via terraform and then adjust the ecr url variable and image tag to point to your test container.

Once deployed you can watch the logs and inspect the admission requests and responses as you spin up pods.

# Deploying to prod

Merge changes into main and make a release, this will build and push the latest code to our ecr repo. Then update the terraform to point to the new release tag.

### Steps to production

1. on deploy increment helm chart app version (stretch)
1. pull list of system namespaces from somewhere (stretch)
