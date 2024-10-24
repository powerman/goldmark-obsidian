# goldmark-obsidian

[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/goldmark-obsidian.svg)](https://pkg.go.dev/github.com/powerman/goldmark-obsidian)
[![CI/CD](https://github.com/powerman/goldmark-obsidian/actions/workflows/CI&CD.yml/badge.svg)](https://github.com/powerman/goldmark-obsidian/actions/workflows/CI&CD.yml)
[![Coverage Status](https://coveralls.io/repos/github/powerman/goldmark-obsidian/badge.svg?branch=master)](https://coveralls.io/github/powerman/goldmark-obsidian?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/powerman/goldmark-obsidian)](https://goreportcard.com/report/github.com/powerman/goldmark-obsidian)
[![Release](https://img.shields.io/github/v/release/powerman/goldmark-obsidian)](https://github.com/powerman/goldmark-obsidian/releases/latest)

An [Obsidian](https://obsidian.md/) extension for the
[goldmark](https://github.com/yuin/goldmark) markdown parser.

## Features

### Obsidian

- [x] Autolinks (GFM): `https://example.com`
- [x] Internal links: `[[Link]]`
- [x] Embed files: `![[Link]]`
- [x] Block references: `![[Link#^id]]`
- [x] Defining a block: `^id`
- [ ] Comments (block, inline): `%%Text%%`
- [x] Strikethroughs (GFM): `~~Text~~`
- [ ] Highlights: `==Text==`
- [x] Code blocks: ` ``` `
- [x] Incomplete task (GFM): `- [ ]`
- [x] Completed task (GFM): `- [x]`
- [ ] Callouts: `> [!note]`
- [x] Tables (GFM)
- [x] Footnotes (reference): `Example[^1].`
- [ ] Footnotes (inline): `Example. ^[This is an inline footnote.]`
- [x] Diagrams (Mermaid): ` ```mermaid `
- [x] LaTeX (MathJax): `$$`
- [x] Tags: `#tag`
- [x] Properties (metadata): `---` at the very beginning of a file

Known inconsistencies with Obsidian:

- Properties may start and end with at least 1 dash, not exactly 3 (---).
- Tags defined in properties are not applied to document in same way as other tags.
- Document aliases defined in properties are not processed as internal link targets.

### Obsidian plugins

- [x] [Tasks](https://github.com/obsidian-tasks-group/obsidian-tasks).

## Installation

```sh
go get github.com/powerman/goldmark-obsidian
```

## Usage

```go
source := []byte(`
- [ ] Happy New Year 📅 2025-01-01 ^first-task
- [x] Happy Old Year 📅 2024-01-01
`)

md := goldmark.New(
    goldmark.WithExtensions(
        obsidian.NewPlugTasks(),
        obsidian.NewObsidian(),
    ),
)
err := md.Convert(source, os.Stdout)
if err != nil {
    fmt.Println(err)
}
// Output:
// <ul class="contains-task-list">
// <li data-task="" class="task-list-item" id="^first-task"><input disabled="" type="checkbox" class="task-list-item-checkbox"> Happy New Year 📅 2025-01-01</li>
// <li data-task="x" class="task-list-item is-checked"><input checked="" disabled="" type="checkbox" class="task-list-item-checkbox"> Happy Old Year 📅 2024-01-01</li>
// </ul>
```
