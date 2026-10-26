game {

  start_zone_id = "zone-one"

  desc = "You find yourself in a test world. Can you get out?"

  zone "zone-one" {
    desc = "a black void we call zone-one"
    exit "east" {
      flags = ["end-flag"]
      to = "zone-two"
    }
  }

  zone "zone-two" {
    desc = "a black void we call zone-two"
  }

  object "button" {
    zone_id = "zone-two"
    desc = "it's a button"
    action "use" {
      rule {
        effect "end_game" {
          is_win = false
          text = "You can never escape!"
        }
      }
    }
  }

  object "brass-key" {
    zone_id     = "zone-one"
    visible     = false
    targetable  = false
    carriable   = true
    desc        = "a small brass key"
  }

  object "desk" {
    zone_id     = "zone-one"
    visible     = true
    targetable  = true
    carriable   = false

    # Override built-in examine to have a side effect
    # Note desc is not required in this case.
    action "examine" {
      rule {
        effect "say" {
          text = "there is a brass-key in the drawer"
        }
        effect "object_targetable" {
          object_id = "brass-key"
          value = true
        }
      }
    }
  }

  object "door" {
    zone_id     = "zone-one"
    desc        = "A solid door. The way out?"
    visible     = true
    targetable  = true
    carriable   = false

    property "lock" {
      value = "locked"
    }

    action "use" {

      # The door is locked and key not used -> say
      rule {
        when {
          property "lock" {
            value = "locked"
          }
        }
        effect "say" {
          text = "the door is locked"
        }
      }

      # The door is locked. key is used. -> door unlocked, flag
      rule {
        when {
          use_with = "brass-key"
          property "lock" {
            value = "locked"
          }
        }
        effect "say" {
          text = "the key clicks in the latch"
        }
        effect "say" {
          text = "the door opens"
        }
        effect "property" {
          object_id   = "door"
          property_id = "lock"
          value       = "unlocked"
        }
        effect "add_flag" {
          flag_id   = "end-flag"
        }
      }

      # The door is unlocked. the key is not used -> say
      rule {
        when {
          property "lock" {
            value = "unlocked"
          }
        }
        effect "say" {
          text = "the door open and the way is clear"
        }
      }

      # The door is unlocked. the key is used -> Locked, flag
      rule {
        when {
          use_with = "brass-key"
          property "lock" {
            value = "unlocked"
          }
        }
        effect "say" {
          text = "the key clicks in the latch"
        }
        effect "say" {
          text = "the door is locked"
        }
        effect "property" {
          object_id   = "door"
          property_id = "lock"
          value       = "locked"
        }
        effect "remove_flag" {
          flag_id = "end-flag"
        }
      }
    }
  }
}
