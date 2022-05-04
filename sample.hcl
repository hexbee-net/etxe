# TODO: define precedence for description comments
# TODO: we want a state file to be able to delete resources without access to the declarations
# TODO: on plan command, let navigate the plan and use vi command for apply, exit, ...
# TODO: real FP map/flatmap for expressions

#  Inline comment
// Inline comment

/*
Multiline comment
*/

########################################

// var description
const CONST_NAME = "const value"

########################################

locals {
  string-local-block = expression
}

local string-local-inline = expression

########################################

variable "string-var-long" {
  type = string
  description = "var description"
  default = "default value"
}

// var description
var string-var-short: string = "default value"

var {
  // var description
  string-var-block: string = "default value"  // alternate var description, TODO: choose precedence
}

########################################

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region  = "eu-west-3"
}

data {}

resource "aws_instance" "app_server" {
  ami           = "ami-062fdd189639d3e93"
  instance_type = "t2.micro"

  tags = {
    Name = "ExampleAppServerInstance"
  }
}

check "error_name" {
  condition = expression
  error = expression
  #TODO Check format
}

if {
  // ...
  // access to properties or object will throw an error
}

resource "foo" {
  if: expression // access to properties will yield null
}
