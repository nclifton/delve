resource "helm_release" "harness" {
  # Only deploy to staging and production environments
  # https://stackoverflow.com/a/58193941
  count = (length(regexall(".*staging.*", terraform.workspace)) > 0 || length(regexall(".*production.*", terraform.workspace)) > 0) ? 1 : 0

  name       = "harness"
  repository = "https://app.harness.io/storage/harness-download/harness-helm-charts/"
  chart      = "harness-delegate"
  version    = "1.0.2"

  values = [
    templatefile("${path.module}/values.tmpl.yaml", { environment_name : terraform.workspace })
  ]
}
