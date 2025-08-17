package command

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

type flags struct {
	Lines     LinesFlag
	Words     WordsFlag
	Chars     CharsFlag
	Bytes     BytesFlag
	MaxLength MaxLengthFlag
}

func (f LinesFlag) Configure(flags *flags)     { flags.Lines = f }
func (f WordsFlag) Configure(flags *flags)     { flags.Words = f }
func (f CharsFlag) Configure(flags *flags)     { flags.Chars = f }
func (f BytesFlag) Configure(flags *flags)     { flags.Bytes = f }
func (f MaxLengthFlag) Configure(flags *flags) { flags.MaxLength = f }
