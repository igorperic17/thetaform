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

# Fetch the list of organizations
data "theta_organizations" "org_list" {}

# Fetch the list of projects for the first organization
data "theta_projects" "projects" {
  organization_id = data.theta_organizations.org_list.organizations[0].id
}

# Output the fetched organizations
output "organizations" {
  value = data.theta_organizations.org_list.organizations
}

# Output the fetched projects
output "projects" {
  value = data.theta_projects.projects
}