resource "digitalocean_firewall" "prod" {
  name = "mta-sts-audit"

  droplet_ids = [data.digitalocean_droplet.prod.id]

  # SSH from my home IP
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["173.79.54.83"]
  }

  # HTTPS in from world
  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # SMTP in from world
  inbound_rule {
    protocol         = "tcp"
    port_range       = "25"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  # Unrestricted outbound
  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

