variable "email" {
  description = "Email for the Theta API"
  type        = string
}

variable "password" {
  description = "Password for the Theta API"
  type        = string
  sensitive   = true
}

variable "hf_token" {
  description = "HugginFace API token"
  type        = string
  sensitive   = true
}
