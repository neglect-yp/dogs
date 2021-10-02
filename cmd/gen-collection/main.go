// Usage:
//     package yourcollection
//
//     import "github.com/genkami/dogs/types/iterator"
//
//     // go:generate go run .../../cmd/gen-collection -pkg yourcollection -name YourCollection -out zz_generated.collection.go
//     type YourCollection[T] struct { ... }
//     func FromIterator[T any](it iterator.Iterator[T]) YourCollection[T] { ... }
//     func (xs YourCollection[T]) Iter() iterator.Iterator[T] { ... }

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
)

func main() {
	var (
		pkgName    string
		typeName   string
		constraint string
		exclude    string
		output     string
	)

	flag.StringVar(&pkgName, "pkg", "", "the name of the package")
	flag.StringVar(&typeName, "name", "", "the name of the type")
	flag.StringVar(&constraint, "constraint", "any", "type constraint that the type argument of this type should satisfy")
	flag.StringVar(&exclude, "exclude", "", "comma-separated names of funcions to exclude")
	flag.StringVar(&output, "out", "", "path to output")
	flag.Parse()

	if pkgName == "" {
		fmt.Fprintf(os.Stderr, "gen-collection: missing -pkg\n")
		os.Exit(1)
	}
	if typeName == "" {
		fmt.Fprintf(os.Stderr, "gen-collection: missing -type-name\n")
		os.Exit(1)
	}
	if output == "" {
		fmt.Fprintf(os.Stderr, "gen-collection: missing -out\n")
		os.Exit(1)
	}

	w, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen-collection: failed to open %s: %s\n", output, err.Error())
		os.Exit(1)
	}
	defer w.Close()

	var denyList []string
	if len(exclude) > 0 {
		denyList = strings.Split(exclude, ",")
	}

	tmpl, err := generateTemplate(denyList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen-collection: failed to prepare template: %s\n", err.Error())
		os.Exit(1)
	}
	err = tmpl.Execute(w, map[string]string{
		"PkgName":    pkgName,
		"TypeName":   typeName,
		"Constraint": constraint,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen-collection: failed to write template: %s\n", err.Error())
		os.Exit(1)
	}
}

func generateTemplate(denyList []string) (*template.Template, error) {
	var buf bytes.Buffer
	buf.WriteString(header)

	denyMap := make(map[string]struct{})
	for _, denied := range denyList {
		if _, ok := allFuncs[denied]; !ok {
			return nil, fmt.Errorf("invalid -exclude option: %s not found\n", denied)
		}
		denyMap[denied] = struct{}{}
	}

	keys := make([]string, 0, len(allFuncs))
outer:
	for k, _ := range allFuncs {
		if _, ok := denyMap[k]; !ok {
			keys = append(keys, k)
			continue outer
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(allFuncs[k])
	}
	return template.Must(template.New("").Parse(buf.String())), nil
}

const header = `// Code generated by gen-collection; DO NOT EDIT.

package {{ .PkgName }}

import (
	"github.com/genkami/dogs/classes/algebra"
	"github.com/genkami/dogs/classes/cmp"
	"github.com/genkami/dogs/types/iterator"
	"github.com/genkami/dogs/types/pair"
)

// Some packages are unused depending on -include CLI option.
// This prevents compile error when corresponding functions are not defined.
var _ = (algebra.Monoid[int])(nil)
var _ = (cmp.Ord[int])(nil)
var _ = (iterator.Iterator[int])(nil)
var _ = (*pair.Pair[int, int])(nil)

`

var allFuncs = map[string]string{
	"Find": `
// Find returns a first element in xs that satisfies the given predicate fn.
// It returns false as a second return value if no elements are found.
func Find[T {{ .Constraint }}](xs {{ .TypeName }}[T], fn func(T) bool) (T, bool) {
	return iterator.Find[T](xs.Iter(), fn)
}
`,
	"FindIndex": `
// FindIndex returns a first index of an element in xs that satisfies the given predicate fn.
// It returns negative value if no elements are found.
func FindIndex[T {{ .Constraint }}](xs {{ .TypeName }}[T], fn func(T) bool) int {
	return iterator.FindIndex[T](xs.Iter(), fn)
}
`,
	"FindElem": `
// FindElem returns a first element in xs that equals to e in the sense of given Eq.
// It returns false as a second return value if no elements are found.
func FindElem[T {{ .Constraint }}](xs {{ .TypeName }}[T], e T, eq cmp.Eq[T]) (T, bool) {
	return iterator.FindElem[T](xs.Iter(), e, eq)
}
`,
	"FindElemIndex": `
// FindElemIndex returns a first index of an element in xs that equals to e in the sense of given Eq.
// It returns negative value if no elements are found.
func FindElemIndex[T {{ .Constraint }}](xs {{ .TypeName }}[T], e T, eq cmp.Eq[T]) int {
	return iterator.FindElemIndex[T](xs.Iter(), e, eq)
}
`,
	"Filter": `
// Filter returns a collection that only returns elements that satisfies given predicate.
func Filter[T {{ .Constraint }}](xs {{ .TypeName }}[T], fn func(T) bool) {{ .TypeName }}[T] {
	return FromIterator[T](iterator.Filter[T](xs.Iter(), fn))
}
`,
	"Map": `
// Map returns a collection that applies fn to each element of xs.
func Map[T, U {{ .Constraint }}](xs {{ .TypeName }}[T], fn func(T) U) {{ .TypeName }}[U] {
	return FromIterator[U](iterator.Map[T, U](xs.Iter(), fn))
}
`,
	"ForEach": `
// ForEach applies fn to each element in xs.
func ForEach[T {{ .Constraint }}](xs {{ .TypeName }}[T], fn func(T)) {
	iterator.ForEach[T](xs.Iter(), fn)
}
`,
	"Fold": `
// Fold accumulates every element in a collection by applying fn.
func Fold[T any, U {{ .Constraint }}](init T, xs {{ .TypeName }}[U], fn func(T, U) T) T {
	return iterator.Fold[T, U](init, xs.Iter(), fn)
}
`,
	"Zip": `
// Zip combines two collections into one that contains pairs of corresponding elements.
func Zip[T, U {{ .Constraint }}](a {{ .TypeName }}[T], b {{ .TypeName }}[U]) {{ .TypeName }}[pair.Pair[T, U]] {
	return FromIterator[pair.Pair[T, U]](iterator.Zip(a.Iter(), b.Iter()))
}
`,
	"SumWithInit": `
// SumWithInit sums up init and all values in xs.
func SumWithInit[T {{ .Constraint }}](init T, xs {{ .TypeName }}[T], s algebra.Semigroup[T]) T {
	return Fold[T, T](init, xs, s.Combine)
}
`,
	"Sum": `
// Sum sums up all values in xs.
// It returns m.Empty() when xs is empty.
func Sum[T {{ .Constraint }}](xs {{ .TypeName }}[T], m algebra.Monoid[T]) T {
	var s algebra.Semigroup[T] = m
	return SumWithInit[T](m.Empty(), xs, s)
}
`,
}
