resource "kubernetes_namespace" "redis" {
  metadata {
    name = "redis"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
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

resource "kubernetes_service" "redis_commander" {
  metadata {
    name      = "redis-commander"
    namespace = "redis"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }

  spec {
    port {
      name        = "http"
      protocol    = "TCP"
      port        = 8081
      target_port = "8081"
    }

    selector = {
      app = "redis-commander"
    }

    type = "NodePort"
  }

  depends_on = [kubernetes_namespace.redis]
}

resource "kubernetes_deployment" "redis_commander" {
  # this takes longer than 10 minutes to succeed
  wait_for_rollout = false
  metadata {
    name      = "redis-commander"
    namespace = "redis"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "redis-commander"
      }
    }

    template {
      metadata {
        labels = {
          app  = "redis-commander"
          tier = "backend"

          "managed-by"   = "terraform"
          "cluster-name" = "${terraform.workspace}"
        }
      }

      spec {
        container {
          name  = "redis-commander"
          image = "rediscommander/redis-commander"

          port {
            name           = "redis-commander"
            container_port = 8081
          }

          env {
            name  = "REDIS_HOSTS"
            value = "redis-master.redis"
          }

          env {
            name  = "K8S_SIGTERM"
            value = "1"
          }

          env {
            name  = "HTTP_USER"
            value = "burstsms"
          }

          env {
            name  = "HTTP_PASSWORD"
            value = "324nhsedf8sdf"
          }

          resources {
            limits {
              cpu    = "500m"
              memory = "512M"
            }
          }

          security_context {
            capabilities {
              drop = ["ALL"]
            }

            run_as_non_root = true
          }
        }
      }
    }
  }

  depends_on = [kubernetes_namespace.redis]
}

