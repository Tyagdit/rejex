package rejex

import (
    "fmt"
    "strconv"
)

func (r *RejexBuilder) checkForSelection(segment string) *RejexBuilder {
    if !r.selectionActive {
        segment = fmt.Sprintf("[%s]", segment)
        unsegment := fmt.Sprintf("[^%s]", segment)
        r.appendSegment(characters, segment, unsegment)
    } else {
        r.appendSegment(characters, segment)
    }
    return r
}

// AnyFrom matches any single character from the provided input
func (r *RejexBuilder) AnyFrom(s string) *RejexBuilder {
    return r.checkForSelection(s)
}

// AnyFromCharRange matches any single character in the range between the 2 characters provided
func (r *RejexBuilder) AnyFromCharRange(from, to string) *RejexBuilder {
    segment := fmt.Sprintf("%s-%s", from, to)
    return r.checkForSelection(segment)
}

// AnyWhitespace matches any single whitespace character
func (r *RejexBuilder) AnyWhitespace() *RejexBuilder {
    return r.appendSegment(characters, "\\s", "\\S")
}

// AnyWordChar matches any single word character
func (r *RejexBuilder) AnyWordChar() *RejexBuilder {
    return r.appendSegment(characters, "\\w", "\\W")
}

// AnyDigit matches any single decimal digit
func (r *RejexBuilder) AnyDigit() *RejexBuilder {
    return r.appendSegment(characters, "\\d", "\\D")
}

// AnyLetter matches any single english letter
func (r *RejexBuilder) AnyLetter() *RejexBuilder {
    return r.checkForSelection("a-zA-Z")
}

// AnyUppercase matches any single uppercase english letter
func (r *RejexBuilder) AnyUppercase() *RejexBuilder {
    return r.checkForSelection("A-Z")
}

// AnyLowercase matches any single lowercase english letter
func (r *RejexBuilder) AnyLowercase() *RejexBuilder {
    return r.checkForSelection("a-z")
}

// AnyAlNumChar matches any single english letter or digit
func (r *RejexBuilder) AnyAlNumChar() *RejexBuilder {
    return r.checkForSelection("0-9a-zA-Z")
}

// AnyPunctuation matches any single Punctuation character
func (r *RejexBuilder) AnyPunctuation() *RejexBuilder {
    return r.checkForSelection("!-/:-@[-`{-~")
}

// AnyGraphicChar matches any visible character
func (r *RejexBuilder) AnyGraphicChar() *RejexBuilder {
    return r.checkForSelection("!-~")
}

// AnyASCIIChar matches any single ASCII character
func (r *RejexBuilder) AnyASCIIChar() *RejexBuilder {
    return r.checkForSelection("\x00-\x7F")
}

// AnyControlChar matches any sigle control character
func (r *RejexBuilder) AnyControlChar() *RejexBuilder {
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

    return r.appendSegment(characters, segment, unsegment)
}

// AnyUnicodeGrapheme matches a single Unicode grapheme, whether encoded as a
// single code point or multiple code points using combining marks. A grapheme
// most closely resembles the everyday concept of a “character”
func (r *RejexBuilder) AnyUnicodeGrapheme() *RejexBuilder {
    return r.appendSegment(characters, "\\X")
}

// AnyUnicodeLetter matches any single unicode letter
func (r *RejexBuilder) AnyUnicodeLetter() *RejexBuilder {
    return r.appendSegment(characters, "\\pL", "\\PL")
}

// AnyUnicodeUppercase matches any single uppercase unicode character
func (r *RejexBuilder) AnyUnicodeUppercase() *RejexBuilder {
    return r.appendSegment(characters, "\\p{Lu}", "\\P{Lu}")
}

// AnyUnicodeLowercase matches any single lowercase unicode character
func (r *RejexBuilder) AnyUnicodeLowercase() *RejexBuilder {
    return r.appendSegment(characters, "\\p{Ll}", "\\P{Ll}")
}

// AnyUnicodeWhitespace matches any single unicode whitespace
func (r *RejexBuilder) AnyUnicodeWhitespace() *RejexBuilder {
    return r.appendSegment(characters, "\\pZ", "\\PZ")
}

// AnyUnicodeSymbol matches any single unicode symbol character
func (r *RejexBuilder) AnyUnicodeSymbol() *RejexBuilder {
    return r.appendSegment(characters, "\\pS", "\\PS")
}

// AnyUnicodeNumber matches any single unicode number
func (r *RejexBuilder) AnyUnicodeNumber() *RejexBuilder {
    return r.appendSegment(characters, "\\pN", "\\PN")
}

// AnyUnicodePunctuation matches any single unicode punctuation character
func (r *RejexBuilder) AnyUnicodePunctuation() *RejexBuilder {
    return r.appendSegment(characters, "\\pP", "\\PP")
}

// Non Negate-able classes

// OctalChar matches the character represented by the provided octal character code
func (r *RejexBuilder) OctalChar(c int) *RejexBuilder {
    if c >= 0 && c < 778 {
        segment := fmt.Sprintf("\\%03d", c)
        r.appendSegment(characters, segment)
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
        r.appendSegment(characters, segment)
    } else {
        segment = fmt.Sprintf("\\x{%s}", s)
        r.appendSegment(characters, segment)
    }

    return r
}

// ControlChar matches the control character represented by the provided control
// character code
func (r *RejexBuilder) ControlChar(s string) *RejexBuilder {
    if len(s) == 1 {
        segment := fmt.Sprintf("\\c%s", s)
        r.appendSegment(characters, segment)
    } else {
        r.addError("Invalid control character code")
    }
    return r
}

