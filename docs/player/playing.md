# Playing a Game

guescape is a text adventure. You read where you are, then type what you want
to do. This guide covers starting a game and the basics of typing commands.

## Starting a Game

A game is a single world file (`.hcl`). Point guescape at it with the `-path`
flag:

```
guescape -path ./my-world.hcl
```

If you are running from source:

```
go run . -path ./my-world.hcl
```

The game prints the opening description and then waits for input at a prompt
that shows your current location:

```
start>
```

## Typing Commands

Type a command and press Enter. Commands are a verb, sometimes followed by a
target:

```
examine desk
take key
move north
```

Input is not case-sensitive, and extra spaces are ignored.

### Short forms for commands

You can type any unambiguous beginning of a command. All of these run
`examine`:

```
examine desk
exam desk
ex desk
```

### Short forms for targets

You do not have to type an object's full name. Any part of the name works, as
long as it matches only one thing you can currently reach:

```
examine engineering-coolant-console
examine coolant
examine cool
```

If what you type matches more than one thing, the game cannot tell which you
mean and will say it found nothing. Type a longer, more specific piece of the
name.

## Finding the Commands

The game knows the full list of commands, use `help`.

- `help` lists every command with a short description.
- `help <command>` explains one command and how to use it, for example
  `help use`.

A few forms worth knowing up front, because they are easy to miss:

- `move <direction>` gets you around, such as `move north`.
- `use <item> on <object>` uses something you are carrying on something in the
  world, such as `use key on door`.
- `look` (or `show`) describes where you are; `inventory` lists what you are
  carrying.

Everything else, and the exact wording for each command, is one `help` away.

## When a Command Does Not Work

The game answers in plain language. A few things it may tell you:

- "Nothing like that here" - it could not find what you named where you are.
- "That way is blocked" - there is a way out, but you cannot pass yet.
- "Cannot carry that" - the thing cannot be picked up.
- "You do not have that" - you tried to use an item you are not carrying.

The exact wording is chosen by the game's designer and may differ from world to
world.
