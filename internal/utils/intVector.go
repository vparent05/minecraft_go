package utils

type IntVector2 struct {
	X int
	Y int
}

type IntVector3 struct {
	X int
	Y int
	Z int
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 IntVector2) Add(v2 IntVector2) IntVector2 {
	return IntVector2{v1.X + v2.X, v1.Y + v2.Y}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 IntVector2) Sub(v2 IntVector2) IntVector2 {
	return IntVector2{v1.X - v2.X, v1.Y - v2.Y}
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 IntVector3) Add(v2 IntVector3) IntVector3 {
	return IntVector3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 IntVector3) Sub(v2 IntVector3) IntVector3 {
	return IntVector3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}
