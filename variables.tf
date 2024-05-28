variable "email" {
  description = "Email for the Theta API"
  type        = string
}

variable "password" {
  description = "Password for the Theta API"
  type        = string
  sensitive   = true
}
