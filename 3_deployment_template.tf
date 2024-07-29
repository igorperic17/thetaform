
# Fetch the list of organizations
data "theta_organizations" "org_list3" {}

# Fetch the list of projects for the first organization
data "theta_projects" "projects3" {
  organization_id = data.theta_organizations.org_list3.organizations[0].id
}

# Fetch the list of deployment templates for the first project
data "theta_deployment_templates" "templates3" {
  project_id = data.theta_projects.projects3.projects[0].id
}

# Output the fetched organizations
output "organizations3" {
  value = data.theta_organizations.org_list3.organizations
}

# Output the fetched projects
output "projects3" {
  value = data.theta_projects.projects3.projects
}

# Output the fetched deployment templates
output "deployment_templates" {
  value = data.theta_deployment_templates.templates3.deployment_templates
}