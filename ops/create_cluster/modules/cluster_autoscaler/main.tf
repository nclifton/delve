
resource "helm_release" "cluster_autoscaler" {
  name       = "cluster-autoscaler"
  namespace  = "kube-system"
  repository = "https://kubernetes.github.io/autoscaler"
  chart      = "cluster-autoscaler-chart"
  version    = "1.0.4"

  set {
    name  = "awsRegion"
    value = "ap-southeast-2"
  }

  set {
    name  = "autoDiscovery.clusterName"
    value = terraform.workspace
  }
}
