package utils

import (
	"iter"
)

func FromOriginIterator1[T any](values []T, origin int) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		length := len(values)

		// left [ originJ - width/2, originJ [
		for virtualI := origin - length/2; virtualI < origin; virtualI++ {
			i := Mod(virtualI, length)
			if !yield(i, values[i]) {
				return
			}
		}

		// right [ originJ + width/2, originJ [
		for virtualI := origin + length/2; virtualI >= origin; virtualI-- {
			i := Mod(virtualI, length)
			if !yield(i, values[i]) {
				return
			}
		}
	}
}

func FromOriginIterator2[T any](values [][]T, originI, originJ int) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, row := range FromOriginIterator1(values, originI) {
			for j, t := range FromOriginIterator1(row, originJ) {
				if !yield(i*len(row)+j, t) {
					return
				}
			}
		}
	}
}
