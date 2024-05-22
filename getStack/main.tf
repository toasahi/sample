data "aws_iam_policy_document" "get_stack_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "get_stack" {
  name               = "get_stack_sfn"
  assume_role_policy = data.aws_iam_policy_document.get_stack_assume_role.json
}

resource "aws_iam_role_policy" "get_stack" {
  name = "get-stack"
  role = aws_iam_role.get_stack.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "cloudformation:DescribeStacks"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "terraform_data" "get_stack" {
  triggers_replace = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "cd ${path.module}/lambda/ && GOOS=linux GOARCH=arm64 go build -o ./build/bootstrap main.go"
  }
}

data "archive_file" "get_stack" {
  type        = "zip"
  source_file = "${path.module}/lambda/build/bootstrap"
  output_path = "${path.module}/lambda/go.zip"

  depends_on = [terraform_data.get_stack]
}

resource "aws_lambda_function" "get_stack" {
  filename         = "${path.module}/lambda/go.zip"
  function_name    = "get_stack"
  role             = aws_iam_role.get_stack.arn
  handler          = "sample"
  architectures    = ["arm64"]
  source_code_hash = data.archive_file.get_stack.output_base64sha256

  runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "get_stack" {
  name = "/aws/lambda/${aws_lambda_function.get_stack.function_name}"
}
