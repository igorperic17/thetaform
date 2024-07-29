terraform {
  required_providers {
    theta = {
      source = "hashicorp.com/edu/theta"
    }
  }
}

provider "theta" {
  email    = var.email
  password = var.password
}

data "theta_organizations" "org_list" {}

output "organizations" {
  value = data.theta_organizations.org_list.organizations
}