go install .
# TF_LOG=INFO terraform plan --var-file=local.tfvars
# TF_LOG=DEBUG terraform destroy --var-file=local.tfvars
TF_LOG=1 terraform apply --var-file=local.tfvars