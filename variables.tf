variable "HCLOUD_TOKEN" {
  type      = string
  sensitive = true
}

variable "location" {
  default = "nbg1"
}

variable "http_protocol" {
  default = "http"
}

variable "http_port" {
  default = "80"
}

variable "instances" {
  default = "1"
}

#CPX31 vCPU 4 RAM 8 NVME GB160 TRAF GB20
variable "server_type" {
  default = "cpx31" #servers names should be lowecase without whitespaces
}

variable "os_type" {
  default = "ubuntu-20.04"
}

variable "disk_size" {
  default = "160"
}

variable "ip_range" {
  default = "10.0.1.0/24"
}
variable "ssh_key_pub" {}