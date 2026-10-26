# World Zone

A `Zone` is a area in the game to which objects and exists are associated.

```hcl
zone "zone-id" {
    desc = "zone description"
    exit "direction" { }
}
```

- **desc**: A string that describes the zone. Returned as a *say* `Effect` to a `look` command.
- **exit**: A `Exit` map and the associated `direction` name.



