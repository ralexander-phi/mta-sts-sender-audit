output "ipv4_address" {
  value = data.digitalocean_droplet.prod.ipv4_address
}

output "ipv6_address" {
  value = data.digitalocean_droplet.prod.ipv6_address
}

output "ipv4_address_reserved" {
  value = data.digitalocean_reserved_ip.prod.ip_address
}
