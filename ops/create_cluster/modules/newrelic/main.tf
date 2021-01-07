# Install New Relic via Helm

resource "kubernetes_namespace" "newrelic" {
  metadata {
    name = "newrelic"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "helm_release" "newrelic" {
  name       = "newrelic-bundle"
  namespace  = "newrelic"
  repository = "https://helm-charts.newrelic.com"
  chart      = "nri-bundle"
  version    = "1.1.0"

  set {
    name  = "global.licenseKey"
    value = "3289b893db560c39096d3c222d4a33037c24f943"
  }

  set {
    name  = "global.cluster"
    value = terraform.workspace
  }

  set {
    name  = "global.namespace"
    value = "newrelic"
  }

  set {
    name  = "newrelic-infrastructure.privileged"
    value = true
  }

  set {
    name  = "ksm.enabled"
    value = true
  }

  set {
    name  = "prometheus.enabled"
    value = true
  }

  set {
    name  = "kubeEvents.enabled"
    value = true
  }

  set {
    name  = "logging.enabled"
    value = true
  }

  depends_on = [kubernetes_namespace.newrelic]
}
