resource "google_compute_network" "git-proxy-network" {
  name                    = "git-proxy-network"
  auto_create_subnetworks = true
}

resource "google_compute_firewall" "allow-http" {
  name    = "allow-http"
  network = google_compute_network.git-proxy-network.id

  allow {
    protocol = "tcp"
    ports    = ["80", "22", "8080"]
  }
  source_ranges = ["0.0.0.0/0"]
}
