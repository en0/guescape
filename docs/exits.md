# Zone Exit

A `Exit` exists on a `Zone`. It provides the available `Directions` the player can `move` to. An
`Exit` can require a `Flag` to be availble. If the user does not have the required `Flag`, the exit
is considered *blocked*.

```hcl
exit "direction" {
    flags = ["win-flag"]
    to = "win"
}
```

- **Flags**: A optional `Identity` list of required flags.
- **to**: The `Identity` of the target `Zone`.
