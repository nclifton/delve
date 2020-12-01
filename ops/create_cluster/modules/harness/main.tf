resource "kubernetes_namespace" "harness" {
  metadata {
    name = "harness"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }
}

resource "helm_release" "harness" {
  name       = "harness"
  repository = "https://app.harness.io/storage/harness-download/harness-helm-charts/"
  chart      = "harness-delegate"
  namespace  = "harness"
  version    = "1.0.2"

  values = [
    "${templatefile("${path.module}/values.tmpl.yaml", { environment_name : terraform.workspace })}"
  ]

  depends_on = [kubernetes_namespace.harness]
}
