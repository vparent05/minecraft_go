package utils

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

var PERMUTATION_ARRAY = generatePermutationArray(512)

func FractalNoise2(x float32, y float32, numberOfOctaves int) float32 {
	var result float32 = 0
	var amplitude float32 = 1
	var frequency float32 = 0.005

	for i := 0; i < numberOfOctaves; i++ {
		result += perlinNoise2(x*frequency, y*frequency) * amplitude
		amplitude /= 2
		frequency *= 2
	}

	return result
}

func perlinNoise2(x float32, y float32) float32 {
	X := int(math.Floor(float64(x))) & 255
	Y := int(math.Floor(float64(y))) & 255

	xf := x - float32(math.Floor(float64(x)))
	yf := y - float32(math.Floor(float64(y)))

	topRight := mgl32.Vec2{xf - 1, yf - 1}
	topLeft := mgl32.Vec2{xf, yf - 1}
	bottomRight := mgl32.Vec2{xf - 1, yf}
	bottomLeft := mgl32.Vec2{xf, yf}

	valueTopRight := PERMUTATION_ARRAY[PERMUTATION_ARRAY[X+1]+Y+1]
	valueTopLeft := PERMUTATION_ARRAY[PERMUTATION_ARRAY[X]+Y+1]
	valueBottomRight := PERMUTATION_ARRAY[PERMUTATION_ARRAY[X+1]+Y]
	valueBottomLeft := PERMUTATION_ARRAY[PERMUTATION_ARRAY[X]+Y]

	dotTopRight := getConstantVector(valueTopRight).Dot(topRight)
	dotTopLeft := getConstantVector(valueTopLeft).Dot(topLeft)
	dotBottomRight := getConstantVector(valueBottomRight).Dot(bottomRight)
	dotBottomLeft := getConstantVector(valueBottomLeft).Dot(bottomLeft)

	u := fade(xf)
	v := fade(yf)

	return lerp(u, lerp(v, dotBottomLeft, dotTopLeft), lerp(v, dotBottomRight, dotTopRight))
}

func shuffleArray(array []int) {
	for i := len(array) - 1; i > 0; i-- {
		index := rand.Int() % i

		temp := array[i]
		array[i] = array[index]
		array[index] = temp
	}
}

func generatePermutationArray(size int) []int {
	array := make([]int, size/2)
	for i := 0; i < size/2; i++ {
		array[i] = i
	}
	shuffleArray(array)
	finalArray := make([]int, size)
	for i := 0; i < size/2; i++ {
		finalArray[i] = array[i]
		finalArray[i+size/2] = array[i]
	}
	return finalArray
}

func lerp(t float32, a float32, b float32) float32 {
	return a + t*(b-a)
}

func fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

func getConstantVector(v int) mgl32.Vec2 {
	h := v & 3
	if h == 0 {
		return mgl32.Vec2{1, 1}
	} else if h == 1 {
		return mgl32.Vec2{-1, 1}
	} else if h == 2 {
		return mgl32.Vec2{-1, -1}
	} else {
		return mgl32.Vec2{1, -1}
	}
}
