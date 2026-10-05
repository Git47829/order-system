resource "helm_release" "worker" {
  name             = "woker"
  chart            = "${path.module}/charts/worker"
  namespace        = "apps"
  create_namespace = true

  values = [
    yamlencode({
      replicaCount = 2
      image = {
        repository = var.worker-artefact-url
        tag        = "1.0.0"
      }
      nodeSelector = {
        workload = "general"
      }
    })
  ]

  depends_on = [google_container_node_pool.gke]
}

resource "helm_release" "order-api" {
  name      = "analytics"
  chart     = "${path.module}/charts/order-api"
  namespace = "apps"

  values = [
    yamlencode({
      replicaCount = 1
      image = {
        repository = var.order-api-artefact-url
        tag        = "1.0.0"
      }
    })
  ]
}


