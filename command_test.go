package command_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/wc"
)

// ==============================================================================
// Test Default Behavior (lines, words, bytes)
// ==============================================================================

func TestWc_Default(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines("hello world", "foo bar").
		Run()

	assertion.NoError(t, result.Err)
	// Default: lines, words, bytes
	// 2 lines, 4 words, 20 bytes (11 + 1 + 7 + 1)
	assertion.Contains(t, result.Stdout, "2")
	assertion.Contains(t, result.Stdout, "4")
	assertion.Contains(t, result.Stdout, "20")
}

func TestWc_EmptyInput(t *testing.T) {
	result := run.Quick(command.Wc())

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "0")
}

func TestWc_SingleLine(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	// 1 line, 1 word, 6 bytes (5 + newline)
	assertion.Contains(t, result.Stdout, "1")
	assertion.Contains(t, result.Stdout, "6")
}

func TestWc_EmptyLine(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines("").
		Run()

	assertion.NoError(t, result.Err)
	// 1 line, 0 words, 1 byte (just newline)
	output := strings.Fields(result.Stdout[0])
	assertion.Equal(t, len(output), 3, "three counts")
	assertion.Equal(t, output[0], "1", "lines")
	assertion.Equal(t, output[1], "0", "words")
	assertion.Equal(t, output[2], "1", "bytes")
}

// ==============================================================================
// Test Individual Flags
// ==============================================================================

func TestWc_LinesOnly(t *testing.T) {
	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines("line1", "line2", "line3").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "3", "three lines")
}

func TestWc_WordsOnly(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("hello world", "foo bar baz").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "five words")
}

func TestWc_BytesOnly(t *testing.T) {
	result := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "6", "six bytes")
}

func TestWc_CharsOnly(t *testing.T) {
	result := run.Command(command.Wc(command.Chars)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "five characters")
}

func TestWc_MaxLengthOnly(t *testing.T) {
	result := run.Command(command.Wc(command.MaxLength)).
		WithStdinLines("short", "medium line", "a").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "11", "max length 11")
}

// ==============================================================================
// Test Flag Combinations
// ==============================================================================

func TestWc_LinesAndWords(t *testing.T) {
	result := run.Command(command.Wc(command.Lines, command.Words)).
		WithStdinLines("hello world", "foo").
		Run()

	assertion.NoError(t, result.Err)
	// 2 lines, 3 words
	assertion.Contains(t, result.Stdout, "2")
	assertion.Contains(t, result.Stdout, "3")
}

func TestWc_AllCounts(t *testing.T) {
	result := run.Command(command.Wc(
		command.Lines,
		command.Words,
		command.Chars,
		command.Bytes,
		command.MaxLength,
	)).WithStdinLines("hello", "world").Run()

	assertion.NoError(t, result.Err)
	// 2 lines, 2 words, 10 chars, 12 bytes, max 5
	output := strings.Fields(result.Stdout[0])
	assertion.Equal(t, len(output), 5, "five counts")
}

// ==============================================================================
// Test Line Counting
// ==============================================================================

func TestWc_MultipleLines(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5"}
	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "five lines")
}

func TestWc_ManyLines(t *testing.T) {
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = "line"
	}

	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "100")
}

func TestWc_EmptyLines(t *testing.T) {
	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines("", "", "").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "3", "three empty lines")
}

// ==============================================================================
// Test Word Counting
// ==============================================================================

func TestWc_SingleWord(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "1", "one word")
}

func TestWc_MultipleWords(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("one two three four five").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "five words")
}

func TestWc_WordsAcrossLines(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("hello world", "foo bar").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "4", "four words")
}

func TestWc_WordsWithWhitespace(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("  hello   world  ").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "2", "two words")
}

func TestWc_NoWords(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("", "   ", "\t\t").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "0", "zero words")
}

// ==============================================================================
// Test Byte Counting
// ==============================================================================

func TestWc_Bytes_ASCII(t *testing.T) {
	result := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	// "hello" = 5 bytes + 1 newline = 6
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "6", "six bytes")
}

