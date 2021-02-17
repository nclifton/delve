resource "kubernetes_namespace" "traefik" {
  metadata {
    name = "traefik"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "helm_release" "traefik" {
  chart      = "traefik"
  name       = "traefik"
  namespace  = "traefik"
  repository = "https://helm.traefik.io/traefik"
  version    = "9.14.2"

  values = [
    templatefile("${path.module}/values.tmpl.yaml", { environment_name : terraform.workspace })
  ]

  depends_on = [kubernetes_namespace.traefik]
}
