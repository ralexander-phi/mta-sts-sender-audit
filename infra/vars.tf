variable "do_token" {
  type = string
}

variable "home_ip" {
  type = string
}

variable "mail_server_labels" {
  type = list(string)
  default = ["a", "b"]
}
