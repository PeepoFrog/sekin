# needed for terraform >= 0.13
terraform {
  required_providers {
    hcloud = {
      source = "hetznercloud/hcloud"
      version = "1.25.2"
    }
  }
}