package rejex

import (
    "fmt"
)

type RejexFlag rune

const (
    CaseInsensitiveFlag RejexFlag = 'i'
    MultilineFlag RejexFlag = 'm'
    SingleLineFlag RejexFlag = 's'
    UngreedyFlag RejexFlag = 'U'
    StickyFlag RejexFlag = 'y'
    UnicodeFlag RejexFlag = 'u'
    GlobalFlag RejexFlag = 'g'
)

func (r *RejexBuilder) changeFlags(f []RejexFlag, state bool) *RejexBuilder {
    if len(f) == 0 {
        r.addError("No flags provided")
    }
    for _, flag := range f {
        if _, ok := r.flags[flag]; ok {
            r.flags[flag] = state
        } else {
            r.addError(fmt.Sprintf("Invalid flag '%v'", flag))
        }
    }
    return r
}

func (r *RejexBuilder) AddFlags(f ...RejexFlag) *RejexBuilder {
    return r.changeFlags(f, true)
}

func (r *RejexBuilder) RemoveFlags(f ...RejexFlag) *RejexBuilder {
    return r.changeFlags(f, true)
}
