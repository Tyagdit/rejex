# Rejex
An intuitive and readable way to create Regular Expressions in Go using method chains.
Find the documentation [here](https://pkg.go.dev/github.com/tyagdit/rejex)

# Usage

### Building regex

Each method in the chain represents one segment of the final regex.
Begin by creating a Rejex object using `NewRejex()` or `NewRejexFromString()` with an existing regex string.
Finish building the regex with the Build() method to obtain the built regex

```Go
reg, _ := rejex.NewRejex().
        Starting().
        Characters("abcd").
        AnyDigit().
        Build()
```

Group constructs are created using a starting method such as `BeginCaptureGroup()` followed by the group contents
and completed with `EndGroup()`. A similar method exists for selection sets (like [a-z0-9]).

```Go
reg, _ := rejex.NewRejex().
        BeginNamedCaptureGroup("groupname").
            AnyUppercase().
            AnyLowercase().
            AnyDigit().
        EndGroup().
        Build()
```

**Note**: Be mindful of passing escaped characters or characters with a preceding `\`, Go requires you to
escape the `\` for it to be passed to the method or it will be treated as a metacharacter for the compiler
to be interpreted before being passed. In other words if you're passing `\d` or directly to a method Go will
error out saying `Unknown escape sequence`, instead pass it as `\\d`. Same rule applies to metacharacters
that Go does recognise like `\n`

### Negation

Negation can only be used on character classses such as `AnyDigit()` or `AnyFromCharRange()`.
To get negated selection sets either use `BeginNonSelectionSet()` or `Not().AnyFrom()`.

```Go
reg, _ := rejex.NewRejex().
        BeginNamedCaptureGroup("groupname").
            AnyUppercase().
            AnyLowercase().
            Not().AnyDigit().
        EndGroup().
        Build()
```

### Quantifiers

Quantifier methods can be provided with a string to quantify or with an empty string following
a character segment to quantify said preceding segment

```Go
reg, _ := rejex.NewRejex().
        BeginNamedCaptureGroup("groupname").
            AnyUppercase().
            AnyLowercase().
            NOrMoreOf("hi", 2).
            Not().AnyDigit().OneOrMoreOf("").
        EndGroup().
        Build()
```

### Flags

Add or remove flags using the `AddFlags()` and `RemoveFlags()` methods. These can be used anywhere
in the chain before `Build()`.

```Go
reg, _ := rejex.NewRejex().
        BeginNamedCaptureGroup("groupname").
            AnyUppercase().
            AnyLowercase().
            NOrMoreOf("hi", 2).
            Not().AnyDigit().OneOrMoreOf("").
        EndGroup().
        AddFlags(
            rejex.CaseInsensitiveFlag,
            rejex.MultilineFlag,
        ).
        Build()
```

### Errors

The `Build()` method returns 2 values, the regex string and errors encountered while building it.
Most errors are handled gracefully by omiting the offending segments and the error is logged.
The returned error value is of type `[]RejexError` which is a list of the errors and their positions.

```Go
reg, e = rejex.NewRejex().
        AnyDigit().
        PreferFewer().
        EndGroup().
        Build()
fmt.Printf("%#q\n%s\n", e, reg)
```
gives the output

```
Error while building regex at position 2: 'PreferFewer()' should only be used after a quantifier
Error while building regex at position 2: Cannot end group, no group open
[{'\x02' `'PreferFewer()' should only be used after a quantifier`} {'\x02' `Cannot end group, no group open`}]
\d
```

Notice the errors reported directly, these will be printed by default. To silence them, the constructor
methods have an `ignoreErrors` param.

```Go
reg, e = rejex.NewRejex(true).
        AnyDigit().
        PreferFewer().
        EndGroup().
        Build()
fmt.Printf("%#q\n%s\n", e, reg)
```
Now this gives the output 

```
[{'\x02' `'PreferFewer()' should only be used after a quantifier`} {'\x02' `Cannot end group, no group open`}]
\d
```

### Flavors

The default flavor is the Go regex syntax specified in the
[regexp/syntax](https://pkg.go.dev/regexp/syntax@go1.18.1) package of the standard library.
To use different flavors, use the corresponding constructors.

```Go
reg, _ := rejex.NewECMARejex().
        BeginPosLookahead().
            AnyUppercase().
            AnyLowercase().
            NOrMoreOf("hi", 2).
        EndGroup().
        Build()
```

The supported flavors are:

- Golang
- ECMAScript

Many flavors of regex have multiple implementations, subflavors and different default options
which makes it impractical to provide a comprehensive way to generate reliable regexes. This means
that the features available, or their particular syntax may be incompatible or behave differently
in your particular case. It is advised that you look at the docs of the specific implementation
you're using in case you need some complex regexes or the generated regex fails. You should also
test that the regex behaves the way you intend to before using it in production.

# Examples

An example to make a regex for UNIX file paths

```Go
reg, _ := rejex.NewRejex().
        BeginNonCaptureGroup().
            Characters("/").
            BeginSelectionSet().
                AnyFromCharRange("a", "z").
                AnyUppercase().
                AnyFrom("\\d.").
            EndSelectionSet().
            BeginSelectionSet().
                AnyLetter().
                AnyDigit().
                EscapedCharacters("-.").
            EndSelectionSet().
            NToMOf("", 0, 61).
        EndGroup().
        OneOrMoreOf("").
        Build()
```
creates `(?:/[a-zA-Z\d.][a-zA-Z\d-\.]{0,61})+`
