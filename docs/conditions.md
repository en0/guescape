# Action Rule Conditions

A `Condition` provides a set of different criterian.

**Note** that all the criterian in a condition must be true for the `Condigion` to be met.

## Property 
```hcl
when {
    property "property-id" { }
    flag "flag-id" { }
    use_with = "item-id"
}
```

- **Property**: A `PropertyCondition` list.
- **Flags**: A `FlagCondition` list.
- **use_with**: An item that must be used on this action.

### ProertyCondition

A list of required `Object` properties and their values that must exist on the current `Object`.

```hcl
property "property-id" {
    value = "property value"
}
```

- **Value**: The expected value of the given property on the current `Object`.

### FlagCondition

A list of required `Flag`s that the player must or must not have.

```hcl
flag "flag-id" {
    has = true
}
```

- **has**: A boolean indicating of the player must or must not have the given flag.

