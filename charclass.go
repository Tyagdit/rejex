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

func (r *RejexBuilder) AnyFrom(s string) *RejexBuilder {
    return r.checkForSelection(s)
}

func (r *RejexBuilder) CharRange(from, to string) *RejexBuilder {
    segment := fmt.Sprintf("%s-%s", from, to)
    return r.checkForSelection(segment)
}

func (r *RejexBuilder) Whitespace() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\s", "\\S")
}

func (r *RejexBuilder) WordChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\w", "\\W")
}

func (r *RejexBuilder) Digit() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\d", "\\D")
}

func (r *RejexBuilder) Letter() *RejexBuilder {
    return r.checkForSelection("a-zA-Z")
}

func (r *RejexBuilder) Uppercase() *RejexBuilder {
    return r.checkForSelection("A-Z")
}

func (r *RejexBuilder) Lowercase() *RejexBuilder {
    return r.checkForSelection("a-z")
}

func (r *RejexBuilder) AlNumChar() *RejexBuilder {
    return r.checkForSelection("0-9a-zA-Z")
}

func (r *RejexBuilder) Punctuation() *RejexBuilder {
    return r.checkForSelection("!-/:-@[-`{-~")
}

func (r *RejexBuilder) GraphicChar() *RejexBuilder {
    return r.checkForSelection("!-~")
}

func (r *RejexBuilder) ASCIIChar() *RejexBuilder {
    return r.checkForSelection("\x00-\x7F")
}

func (r *RejexBuilder) ControlChar() *RejexBuilder {
    return r.checkForSelection("\x00-\x1F\x7F")
}

// Unicode Classes

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

func (r *RejexBuilder) LetterUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pL", "\\PL")
}

func (r *RejexBuilder) UppercaseUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\p{Lu}", "\\P{Lu}")
}

func (r *RejexBuilder) LowercaseUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\p{Ll}", "\\P{Ll}")
}

func (r *RejexBuilder) WhitespaceUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pZ", "\\PZ")
}

func (r *RejexBuilder) SymbolUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pS", "\\PS")
}

func (r *RejexBuilder) NumberUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pN", "\\PN")
}

func (r *RejexBuilder) PunctuationUnicode() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\pP", "\\PP")
}

// Non Negate-able classes

func (r *RejexBuilder) OctalChar(c int) *RejexBuilder {
    if c >= 0 && c < 778 {
        segment := fmt.Sprintf("\\%03d", c)
        r.appendSegment(CHARACTERS, segment)
    } else {
        r.addError("Invalid octal character code")
    }
    return r
}

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

