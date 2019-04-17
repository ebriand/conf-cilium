variable "gce_region" {
  type    = "string"
  default = "europe-west1"
}

variable "gce_zone" {
  type = "string"
  default = "europe-west1-c"
}

variable "gce_ssh_user" {
  type = "string"
  default = "ubuntu"
}

variable "gce_ssh_pub_key_file" {
  type = "string"
  default = "kube_rsa.pub"
}

variable "worker_count" {
  type = "string"
  default = "3"
}

resource "google_compute_network" "kubeadm-network" {
  name                    = "kubeadm-network"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "kubeadm-network-sub0" {
  name          = "kubeadm-network-sub0"
  region		= "${var.gce_region}"
  network       = "${google_compute_network.kubeadm-network.name}"
  ip_cidr_range = "10.240.0.0/20"
}

resource "google_compute_firewall" "internal" {
  name    = "kubeadm-network--allow-internal"
  network = "${google_compute_network.kubeadm-network.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "udp"
  }

  allow {
    protocol = "tcp"
  }

  source_ranges = ["10.240.0.0/20", "10.200.0.0/16"]
}

resource "google_compute_firewall" "external" {
  name    = "kubeadm-network--allow-external"
  network = "${google_compute_network.kubeadm-network.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "udp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "6443", "30000-32767"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "worker" {
  count = "${var.worker_count}"
  name            = "worker-${count.index}"
  machine_type    = "n1-standard-1"
  zone            = "${var.gce_zone}"
  can_ip_forward  = true

  tags = ["kubeadm","worker"]
  
  scheduling {
    automatic_restart   = false
  }
  
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.kubeadm-network-sub0.self_link}"
    network_ip = "10.240.0.2${count.index}"

    access_config {
      // Ephemeral IP
    }
  }
  
  service_account {
    scopes = ["compute-rw","storage-ro","service-management","service-control","logging-write","monitoring"]
  }

  metadata {
    sshKeys = "${var.gce_ssh_user}:${file(var.gce_ssh_pub_key_file)}"
  }
}

resource "google_compute_instance_group" "worker_group" {
  name               = "worker-group"
  zone               = "${var.gce_zone}"
  instances          = ["${google_compute_instance.worker.*.self_link}"]
}

resource "google_compute_instance" "controller" {
  count    = 1
  name     = "controller-${count.index}"
  machine_type    = "n1-standard-2"
  zone            = "${var.gce_zone}"
  can_ip_forward  = true

  tags = ["kubeadm","controller"]
  
  scheduling {
    automatic_restart   = false
  }
  
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    subnetwork = "${google_compute_subnetwork.kubeadm-network-sub0.self_link}"
    network_ip = "10.240.0.1${count.index}"

    access_config {
      // Ephemeral IP
    }
  }
  
  service_account {
    scopes = ["compute-rw","storage-ro","service-management","service-control","logging-write","monitoring"]
  }
  metadata {
    sshKeys = "${var.gce_ssh_user}:${file(var.gce_ssh_pub_key_file)}"
  }
}

resource "google_compute_instance_group" "controller_group" {
  name               = "controller-group"
  zone               = "${var.gce_zone}"
  instances          = ["${google_compute_instance.controller.0.self_link}"]
  named_port {
  	name = "api-server"
  	port = 6443
  }
}
