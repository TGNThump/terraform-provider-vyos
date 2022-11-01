resource "vyos_config" "test" {
  path = "firewall name WAN_LOCAL"
  value = jsonencode({
    default-action = "drop"
    rule = {
      10 = {
        action = "accept"
        state = {
          established = "enable"
          related     = "enable"
        }
      }
    }
  })
}