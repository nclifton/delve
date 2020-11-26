resource "kubernetes_namespace" "keda" {
  metadata {
    name = "keda"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }
}

resource "helm_release" "keda" {
  name      = "keda"
  chart     = "kedacore/keda"
  namespace = "keda"
  version   = "1.5.0"

  depends_on = [kubernetes_namespace.keda]
}
