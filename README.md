# Rejex
 An intuitive and readable way to create Regular Expressions in Go using method chains

# Usage

An example to make a regex for UNIX file paths

```Go
reg := rejex.NewRejex().
        BeginNonCaptureGroup().
            Characters("/").
            BeginSelectionSet().
                CharRange("a", "z").
                Uppercase().
                AnyFrom("\\d.").
            EndSelectionSet().
            BeginSelectionSet().
                Letter().
                Digit().
                EscapedCharacters("-.").
            EndSelectionSet().
            NToMOf("", 0, 61).
        EndGroup().
        OneOrMoreOf("").
        Build()
```
creates `(?:/[a-zA-Z\d.][a-zA-Z\d-\.]{0,61})+`
