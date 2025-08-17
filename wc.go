package wc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/wc/opt"
)

// Statistics holds the count results
type Stats struct {
	Lines     int
	Words     int
	Chars     int
	Bytes     int
	MaxLength int
}

// Flags represents the configuration options for the wc command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Wc creates a new wc command with the given parameters
func Wc(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// If no specific flags are set, show lines, words, and bytes (default behavior)
	showDefault := !bool(c.Flags.Lines) && !bool(c.Flags.Words) && !bool(c.Flags.Chars) && !bool(c.Flags.Bytes) && !bool(c.Flags.MaxLength)

	// If no files specified, read from stdin
	if len(c.Positional) == 0 {
		stats, err := c.countReader(ctx, stdin)
		if err != nil {
			return err
		}
		c.printStats(stdout, stats, "", showDefault)
		return nil
	}

	var totalStats Stats
	multipleFiles := len(c.Positional) > 1

	err := yup.ProcessFilesWithContext(
		ctx, c.Positional, stdin, stdout, stderr,
		yup.FileProcessorOptions{
			CommandName:     "wc",
			ContinueOnError: true,
		},
		func(ctx context.Context, source yup.InputSource, output io.Writer) error {
			stats, err := c.countReader(ctx, source.Reader)
			if err != nil {
				return err
			}

			c.printStats(output, stats, source.Filename, showDefault)

			// Add to totals
			totalStats.Lines += stats.Lines
			totalStats.Words += stats.Words
			totalStats.Chars += stats.Chars
			totalStats.Bytes += stats.Bytes
			if stats.MaxLength > totalStats.MaxLength {
				totalStats.MaxLength = stats.MaxLength
			}

			return nil
		},
	)

	// Print totals if multiple files
	if multipleFiles && err == nil {
		c.printStats(stdout, totalStats, "total", showDefault)
	}

	return err
}

func (c command) countReader(ctx context.Context, reader io.Reader) (Stats, error) {
	var stats Stats
	scanner := bufio.NewScanner(reader)

	for yup.ScanWithContext(ctx, scanner) {
		line := scanner.Text()
		stats.Lines++

		// Count characters (runes)
		runes := []rune(line)
		stats.Chars += len(runes) + 1 // +1 for newline

		// Count bytes
		stats.Bytes += len(scanner.Bytes()) + 1 // +1 for newline

		// Count words
		words := strings.FieldsFunc(line, func(r rune) bool {
			return unicode.IsSpace(r)
		})
		stats.Words += len(words)

		// Track max line length
		if len(runes) > stats.MaxLength {
			stats.MaxLength = len(runes)
		}
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return stats, err
	}

	return stats, scanner.Err()
}

func (c command) printStats(output io.Writer, stats Stats, filename string, showDefault bool) {
	var parts []string

	if bool(c.Flags.Lines) || showDefault {
		parts = append(parts, fmt.Sprintf("%8d", stats.Lines))
	}

	if bool(c.Flags.Words) || showDefault {
		parts = append(parts, fmt.Sprintf("%8d", stats.Words))
	}

	if bool(c.Flags.Chars) && !bool(c.Flags.Bytes) {
		parts = append(parts, fmt.Sprintf("%8d", stats.Chars))
	}

	if bool(c.Flags.Bytes) || (showDefault && !bool(c.Flags.Chars)) {
		parts = append(parts, fmt.Sprintf("%8d", stats.Bytes))
	}

	if bool(c.Flags.MaxLength) {
		parts = append(parts, fmt.Sprintf("%8d", stats.MaxLength))
	}

	result := strings.Join(parts, "")
	if filename != "" {
		result += " " + filename
	}

	fmt.Fprintln(output, result)
}
