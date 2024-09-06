# API server
resource "digitalocean_record" "api" {
  domain = data.digitalocean_domain.alexsci-com.id
  type   = "A"
  name   = "api.audit"
  value  = data.digitalocean_reserved_ip.prod.ip_address
}


# Mail servers
resource "digitalocean_record" "A-mail-a" {
  domain = data.digitalocean_domain.alexsci-com.id
  type   = "A"
  name   = "mail-a.audit"
  value  = data.digitalocean_reserved_ip.prod.ip_address
}


resource "digitalocean_record" "A-mail-b" {
  domain = data.digitalocean_domain.alexsci-com.id
  type   = "A"
  name   = "mail-b.audit"
  value  = data.digitalocean_droplet.prod.ipv4_address
}


resource "digitalocean_record" "MX-mail" {
  domain = data.digitalocean_domain.alexsci-com.id
  type   = "MX"
  priority = 50

  for_each = toset(var.mail_server_labels)
  name   = "${each.key}.audit"
  value  = "mail-${each.key}.audit.alexsci.com."
}


# MTA-STS policy static hosting
resource "digitalocean_record" "A-a" {
  domain = data.digitalocean_domain.alexsci-com.id
  type   = "A"

  for_each = toset(var.mail_servers_with_sts_labels)
  name   = "mta-sts.${each.key}.audit"

  # Co-hosted with API server
  value  = digitalocean_record.api.value
}

# MTA-STS record with policy ID
resource "digitalocean_record" "TXT-mta-sts" {
  domain = data.digitalocean_domain.alexsci-com.id
  type = "TXT"

  for_each = toset(var.mail_servers_with_sts_labels)
  name = "_mta-sts.${each.key}.audit"

  value = var.mta_sts_policy_record
}
