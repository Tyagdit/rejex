package rejex

type RejexFlavor string

const (
    GoFlavor RejexFlavor = "GO"
    ECMAFlavor RejexFlavor = "ECMA"
    PerlFlavor RejexFlavor = "PERL"
)

var goFlavorFlags = map[RejexFlag]bool{
    'i': false, // Case Insensitive
    'm': false, // Multiline
    's': false, // Single Line
    'U': false, // Ungreedy
}

// GoFlavorInterface represents regex of the Go standard syntax
type GoFlavorInterface interface {
    Build() (string, []RejexError)

    // General
    Not() *RejexBuilder
    Characters(string) *RejexBuilder
    EscapedCharacters(string) *RejexBuilder
    AnyChar() *RejexBuilder
    Literally(string) *RejexBuilder

    // Anchors
    Starting() *RejexBuilder
    AbsoluteStarting() *RejexBuilder
    Ending() *RejexBuilder
    AbsoluteEnding() *RejexBuilder
    WordBoundary() *RejexBuilder

    // Quantifiers
    ZeroOrOneOf(string) *RejexBuilder
    ZeroOrMoreOf(string) *RejexBuilder
    OneOrMoreOf(string) *RejexBuilder
    NOf(string, int) *RejexBuilder
    NOrMoreOf(string, int) *RejexBuilder
    NToMOf(string, int, int) *RejexBuilder

    // Meta
    PreferFewer() *RejexBuilder
    Or() *RejexBuilder
    EitherOr(...string) *RejexBuilder

    // Group Constructs
    BeginCaptureGroup() *RejexBuilder
    BeginNamedCaptureGroup(string) *RejexBuilder
    BeginNonCaptureGroup() *RejexBuilder
    BeginGroupWithFlags([]RejexFlag) *RejexBuilder
    EndGroup() *RejexBuilder
    BeginSelectionSet() *RejexBuilder
    BeginNonSelectionSet() *RejexBuilder
    EndSelectionSet() *RejexBuilder

    // Char Classes
    AnyFrom(string) *RejexBuilder
    AnyFromCharRange(string, string) *RejexBuilder
    AnyWhitespace() *RejexBuilder
    AnyWordChar() *RejexBuilder
    AnyDigit() *RejexBuilder
    AnyLetter() *RejexBuilder
    AnyUppercase() *RejexBuilder
    AnyLowercase() *RejexBuilder
    AnyAlNumChar() *RejexBuilder
    AnyPunctuation() *RejexBuilder
    AnyGraphicChar() *RejexBuilder
    AnyASCIIChar() *RejexBuilder
    AnyControlChar() *RejexBuilder
    AnyUnicodeLetter() *RejexBuilder
    AnyUnicodeUppercase() *RejexBuilder
    AnyUnicodeLowercase() *RejexBuilder
    AnyUnicodeWhitespace() *RejexBuilder
    AnyUnicodeSymbol() *RejexBuilder
    AnyUnicodeNumber() *RejexBuilder
    AnyUnicodePunctuation() *RejexBuilder
    UnicodeClass(string) *RejexBuilder
    OctalChar(int) *RejexBuilder
    HexChar(string) *RejexBuilder

    // Flags
    AddFlags(...RejexFlag) *RejexBuilder
    RemoveFlags(...RejexFlag) *RejexBuilder

    // Utils
    LineEnding() *RejexBuilder
}

var ecmaFlavorFlags = map[RejexFlag]bool{
    'g': false, // Global
    'i': false, // Case Insensitive
    'm': false, // Multiline
    'y': false, // Sticky
    'u': false, // Unicode
}

// ECMAFlavorInterface represents regex of the ECMAScript standard syntax
type ECMAFlavorInterface interface {
    Build() (string, []RejexError)

    // General
    Not() *RejexBuilder
    Characters(string) *RejexBuilder
    EscapedCharacters(string) *RejexBuilder
    AnyChar() *RejexBuilder

    // Anchors
    Starting() *RejexBuilder
    Ending() *RejexBuilder
    WordBoundary() *RejexBuilder

    // Quantifiers
    ZeroOrOneOf(string) *RejexBuilder
    ZeroOrMoreOf(string) *RejexBuilder
    OneOrMoreOf(string) *RejexBuilder
    NOf(string, int) *RejexBuilder
    NOrMoreOf(string, int) *RejexBuilder
    NToMOf(string, int, int) *RejexBuilder

    // Meta
    PreferFewer() *RejexBuilder
    Or() *RejexBuilder
    EitherOr(...string) *RejexBuilder
    CapturedPatternByNum(int) *RejexBuilder
    CapturedPatternByName(string) *RejexBuilder

    // Group Constructs
    BeginCaptureGroup() *RejexBuilder
    BeginNamedCaptureGroup(string) *RejexBuilder
    BeginNonCaptureGroup() *RejexBuilder
    BeginPosLookahead() *RejexBuilder
    BeginNegLookahead() *RejexBuilder
    BeginPosLookbehind() *RejexBuilder
    BeginNegLookbehind() *RejexBuilder
    EndGroup() *RejexBuilder
    BeginSelectionSet() *RejexBuilder
    BeginNonSelectionSet() *RejexBuilder
    EndSelectionSet() *RejexBuilder

    // Char Classes
    AnyFrom(string) *RejexBuilder
    AnyFromCharRange(string, string) *RejexBuilder
    AnyWhitespace() *RejexBuilder
    AnyWordChar() *RejexBuilder
    AnyDigit() *RejexBuilder
    AnyLetter() *RejexBuilder
    AnyUppercase() *RejexBuilder
    AnyLowercase() *RejexBuilder
    AnyAlNumChar() *RejexBuilder
    AnyPunctuation() *RejexBuilder
    AnyGraphicChar() *RejexBuilder
    AnyASCIIChar() *RejexBuilder
    AnyControlChar() *RejexBuilder
    AnyUnicodeLetter() *RejexBuilder
    AnyUnicodeUppercase() *RejexBuilder
    AnyUnicodeLowercase() *RejexBuilder
    AnyUnicodeWhitespace() *RejexBuilder
    AnyUnicodeSymbol() *RejexBuilder
    AnyUnicodeNumber() *RejexBuilder
    AnyUnicodePunctuation() *RejexBuilder
    UnicodeClass(string) *RejexBuilder
    OctalChar(int) *RejexBuilder
    HexChar(string) *RejexBuilder
    ControlChar(string) *RejexBuilder

    // Flags
    AddFlags(...RejexFlag) *RejexBuilder
    RemoveFlags(...RejexFlag) *RejexBuilder

    // Utils
    LineEnding() *RejexBuilder
}

