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

### Steps to production

1. on deploy increment helm chart app version (stretch)
1. pull list of system namespaces from somewhere (stretch)
