package mtransform

import (
	"errors"
	"fmt"
	"math"
)

type Transform [3][3]float64

// Point represents a 2D point
type Point struct {
	X, Y float64
}

func (t *Transform) Apply(x float64, y float64) (float64, float64) {
	var X, Y float64
	X = t[0][0]*x + t[0][1]*y + t[0][2]
	Y = t[1][0]*x + t[1][1]*y + t[1][2]
	return X, Y
}

func Identity() Transform {
	var t Transform
	t[0][0] = 1
	t[1][1] = 1
	t[2][2] = 1
	return t
}
func NewTransform() *Transform {
	var t Transform
	t = Identity()
	return &t
}

func MultiplyTransforms(a Transform, b Transform) Transform {
	return Transform{
		{
			a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0],
			a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1],
			a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2],
		},
		{
			a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0],
			a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1],
			a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2],
		},
		{
			a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0],
			a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1],
			a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2],
		},
	}
}

func (a *Transform) MultiplyWith(b Transform) {
	*a = MultiplyTransforms(*a, b)
}

func (t *Transform) Scale(x float64, y float64) {
	a := Identity()
	a[0][0] = x
	a[1][1] = y
	t.MultiplyWith(a)
}
func (t *Transform) Translate(x float64, y float64) {
	a := Identity()

	a[0][2] = x
	a[1][2] = y
	t.MultiplyWith(a)
}

func (t *Transform) RotateOrigin(angle float64) {
	a := Identity()
	a[0][0] = math.Cos(angle)
	a[0][1] = -math.Sin(angle)
	a[1][0] = math.Sin(angle)
	a[1][1] = a[0][0]
	t.MultiplyWith(a)
}

func (t *Transform) RotatePoint(angle float64, x float64, y float64) {
	t.Translate(x, y)
	t.RotateOrigin(angle)
	t.Translate(-x, -y)
}

func (t *Transform) SkewX(angle float64) {
	a := Identity()
	a[0][1] = math.Tan(angle)
	t.MultiplyWith(a)
}

func (t *Transform) SkewY(angle float64) {
	a := Identity()
	a[1][0] = math.Tan(angle)
	t.MultiplyWith(a)
}

func (t *Transform) Equals(t2 *Transform) bool {
	return t[0][0] == t2[0][0] && t[0][1] == t2[0][1] && t[0][2] == t2[0][2] &&
		t[1][0] == t2[1][0] && t[1][1] == t2[1][1] && t[1][2] == t2[1][2] &&
		t[2][0] == t2[2][0] && t[2][1] == t2[2][1] && t[2][2] == t2[2][2]
}

// Reflection transformations
func (t *Transform) ReflectX() {
	// Reflect across X-axis
	a := Identity()
	a[1][1] = -1
	t.MultiplyWith(a)
}

func (t *Transform) ReflectY() {
	// Reflect across Y-axis
	a := Identity()
	a[0][0] = -1
	t.MultiplyWith(a)
}

func (t *Transform) ReflectOrigin() {
	// Reflect through origin (180° rotation)
	a := Identity()
	a[0][0] = -1
	a[1][1] = -1
	t.MultiplyWith(a)
}

// Matrix analysis functions
func (t *Transform) Determinant() float64 {
	// Calculate determinant for 2D part
	return t[0][0]*t[1][1] - t[0][1]*t[1][0]
}

func (t *Transform) IsInvertible() bool {
	return math.Abs(t.Determinant()) > 1e-10
}

func (t *Transform) Invert() (*Transform, error) {
	det := t.Determinant()
	if math.Abs(det) < 1e-10 {
		return nil, errors.New("matrix is not invertible")
	}

	invDet := 1.0 / det
	result := Transform{
		{t[1][1] * invDet, -t[0][1] * invDet, (t[0][1]*t[1][2] - t[1][1]*t[0][2]) * invDet},
		{-t[1][0] * invDet, t[0][0] * invDet, (t[1][0]*t[0][2] - t[0][0]*t[1][2]) * invDet},
		{0, 0, 1},
	}
	return &result, nil
}

// Decomposition and analysis functions
func (t *Transform) GetTranslation() (float64, float64) {
	return t[0][2], t[1][2]
}