func TestWc_Bytes_MultipleLines(t *testing.T) {
	result := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("abc", "def").
		Run()

	assertion.NoError(t, result.Err)
	// "abc\n" = 4, "def\n" = 4, total = 8
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "8", "eight bytes")
}

func TestWc_Bytes_Empty(t *testing.T) {
	result := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("").
		Run()

	assertion.NoError(t, result.Err)
	// Empty line = just newline = 1 byte
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "1", "one byte")
}

func TestWc_Bytes_Unicode(t *testing.T) {
	result := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("日本語").
		Run()

	assertion.NoError(t, result.Err)
	// "日本語" in UTF-8 = 9 bytes + 1 newline = 10
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "10", "ten bytes")
}

// ==============================================================================
// Test Character Counting
// ==============================================================================

func TestWc_Chars_ASCII(t *testing.T) {
	result := run.Command(command.Wc(command.Chars)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	// "hello" = 5 characters (newline not counted in char count)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "five characters")
}

func TestWc_Chars_Unicode(t *testing.T) {
	result := run.Command(command.Wc(command.Chars)).
		WithStdinLines("日本語").
		Run()

	assertion.NoError(t, result.Err)
	// "日本語" = 3 characters (not 9 bytes)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "3", "three characters")
}

func TestWc_Chars_Mixed(t *testing.T) {
	result := run.Command(command.Wc(command.Chars)).
		WithStdinLines("hello世界").
		Run()

	assertion.NoError(t, result.Err)
	// "hello" (5) + "世界" (2) = 7 characters
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "7", "seven characters")
}

func TestWc_Chars_Emoji(t *testing.T) {
	result := run.Command(command.Wc(command.Chars)).
		WithStdinLines("😀👋").
		Run()

	assertion.NoError(t, result.Err)
	// 2 emoji = 2 characters
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "2", "two characters")
}

// ==============================================================================
// Test Max Length
// ==============================================================================

func TestWc_MaxLength_SingleLine(t *testing.T) {
	result := run.Command(command.Wc(command.MaxLength)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "max length 5")
}

func TestWc_MaxLength_MultipleLines(t *testing.T) {
	result := run.Command(command.Wc(command.MaxLength)).
		WithStdinLines("short", "this is a much longer line", "tiny").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "26", "max length 26")
}

func TestWc_MaxLength_Empty(t *testing.T) {
	result := run.Command(command.Wc(command.MaxLength)).
		WithStdinLines("").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "0", "max length 0")
}

func TestWc_MaxLength_EmptyAndNonEmpty(t *testing.T) {
	result := run.Command(command.Wc(command.MaxLength)).
		WithStdinLines("", "hello", "").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "5", "max length 5")
}

// ==============================================================================
// Test With Whitespace
// ==============================================================================

func TestWc_Whitespace_LeadingSpaces(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("   hello world").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "2", "two words")
}

func TestWc_Whitespace_TrailingSpaces(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("hello world   ").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "2", "two words")
}

func TestWc_Whitespace_Tabs(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("hello\tworld").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "2", "two words")
}

func TestWc_Whitespace_OnlySpaces(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines("   ").
		Run()

	assertion.NoError(t, result.Err)
	// 1 line, 0 words
	assertion.Contains(t, result.Stdout, "1")
	assertion.Contains(t, result.Stdout, "0")
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestWc_InputError(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinError(errors.New("read failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "read failed")
}

func TestWc_OutputError(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines("test").
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestWc_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		flag          any
		expectedValue string
	}{
		{
			name:          "lines - three",
			input:         []string{"a", "b", "c"},
			flag:          command.Lines,
			expectedValue: "3",
		},
		{
			name:          "words - five",
			input:         []string{"one two", "three four five"},
			flag:          command.Words,
			expectedValue: "5",
		},
		{
			name:          "bytes - abc",
			input:         []string{"abc"},
			flag:          command.Bytes,
			expectedValue: "4",
		},
		{
			name:          "chars - abc",
			input:         []string{"abc"},
			flag:          command.Chars,
			expectedValue: "3",
		},
		{
			name:          "max length",
			input:         []string{"hi", "hello", "yo"},
			flag:          command.MaxLength,
			expectedValue: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Wc(tt.flag)).
				WithStdinLines(tt.input...).
				Run()

			assertion.NoError(t, result.Err)
			output := strings.TrimSpace(result.Stdout[0])
			assertion.Equal(t, output, tt.expectedValue, "expected value")
		})
	}
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestWc_CodeFile(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines(
			"package main",
			"",
			"func main() {",
			"}",
		).Run()

	assertion.NoError(t, result.Err)
	// 4 lines, 4 words
	assertion.Contains(t, result.Stdout, "4")
}

func TestWc_Document(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines(
			"The quick brown fox",
			"jumps over the",
			"lazy dog",
		).Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "9", "nine words")
}

