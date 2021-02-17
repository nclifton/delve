# Install Pgadmin via Helm

resource "random_string" "password" {
  length  = 16
  special = false # can accidentally generate strings starting with { that can't be parsed correctly
}

resource "kubernetes_namespace" "pgadmin" {
  metadata {
    name = "pgadmin"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }
}

resource "aws_secretsmanager_secret" "pgadmin4-password" {
  name                    = "${terraform.workspace}-pgadmin4-password"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "pgadmin4-password-str" {
  secret_id     = aws_secretsmanager_secret.pgadmin4-password.id
  secret_string = random_string.password.result

  lifecycle {
    ignore_changes = [secret_string]
  }
}

resource "kubernetes_secret" "pgpassfile" {
  metadata {
    name      = "pgpassfile"
    namespace = "pgadmin"
    labels = {
      "managed-by"   = "terraform"
      "cluster-name" = terraform.workspace
    }
  }

  data = {
    pgpassfile = <<EOT
${var.postgres_endpoint}:5432:sendsei:foo:barbut8chars
    EOT
  }

  depends_on = [kubernetes_namespace.pgadmin]
}

resource "helm_release" "pgadmin" {
  name       = "pgadmin"
  repository = "https://helm.runix.net/"
  chart      = "pgadmin4"
  namespace  = "pgadmin"
  version    = "1.3.2"

  values = [
    templatefile("${path.module}/values.tmpl.yaml", { postgres_endpoint : var.postgres_endpoint })
  ]
  set {
    name  = "env.email"
    value = "burstsms@burstsms.com"
  }

  set {
    name  = "env.password"
    value = random_string.password.result
  }

  depends_on = [kubernetes_namespace.pgadmin]
}

resource "kubernetes_ingress" "pgadmin_ingress" {
  metadata {
    name      = "pgadmin-ingress"
    namespace = "pgadmin"
    annotations = {
      "kubernetes.io/ingress.class"                      = "traefik"
      "cert-manager.io/issuer"                           = "letsencrypt-prod"
    }
    labels = {
      "app" = "pgadmin"
    }
  }

  spec {
    tls {
      hosts       = ["pgadmin.${var.env_dns}"]
      secret_name = "pgadmin-tls"
    }

    rule {
      host = "pgadmin.${var.env_dns}"

      http {
        path {
          path = "/"

          backend {
            service_name = "pgadmin-pgadmin4"
            service_port = 80
          }
        }
      }
    }
  }

  depends_on = [kubernetes_namespace.pgadmin]
}
