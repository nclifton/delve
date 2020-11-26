resource "kubernetes_namespace" "kubernetes_dashboard" {
  metadata {
    name = "kubernetes-dashboard"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }
}

resource "kubernetes_service_account" "eks_admin" {
  metadata {
    name = "eks-admin"
    namespace = "kube-system"

    labels = {
      "managed-by"             = "terraform"
      "cluster-name"           = "${terraform.workspace}"
      "app.kubernetes.io/name" = "eks-admin"
    }
  }
}

resource "kubernetes_cluster_role_binding" "eks_admin" {
  metadata {
    name = "eks-admin"

    labels = {
      "managed-by"             = "terraform"
      "cluster-name"           = "${terraform.workspace}"
      "app.kubernetes.io/name" = "eks-admin"
    }
  }

  subject {
    kind      = "ServiceAccount"
    name      = "eks-admin"
    namespace = "kube-system"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "cluster-admin"
  }
}

resource "helm_release" "kubernetes_dashboard" {
  name       = "kubernetes-dashboard"
  repository = "https://kubernetes.github.io/dashboard/"
  chart      = "kubernetes-dashboard"
  version    = "3.0.0"
  namespace  = "kubernetes-dashboard"

  set {
    name  = "metricsScraper.enabled"
    value = "true"
  }

  # Required for kubectl proxy to work
  set {
    name  = "fullnameOverride"
    value = "kubernetes-dashboard"
  }

  depends_on = [kubernetes_namespace.kubernetes_dashboard]
}
