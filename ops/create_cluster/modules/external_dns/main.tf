resource "kubernetes_namespace" "external_dns" {
  metadata {
    name = "external-dns"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }
}

resource "helm_release" "external_dns" {
  name       = "external-dns"
  chart      = "external-dns"
  namespace  = "external-dns"
  repository = "https://charts.bitnami.com/bitnami"
  version    = "3.4.1"

  set {
    name  = "provider"
    value = "aws"
  }

  set {
    name  = "aws.region"
    value = "ap-southeast-2"
  }

  depends_on = [kubernetes_namespace.external_dns]
}
