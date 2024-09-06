variable "do_token" {
  type = string
}

variable "home_ip" {
  type = string
}

# All mail servers
variable "mail_server_labels" {
  type = list(string)
  default = ["a", "b"]
}

# Which mail servers should have MTA-STS DNS records
variable "mail_servers_with_sts_labels" {
  type = list(string)
  default = ["a", "b"]
}
