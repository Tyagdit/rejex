package rejex

import (
    "fmt"
)

func (r *RejexBuilder) AnyFrom(s string) *RejexBuilder {
    anyFrom := fmt.Sprintf("[%s]", s)
    noneFrom := fmt.Sprintf("[^%s]", s)

    return r.appendSegment(CHARACTERS, anyFrom, noneFrom)
}

func (r *RejexBuilder) CharRange(from, to string) *RejexBuilder {
    charRange := fmt.Sprintf("[%s-%s]", from, to)
    noCharRange := fmt.Sprintf("[^%s-%s]", from, to)

    return r.appendSegment(CHARACTERS, charRange, noCharRange)
}

func (r *RejexBuilder) Literally(s string) *RejexBuilder {
    segment := fmt.Sprintf("\\Q%s\\E", s)
    return r.appendSegment(CHARACTERS, segment)
}

func (r *RejexBuilder) AnyChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, ".")
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
    return r.appendSegment(CHARACTERS, "[a-zA-Z]", "[^a-zA-Z]")
}

func (r *RejexBuilder) AlNumChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "[0-9a-zA-Z]", "[^0-9a-zA-Z]")
}

func (r *RejexBuilder) Punctuation() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "[!-/:-@[-`{-~]", "[^!-/:-@[-`{-~]")
}

func (r *RejexBuilder) GraphicChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "[!-~]", "[^!-~]")
}

func (r *RejexBuilder) ASCIIChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "[\x00-\x7F]", "[^\x00-\x7F]")
}

func (r *RejexBuilder) OctalChar(s string) *RejexBuilder {
    segment := fmt.Sprintf("\\%s", s)
    return r.appendSegment(CHARACTERS, segment)
}

func (r *RejexBuilder) HexChar(s string) *RejexBuilder {
    var segment string
    if len(s) == 2 {
        segment = fmt.Sprintf("\\x%s", s)
    } else {
        segment = fmt.Sprintf("\\x{%s}", s)
    }

    return r.appendSegment(CHARACTERS, segment)
}

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
