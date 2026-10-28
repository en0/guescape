raw input → tokenize → normalize first token → world command? → nav command? → object command


---

Mode ideas

We add a new mode slice, `[]Mode` to the world state. the lua vm runs on `plan_effects` with a new
action type, `ModeInputAction` that takes the raw input and the token array. The lua script emits
effects through some custom functions like `say()`. 

Things to work out:

1. What world state do we need to expose to lua and whats the grammer look like for that world state
   and the effect methods and the prompt/state mutation requests (still just effects)?
2. how do i actually run the lua and create the custom functions?

Examples below:

```go
for r.running {
    p := r.ws.Prompt()          // game answers: mode's prompt, or default "zone-id> "
    tokens, err := r.prompt(p.Text)
    ...
    var act game.Action
    if p.ModeActive {           // same query tells the REPL which action to build
        act = game.NewModeInputAction(r.in.Text(), tokens)
    } else {
        // existing findCommand dispatch -> builds a normal Action
    }
    for _, ef := range game.PlanEffects(r.ws, act) { ... apply ... }
}

// So ws.Prompt() (or call it ws.InputContext()) returns a small value that answers both questions the REPL has at the top of the loop:

type InputContext struct {
    Prompt     string
    ModeActive bool   // the REPL's only signal, read fresh each turn
}
```

```hcl
game {
    ...
    mode "some-mode-id" {
        prompt = "Enter Code: "
        state {
            some_state_variable = "locked"
        }
        script = <<-LUA
            local state = get_state("some_state_variable")
            if state == "locked" then
                if arg[0] == "1234" then
                    say("Accepted")
                    set_prompt("UNLOCKED")
                    clear_mode()
                end
            end
        LUA
    }
}
```
