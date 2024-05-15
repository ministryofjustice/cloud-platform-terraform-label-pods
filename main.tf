locals {
  ns = "cloud-platform-label-pods"
}

resource "kubernetes_namespace" "label-pods" {
  metadata {
    name = local.ns

    labels = {
      "name"                               = local.ns
      "pod-security.kubernetes.io/enforce" = "restricted"
    }

    annotations = {
      "cloud-platform.justice.gov.uk/application"   = "Cloud Platform label pods controller"
      "cloud-platform.justice.gov.uk/business-unit" = "cloud-platform"
      "cloud-platform.justice.gov.uk/owner"         = "Cloud Platform: platforms@digital.justice.gov.uk"
      "cloud-platform.justice.gov.uk/source-code"   = "https://github.com/ministryofjustice/cloud-platform-label-pods"
    }
  }
}

resource "helm_release" "label-pods" {
  name       = "label-pods-controller"
  namespace  = local.ns
  chart      = "cloud-platform-label-pods"
  repository = "https://ministryofjustice.github.io/cloud-platform-helm-charts"
  version    = var.chart_version

  values = [templatefile("${path.module}/templates/values.yaml.tpl", {
    ecrUrl   = var.ecr_url
    imageTag = var.image_tag
  })]
}

