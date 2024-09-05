variable "do_token" {
  type = string
}

variable "mail_server_labels" {
  type = list(string)
  default = ["a", "b"]
}
