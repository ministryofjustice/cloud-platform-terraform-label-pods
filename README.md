# Cloud Platform Label Pods

The service is a mutating admission controller webhook.

## Purpose

> Label all user pods with their github team name.

## Why? 

> We need to identify who owns which logs, so that we can prevent users seeing logs they shouldn't

### Steps to production

1. add unit tests and remove logging
1. update docker tests to check for ca.cert
1. lock down rbac
1. reject admissions properly
1. exclude current namespace from webhook
1. pass proper prod vars to prod deploy
1. fill out readme
1. on deploy increment helm chart app version
1. pull list of system namespaces from somewhere
