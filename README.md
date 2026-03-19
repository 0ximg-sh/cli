# 0ximg CLI

`0ximg` is a Go-based command-line tool for rendering source code into shareable images through the `0ximg.sh` API.

It supports reading code from:

- a file path
- standard input
- the system clipboard

The CLI also supports syntax highlighting, theme selection, line slicing, highlighted line ranges, layout customization, and automatic image download after rendering.

## Features

- Render code snippets to PNG images
- Read input from file, `stdin`, or clipboard
- Auto-detect language from file extension or shebang
- Render only a selected line range
- Highlight specific lines or ranges
- Customize theme, font, background, padding, shadows, and window chrome
- Print or copy a preview URL after rendering

## Requirements

- Go `1.25.7` or newer
- Network access to `https://0ximg.sh`

## Installation

### Build from source

```bash
go build -o 0ximg .
```

### Install with Go

```bash
go install 0ximg.sh/cli@latest
```

## Usage

```bash
0ximg [file] [flags]
```

If no file is provided, `0ximg` can read from `stdin`. You can also use `--from-clipboard` to render code directly from your clipboard.

## Basic Examples

### Render a file

```bash
0ximg main.go --theme Dracula --output main.go.png
```

### Render from standard input

```bash
cat main.go | 0ximg --language go --theme Dracula --output main.go.png
```

### Render from clipboard

```bash
0ximg --from-clipboard --language go --theme Nord --output snippet.png
```

### Render only selected lines

```bash
0ximg main.go --lines 10-20 --output snippet.png
```

### Highlight specific lines

```bash
0ximg main.go --lines 10-20 --highlight-lines 12-14 --output snippet.png
```

### List available themes

```bash
0ximg --list-themes
```

## Common Flags

| Flag | Description |
| --- | --- |
| `--from-clipboard` | Read code from the clipboard |
| `--list-themes` | Print available themes |
| `-o, --output` | Output image path |
| `-l, --language` | Language for syntax highlighting |
| `-t, --theme` | Theme used by the renderer |
| `-b, --background` | Background color |
| `-f, --font` | Font family |
| `--title` | Title metadata for rendering |
| `--window-title` | Custom window title |
| `--lines` | Render only a line range such as `10-20` |
| `--highlight-lines` | Highlight lines such as `1;5-12` |
| `--line-offset` | Starting line number |
| `--line-pad` | Padding between lines |
| `--no-line-number` | Hide line numbers |
| `--no-round-corner` | Disable rounded corners |
| `--no-window-controls` | Hide window controls |
| `--pad-horiz` | Horizontal padding |
| `--pad-vert` | Vertical padding |
| `--code-pad-right` | Right-side code padding |
| `--shadow-blur-radius` | Shadow blur radius |
| `--shadow-color` | Shadow color |
| `--shadow-offset-x` | Shadow X offset |
| `--shadow-offset-y` | Shadow Y offset |
| `--tab-width` | Tab width |

## Notes

- When `--lines` is used without `--line-offset`, the CLI automatically adjusts the displayed line numbers to match the original source file.
- If `--output` is not explicitly set, the tool derives an output path automatically.
- When rendering succeeds without an explicit `--output`, the CLI prints a preview URL and attempts to copy it to the clipboard.

## Version

```bash
0ximg version
```

Example output:

```text
version: dev
commit: none
date: unknown
```

## Development

Run tests with:

```bash
go test ./...
```

## Support

If this tool is useful to you, support the project here:

`https://buymeacoffee.com/levinhne`
