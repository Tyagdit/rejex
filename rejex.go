package rejex

import (
	"fmt"
    "strings"
    "regexp"
)

const (
    ANCHOR = "ANCHOR"
    QUANTIFIER = "QUANTIFIER"
    CHARACTERS = "CHARACTERS"
    META = "META"
)

type RejexError struct {
    Position int
    Err string
}

func (e *RejexError) Error() string {
    return fmt.Sprintf("Error while building regex at position %d: %s", e.Position, e.Err)
}

type RejexBuilder struct {
    strings.Builder

    flags map[RejexFlag]bool

    negateNext bool
    lastSegmentType string

    groupActive bool
    groupContent []string
    groupNestingLevel int

    selectionActive bool
    selectionContent string

    ignoreErrors bool
    Errors []RejexError

    // bufferedQuantifier string
}

func createRejexBuilder(ignoreErrors []bool) *RejexBuilder {
    var ie bool
    if len(ignoreErrors) > 0 {
        ie = ignoreErrors[0]
    } else {
        ie = false
    }
    return &RejexBuilder{
        flags: map[RejexFlag]bool{
            'i': false, // Case Insensitive
            'm': false, // Multiline
            's': false, // Single Line
            'U': false, // Ungreedy
        },
        groupContent: make([]string, 2),
        ignoreErrors: ie,
    }
}

func NewRejex(ignoreErrors ...bool) *RejexBuilder {
    return createRejexBuilder(ignoreErrors)
}

func NewRejexFromString(s string, ignoreErrors ...bool) *RejexBuilder {
    r := createRejexBuilder(ignoreErrors)
    r.WriteString(s)
    return r
}

func (r *RejexBuilder) Build() (string, []RejexError) {
    flagStr := "(?"
    for f, b := range r.flags {
        if b { flagStr += string(f) }
    }
    flagStr += ")"

    r.negateNext = false

    if r.selectionActive {
        r.addError("Building without closing selection set")
    }
    if r.groupActive {
        r.addError("Building without closing group")
    }

    if !r.ignoreErrors {
        for _, err := range r.Errors {
            fmt.Println(err.Error())
        }
    }

    if flagStr == "(?)" {
        return r.String(), r.Errors
    } else {
        return flagStr + r.String(), r.Errors
    }
}

func (r *RejexBuilder) appendSegment(segmentType string, segment string, alt ...string) *RejexBuilder {
    var toWrite string
    if !r.negateNext || len(alt) == 0 {
        toWrite = segment
    } else {
        toWrite = alt[0]
    }

    if r.selectionActive {
        r.selectionContent += toWrite
    } else if r.groupActive {
        r.groupContent[r.groupNestingLevel] += toWrite
    } else {
        r.WriteString(toWrite)
    }

    r.negateNext = false
    r.lastSegmentType = segmentType
    return r
}

func (r *RejexBuilder) addError(err string) {
    r.Errors = append(r.Errors,
        RejexError{
            r.Len(),
            err,
        },
    )
}

// General

func (r *RejexBuilder) Not() *RejexBuilder {
    if r.selectionActive {
        r.addError(
            "Negation cannot be used in a selection set, use `BeginNonSelectionSet()` instead",
        )
    } else {
        r.negateNext = !r.negateNext
    }
    return r
}

func (r *RejexBuilder) Characters(s string) *RejexBuilder {
    return r.appendSegment(CHARACTERS, s)
}

func (r *RejexBuilder) EscapedCharacters(s string) *RejexBuilder {
    segment := regexp.QuoteMeta(s)
    return r.appendSegment(CHARACTERS, segment)
}

func (r *RejexBuilder) AnyChar() *RejexBuilder {
    return r.appendSegment(CHARACTERS, ".")
}

func (r *RejexBuilder) Literally(s string) *RejexBuilder {
    segment := fmt.Sprintf("\\Q%s\\E", s)
    return r.appendSegment(CHARACTERS, segment)
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
        return fmt.Sprintf("(?:%s)%s", s, q)
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
        r.appendSegment(META, "?")
    } else {
        r.addError(
            "'PreferFewer()' should only be used after a quantifier",
        )
    }
    return r
}

func (r *RejexBuilder) Or() *RejexBuilder {
    return r.appendSegment(META, "|")
}

func (r *RejexBuilder) EitherOr(s ...string) *RejexBuilder {
    var segment string
    if len(s) > 1 {
        segment = fmt.Sprintf("(?:%s)", strings.Join(s, "|"))
        r.appendSegment(CHARACTERS, segment)
    } else {
        r.addError(
            "Not enough options specified in 'EitherOr()'",
        )
    }
    return r
}

// Grouping

func (r *RejexBuilder) startNewGroup(s string) *RejexBuilder {
    if !r.selectionActive {
        r.groupActive = true
        r.groupNestingLevel++
        r.groupContent[r.groupNestingLevel] = s
    } else {
        r.addError(
            "Group constructs do not work inside a selection set",
        )
    }
    return r
}

func (r *RejexBuilder) BeginCaptureGroup() *RejexBuilder {
    return r.startNewGroup("(")
}

func (r *RejexBuilder) BeginNamedCaptureGroup(name string) *RejexBuilder {
    segment := fmt.Sprintf("(?P<%s>", name)
    return r.startNewGroup(segment)
}

func (r *RejexBuilder) BeginNonCaptureGroup() *RejexBuilder {
    return r.startNewGroup("(?:")
}

func (r *RejexBuilder) BeginGroupWithFlags(f RejexFlag) *RejexBuilder {
    segment := fmt.Sprintf("(?%s:", string(f))
    return r.startNewGroup(segment)
}

func (r *RejexBuilder) EndGroup() *RejexBuilder {
    if r.groupActive {
        segment := r.groupContent[r.groupNestingLevel] + ")"
        r.groupContent = r.groupContent[:r.groupNestingLevel]
        r.groupNestingLevel--
        if r.groupNestingLevel == 0 {
            r.groupActive = false
        }
        r.appendSegment(CHARACTERS, segment)
    } else {
        r.addError(
            "Cannot end group, no group open",
        )
    }
    return r
}

func (r *RejexBuilder) BeginSelectionSet() *RejexBuilder {
    if !r.selectionActive {
        r.selectionActive = true
        r.selectionContent = "["
    } else {
        r.addError("Cannot nest selection sets")
    }
    return r
}

func (r *RejexBuilder) BeginNonSelectionSet() *RejexBuilder {
    if !r.selectionActive {
        r.selectionActive = true
        r.selectionContent = "[^"
    } else {
        r.addError("Cannot nest selection sets")
    }
    return r
}

func (r *RejexBuilder) EndSelectionSet() *RejexBuilder {
    if r.selectionActive {
        segment := r.selectionContent + "]"
        r.selectionContent = ""
        r.selectionActive = false
        r.appendSegment(CHARACTERS, segment)
    } else {
        r.addError(
            "Cannot end selection set, no set open",
        )
    }
    return r
}
