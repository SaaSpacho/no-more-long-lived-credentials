resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.this.function_name}"
  retention_in_days = 7
}

data "aws_iam_policy_document" "assume" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "execution" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      aws_cloudwatch_log_group.lambda_logs.arn,
      "${aws_cloudwatch_log_group.lambda_logs.arn}:*"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "sts:GetWebIdentityToken",
    ]
    resources = [
      "*"
    ]
  }
}

resource "aws_iam_policy" "execution" {
  name   = "no-more-long-lived-credentials-lambda-execution-policy"
  policy = data.aws_iam_policy_document.execution.json
}


resource "aws_iam_role_policy_attachment" "this" {
  role       = aws_iam_role.this.name
  policy_arn = aws_iam_policy.execution.arn
}

resource "aws_iam_role" "this" {
  assume_role_policy = data.aws_iam_policy_document.assume.json
  name               = "no-more-long-lived-credentials-lambda-role"
}

resource "archive_file" "code" {
  type        = "zip"
  source_file = "${path.module}/../src/bootstrap"
  output_path = "${path.module}/code/bootstrap.zip"
}

resource "aws_lambda_function" "this" {
  function_name = "no-more-long-lived-credentials"
  role          = aws_iam_role.this.arn

  filename         = archive_file.code.output_path
  source_code_hash = archive_file.code.output_base64sha256

  handler = "bootstrap"
  runtime = "provided.al2023"

  architectures = ["arm64"]
}
