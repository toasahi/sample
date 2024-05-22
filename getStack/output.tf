output "get_stack_arn" {
  value = aws_lambda_function.get_stack.invoke_arn
}
