resource "google_service_account" "gke_nodes" {
  account_id   = "gke_nodes"
  display_name = "GKE NODE Service Account"
}

resource "google_project_iam_member" "gke_nodes" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
    "roles/artefactregistry.reader",
  ])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}


resource "google_compute_network" "default" {
  name = "worker-network-prod"

  auto_create_subnetworks  = false
  enable_ula_internal_ipv6 = true

}

resource "google_compute_subnetwork" "default" {
  name = "worker-subnet-prod"

  ip_cidr_range = "10.0.0.0/16"
  region        = "eu-central-1"

  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "INTERNAL"

  network = google_compute_network.default.id
  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.0.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.1.0/24"
  }

}

resource "google_container_cluster" "gke" {
  name                     = "woker-system-cluster"
  location                 = "eu-central-1"
  remove_default_node_pool = true
  initial_node_count       = 1
  enable_autopilot         = false
  enable_l4_ilb_subsetting = true

  network    = google_compute_network.default.id
  subnetwork = google_compute_subnetwork.default.id

  ip_allocation_policy {
    stack_type                    = "IPV4_IPV6"
    services_secondary_range_name = google_compute_subnetwork.default.secondary_ip_range[0].range_name
    cluster_secondary_range_name  = google_compute_subnetwork.default.secondary_ip_range[1].range_name

  }

  deletion_protection = false

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  release_channel {
    channel = "REGULAR"
  }
}

resource "google_container_node_pool" "general" {
  name     = "general"
  cluster  = google_container_cluster.gke.name
  location = var.cluster_location

  autoscaling {
    min_node_count = 1
    max_node_count = 1
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    machine_type    = "e2-standard-4"
    disk_size_gb    = 5
    service_account = google_service_account.gke_nodes.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]

    labels = {
      workload = "general"
    }

    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
}

