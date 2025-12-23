resource "scaleway_function_namespace" "namespace" {
  name        = "no-more-long-lived-credentials"
  description = "Main function namespace"
}

resource "archive_file" "code" {
  type        = "zip"
  source_dir  = "${path.module}/../src/"
  output_path = "${path.module}/code/bootstrap.zip"
}

resource "scaleway_function" "this" {
  namespace_id = scaleway_function_namespace.namespace.id
  name         = "no-more-long-lived-credentials"
  runtime      = "go124"
  timeout      = 10

  min_scale = 0

  handler = "Handle"
  privacy = "public"

  zip_file = archive_file.code.output_path
  zip_hash = archive_file.code.output_base64sha256
  deploy   = true
}
