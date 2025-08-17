package command

import (
	"fmt"
	"io"
	"strings"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Wc(parameters ...any) yup.Command {
	return command(yup.Initialize[yup.File, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return yup.Inputs[yup.File, flags](p).Wrap(
		yup.AccumulateAndOutput(func(lines []string, stdout io.Writer) error {
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
