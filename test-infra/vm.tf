resource "google_compute_instance" "git-proxy-vm" {
  boot_disk {
    auto_delete = true
    device_name = "instance-20240223-130250"

    initialize_params {
      image = "projects/ubuntu-os-cloud/global/images/ubuntu-2310-mantic-amd64-v20240223"
      size  = 10
      type  = "pd-balanced"
    }

    mode = "READ_WRITE"
  }

  can_ip_forward      = false
  deletion_protection = false
  enable_display      = false

  labels = {
    goog-ec-src = "vm_add-tf"
  }

  machine_type = "e2-micro"
  name         = "vm-1"

  metadata_startup_script = file("./init_git_proxy.sh")

  network_interface {
    network = google_compute_network.bramble.id

    access_config {
      network_tier = "PREMIUM"
    }

    queue_count = 0
    stack_type  = "IPV4_ONLY"
  }

  scheduling {
    automatic_restart   = false
    on_host_maintenance = "TERMINATE"
    preemptible         = true
    provisioning_model  = "SPOT"
  }

  shielded_instance_config {
    enable_integrity_monitoring = true
    enable_secure_boot          = false
    enable_vtpm                 = true
  }

  tags = ["http-server", "https-server"]
  zone = "us-central1-a"
}
