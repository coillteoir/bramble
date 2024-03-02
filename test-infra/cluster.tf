resource "google_container_cluster" "bramble_test_cluster" {
  name                     = "bramble-test-cluster"
  location                 = "us-central1-a"
  remove_default_node_pool = true
  initial_node_count       = 1
  network                  = google_compute_network.bramble-cluster-network.id
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "bramble-node-pool"
  location   = "us-central1-a"
  cluster    = google_container_cluster.bramble_test_cluster.id
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-micro"
  }
}
