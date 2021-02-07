resource "helm_release" "harness" {
  name       = "harness"
  repository = "https://app.harness.io/storage/harness-download/harness-helm-charts/"
  chart      = "harness-delegate"
  version    = "1.0.2"

  values = [
    templatefile("${path.module}/values.tmpl.yaml", { environment_name : terraform.workspace })
  ]
}
