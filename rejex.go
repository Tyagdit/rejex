package rejex

import (
	"fmt"
    "strings"
)

const (
    ANCHOR = "ANCHOR"
    QUANTIFIER = "QUANTIFIER"
    CHARACTERS = "CHARACTERS"
    META = "META"
)

type RejexBuilder struct {
    strings.Builder
    flags map[RejexFlag]bool
    negateNext bool
    lastSegmentType string
    groupActive bool
    groupContent string
    savedQuantifier string
}

func NewRejex() *RejexBuilder {
    return &RejexBuilder{flags: map[RejexFlag]bool{
        'i': false, // Case Insensitive
        'm': false, // Multiline
        's': false, // Single Line
        'U': false, // Ungreedy
    }}
}

func NewRejexFromString(s string) *RejexBuilder {
    r := RejexBuilder{flags: map[RejexFlag]bool{
        'i': false, // Case Insensitive
        'm': false, // Multiline
        's': false, // Single Line
        'U': false, // Ungreedy
    }}

    r.WriteString(s)
    return &r
}

func (r *RejexBuilder) Build() string {
    flagStr := "(?"
    for f, b := range r.flags {
        if b { flagStr += string(f) }
    }
    flagStr += ")"

    r.negateNext = false

    if flagStr == "(?)" {
        return r.String()
    } else {
        return flagStr + r.String()
    }
}

func (r *RejexBuilder) appendSegment(segmentType string, segment string, alt ...string) *RejexBuilder {
    var toWrite string
    if !r.negateNext || len(alt) == 0 {
        toWrite = segment
    } else {
        toWrite = alt[0]
    }

    if r.groupActive {
        r.groupContent += segment
    } else {
        r.WriteString(toWrite)
    }

    r.negateNext = false
    r.lastSegmentType = segmentType
    return r
}

func (r *RejexBuilder) Not() *RejexBuilder {
    r.negateNext = !r.negateNext
    return r
}

func (r *RejexBuilder) Exactly(s string) *RejexBuilder {
    return r.appendSegment(CHARACTERS, s)
}

// Anchors

func (r *RejexBuilder) Starting() *RejexBuilder {
    return r.appendSegment(ANCHOR, "^")
}

func (r *RejexBuilder) AbsoluteStarting() *RejexBuilder {
    return r.appendSegment(ANCHOR, "\\A")
}

func (r *RejexBuilder) Ending() *RejexBuilder {
    return r.appendSegment(ANCHOR, "$")
}

func (r *RejexBuilder) AbsoluteEnding() *RejexBuilder {
    return r.appendSegment(ANCHOR, "\\z")
}

func (r *RejexBuilder) WordBoundary() *RejexBuilder {
    return r.appendSegment(ANCHOR, "\\b", "\\B")
}


// Quantifiers

func checkForGroup(s, q string) string {
    if len(s) > 1 {
        return fmt.Sprintf("(%s)%s", s, q)
    } else {
        return fmt.Sprintf("%s%s", s, q)
    }
}

func (r *RejexBuilder) ZeroOrOneOf(s string) *RejexBuilder {
    segment := checkForGroup(s, "?")
    return r.appendSegment(QUANTIFIER, segment)
}

func (r *RejexBuilder) ZeroOrMoreOf(s string) *RejexBuilder {
    segment := checkForGroup(s, "*")
    return r.appendSegment(QUANTIFIER, segment)
}

func (r *RejexBuilder) OneOrMoreOf(s string) *RejexBuilder {
    segment := checkForGroup(s, "+")
    return r.appendSegment(QUANTIFIER, segment)
}

func (r *RejexBuilder) NOf(s string, n int) *RejexBuilder {
    segment := checkForGroup(s, fmt.Sprintf("{%d}", n))
    return r.appendSegment(QUANTIFIER, segment)
}

func (r *RejexBuilder) NOrMoreOf(s string, n int) *RejexBuilder {
    segment := checkForGroup(s, fmt.Sprintf("{%d,}", n))
    return r.appendSegment(QUANTIFIER, segment)
}

func (r *RejexBuilder) NToMOf(s string, n, m int) *RejexBuilder {
    segment := checkForGroup(s, fmt.Sprintf("{%d,%d}", n, m))
    return r.appendSegment(QUANTIFIER, segment)
}

// Meta

func (r *RejexBuilder) PreferFewer() *RejexBuilder {
    if r.lastSegmentType == QUANTIFIER {
        return r.appendSegment(META, "?")
    } else {
        return r
    }
}

func (r *RejexBuilder) Or() *RejexBuilder {
    return r.appendSegment(META, "|")
}

func (r *RejexBuilder) EitherOr(s ...string) *RejexBuilder {
    var segment string
    if len(s) > 0 {
        segment = fmt.Sprintf("(?:%s)", strings.Join(s, "|"))
    } else {
        segment = ""
    }

    return r.appendSegment(CHARACTERS, segment)
}

// Grouping

func (r *RejexBuilder) StartCaptureGroup() *RejexBuilder {
    r.groupActive = true
    r.groupContent = "("
    return r
}

func (r *RejexBuilder) StartNamedCaptureGroup(name string) *RejexBuilder {
    r.groupActive = true
    r.groupContent = fmt.Sprintf("(?P<%s>", name)
    return r
}

func (r *RejexBuilder) StartNonCaptureGroup() *RejexBuilder {
    r.groupActive = true
    r.groupContent = "(?:"
    return r
}

func (r *RejexBuilder) StopGroup() *RejexBuilder {
    if r.groupActive {
        r.groupActive = false
        r.appendSegment(CHARACTERS, r.groupContent + ")")
    }
    return r
}
