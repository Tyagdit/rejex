package rejex

type RejexFlag rune

const (
    CaseInsensitiveFlag RejexFlag = 'i'
    MultilineFlag RejexFlag = 'm'
    SingleLineFlag RejexFlag = 's'
    UngreedyFlag RejexFlag = 'U'
)

func (r *RejexBuilder) AddFlag(f RejexFlag) {
    r.flags[f] = true
}

func (r *RejexBuilder) RemoveFlag(f RejexFlag) {
    r.flags[f] = false
}
