resource "kubernetes_namespace" "cert_manager" {
  metadata {
    name = "cert-manager"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "helm_release" "cert_manager" {
  chart      = "cert-manager"
  name       = "cert-manager"
  namespace  = "cert-manager"
  repository = "https://charts.jetstack.io"
  version    = "1.1.0"

  set {
    name  = "installCRDs"
    value = true
  }

  # Below are set for AWS EKS compatibility
  # https://cert-manager.io/docs/installation/compatibility/

  set {
    name  = "webhook.hostNetwork"
    value = true
  }

  set {
    name  = "webhook.securePort"
    value = 10260
  }

  depends_on = [kubernetes_namespace.cert_manager]
}

# There is an issue with applying CertIssuer's straight after cert-manager is deployed
# Workaround is to wait a minute then proceed
# https://github.com/jetstack/cert-manager/issues/3338#issuecomment-707579834
resource "time_sleep" "wait_60_seconds" {
  depends_on = [helm_release.cert_manager]

  create_duration = "60s"
}

# TODO: In the future, Kubernetes Provider will have a way to apply
# CRD's built in but for now use kubectl apply
resource "null_resource" "kubernetes_apply" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.module}/letsencrypt-prod.yaml -n cert-manager"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "kubectl delete -f ${path.module}/letsencrypt-prod.yaml -n cert-manager"
  }

  depends_on = [time_sleep.wait_60_seconds]
}
