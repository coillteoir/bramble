resource "google_compute_network" "tf-network" {
  name = "tf-network"
  auto_create_subnetworks = true
}

resource "google_compute_firewall" "allow-http" {
  name    = "allow-http"
  network = "tf-network"

  allow {
    protocol = "tcp"
    ports    = ["80", "22", "8080"]
  }
  source_ranges = ["0.0.0.0/0"]
}
