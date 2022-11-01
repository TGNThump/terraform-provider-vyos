terraform {
  required_providers {
    vyos = {
      source = "TGNThump/vyos"
    }
  }
}

provider "vyos" {
  endpoint = "https://vyos"
  api_key = "abcdefg"
}
