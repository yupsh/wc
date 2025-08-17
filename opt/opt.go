package opt

// Boolean flag types with constants
type LinesFlag bool
const (
	Lines   LinesFlag = true
	NoLines LinesFlag = false
)

type WordsFlag bool
const (
	Words   WordsFlag = true
	NoWords WordsFlag = false
)

type CharsFlag bool
const (
	Chars   CharsFlag = true
	NoChars CharsFlag = false
)

type BytesFlag bool
const (
	Bytes   BytesFlag = true
	NoBytes BytesFlag = false
)

type MaxLengthFlag bool
const (
	MaxLength   MaxLengthFlag = true
	NoMaxLength MaxLengthFlag = false
)

// Flags represents the configuration options for the wc command
type Flags struct {
	Lines     LinesFlag     // Count lines only
	Words     WordsFlag     // Count words only
	Chars     CharsFlag     // Count characters only
	Bytes     BytesFlag     // Count bytes only
	MaxLength MaxLengthFlag // Show maximum line length
}

// Flag configuration methods
func (f LinesFlag) Configure(flags *Flags) { flags.Lines = f }
func (f WordsFlag) Configure(flags *Flags) { flags.Words = f }
func (f CharsFlag) Configure(flags *Flags) { flags.Chars = f }
func (f BytesFlag) Configure(flags *Flags) { flags.Bytes = f }
func (f MaxLengthFlag) Configure(flags *Flags) { flags.MaxLength = f }
