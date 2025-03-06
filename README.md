<div align="center">

[![Slug Craft](https://github.com/mnuddindev/slugcraft/blob/main/slugcraft_logo.png)](https://github.com/mnuddindev/slugcraft.git)

# Slug Craft

[![Go version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/mnuddindev/slugcraft)
[![Go report](https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none)](https://goreportcard.com/report/github.com/mnuddindev/slugcraft)
[![Code coverage](https://img.shields.io/badge/code_coverage-88%25-success?style=for-the-badge&logo=none)](https://github.com/mnuddindev/slugcraft.git)<br/>
[![Wiki](https://img.shields.io/badge/docs-wiki_page-blue?style=for-the-badge&logo=none)](https://github.com/create-go-app/cli/wiki)
[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge&logo=none)](https://github.com/mnuddindev/slugcraft/blob/main/LICENSE)
[![Build Status](https://github.com/mnuddindev/slugcraft/actions/workflows/go.yml/badge.svg)](https://github.com/mnuddindev/slugcraft/actions)

**The ultimate Go package for crafting URL-friendly slugs.** Fast, flexible, and built for the real world—SlugCraft handles multilingual text, avoids collisions, and optimizes for SEO and UX like no other.
</div>

## Why SlugCraft?

- 🌍 **Multilingual Magic**: Smart transliteration for non-Latin scripts (e.g., "привет" → "privet").
- ⚡ **Blazing Fast**: Zero-allocation options and bulk processing for scale.
- 🛠️ **Configurable**: Build your own slug pipeline with ease.
- 🚀 **Unique Slugs**: Collision avoidance with timestamps, UUIDs, or custom suffixes.
- ✨ **SEO & UX Ready**: Human-readable, meaningful slugs out of the box.

Say goodbye to boring, brittle slug libraries. SlugCraft is here to level up your Go projects.

---

## Package Installation
To use SlugCraft as a library in your Go project, install it with:

```bash
go get github.com/mnuddindev/slugcraft@latest
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/mnuddindev/slugcraft"
)

func main() {
	s := slugcraft.New(
		slugcraft.WithLanguage("bn"),
		slugcraft.WithStopwords("en"),
		slugcraft.WithAbbreviation("বাংলা", "BN"),
	)
	slug, err := s.Make(context.Background(), "বাংলা the World")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(slug) // Output: "bn-world"
}
```
## CLI Installation
To install the SlugCraft CLI tool globally on your machine, use:

```bash
go install github.com/mnuddindev/slugcraft/cmd/slugcraft@latest
```

## CLI Usage
Generate slugs directly from the command line:

```bash
slugcraft -input "বাংলা প্রিয়" -lang=bn

output: bangla-priyo

```
## Available Flags
```shell
Available Flags
    -input string: Text to slugify (required)
    -lang string: Language (e.g., bn, ru; optional)
    -cache bool: Enable cache for uniqueness (default: false)
    -suffix string: Suffix style (numeric, version, revision; default: numeric)
    -max int: Maximum slug length (default: 100)
    -stopwords string: Language for stopwords (e.g., en; optional)
    -regex string: Regex filter pattern (e.g., [^a-z0-9-]) (optional)
    -replace string: Regex replacement (default: "")
    -abbr string: Abbreviations (e.g., বাংলা=BN,আমি=ME) (optional)
    -help: Show usage info
```

## CLI Examples
```shell
# English with stopwords and regex
slugcraft -input "Hello the World!" -stopwords=en -regex="[^a-z0-9-]" -replace=""
# Output: hello-world

# Bangla with abbreviations
slugcraft -input "বাংলা আমি" -lang=bn -abbr="বাংলা=BN,আমি=ME"
# Output: bn-me
```

## License

The source files are distributed under the
[the Massachusetts Institute of Technology](https://github.com/mnuddindev/slugcraft/blob/main/LICENSE),
unless otherwise noted.