func (t *Transform) GetScale() (float64, float64) {
	// Extract scale factors from matrix
	sx := math.Sqrt(t[0][0]*t[0][0] + t[1][0]*t[1][0])
	sy := math.Sqrt(t[0][1]*t[0][1] + t[1][1]*t[1][1])

	// Handle negative determinant (reflection)
	if t.Determinant() < 0 {
		sy = -sy
	}

	return sx, sy
}

func (t *Transform) GetRotation() float64 {
	// Extract rotation angle
	return math.Atan2(t[1][0], t[0][0])
}

func (t *Transform) IsIdentity() bool {
	id := Identity()
	return t.Equals(&id)
}

// Utility functions
func (t *Transform) Reset() {
	*t = Identity()
}

func (t *Transform) Clone() *Transform {
	clone := *t
	return &clone
}

func (t *Transform) String() string {
	return fmt.Sprintf("Transform[%.3f %.3f %.3f; %.3f %.3f %.3f; %.3f %.3f %.3f]",
		t[0][0], t[0][1], t[0][2],
		t[1][0], t[1][1], t[1][2],
		t[2][0], t[2][1], t[2][2])
}

// Point operations
func (t *Transform) ApplyToPoint(p Point) Point {
	x, y := t.Apply(p.X, p.Y)
	return Point{X: x, Y: y}
}

func (t *Transform) ApplyToPoints(points []Point) []Point {
	result := make([]Point, len(points))
	for i, p := range points {
		result[i] = t.ApplyToPoint(p)
	}
	return result
}

// Advanced transformations
func (t *Transform) Shear(shx, shy float64) {
	// General shearing transformation
	a := Identity()
	a[0][1] = shx
	a[1][0] = shy
	t.MultiplyWith(a)
}

func (t *Transform) ScaleAroundPoint(sx, sy, cx, cy float64) {
	// Scale around a specific center point
	// Create the composite transformation: T(cx,cy) * S(sx,sy) * T(-cx,-cy)
	translateBack := Identity()
	translateBack[0][2] = cx
	translateBack[1][2] = cy

	scale := Identity()
	scale[0][0] = sx
	scale[1][1] = sy

	translateToOrigin := Identity()
	translateToOrigin[0][2] = -cx
	translateToOrigin[1][2] = -cy

	// Build the composite transformation
	temp := MultiplyTransforms(scale, translateToOrigin)
	composite := MultiplyTransforms(translateBack, temp)

	t.MultiplyWith(composite)
}

func (t *Transform) RotateAroundPoint(angle, cx, cy float64) {
	// Rotate around a specific center point (alternative to RotatePoint)
	// Create the composite transformation: T(cx,cy) * R(angle) * T(-cx,-cy)
	translateBack := Identity()
	translateBack[0][2] = cx
	translateBack[1][2] = cy

	rotate := Identity()
	rotate[0][0] = math.Cos(angle)
	rotate[0][1] = -math.Sin(angle)
	rotate[1][0] = math.Sin(angle)
	rotate[1][1] = rotate[0][0]

	translateToOrigin := Identity()
	translateToOrigin[0][2] = -cx
	translateToOrigin[1][2] = -cy

	// Build the composite transformation
	temp := MultiplyTransforms(rotate, translateToOrigin)
	composite := MultiplyTransforms(translateBack, temp)

	t.MultiplyWith(composite)
}

// Validation and comparison
func (t *Transform) IsNearlyEqual(other *Transform, epsilon float64) bool {
	// Compare with tolerance for floating point precision
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(t[i][j]-other[i][j]) > epsilon {
				return false
			}
		}
	}
	return true
}

func (t *Transform) IsOrthogonal() bool {
	// Check if transformation preserves angles (rotation + reflection only)
	det := t.Determinant()
	return math.Abs(math.Abs(det)-1.0) < 1e-10
}

// Interpolation
func (t *Transform) Lerp(other *Transform, factor float64) Transform {
	// Linear interpolation between two transforms
	var result Transform
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result[i][j] = t[i][j]*(1-factor) + other[i][j]*factor
		}
	}
	return result
}

// SVG integration
func (t *Transform) ToSVGMatrix() string {
	return fmt.Sprintf("matrix(%g,%g,%g,%g,%g,%g)",
		t[0][0], t[1][0], t[0][1], t[1][1], t[0][2], t[1][2])
}
