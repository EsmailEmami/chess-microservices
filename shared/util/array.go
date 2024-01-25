package util

import "reflect"

func ArrayRemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func ArrayRemoveItem[T any](s []T, item T) []T {
	for i, v := range s {
		if reflect.DeepEqual(v, item) {
			s[i] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}
