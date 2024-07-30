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
  name            = "stable-diffusion-service"
  project_id      = data.theta_projects.projects.projects[0].id
  description     = "Image generation service "
  container_images = ["thetalabsorg/sketch_to_3d:v0.0.1"]
  container_port  = 7860
  container_args  = []
  env_vars = {
    HUGGING_FACE_HUB_TOKEN = "hf_xxx"
  }
  tags     = ["ImageGen", "CodeGen"]
  icon_url = ""
}

output "deployment_template_id" {
  value = theta_deployment_template.my_first_tf_managed_template.id
}

# Fetch the list of deployment templates for the first project
data "theta_deployment_templates" "templates" {
  project_id = data.theta_projects.projects.projects[0].id
}
