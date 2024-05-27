terraform {
  required_providers {
    theta = {
      source = "hashicorp.com/edu/theta"
    }
  }
}

provider "theta" {
  api_key    = var.api_key
  api_secret = var.api_secret
}

resource "theta_endpoint" "example" {
  name                = "llama38b31eu8r5y5d"
  project_id          = "prj_8qf89pmjgdqurbaqfpdu3u854s6p"
  deployment_image_id = "img_rrdau7uikg8rhurf7cbei8j77nbp"
  container_image     = "vllm/vllm-openai"
  min_replicas        = 1
  max_replicas        = 1
  vm_id               = "vm_c1"
  annotations = {
    tags = "[\"LLM\",\"API\"]"
  }
  env_vars = {
    HUGGING_FACE_HUB_TOKEN = "12345"
  }
}
