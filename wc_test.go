package wc_test

import (
	"context"
	"os"
	"strings"

	"github.com/yupsh/wc"
	"github.com/yupsh/wc/opt"
)

func ExampleWc() {
	ctx := context.Background()
	input := strings.NewReader("Hello world\nThis is a test\nWith multiple lines\n")

	cmd := wc.Wc() // Default: lines, words, bytes
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	//        3       9      47
}

func ExampleWc_linesOnly() {
	ctx := context.Background()
	input := strings.NewReader("Line 1\nLine 2\nLine 3\n")

	cmd := wc.Wc(opt.Lines)
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	//        3
}

func ExampleWc_wordsOnly() {
	ctx := context.Background()
	input := strings.NewReader("Hello world this is a test")

	cmd := wc.Wc(opt.Words)
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	//        6
}

func ExampleWc_bytesOnly() {
	ctx := context.Background()
	input := strings.NewReader("Hello")

	cmd := wc.Wc(opt.Bytes)
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	//        6
}

func ExampleWc_charsOnly() {
	ctx := context.Background()
	input := strings.NewReader("Hello")

	cmd := wc.Wc(opt.Chars)
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	//        6
}
