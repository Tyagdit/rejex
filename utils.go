package rejex

// LineEnding matches any single character that starts a new line
func (r *RejexBuilder) LineEnding() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\n\\r\\v\\f")
}
