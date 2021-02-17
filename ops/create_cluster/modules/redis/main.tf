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
      "kubernetes.io/ingress.class"                      = "traefik"
      "cert-manager.io/issuer"                           = "letsencrypt-prod"
    }
    labels = {
      "app" = "redis-commander"
    }
  }

  spec {
    tls {
      hosts       = ["redis.${var.env_dns}"]
      secret_name = "redis-commander-tls"
    }

    rule {
      host = "redis.${var.env_dns}"

      http {
        path {
          path = "/"

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
