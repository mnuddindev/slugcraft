package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mnuddindev/slugcraft/"
)

func main() {
	// Define flags
	input := flag.String("input", "", "Text to slugify")
	lang := flag.String("lang", "", "Language (bn, default: en)")
	cache := flag.Bool("cache", false, "Enable in-memory cache for uniqueness")
	suffix := flag.String("suffix", "numeric", "Suffix style: numeric, version, revision")
	maxLength := flag.Int("max", 100, "Maximum slug length")
	stopwords := flag.String("stopwords", "", "Language for stopwords (e.g., en)")
	regex := flag.String("regex", "", "Regex pattern to filter (e.g., [^a-z0-9-])")
	regexReplace := flag.String("replace", "", "Replacement for regex filter")
	abbr := flag.String("abbr", "", "Abbreviations (format: key1=value1,key2=value2)")
	help := flag.Bool("help", false, "Show usage information")

	flag.Parse()

	// Show help or validate input
	if *help || *input == "" {
		printUsage()
		os.Exit(0)
	}

	// Parse abbreviations
	abbreviations := make(map[string]string)
	if *abbr != "" {
		pairs := strings.Split(*abbr, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				abbreviations[kv[0]] = kv[1]
			}
		}
	}

	// Create Slugger with options
	opts := []slugcraft.Option{}
	if *lang != "" {
		if *lang == "bn" || *lang == "en" {
			opts = append(opts, slugcraft.WithLanguage(*lang))
		} else {
			fmt.Println("Star the repository and wait for more language support. \n https://github.com/mnuddindev/slugcraft")
		}
	}
	if *cache {
		opts = append(opts, slugcraft.WithUseCache(true))
	}
	if *suffix != "" {
		opts = append(opts, slugcraft.WithSuffixStyle(*suffix))
	}
	if *maxLength > 0 {
		opts = append(opts, slugcraft.WithMaxLength(*maxLength))
	}
	if *stopwords != "" {
		opts = append(opts, slugcraft.WithStopwords(*stopwords))
	}
	if *regex != "" {
		opts = append(opts, slugcraft.WithRegexFilter(*regex, *regexReplace))
	}
	for k, v := range abbreviations {
		opts = append(opts, slugcraft.WithAbbreviation(k, v))
	}

	s := slugcraft.New(opts...)

	// Generate slug
	slug, err := s.Make(context.Background(), *input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(slug)
}

func printUsage() {
	fmt.Println("SlugCraft CLI - Generate slugs from text")
	fmt.Println("Usage: slugcraft [flags]")
	fmt.Println("Flags:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  -%s\t%s (default: %v)\n", f.Name, f.Usage, f.DefValue)
	})
	fmt.Println("Examples:")
	fmt.Println(`  slugcraft -input "বাংলা প্রিয়" -lang=bn`)
	fmt.Println(`  slugcraft -input "Привет, мир!" -lang=ru -cache=true -suffix=version`)
	fmt.Println(`  slugcraft -input "Hello the World" -stopwords=en -regex="[^a-z0-9-]" -replace=""`)
	fmt.Println(`  slugcraft -input "বাংলা আমি" -lang=bn -abbr="বাংলা=BN,আমি=ME"`)
}
