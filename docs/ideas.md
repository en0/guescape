## Ideas

### Save/Load state effects

Add an effect that the game can emit to cause a save checkpoint. This might be useful when entering
a particularly hard room.

Add an effect that the game can emit to cause a save-point to restore. This would be an alternative
to a hard death, game over option.

**open question** Should the player be able to trigger saves and restores themselves?

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

### Expose the messaging system to the REPL

The REPL now resolves targets (objects, items, exits) through the resolver before building an
action. That resolution step runs in the CLI, so when it fails the CLI has to emit its own message.
Today those are hardcoded developer strings ("object not found: ...", "exit not found: ...") wrapped
in "Got an error: ...". This breaks immersion and, worse, ignores the designer's `message`
overrides.

It also created an inconsistency. Resolution-time failures (bad object/item/exit name) now bypass
the designer's voice entirely, while planning-time failures (way blocked, not carriable, already
taken, no handler) still flow through `message` as before. Some designer messages are now effectively
dead on the normal path. For example, `CannotMove` can no longer surface, because the resolver
rejects an unknown direction before the planner's exit check ever runs.

The idea is to open the messaging system up so the REPL can ask the world for a message while
keeping control of IO:

    ws.GameMessage(MsgCannotMove)

This returns the designer's override if present, or the default otherwise (the same resolution
`message()` already does internally). The REPL keeps owning when and how it prints; the designer
keeps owning the wording.

**open question** What is the right key vocabulary? The resolver failures do not map cleanly onto
the existing keys. A missing object could reuse `ObjectNotFound`, but "ambiguous, more than one
match" has no key today and is a distinct player-facing situation worth its own message.

**open question** Should resolution-time and planning-time failures share keys, or should resolution
get its own set? Sharing keeps the vocabulary small; separating lets a designer distinguish "no such
thing" from "that thing cannot be used."

**open question** Does exposing `GameMessage` leak a concern the domain would rather keep internal,
or is "give me the text for this situation" a legitimate public query alongside the other read-only
view methods?

#### Other REPL sites that bypass the designer

Beyond the resolvers, other hardcoded REPL strings skip the designer's voice. Candidates to route
through the messaging system:

- Unknown command ("Those words mean nothing here.").
- The quit flow ("Leave the story here? (y/n)", "The story continues!", and the "I will take that
  as a no" fallback).
- Empty inventory ("inventory empty"), which reads as a narrative beat rather than chrome.

Explicitly out of scope: the help listing and per-command help are UI and stay internal; the engine
should not expect a REPL-style interface. Win/lose banners are temporary and not part of this.
Engine errors ("Got an error", "this story has a hiccup") stay developer-facing and should not be
dressed in the designer's voice.

**open question** Should the `look`/inventory section headers ("exits:", "objects:") be
designer-controlled messages, or are they UI chrome that stays internal? The headers are only text
embedded in a structure the REPL owns (bullets, spacing, ordering), and the UI may change
drastically, which argues for chrome. But a designer may reasonably want "Ways out:" or "You see:".
The empty-inventory string ("inventory empty") is a cleaner case: it reads as a narrative beat
rather than a label, so it may belong to the designer even if the headers do not. Note this is
entangled with the targetable-markers idea below: if `look` collapses to description-only, the
headers disappear and this question is moot.

**open question** The invalid-input case ("doesn't take those words that way") should not key on CLI
grammar (argument counts, command names), since the engine should not presume a REPL. Reframe it as
a flat "confused" message: the player asked for something the game cannot make sense of
(e.g. "I don't understand"). A single `Confused` key could absorb every malformed-input path
(unknown command, bad argument shape, help-context unknown command), letting the CLI change its
command surface freely without minting new keys. Only flat messages; no templating.

**open question** Is there a meaningful split between "confused" (input not understood) and the
existing `NoHandler` (understood, but not valid here)? A two-pole model (`Confused` vs `NoHandler`)
may cover most situations without a sprawl of per-situation keys.

**open question** Does a well-formed-but-incomplete command (e.g. `use` with no object) deserve
different feedback than pure confusion? "Use what?" is command-aware and parameterized, which
conflicts with flat-only, so it may have to collapse into `Confused` and accept blunter feedback.

**open question** More broadly: reconsider the entire set of currently defined message keys and
reduce it where possible. The resolver and confused-message work will likely reveal keys that
overlap or that can be collapsed into a smaller, situation-based vocabulary.

### Targetable markers in description text

Introduce a marker in description text that flags a word as referring to a targetable thing, using
a sigil such as `@`:

    "You see a @table in the corner with a @dusty-book atop."

The marker serves two purposes. First, presentation: the UI could highlight targetable nouns so the
player knows what they can act on. Second, validation: at load time the engine parses each marker
and confirms it resolves to a real, targetable object in scope. That gives the designer compile-time
safety that every object they *mention* is one the player can actually *act on*, and vice versa.
The description becomes the source of truth for what is targetable in the player's mental model.

#### Relation to the look output (#5/#6)

If descriptions carry their own targetable nouns, the generated `exits:`/`objects:` lists in `look`
become redundant: the prose already tells the player what is here and what they can touch. That
would collapse `look` to just the description and dissolve the section-header question entirely,
since there would be no lists to label.

**open question** The generated lists are currently a safety net. `ViewCurrentZone` derives them
from live state (`visible` and current zone), so the player is always told about every visible
object, recomputed every `look`. Removing the lists moves that guarantee onto the designer, who
would then hand-maintain agreement between the description's marked nouns and the actual visible set.

**open question** Static markers do not track dynamic visibility. A marker validates that a noun
exists and is targetable at load time, but an object revealed or hidden at runtime (via an
`object_visible` effect) would make static description text stale. The lists never had this problem
because they are recomputed. Is static validation enough, or does the player need live presence?

**open question** A middle path: make marker rendering conditional on current visibility. The engine
renders a marked noun as plain text (or omits its clause) when the object is not currently visible,
and shows/highlights it when it is. This preserves the safety net and keeps a single authored
description, at the cost of more engine machinery. Worth the complexity?

**open question** Is the sync burden simply too much to put on designers, regardless of mechanism?
The appeal of the current lists is that the designer manages only the `visible` flag and the engine
handles discovery.

### HCL: Move objects into zone block

When you have a lot of zones and a lot of objects, it gets confusing.

```hcl
game {
    zone "foo" {
        object "bar" {
        }
    }
}
```
