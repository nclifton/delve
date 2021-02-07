resource "kubernetes_namespace" "redis" {
  metadata {
    name = "redis"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "helm_release" "redis" {
  name      = "redis"
  chart     = "bitnami/redis"
  namespace = "redis"
  version   = "10.6.16"

  values = [
    file("${path.module}/values.yaml")
  ]

  depends_on = [kubernetes_namespace.redis]
}

resource "null_resource" "redis_commander" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.module}/redis-commander.yaml"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "kubectl delete -f ${path.module}/redis-commander.yaml"
  }

  depends_on = [kubernetes_namespace.redis]
}

resource "kubernetes_ingress" "redis_commander_ingress" {
  metadata {
    name      = "redis-commander-ingress"
    namespace = "redis"
    annotations = {
      "kubernetes.io/ingress.class"                    = "alb"
      "alb.ingress.kubernetes.io/scheme"               = "internet-facing"
      "alb.ingress.kubernetes.io/actions.ssl-redirect" = "{\"Type\": \"redirect\", \"RedirectConfig\": { \"Protocol\": \"HTTPS\", \"Port\": \"443\", \"StatusCode\": \"HTTP_301\"}}"
      "alb.ingress.kubernetes.io/listen-ports"         = "[{\"HTTP\": 80}, {\"HTTPS\":443}]"
      "alb.ingress.kubernetes.io/success-codes"        = "200,404"
      "alb.ingress.kubernetes.io/target-type"          = "ip"
    }
    labels = {
      "app" = "redis-commander"
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
      host = "redis-commander.${var.env_dns}"

      http {
        path {
          path = "/*"

          backend {
            service_name = "redis-commander"
            service_port = 8081
          }
        }
      }
    }
  }

  depends_on = [kubernetes_namespace.redis]
}
