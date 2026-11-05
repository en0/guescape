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

  mode "mode-one" {
    prompt = "username: "
    state = {
      "s" = "login"
      "user" = "zerocool"
      "pw" = "abc123"
      "cwd" = "~"
      "last" = ""
    }
    lua = <<-LUA
      local function input() return table.concat(args) end

      local function doLogin()
          setState("last", input())
          setState("s", "login-password")
          setPrompt("password: ")
      end

      local function doPassword()
          if getState("last") == getState("user") and input() == getState("pw") then
              setState("s", "term")
              setPrompt(getState("user").."@mode "..getState("cwd").." $ ")
          else
              clearMode()
              setPrompt("username: ")
              setState("s", "login")
              say("Wrong username or password")
          end
      end

      local function doTerm()
          local cmd = args[1]
          if cmd == "ls" then
              say("secret.txt")
          elseif cmd == "cat" then
              if args[2] == "secret.txt" then say("there is no way out")
              else say("No such file") end
          elseif cmd == "exit" then
              clearMode()
          else
              say("command not found: "..tostring(cmd))
          end
      end

      local states = {
          ["login"] = doLogin,
          ["login-password"] = doPassword,
          ["term"] = doTerm,
      }

      local h = states[getState("s")]
      if h then h() end
    LUA
  }

  object "computer" {
    zone_id = "zone-two"
    desc = "an old computer"
    action "use" {
      rule {
        effect "set_mode" {
          mode_id = "mode-one"
        }
      }
    }
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

    properties = {
      "lock" = "locked"
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
