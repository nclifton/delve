resource "kubernetes_namespace" "keda" {
  metadata {
    name = "keda"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "helm_release" "keda" {
  name       = "keda"
  chart      = "keda"
  namespace  = "keda"
  repository = "https://kedacore.github.io/charts"
  version    = "1.5.0"

  depends_on = [kubernetes_namespace.keda]
}
