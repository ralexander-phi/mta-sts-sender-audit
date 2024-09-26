variable "do_token" {
  type = string
}

variable "home_ip" {
  type = string
}

# All mail servers
variable "mail_server_labels" {
  type = list(string)
  default = ["a", "b", "c", "d", "e", "f", "g"]
}

# Which mail servers should have MTA-STS DNS records
variable "mail_servers_with_sts_labels" {
  type = list(string)
  # Excludes:
  # f - Legacy non-TLS server
  # g - DANE TLS server
  default = ["a", "b", "c", "d", "e"]
}

variable "mta_sts_policy_record" {
  type = string
  default = "v=STSv1; id=20240906T1730;"
}

variable "tlsrpt_address" {
  type = string
  default = "tlsrpt-audit@robalexdev.com"
}

variable "ipv6_subnet_size" {
  type = number
  default = 124
}

