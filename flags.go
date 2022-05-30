package rejex

type RejexFlag rune

const (
    CaseInsensitiveFlag RejexFlag = 'i'
    MultilineFlag RejexFlag = 'm'
    SingleLineFlag RejexFlag = 's'
    UngreedyFlag RejexFlag = 'U'
)

func (r *RejexBuilder) AddFlags(f ...RejexFlag) *RejexBuilder {
    if len(f) == 0 {
        r.addError("No flags provided")
    }
    for _, flag := range f {
        r.flags[flag] = true
    }
    return r
}

func (r *RejexBuilder) RemoveFlags(f ...RejexFlag) *RejexBuilder { 
    if len(f) == 0 {
        r.addError("No flags provided to be removed")
    }
    for _, flag := range f {
        r.flags[flag] = false
    }
    return r
}
