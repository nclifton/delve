resource "kubernetes_namespace" "rabbitmq" {
  metadata {
    name = "rabbitmq"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }
}

resource "kubernetes_secret" "rabbitmq" {
  metadata {
    name      = "rabbitmq-load-definition"
    namespace = "rabbitmq"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = "${terraform.workspace}"
    }
  }

  type = "Opaque"

  data = {
    "load_definition.json" = "${file("${path.module}/load-definition.json")}"
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