func TestWc_CSVFile(t *testing.T) {
	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines(
			"Name,Age,City",
			"Alice,30,NYC",
			"Bob,25,LA",
		).Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "3", "three lines")
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestWc_VeryLongLine(t *testing.T) {
	longLine := strings.Repeat("a", 10000)
	result := run.Command(command.Wc(command.MaxLength)).
		WithStdinLines(longLine).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "10000")
}

func TestWc_ManyShortLines(t *testing.T) {
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = "x"
	}

	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Contains(t, result.Stdout, "1000")
}

func TestWc_UnicodeWords(t *testing.T) {
	result := run.Command(command.Wc(command.Words)).
		WithStdinLines("こんにちは 世界").
		Run()

	assertion.NoError(t, result.Err)
	output := strings.TrimSpace(result.Stdout[0])
	assertion.Equal(t, output, "2", "two words")
}

func TestWc_MixedContent(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines(
			"normal line",
			"",
			"line with\ttabs",
			"  spaces  ",
			"unicode: 日本語",
		).Run()

	assertion.NoError(t, result.Err)
	// 5 lines
	assertion.Contains(t, result.Stdout, "5")
}

// ==============================================================================
// Test Bytes vs Characters Difference
// ==============================================================================

func TestWc_BytesVsChars_ASCII(t *testing.T) {
	resultBytes := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("hello").
		Run()

	resultChars := run.Command(command.Wc(command.Chars)).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, resultBytes.Err)
	assertion.NoError(t, resultChars.Err)

	// ASCII: bytes includes newline, chars doesn't
	bytesOut := strings.TrimSpace(resultBytes.Stdout[0])
	charsOut := strings.TrimSpace(resultChars.Stdout[0])

	assertion.Equal(t, bytesOut, "6", "6 bytes (with newline)")
	assertion.Equal(t, charsOut, "5", "5 chars (without newline)")
}

func TestWc_BytesVsChars_Unicode(t *testing.T) {
	resultBytes := run.Command(command.Wc(command.Bytes)).
		WithStdinLines("日本語").
		Run()

	resultChars := run.Command(command.Wc(command.Chars)).
		WithStdinLines("日本語").
		Run()

	assertion.NoError(t, resultBytes.Err)
	assertion.NoError(t, resultChars.Err)

	// Unicode: more bytes than characters
	bytesOut := strings.TrimSpace(resultBytes.Stdout[0])
	charsOut := strings.TrimSpace(resultChars.Stdout[0])

	assertion.Equal(t, bytesOut, "10", "10 bytes")
	assertion.Equal(t, charsOut, "3", "3 chars")
}

// ==============================================================================
// Test Output Formatting
// ==============================================================================

func TestWc_OutputFormat_Default(t *testing.T) {
	result := run.Command(command.Wc()).
		WithStdinLines("hello").
		Run()

	assertion.NoError(t, result.Err)
	// Should have 3 numbers in output (lines, words, bytes)
	output := strings.Fields(result.Stdout[0])
	assertion.Equal(t, len(output), 3, "three fields")
}

func TestWc_OutputFormat_SingleFlag(t *testing.T) {
	result := run.Command(command.Wc(command.Lines)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Should have just 1 number
	output := strings.TrimSpace(result.Stdout[0])
	// Should be just a number, no extra spaces
	_, err := fmt.Sscanf(output, "%d", new(int))
	assertion.NoError(t, err)
}

