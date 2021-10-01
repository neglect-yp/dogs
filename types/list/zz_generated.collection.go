// Code generated by gen-collection; DO NOT EDIT.

package list

import (
	"github.com/genkami/dogs/classes/algebra"
	"github.com/genkami/dogs/classes/cmp"
	"github.com/genkami/dogs/types/iterator"
	"github.com/genkami/dogs/types/pair"
)

// Find returns a first element in xs that satisfies the given predicate fn.
// It returns false as a second return value if no elements are found.
func Find[T any](xs *List[T], fn func(T) bool) (T, bool) {
	return iterator.Find[T](xs.Iter(), fn)
}

// FindIndex returns a first index of an element in xs that satisfies the given predicate fn.
// It returns negative value if no elements are found.
func FindIndex[T any](xs *List[T], fn func(T) bool) int {
	return iterator.FindIndex[T](xs.Iter(), fn)
}

// FindElem returns a first element in xs that equals to e in the sense of given Eq.
// It returns false as a second return value if no elements are found.
func FindElem[T any](xs *List[T], e T, eq cmp.Eq[T]) (T, bool) {
	return iterator.FindElem[T](xs.Iter(), e, eq)
}

// FindElemIndex returns a first index of an element in xs that equals to e in the sense of given Eq.
// It returns negative value if no elements are found.
func FindElemIndex[T any](xs *List[T], e T, eq cmp.Eq[T]) int {
	return iterator.FindElemIndex[T](xs.Iter(), e, eq)
}

// Map(xs, f) returns a collection that applies f to each element of xs.
func Map[T, U any](xs *List[T], fn func(T) U) *List[U] {
	return FromIterator[U](iterator.Map[T, U](xs.Iter(), fn))
}

// Fold accumulates every element in a collection by applying fn.
func Fold[T any, U any](init T, xs *List[U], fn func(T, U) T) T {
	return iterator.Fold[T, U](init, xs.Iter(), fn)
}

// Zip combines two collections into one that contains pairs of corresponding elements.
func Zip[T, U any](a *List[T], b *List[U]) *List[pair.Pair[T, U]] {
	return FromIterator[pair.Pair[T, U]](iterator.Zip(a.Iter(), b.Iter()))
}

// SumWithInit sums up init and all values in xs.
func SumWithInit[T any](init T, xs *List[T], s algebra.Semigroup[T]) T {
	return Fold[T, T](init, xs, s.Combine)
}

// Sum sums up all values in xs.
// It returns m.Empty() when xs is empty.
func Sum[T any](xs *List[T], m algebra.Monoid[T]) T {
	var s algebra.Semigroup[T] = m
	return SumWithInit[T](m.Empty(), xs, s)
}
