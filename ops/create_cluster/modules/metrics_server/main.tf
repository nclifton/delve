resource "helm_release" "metric-server" {
  name       = "metric-server"
  repository = "https://charts.bitnami.com/bitnami" 
  chart      = "metrics-server"
  version    = "5.0.1"

  namespace = "kube-system"

  set {
    name  = "apiService.create"
    value = "true"
  }
}