var perlFlavorFlags = map[RejexFlag]bool{
    'g': false, // Global
    'i': false, // Case Insensitive
    'm': false, // Multiline
    's': false, // Single Line
}

// PerlFlavorInterface represents regex of the Perl standard syntax
type PerlFlavorInterface interface {
    // There's a bunch of features missing here like recursion, conditionals,
    // subroutines etc. but like who tf needs these?? and why??? why is there
    // programming logic in regular expressions????

    Build() (string, []RejexError)

    // General
    Not() *RejexBuilder
    Characters(string) *RejexBuilder
    EscapedCharacters(string) *RejexBuilder
    AnyChar() *RejexBuilder

    // Anchors
    Starting() *RejexBuilder
    AbsoluteStarting() *RejexBuilder
    Ending() *RejexBuilder
    AbsoluteEnding() *RejexBuilder
    WordBoundary() *RejexBuilder
    EndOfLastMatch() *RejexBuilder
    // AbsoluteEndingWithNewline() *RejexBuilder // "\Z"

    // Quantifiers
    ZeroOrOneOf(string) *RejexBuilder
    ZeroOrMoreOf(string) *RejexBuilder
    OneOrMoreOf(string) *RejexBuilder
    NOf(string, int) *RejexBuilder
    NOrMoreOf(string, int) *RejexBuilder
    NToMOf(string, int, int) *RejexBuilder

    // Meta
    PreferFewer() *RejexBuilder
    PossessiveQuantifier() *RejexBuilder
    Or() *RejexBuilder
    EitherOr(...string) *RejexBuilder
    CapturedPatternByNum(int) *RejexBuilder
    CapturedPatternByName(string) *RejexBuilder

    // Group Constructs
    BeginCaptureGroup() *RejexBuilder
    BeginNamedCaptureGroup(string) *RejexBuilder
    BeginNonCaptureGroup() *RejexBuilder
    BeginGroupWithFlags([]RejexFlag) *RejexBuilder
    BeginPosLookahead() *RejexBuilder
    BeginNegLookahead() *RejexBuilder
    BeginPosLookbehind() *RejexBuilder
    BeginNegLookbehind() *RejexBuilder
    BeginAtomicGroup() *RejexBuilder
    BeginBranchResetGroup() *RejexBuilder
    EndGroup() *RejexBuilder
    BeginSelectionSet() *RejexBuilder
    BeginNonSelectionSet() *RejexBuilder
    EndSelectionSet() *RejexBuilder

    // Char Classes
    AnyFrom(string) *RejexBuilder
    AnyFromCharRange(string, string) *RejexBuilder
    AnyWhitespace() *RejexBuilder
    AnyWordChar() *RejexBuilder
    AnyDigit() *RejexBuilder
    AnyLetter() *RejexBuilder
    AnyUppercase() *RejexBuilder
    AnyLowercase() *RejexBuilder
    AnyAlNumChar() *RejexBuilder
    AnyPunctuation() *RejexBuilder
    AnyGraphicChar() *RejexBuilder
    AnyASCIIChar() *RejexBuilder
    AnyControlChar() *RejexBuilder
    AnyUnicodeGrapheme() *RejexBuilder
    AnyUnicodeLetter() *RejexBuilder
    AnyUnicodeUppercase() *RejexBuilder
    AnyUnicodeLowercase() *RejexBuilder
    AnyUnicodeWhitespace() *RejexBuilder
    AnyUnicodeSymbol() *RejexBuilder
    AnyUnicodeNumber() *RejexBuilder
    AnyUnicodePunctuation() *RejexBuilder
    UnicodeClass(string) *RejexBuilder
    OctalChar(int) *RejexBuilder
    HexChar(string) *RejexBuilder

    // Flags
    AddFlags(...RejexFlag) *RejexBuilder
    RemoveFlags(...RejexFlag) *RejexBuilder

    // Utils
    LineEnding() *RejexBuilder
}

// EgrepFlavorInterface represents regex of ERE syntax used by GNU egrep (or with grep -E)
// type EgrepFlavorInterface interface {}
// EgrepFlavorInterface represents regex of POSIX ERE syntax used by grep
// type EgrepPOSIXInterface interface {}
// type VimFlavorInterface interface {}
