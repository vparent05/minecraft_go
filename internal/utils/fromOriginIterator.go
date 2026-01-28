package utils

import (
	"iter"
	"unsafe"
)

func FromOriginIterator1[T any](values []T, origin int) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		length := len(values)

		// left [ origin - width/2, origin [
		for virtualI := origin - length/2; virtualI < origin; virtualI++ {
			i := Mod(virtualI, length)
			if !yield(i, values[i]) {
				return
			}
		}

		// right [ origin + width/2, origin ]
		for virtualI := origin + length/2; virtualI >= origin; virtualI-- {
			i := Mod(virtualI, length)
			if !yield(i, values[i]) {
				return
			}
		}
	}
}

func FromOriginIterator2[T any](values [][]T, origin IntVector2) iter.Seq2[IntVector2, T] {
	return func(yield func(IntVector2, T) bool) {
		for i, row := range FromOriginIterator1(values, origin.X) {
			for j, t := range FromOriginIterator1(row, origin.Y) {
				if !yield(IntVector2{i, j}, t) {
					return
				}
			}
		}
	}
}

func UnsafeFromOriginIterator1[T any](values *T, elementSize uintptr, start int, end int, size int, origin int) iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		base := uintptr(unsafe.Pointer(values))

		// left [ start, origin [
		for virtualI := start; virtualI < origin; virtualI++ {
			i := Mod(virtualI, size)
			if !yield(i, (*T)(unsafe.Pointer(base+uintptr(i)*elementSize))) {
				return
			}
		}

		// right [ end, origin ]
		for virtualI := end - 1; virtualI >= origin; virtualI-- {
			i := Mod(virtualI, size)
			if !yield(i, (*T)(unsafe.Pointer(base+uintptr(i)*elementSize))) {
				return
			}
		}
	}
}

func UnsafeFromOriginIterator3[T any](values *T, start, end, size, origin IntVector3) iter.Seq2[IntVector3, T] {
	return func(yield func(IntVector3, T) bool) {
		for i, slice := range UnsafeFromOriginIterator1(values, uintptr(size.Y)*uintptr(size.Z)*unsafe.Sizeof(*values), start.X, end.X, size.X, origin.X) {
			for j, row := range UnsafeFromOriginIterator1(slice, uintptr(size.Z)*unsafe.Sizeof(*values), start.Y, end.Y, size.Y, origin.Y) {
				for k, t := range UnsafeFromOriginIterator1(row, unsafe.Sizeof(*values), start.Z, end.Z, size.Z, origin.Z) {
					if !yield(IntVector3{i, j, k}, *t) {
						return
					}
				}
			}
		}
	}
}
