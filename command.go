package command

import (
	"fmt"
	"io"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[gloo.File, flags]

func Wc(parameters ...any) gloo.Command {
	return command(gloo.Initialize[gloo.File, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return gloo.Inputs[gloo.File, flags](p).Wrap(
		gloo.AccumulateAndOutput(func(lines []string, stdout io.Writer) error {
		var lineCount, wordCount, charCount, byteCount, maxLength int

		for _, line := range lines {
			lineCount++
			charCount += len([]rune(line))
			byteCount += len(line) + 1 // +1 for newline
			wordCount += len(strings.Fields(line))

			if len(line) > maxLength {
				maxLength = len(line)
			}
		}

		// Output based on flags (default: all)
		showAll := !bool(p.Flags.Lines) && !bool(p.Flags.Words) &&
		           !bool(p.Flags.Chars) && !bool(p.Flags.Bytes) &&
		           !bool(p.Flags.MaxLength)

		var output string
		if bool(p.Flags.Lines) || showAll {
			output += fmt.Sprintf("%7d ", lineCount)
		}
		if bool(p.Flags.Words) || showAll {
			output += fmt.Sprintf("%7d ", wordCount)
		}
		if bool(p.Flags.Chars) {
			output += fmt.Sprintf("%7d ", charCount)
		}
		if bool(p.Flags.Bytes) || showAll {
			output += fmt.Sprintf("%7d ", byteCount)
		}
		if bool(p.Flags.MaxLength) {
			output += fmt.Sprintf("%7d ", maxLength)
		}

		_, err := fmt.Fprintln(stdout, strings.TrimSpace(output))
		return err
	}).Executor(),
	)
}
