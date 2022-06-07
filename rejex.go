package rejex

import (
	"fmt"
    "strings"
    "regexp"
)

const (
    anchor = "ANCHOR"
    quantifier = "QUANTIFIER"
    characters = "CHARACTERS"
    meta = "META"
)

// RejexError is an error reported while constructing a regex
type RejexError struct {
    Position int
    Err string
}

// Error returns a formatted error message
func (e *RejexError) Error() string {
    return fmt.Sprintf("Error while building regex at position %d: %s", e.Position, e.Err)
}

// RejexBuilder defines a regex string before being fully constructed
type RejexBuilder struct {
    strings.Builder

    flags map[RejexFlag]bool
    flavor RejexFlavor

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

func createRejexBuilder(flavor RejexFlavor, ignoreErrors []bool) *RejexBuilder {
    var ie bool
    if len(ignoreErrors) > 0 {
        ie = ignoreErrors[0]
    } else {
        ie = false
    }

    r := RejexBuilder{
        groupContent: make([]string, 2),
        ignoreErrors: ie,
    }

    r.flavor = flavor
    switch flavor {
    case GoFlavor:
        r.flags = goFlavorFlags
    case ECMAFlavor:
        r.flags = ecmaFlavorFlags
    case PerlFlavor:
        r.flags = perlFlavorFlags
    }

    return &r
}

// NewRejex creates a new RejexBuilder object used to construct a regex. This uses
// the Go flavored syntax.
func NewRejex(ignoreErrors ...bool) GoFlavorInterface {
    r := createRejexBuilder(GoFlavor, ignoreErrors)
    return GoFlavorInterface(r)
}

// NewRejexFromString creates a new RejexBuilder object used to construct a regex and
// populates it with a provided regex string, this string is not validated to be valid
// syntax. This uses the Go flavored syntax.
func NewRejexFromString(s string, ignoreErrors ...bool) GoFlavorInterface {
    r := createRejexBuilder(GoFlavor, ignoreErrors)
    r.WriteString(s)
    return GoFlavorInterface(r)
}

// NewECMARejex creates a new RejexBuilder object used to construct a regex. This uses
// the ECMAScript flavored syntax.
func NewECMARejex(ignoreErrors ...bool) ECMAFlavorInterface {
    r := createRejexBuilder(ECMAFlavor, ignoreErrors)
    return ECMAFlavorInterface(r)
}

// NewECMARejexFromString creates a new RejexBuilder object used to construct a regex and
// populates it with a provided regex string, this string is not validated to be valid
// syntax. This uses the ECMAScript flavored syntax.
func NewECMARejexFromString(s string, ignoreErrors ...bool) ECMAFlavorInterface {
    r := createRejexBuilder(ECMAFlavor, ignoreErrors)
    r.WriteString(s)
    return ECMAFlavorInterface(r)
}

// NewPerlRejex creates a new RejexBuilder object used to construct a regex. This uses
// the Perl flavored syntax.
func NewPerlRejex(ignoreErrors ...bool) PerlFlavorInterface {
    r := createRejexBuilder(GoFlavor, ignoreErrors)
    return PerlFlavorInterface(r)
}

// NewPerlRejexFromString creates a new RejexBuilder object used to construct a regex and
// populates it with a provided regex string, this string is not validated to be valid
// syntax. This uses the Perl flavored syntax.
func NewPerlRejexFromString(s string, ignoreErrors ...bool) PerlFlavorInterface {
    r := createRejexBuilder(GoFlavor, ignoreErrors)
    r.WriteString(s)
    return PerlFlavorInterface(r)
}

// Build constructs the final regex string and returns it along with a list of errors
func (r *RejexBuilder) Build() (string, []RejexError) {
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

    var flagStr, builtRejex string
    switch r.flavor {
    case GoFlavor:
        flagStr = "(?"
        for f, b := range r.flags {
            if b { flagStr += string(f) }
        }
        flagStr += ")"
        if flagStr == "(?)" {
            builtRejex = r.String()
        } else {
            builtRejex = flagStr + r.String()
        }
    case ECMAFlavor:
        flagStr = ""
        for f, b := range r.flags {
            if b { flagStr += string(f) }
        }
        builtRejex = fmt.Sprintf("/%s/%s", r.String(), flagStr)
    case PerlFlavor:
        flagStr = ""
        for f, b := range r.flags {
            if b { flagStr += string(f) }
        }
        builtRejex = fmt.Sprintf("/%s/%s", r.String(), flagStr)
    }

    return builtRejex, r.Errors
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

// Not queues the following segment to be negated, converting '\d' to '\D' for instance
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

// Characters matches the exact input provided to it
func (r *RejexBuilder) Characters(s string) *RejexBuilder {
    return r.appendSegment(characters, s)
}

// EscapedCharacters matches the input provided after escaping the
// regex special characters from it
func (r *RejexBuilder) EscapedCharacters(s string) *RejexBuilder {
    segment := regexp.QuoteMeta(s)
    return r.appendSegment(characters, segment)
}

// AnyChar matches any single character
func (r *RejexBuilder) AnyChar() *RejexBuilder {
    return r.appendSegment(characters, ".")
}

// Literally matches the provided input enclosed in an escape sequence (\Q...\E)
func (r *RejexBuilder) Literally(s string) *RejexBuilder {
    segment := fmt.Sprintf("\\Q%s\\E", s)
    return r.appendSegment(characters, segment)
}

// Anchors

// Starting matches the beginning of a string or the beginning of a line
// when the multiline flag is set. It does not match any character
func (r *RejexBuilder) Starting() *RejexBuilder {
    return r.appendSegment(anchor, "^")
}

// AbsoluteStarting represents the absolute beginning of a string
// unlike Starting, the multiline flag doesn't affect this, it always matches
// the very beginning of a string. It does not match any character
func (r *RejexBuilder) AbsoluteStarting() *RejexBuilder {
    return r.appendSegment(anchor, "\\A")
}

// Ending matches the end of a string or the end of a line
// when the multiline flag is set. It does not match any character
func (r *RejexBuilder) Ending() *RejexBuilder {
    return r.appendSegment(anchor, "$")
}

// AbsoluteEnding represents the absolute end of a string
// unlike Ending, the multiline flag doesn't affect this, it always matches
// the very end of a string. It does not match any character
func (r *RejexBuilder) AbsoluteEnding() *RejexBuilder {
    return r.appendSegment(anchor, "\\z")
}

// WordBoundary matches the end or beginning of any word, it does not match
// any character but is an anchor between a word character and a non
// word character
func (r *RejexBuilder) WordBoundary() *RejexBuilder {
    return r.appendSegment(anchor, "\\b", "\\B")
}

// EndOfLastMatch matches at the end of the previous match during the second
// and following match attempts. Matches at the start of the string during the
// first match attempt
func (r *RejexBuilder) EndOfLastMatch() *RejexBuilder {
    return r.appendSegment(anchor, "\\G")
}


// Quantifiers

func checkForGroup(s, q string) string {
    if len(s) > 1 {
        return fmt.Sprintf("(?:%s)%s", s, q)
    } else {
        return fmt.Sprintf("%s%s", s, q)
    }
}

// ZeroOrOneOf matches exactly 0 or 1 occurance of the provided input
// the input can be an empty string for this to affect the segment preceding it.
// Matches as many characters as it can
func (r *RejexBuilder) ZeroOrOneOf(s string) *RejexBuilder {
    segment := checkForGroup(s, "?")
    return r.appendSegment(quantifier, segment)
}

// ZeroOrMoreOf matches any number of occurances of the provided input
// the input can be an empty string for this to affect the segment preceding it.
// Matches as many characters as it can
func (r *RejexBuilder) ZeroOrMoreOf(s string) *RejexBuilder {
    segment := checkForGroup(s, "*")
    return r.appendSegment(quantifier, segment)
}

// OneOrMoreOf matches more than 1 occurances of the provided input
// the input can be an empty string for this to affect the segment preceding it.
// Matches as many characters as it can
func (r *RejexBuilder) OneOrMoreOf(s string) *RejexBuilder {
    segment := checkForGroup(s, "+")
    return r.appendSegment(quantifier, segment)
}

// NOf matches exactly n occurances of the provided input
// the input can be an empty string for this to affect the segment preceding it.
// Matches as many characters as it can
func (r *RejexBuilder) NOf(s string, n int) *RejexBuilder {
    segment := checkForGroup(s, fmt.Sprintf("{%d}", n))
    return r.appendSegment(quantifier, segment)
}

// NOrMoreOf matches more than n occurances of the provided input
// the input can be an empty string for this to affect the segment preceding it.
// Matches as many characters as it can
func (r *RejexBuilder) NOrMoreOf(s string, n int) *RejexBuilder {
    segment := checkForGroup(s, fmt.Sprintf("{%d,}", n))
    return r.appendSegment(quantifier, segment)
}

// NToMOf matches more than n and upto m occurances of the provided input
// the input can be an empty string for this to affect the segment preceding it.
// Matches as many characters as it can, fewer than m
func (r *RejexBuilder) NToMOf(s string, n, m int) *RejexBuilder {
    segment := checkForGroup(s, fmt.Sprintf("{%d,%d}", n, m))
    return r.appendSegment(quantifier, segment)
}

// Meta

// PreferFewer when used after a quantifier (such as OneOrMoreOf or NOf) makes
// the segment match as few characters as it can, opposite of their default behaviour
func (r *RejexBuilder) PreferFewer() *RejexBuilder {
    if r.lastSegmentType == quantifier {
        r.appendSegment(meta, "?")
    } else {
        r.addError(
            "'PreferFewer()' should only be used after a quantifier",
        )
    }
    return r
}

// PossessiveQuantifier when used after a quantifier (such as OneOrMoreOf or NOf)
// makes the segment match as many items as possible, without trying any permutations
// with less matches even if the remainder of the regex fails.
func (r *RejexBuilder) PossessiveQuantifier() *RejexBuilder {
    if r.lastSegmentType == quantifier {
        r.appendSegment(meta, "+")
    } else {
        r.addError(
            "'PossessiveQuantifier()' should only be used after a quantifier",
        )
    }
    return r
}

// Or represents an alternative between whatever precedes it and whatever follows it.
// Can be used within a group construct. Can be repeated to provide more than 2 alternatives
func (r *RejexBuilder) Or() *RejexBuilder {
    return r.appendSegment(meta, "|")
}

// EitherOr matches any of the provided input strings by chaining together segments using
// the Or syntax. This uses a non-capturing group by default
func (r *RejexBuilder) EitherOr(s ...string) *RejexBuilder {
    var segment string
    if len(s) > 1 {
        segment = fmt.Sprintf("(?:%s)", strings.Join(s, "|"))
        r.appendSegment(characters, segment)
    } else {
        r.addError(
            "Not enough options specified in 'EitherOr()'",
        )
    }
    return r
}

// CapturedPatternByNum matches a previusly captured group with the provided
// group number
func (r *RejexBuilder) CapturedPatternByNum(n int) *RejexBuilder {
    if n > 0 && n < 100 {
        segment := fmt.Sprintf("\\%d", n)
        r.appendSegment(meta, segment)
    } else if n < 0 && r.flavor == PerlFlavor {
        segment := fmt.Sprintf("\\g%d", n)
        r.appendSegment(meta, segment)
    } else {
        r.addError("Pattern number out of bounds")
    }

    return r
}

// CapturedPatternByName matches a previusly captured group with the provided
// group name
func (r *RejexBuilder) CapturedPatternByName(s string) *RejexBuilder {
    segment := fmt.Sprintf("\\k<%s>", s)
    return r.appendSegment(meta, segment)
}

// Group Constructs

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

// BeginCaptureGroup represents the start of a new capture group with a group number
func (r *RejexBuilder) BeginCaptureGroup() *RejexBuilder {
    return r.startNewGroup("(")
}

// BeginNamedCaptureGroup represents the start of a new capture group with a group name
func (r *RejexBuilder) BeginNamedCaptureGroup(name string) *RejexBuilder {
    segment := fmt.Sprintf("(?P<%s>", name)
    return r.startNewGroup(segment)
}

// BeginNonCaptureGroup represents the start of a new group with no group number or name
func (r *RejexBuilder) BeginNonCaptureGroup() *RejexBuilder {
    return r.startNewGroup("(?:")
}

// BeginGroupWithFlags represents the start of a new group which use the provided flags.
// These flags only affect the pattern within this group
func (r *RejexBuilder) BeginGroupWithFlags(f []RejexFlag) *RejexBuilder {
    segment := fmt.Sprintf("(?%s:", string(f))
    return r.startNewGroup(segment)
}

// BeginAtomicGroup represents the start of a new group which prevents the
// regex engine from backtracking back into the group after a match has been
// found for the group. If the remainder of the regex fails, the engine may
// backtrack over the group if a quantifier or alternation makes it optional.
// But it will not backtrack into the group to try other permutations of the group
func (r *RejexBuilder) BeginAtomicGroup() *RejexBuilder {
    return r.startNewGroup("(?>")
}

// BeginBranchResetGroup represents the start of a new group which if it has
// multiple alternatives with capturing groups, then the capturing group
// numbers are the same in all the alternatives
func (r *RejexBuilder) BeginBranchResetGroup() *RejexBuilder {
    return r.startNewGroup("(?|")
}

// BeginPosLookahead represents the start of a new group which only allows the preceding
// segment to match when the pattern in this group follows it but without actualy matching
// this pattern
func (r *RejexBuilder) BeginPosLookahead() *RejexBuilder {
    return r.startNewGroup("(?=")
}

// BeginNegLookahead represents the start of a new group which only allows the preceding
// segment to match when the pattern in this group does not follow it but without actualy matching
// this pattern
func (r *RejexBuilder) BeginNegLookahead() *RejexBuilder {
    return r.startNewGroup("(?!")
}

// BeginPosLookbehind represents the start of a new group which only allows the following
// segment to match when the pattern in this group precedes it but without actualy matching
// this pattern
func (r *RejexBuilder) BeginPosLookbehind() *RejexBuilder {
    return r.startNewGroup("(?<=")
}

// BeginNegLookbehind represents the start of a new group which only allows the following
// segment to match when the pattern in this group does not precede it but without actualy matching
// this pattern
func (r *RejexBuilder) BeginNegLookbehind() *RejexBuilder {
    return r.startNewGroup("(?<!")
}

// EndGroup represents the end of the last opened group
func (r *RejexBuilder) EndGroup() *RejexBuilder {
    if r.groupActive {
        segment := r.groupContent[r.groupNestingLevel] + ")"
        r.groupContent = r.groupContent[:r.groupNestingLevel]
        r.groupNestingLevel--
        if r.groupNestingLevel == 0 {
            r.groupActive = false
        }
        r.appendSegment(characters, segment)
    } else {
        r.addError(
            "Cannot end group, no group open",
        )
    }
    return r
}

// BeginSelectionSet represents the start of a set of characters out of which only one needs be matched
func (r *RejexBuilder) BeginSelectionSet() *RejexBuilder {
    if !r.selectionActive {
        r.selectionActive = true
        r.selectionContent = "["
    } else {
        r.addError("Cannot nest selection sets")
    }
    return r
}

// BeginNonSelectionSet represents the start of a set of characters out of which none should be matched
func (r *RejexBuilder) BeginNonSelectionSet() *RejexBuilder {
    if !r.selectionActive {
        r.selectionActive = true
        r.selectionContent = "[^"
    } else {
        r.addError("Cannot nest selection sets")
    }
    return r
}

// EndGroup represents the end of the last opened selection set
func (r *RejexBuilder) EndSelectionSet() *RejexBuilder {
    if r.selectionActive {
        segment := r.selectionContent + "]"
        r.selectionContent = ""
        r.selectionActive = false
        r.appendSegment(characters, segment)
    } else {
        r.addError(
            "Cannot end selection set, no set open",
        )
    }
    return r
}
