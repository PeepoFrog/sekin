resource "hcloud_ssh_key" "default" {
  name       = "hetzner_key"
  public_key = "${ secrets.SSH_PRIVATE_KEY }"
}