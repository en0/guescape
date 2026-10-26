# In-Game Commands for MVP

This document outlines the commands supported for *MVP*. Commands come in one of three types.

- **COMMAND**: A world-level command that operates on the current `Zone` or game.
- **VERB TARGET|DIRECTION**: A action taken on/in a specific `Object`/`Direction`.
- **VERB ITEM PREPOSITION TARGET**: A action applied to an `Object` with another.

## World-Level Commands

- **help**: Display the list of available commands
- **quit**: Exit the game
- **look**: Describe the current `Zone`, the directions you can go, and the *visible* `Object`s in
  the current `Zone`.
- **inventory**: List each `Object` in the player's inventory `Zone`.

## Navigation Commands

- **move `[direction]`**: Move the player from the current `Zone` to a connecting `Zone` associated
  with the current `Zone`'s attached `exit` in that identified direction, if any.

## Verbs

The target of an object must be *targetable*

- **examine `[object]`**: *Say* the content of the object's `desc` field
- **use `[object]`**: Execute the `action "use" {...}` handler
- **use `[item]` on `[object]`**: If `[item]` is in the player's inventory, call
  the `action "use" {...}` handler with the given `[item]`
- **take `[object]`**: if `carriable`, move the object to the player's inventory
- **drop `[object]`**: move the object from the player's inventory to the current zone
