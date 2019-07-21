terraform {
  required_version = ">= 0.11.0"
}

# Generate random ID (used for tagging devices)
resource "random_id" "rval" {
  byte_length = 8
}

# Specify the provider and access details
provider "aws" {
  region     = var.aws_region
  access_key = var.access_key
  secret_key = var.secret_key
}

