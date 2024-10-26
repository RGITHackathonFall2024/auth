package utils

import "iter"

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func Map[T any, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for x := range seq {
			if !yield(f(x)) {
				break
			}
		}
	}
}
