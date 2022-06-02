package rejex

type RejexFlavor string

const (
    GoFlavor RejexFlavor = "go"
    ECMAFlavor RejexFlavor = "ecma"
)

var goFlavorFlags = map[RejexFlag]bool{
    'i': false, // Case Insensitive
    'm': false, // Multiline
    's': false, // Single Line
    'U': false, // Ungreedy
}

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
    BeginGroupWithFlags(RejexFlag) *RejexBuilder 
    EndGroup() *RejexBuilder 
    BeginSelectionSet() *RejexBuilder 
    BeginNonSelectionSet() *RejexBuilder 
    EndSelectionSet() *RejexBuilder 

    // Char Classes
    AnyFrom(string) *RejexBuilder 
    CharRange(string, string) *RejexBuilder 
    Whitespace() *RejexBuilder 
    WordChar() *RejexBuilder 
    Digit() *RejexBuilder 
    Letter() *RejexBuilder 
    Uppercase() *RejexBuilder 
    Lowercase() *RejexBuilder 
    AlNumChar() *RejexBuilder 
    Punctuation() *RejexBuilder 
    GraphicChar() *RejexBuilder 
    ASCIIChar() *RejexBuilder 
    ControlChar() *RejexBuilder 
    UnicodeClass(string) *RejexBuilder 
    LetterUnicode() *RejexBuilder 
    UppercaseUnicode() *RejexBuilder 
    LowercaseUnicode() *RejexBuilder 
    WhitespaceUnicode() *RejexBuilder 
    SymbolUnicode() *RejexBuilder 
    NumberUnicode() *RejexBuilder 
    PunctuationUnicode() *RejexBuilder 
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
    CapturedPattern(string) *RejexBuilder

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
    CharRange(string, string) *RejexBuilder 
    Whitespace() *RejexBuilder 
    WordChar() *RejexBuilder 
    Digit() *RejexBuilder 
    Letter() *RejexBuilder 
    Uppercase() *RejexBuilder 
    Lowercase() *RejexBuilder 
    AlNumChar() *RejexBuilder 
    Punctuation() *RejexBuilder 
    GraphicChar() *RejexBuilder 
    ASCIIChar() *RejexBuilder 
    // ControlChar() *RejexBuilder 
    UnicodeClass(string) *RejexBuilder 
    LetterUnicode() *RejexBuilder 
    UppercaseUnicode() *RejexBuilder 
    LowercaseUnicode() *RejexBuilder 
    WhitespaceUnicode() *RejexBuilder 
    SymbolUnicode() *RejexBuilder 
    NumberUnicode() *RejexBuilder 
    PunctuationUnicode() *RejexBuilder 
    OctalChar(int) *RejexBuilder 
    HexChar(string) *RejexBuilder 

    // Flags
    AddFlags(...RejexFlag) *RejexBuilder 
    RemoveFlags(...RejexFlag) *RejexBuilder

    // Utils
    LineEnding() *RejexBuilder 
}

type VimFlavorInterface interface {}
type GrepFlavorInterface interface {}
