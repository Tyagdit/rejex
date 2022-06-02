package rejex

import (
    "fmt"
    "strconv"
)

func (r *RejexBuilder) checkForSelection(segment string) *RejexBuilder {
    if !r.selectionActive {
        segment = fmt.Sprintf("[%s]", segment)
        unsegment := fmt.Sprintf("[^%s]", segment)
        r.appendSegment(CHARACTERS, segment, unsegment)
    } else {
        r.appendSegment(CHARACTERS, segment)
    }
    return r
}

// AnyFrom matches any single character from the provided input
func (r *RejexBuilder) AnyFrom(s string) *RejexBuilder {
    return r.checkForSelection(s)
}

// CharRange matches any single character in the range between the 2 characters provided
func (r *RejexBuilder) CharRange(from, to string) *RejexBuilder {
    segment := fmt.Sprintf("%s-%s", from, to)
    return r.checkForSelection(segment)
}

// Whitespace matches any single whitespace character
func (r *RejexBuilder) Whitespace() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\s", "\\S")
}

// WordChar matches any single word character
func (r *RejexBuilder) WordChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\w", "\\W")
}

// Digit matches any single decimal digit
func (r *RejexBuilder) Digit() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\d", "\\D")
}

// Letter matches any single english letter
func (r *RejexBuilder) Letter() *RejexBuilder {
    return r.checkForSelection("a-zA-Z")
}

// Uppercase matches any single uppercase english letter
func (r *RejexBuilder) Uppercase() *RejexBuilder {
    return r.checkForSelection("A-Z")
}

// Lowercase matches any single lowercase english letter
func (r *RejexBuilder) Lowercase() *RejexBuilder {
    return r.checkForSelection("a-z")
}

// AlNumChar matches any single english letter or digit
func (r *RejexBuilder) AlNumChar() *RejexBuilder {
    return r.checkForSelection("0-9a-zA-Z")
}

// Punctuation matches any single Punctuation character
func (r *RejexBuilder) Punctuation() *RejexBuilder {
    return r.checkForSelection("!-/:-@[-`{-~")
}

// GraphicChar matches any visible character
func (r *RejexBuilder) GraphicChar() *RejexBuilder {
    return r.checkForSelection("!-~")
}

// ASCIIChar matches any single ASCII character
func (r *RejexBuilder) ASCIIChar() *RejexBuilder {
    return r.checkForSelection("\x00-\x7F")
}

// ControlChar matches any sigle control character
func (r *RejexBuilder) ControlChar() *RejexBuilder {
    return r.checkForSelection("\x00-\x1F\x7F")
}

// Unicode Classes

// UnicodeClass matches any character from the provided unicode class
func (r *RejexBuilder) UnicodeClass(s string) *RejexBuilder {
    var segment, unsegment string
    if len(s) == 1 {
        segment = fmt.Sprintf("\\p%s", s)
        unsegment = fmt.Sprintf("\\P%s", s)
    } else {
        segment = fmt.Sprintf("\\p{%s}", s)
        unsegment = fmt.Sprintf("\\P{%s}", s)
    }

    return r.appendSegment(CHARACTERS, segment, unsegment)
}

// LetterUnicode matches any single unicode letter
func (r *RejexBuilder) LetterUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pL", "\\PL")
}

// UppercaseUnicode matches any single uppercase unicode character
func (r *RejexBuilder) UppercaseUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\p{Lu}", "\\P{Lu}")
}

// LowercaseUnicode matches any single lowercase unicode character
func (r *RejexBuilder) LowercaseUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\p{Ll}", "\\P{Ll}")
}

// WhitespaceUnicode matches any single unicode whitespace
func (r *RejexBuilder) WhitespaceUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pZ", "\\PZ")
}

// SymbolUnicode matches any single unicode symbol character
func (r *RejexBuilder) SymbolUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pS", "\\PS")
}

// NumberUnicode matches any single unicode number
func (r *RejexBuilder) NumberUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pN", "\\PN")
}

// PunctuationUnicode matches any single unicode PunctuationUnicode character
func (r *RejexBuilder) PunctuationUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pP", "\\PP")
}

// Non Negate-able classes

// OctalChar matches the character represented by the provided octal character code
func (r *RejexBuilder) OctalChar(c int) *RejexBuilder {
    if c >= 0 && c < 778 {
        segment := fmt.Sprintf("\\%03d", c)
        r.appendSegment(CHARACTERS, segment)
    } else {
        r.addError("Invalid octal character code")
    }
    return r
}

// HexChar matches the character represented by the provided hex character code
func (r *RejexBuilder) HexChar(s string) *RejexBuilder {
    var segment string
    if c, e := strconv.ParseInt(s, 16, 64); e != nil || c < 0 || c > 1114111 {
        r.addError("Invalid hex character code")
    } else if len(s) == 2 {
        segment = fmt.Sprintf("\\x%s", s)
        r.appendSegment(CHARACTERS, segment)
    } else {
        segment = fmt.Sprintf("\\x{%s}", s)
        r.appendSegment(CHARACTERS, segment)
    }

    return r
}

