terraform {
  required_providers {
    theta = {
      source  = "hashicorp.com/edu/theta"
      version = "0.1.0"
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

# Use a valid project ID from the fetched projects to create a deployment template
resource "theta_deployment_template" "my_first_tf_managed_template" {
  name            = "example-template"
  project_id      = data.theta_projects.projects.projects[0].id
  description     = "An example deployment template"
  container_image = "example/image:latest"
  container_port  = 8080
  env_vars = {
    EXAMPLE_VAR = "example-value"
  }
  tags     = ["example", "template"]
  icon_url = "https://example.com/icon.png"
}

output "deployment_template_id" {
  value = theta_deployment_template.my_first_tf_managed_template.id
}

# Fetch the list of deployment templates for the first project
data "theta_deployment_templates" "templates" {
  project_id = data.theta_projects.projects.projects[0].id
}

# Output the fetched organizations
output "organizations" {
  value = data.theta_organizations.org_list.organizations
}

# Output the fetched projects
output "projects" {
  value = data.theta_projects.projects.projects
}

# Output the fetched deployment templates
output "deployment_templates" {
  value = data.theta_deployment_templates.templates
}
