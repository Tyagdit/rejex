package rejex

type RejexFlag rune

const (
    CaseInsensitiveFlag RejexFlag = 'i'
    MultilineFlag RejexFlag = 'm'
    SingleLineFlag RejexFlag = 's'
    UngreedyFlag RejexFlag = 'U'
)

func (r *RejexBuilder) AddFlags(f []RejexFlag) {
    for _, flag := range f {
        r.flags[flag] = true
    }
}

func (r *RejexBuilder) RemoveFlags(f []RejexFlag) {
    for _, flag := range f {
        r.flags[flag] = false
    }
}
