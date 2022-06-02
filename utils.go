package rejex

// LineEnding matches any single character that starts a new line
func (r *RejexBuilder) LineEnding() *RejexBuilder {
    return r.appendSegment(characters, "\\n\\r\\v\\f")
}
