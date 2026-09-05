# Obsidivim

This is a tool I wrote for myself in order to easily compose and edit notes
using a popup floating terminal window running this program.

## Requirements

- fzf
- go
- obsidian (not strictly necessary, I guess)

## Installation

Please note might not work out of the box on your system! You're welcome to
clone/fork the repo, change things to your liking, and build it yourself.

Export the following environment variables:

```bash
EDITOR
OBSIDIAN_VAULT_DIR
OBSIDIAN_TEMPLATES_DIR
```

Create the following template notes wherever `OBSIDIAN_TEMPLATES_DIR`
points to:

- fleeting_note_template.md
- research_note_template.md

Install the binary:
`go install github.com/rqpt/obsidivim`

## Usage

This is meant to be used from a floating terminal window, accessed by some
keybind you've set up.

Launching this, you're first met by a fzf menu with two options:

- New
- Existing

`Existing` runs another fzf menu with options of every single note in your vault.
Select one, and it opens the note in your preferred editor.

`New` runs another fzf menu with options of `Fleeting` and `Research`.

If you selected `Fleeting`, you would be prompted for a note title.
That's it... it is supposed to capture some fleeting thought you have that you
need to persist in your vault.

Selecting `Research`, also prompts you for a note title, then drops you into
your editor with the note open.
