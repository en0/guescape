## Ideas

### Save/Load state effects

Add an effect that the game can emit to cause a save checkpoint. This might be useful when entering
a particularly hard room.

Add an effect that the game can emit to cause a save-point to restore. This would be an alternative
to a hard death, game over option.

**open question** Should the player be able to trigger saves and restores themselves?

### Multi-room objects

**Consider** allowing a object to have a list of zones in which they exist. This specifically
addresses a locked-door problem.

**Counterpoint** A flag actually sovles this instead of using a object state. the logic has to match
but if you support rule referencing, that would solve this issue completely. Lastly, with no
changes, you could set the property on both doors in both rules to keep them in sync.

### HCL Object Subtypes

Subtypes will be useful for object-shorthand to make building worlds more ergonomic. They will
handle comon interactions through a standard interface.

- Flag: A Flag is a takeable item that grants a flag when carried. The flag is removed when dropped.
- Portal: A portal can add and remove exits between rooms.
- Surface: A surface reveals objects when examined.
- Container: A container, like a surface, that can be closed and locked.

### Rule Chaining

Consider adding a name to rules and accumulate matched rules into a set. Then, a rule can check if a
previous rule matched as part of it's condition. Or, ensure no previous rule matched as part of it's
condition.
