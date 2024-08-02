<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_theta"></a> [theta](#requirement\_theta) | 0.1.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_theta"></a> [theta](#provider\_theta) | 0.1.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| theta_deployment.notebook_deployment | resource |
| theta_deployment_template.my_first_tf_managed_template | resource |
| theta_organizations.org_list | data source |
| theta_projects.projects | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_email"></a> [email](#input\_email) | Email for the Theta API | `string` | n/a | yes |
| <a name="input_hf_token"></a> [hf\_token](#input\_hf\_token) | HugginFace API token | `string` | n/a | yes |
| <a name="input_password"></a> [password](#input\_password) | Password for the Theta API | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_deployment_url"></a> [deployment\_url](#output\_deployment\_url) | n/a |
<!-- END_TF_DOCS -->