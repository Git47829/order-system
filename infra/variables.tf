variable "worker-artefact-url" {
  type        = string
  description = "Link for the Container Artefact"
}

variable "order-api-artefact-url" {
  type        = string
  description = "Link for the Order Api Artefact"
}

variable "cluster_location" {
  type        = string
  description = "Region of the Cluster"
  default     = "eu-central-1"
}

variable "project_id" {
  type        = string
  description = "Project Id of the Cluster"
  default     = "floci-local"
}
