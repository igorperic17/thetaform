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

# Create a deployment template
resource "theta_deployment_template" "my_first_tf_managed_template" {
  name              = "stable-diffusion-service"
  project_id        = data.theta_projects.projects.projects[0].id
  description       = "Image generation service"
  container_images  = ["goolashe/comfy-sd-webui:latest"]
  container_port    = 7860
  container_args    = []
  env_vars = {
    HUGGING_FACE_HUB_TOKEN = "hf_xxx"
  }
  tags         = ["ImageGen", "CodeGen"]
  icon_url     = ""
}

# Create a deployment using the deployment template created above
resource "theta_deployment" "my_deployment" {
  name                = "thetaformd1233"
  project_id          = data.theta_projects.projects.projects[0].id
  deployment_image_id = theta_deployment_template.my_first_tf_managed_template.id
  container_image    = theta_deployment_template.my_first_tf_managed_template.container_images[0]
  min_replicas       = 1
  max_replicas       = 1
  vm_id              = "vm_c1"  # Replace with an actual VM ID if needed
  annotations        = {
    tags     = "[\"ImageGen\",\"CodeGen\"]"
    nickname = "dep_name"
  }
  auth_username      = "my_user"
  auth_password      = "my_hackathon_winning_password"
}

output "deployment_url" {
  value = theta_deployment.my_deployment.deployment_url
}
