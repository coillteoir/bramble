resource "google_container_cluster" "bramble_test" {
  name                     = "bramble"
  location                 = "us-central1-a"
  remove_default_node_pool = true
  initial_node_count       = 1
  network                  = google_compute_network.bramble.id
}

resource "google_container_node_pool" "primary-nodes" {
  name       = "bramble"
  location   = "us-central1-a"
  cluster    = google_container_cluster.bramble_test.id
  node_count = 2

  node_config {
    preemptible  = true
    machine_type = "e2-micro"
  }
}
