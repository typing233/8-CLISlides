---
title: "CLISlides Demo"
author: "Demo Author"
date: "2026-06-07"
theme: "dark"
pager: true
---

# Welcome to CLISlides

A terminal-based presentation tool built in Go.

**Features:**
- Markdown rendering with glamour
- Full keyboard navigation
- Code execution
- SSH remote sharing

---

## Navigation

| Key | Action |
|-----|--------|
| `→` / `l` / `n` | Next slide |
| `←` / `h` / `p` | Previous slide |
| `gg` | First slide |
| `G` | Last slide |
| `5G` | Jump to slide 5 |
| `/` | Search (regex) |
| `e` / `x` | Execute code |
| `q` / `Ctrl+C` | Quit |

---

## Markdown Rendering

### Headings work great

- Bullet points
- **Bold** and *italic*
- `inline code`

> Blockquotes too!

1. Numbered lists
2. Also work well

---

## Code Blocks

```bash
echo "Hello from CLISlides!"
date
uname -a
```

Press `e` or `x` to execute the code block above.

---

## Preprocessor Example

The following block uses `~~~` preprocessor syntax.
When the presentation loads, it executes the command
and replaces the block with its output:

~~~
echo "This was generated at load time: $(date)"
~~~

---

# Thank You!

That's all folks. Press `q` to quit.
