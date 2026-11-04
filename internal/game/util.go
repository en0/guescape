package game

import "strings"

func filter[T any](items []T, predicate func (T) bool) []T {
	result := []T{}
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func find[T any](items []T, predicate func (T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

func exists[T any](items []T, predicate func (T) bool) bool {
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

func existsAll[T any](items []T, predicate func (T) bool) bool {
	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}
	return true
}


func transform[T any, R any] (s []T, fn func (T) R) []R {
	var r []R
	for _, i := range s {
		r = append(r, fn(i))
	}
	return r
}

func matchIdentity(candidates []Identity, partial string) (Identity, bool) {
	if partial == "" { return Identity{}, false}
	partial = strings.ToLower(partial)
	m := filter(candidates, func(id Identity) bool {
		return strings.Contains(id.String(), partial)
	})
	if len(m) == 1 { return m[0], true }
	return Identity{}, false
}
