---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vyos Provider"
subcategory: ""
description: |-
  
---

# vyos Provider



## Example Usage

```terraform
terraform {
  required_providers {
    vyos = {
      source = "TGNThump/vyos"
    }
  }
}

provider "vyos" {
  endpoint = "https://vyos"
  api_key  = "abcdefg"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String, Sensitive) API Key for the VyOS HTTP API
- `endpoint` (String) Endpoint of the VyOS HTTP API
- `save_file` (String) Remote file path to save the config too.
- `skip_saving` (Boolean) Set to true to skip saving the config to disk.
