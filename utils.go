package rejex

func (r *RejexBuilder) LineEnding() *RejexBuilder {
    return r.appendSegment(CHARACTERS, "\\n\\r\\v\\f")
}
