resource "digitalocean_project" "mta-sts-audit" {
  name        = "MTA-STS Audit"
  description = "Audit MTA-STS support for mail senders"
  purpose     = "Web Application"
  environment = "Production"
}
