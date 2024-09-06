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

variable "mta_sts_policy_record" {
  type = string
  default = "v=STSv1; id=20240906T1730;"
}

variable "tlsrpt_address" {
  type = string
  default = "tlsrpt-audit@robalexdev.com"
}

