# Cloud Platform Label Pods

The service is a mutating admission controller webhook.

## Purpose

> Label all user pods with their github team name.

## Why? 

> We need to identify who owns which logs, so that we can prevent users seeing logs they shouldn't

### Steps to production

1. add unit tests and remove logging
1. reject admissions properly (this is done waiting unit tests)
1. lock down rbac (create a sa and pass to the deployment, and create a specific cluster role)
1. exclude current namespace from webhook
1. fill out readme (with diagram?)

1. on deploy increment helm chart app version (stretch)
1. pull list of system namespaces from somewhere (stretch)
