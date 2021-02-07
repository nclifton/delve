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
      "kubernetes.io/ingress.class"                    = "alb"
      "alb.ingress.kubernetes.io/scheme"               = "internet-facing"
      "alb.ingress.kubernetes.io/actions.ssl-redirect" = "{\"Type\": \"redirect\", \"RedirectConfig\": { \"Protocol\": \"HTTPS\", \"Port\": \"443\", \"StatusCode\": \"HTTP_301\"}}"
      "alb.ingress.kubernetes.io/listen-ports"         = "[{\"HTTP\": 80}, {\"HTTPS\":443}]"
      "alb.ingress.kubernetes.io/success-codes"        = "200,404"
      "alb.ingress.kubernetes.io/target-type"          = "ip"
    }
    labels = {
      "app" = "rabbitmq"
    }
  }

  spec {
    rule {
      http {
        path {
          path = "/*"

          backend {
            service_name = "ssl-redirect"
            service_port = "use-annotation"
          }
        }
      }
    }

    rule {
      host = "rabbitmq-management.${var.env_dns}"

      http {
        path {
          path = "/*"

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
