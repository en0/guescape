# Effect System

Effects are emited as the result of in-game interactions. The effects are collected and "reduced"
(applied) by the game engine to the world `State`. This document outlines the `Effect` types
available in the MVP release.


## SayEffect

Present text to the user.

```hcl
effect "say" { text =  "Hello, world!" }
```

- Text: string -> The text to *say* to the player.

## CurrentZoneEffect

Set the player's current `Zone`.

```hcl
effect "current_zone" { zone_id = "library" }
```

- zone_id: Identity -> The id of the new current `Zone`.

## FlagEffect

Add or remove a flag.

```hcl
effect "flag" {
  flag      = "door-unlocked"
  operation = "add"
}
```

- flag: Identity -> The `Identity` of the flag to operate on.
- operation: string -> Specifies the flag operation. One of: *add*, *remove*.

## PropertyEffect

Set a property value on an `Object`.

```hcl
effect "property" {
  object_id   = "chest"
  property_id = "state"
  value       = "open"
}
```

- object_id: Identity -> The *optional* `Identity` of the `Object` to change. If not given, the
  current `object` id is used.
- property_id: Identity -> The `Identity` of the property to set.
- value: string -> The string value to set on the identified property.

## ObjectVisibleEffect

Change the visibility of an `Object`.

```hcl
effect "object_visible" {
  object_id = "key"
  visible   = true
}
```

- object_id: Identity -> The `Identity` of the `Object` to change.
- visible: bool -> A boolean value indicating the `Object` visibility.

## ObjectTargetableEffect

Change the targetability of an `Object`.

```hcl
effect "object_targetable" {
  object_id  = "key"
  targetable = true
}
```

- object_id: Identity -> The `Identity` of the `Object` to change.
- targetable: bool -> A boolean value indicating the `Object` targetability.

## ObjectLocationEffect

Move an `Object` to a `Zone`.

```hcl
effect "object_location" {
  object_id = "key"
  zone_id   = "library"
}
```

- object_id: Identity -> The `Identity` of the `Object` to change.
- zone_id: Identity -> The `Identity` of the zone to set on the target `Object`.

## EndGameEffect

End the game, presenting a message to the user.

```hcl
effect "end_game" {
  is_win = true
  text   = "You escaped!"
}
```

- is_win: bool -> A boolean value indicating a *winning* end-game state.
- text: string -> A message to present to the user.
