resource "helm_release" "metric-server" {
  name       = "metrics-server"
  repository = "https://charts.bitnami.com/bitnami" 
  chart      = "metrics-server"
  version    = "5.4.0"

  namespace = "kube-system"

  set {
    name  = "apiService.create"
    value = "true"
  }

  set {
    name  = "rbac.create"
    value = "true"
  }
}
