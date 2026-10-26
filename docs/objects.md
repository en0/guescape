# Built-in Objects

The system supports a set of buit-in objects that cary specific types of functionality.

```hcl
object "object-id" {
    zone_id     = "zone-id"
    visible     = true || false
    targetable  = true || false
    carriable   = true || false
    desc        = "object description"

    action "action-name" {
    }
}
```

## Properties

- **Id**: Unique identification for the object
- **zone_id**": The location of this object: A zone or inventory
- **visible**: A boolean indicating the object is visible in the associated zone
- **targetable**: A boolean indicating the object can be used as the target of an object verb
- **carriable**: A boolean indicating the object can be put in the player's inventory
- **desc**: A string returned when the object is the target of an `examine` verb
- **action <VERB>**: Create an action handler for the given verb. Contains `rules`
