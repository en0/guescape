# Action Rules

Rules provide a list of conditional `Effect`s for custom *Action* handers. They describe the
execpted state of particular elements of the game and provide `Effect`s that should be emitted when
any one of the given `Condition`s are `true`.

**Note** that the `Rule`'s `Effect`s are emitted as long as **any** `Condition` is met.

```hcl
rule {
    when { }
    effect { }
}
```

- **When**: A list of `Condition`s
- **effect**: A list of `Effect`s to emit.

