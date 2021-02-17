resource "kubernetes_namespace" "rabbitmq" {
  metadata {
    name = "rabbitmq"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "kubernetes_secret" "rabbitmq" {
  metadata {
    name      = "rabbitmq-load-definition"
    namespace = "rabbitmq"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }

  type = "Opaque"

  data = {
    "load_definition.json" = file("${path.module}/load-definition.json")
  }

  depends_on = [kubernetes_namespace.rabbitmq]
}

resource "helm_release" "rabbitmq" {
  name      = "rabbitmq"
  chart     = "bitnami/rabbitmq"
  namespace = "rabbitmq"
  version   = "7.6.8"

  values = [
    file("${path.module}/values.yaml")
  ]

  depends_on = [kubernetes_namespace.rabbitmq]
}

resource "kubernetes_ingress" "rabbitmq_ingress" {
  metadata {
    name      = "rabbitmq-ingress"
    namespace = "rabbitmq"
    annotations = {
      "kubernetes.io/ingress.class"                      = "traefik"
      "cert-manager.io/issuer"                           = "letsencrypt-prod"
    }
    labels = {
      "app" = "rabbitmq"
    }
  }

  spec {
    tls {
      hosts       = ["rabbitmq.${var.env_dns}"]
      secret_name = "rabbitmq-tls"
    }

    rule {
      host = "rabbitmq.${var.env_dns}"

      http {
        path {
          path = "/"

          backend {
            service_name = "rabbitmq"
            service_port = 15672
          }
        }
      }
    }
  }

  depends_on = [kubernetes_namespace.rabbitmq]
}
