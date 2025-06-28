package mtransform

import (
	"math"
	"testing"
)

func TestMultiply(t *testing.T) {
	a := Transform{{1, 2, 3}, {-1, -2, -3}, {4, 5, 6}}
	b := Transform{{0, 1, 2}, {0, -1, -2}, {2, 1, 0}}
	want := Transform{{6, 2, -2}, {-6, -2, 2}, {12, 5, -2}}
	got := MultiplyTransforms(a, b)
	if got != want {
		t.Errorf("Multiplying: got %v, want %v", got, want)
	}
}

// Test reflection functions
func TestReflectX(t *testing.T) {
	tr := NewTransform()
	tr.ReflectX()

	x, y := tr.Apply(1, 2)
	if x != 1 || y != -2 {
		t.Errorf("ReflectX: expected (1, -2), got (%f, %f)", x, y)
	}
}

func TestReflectY(t *testing.T) {
	tr := NewTransform()
	tr.ReflectY()

	x, y := tr.Apply(1, 2)
	if x != -1 || y != 2 {
		t.Errorf("ReflectY: expected (-1, 2), got (%f, %f)", x, y)
	}
}

func TestReflectOrigin(t *testing.T) {
	tr := NewTransform()
	tr.ReflectOrigin()

	x, y := tr.Apply(1, 2)
	if x != -1 || y != -2 {
		t.Errorf("ReflectOrigin: expected (-1, -2), got (%f, %f)", x, y)
	}
}

// Test matrix analysis functions
func TestDeterminant(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)

	det := tr.Determinant()
	expected := 6.0
	if math.Abs(det-expected) > 1e-10 {
		t.Errorf("Determinant: expected %f, got %f", expected, det)
	}
}

func TestIsInvertible(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)

	if !tr.IsInvertible() {
		t.Error("IsInvertible: scalable matrix should be invertible")
	}

	// Test non-invertible matrix
	tr2 := &Transform{{1, 2, 0}, {2, 4, 0}, {0, 0, 1}}
	if tr2.IsInvertible() {
		t.Error("IsInvertible: singular matrix should not be invertible")
	}
}

func TestInvert(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)
	tr.Translate(5, 7)

	inv, err := tr.Invert()
	if err != nil {
		t.Errorf("Invert: unexpected error: %v", err)
	}

	// Test that applying transform then inverse gives identity
	combined := MultiplyTransforms(*tr, *inv)
	identity := Identity()

	if !combined.IsNearlyEqual(&identity, 1e-10) {
		t.Errorf("Invert: T * T^-1 should equal identity")
	}
}

// Test decomposition functions
func TestGetTranslation(t *testing.T) {
	tr := NewTransform()
	tr.Translate(5, 7)

	tx, ty := tr.GetTranslation()
	if tx != 5 || ty != 7 {
		t.Errorf("GetTranslation: expected (5, 7), got (%f, %f)", tx, ty)
	}
}

func TestGetScale(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)

	sx, sy := tr.GetScale()
	if math.Abs(sx-2) > 1e-10 || math.Abs(sy-3) > 1e-10 {
		t.Errorf("GetScale: expected (2, 3), got (%f, %f)", sx, sy)
	}
}

func TestGetRotation(t *testing.T) {
	tr := NewTransform()
	angle := math.Pi / 4 // 45 degrees
	tr.RotateOrigin(angle)

	gotAngle := tr.GetRotation()
	if math.Abs(gotAngle-angle) > 1e-10 {
		t.Errorf("GetRotation: expected %f, got %f", angle, gotAngle)
	}
}

func TestIsIdentity(t *testing.T) {
	tr := NewTransform()
	if !tr.IsIdentity() {
		t.Error("IsIdentity: new transform should be identity")
	}

	tr.Scale(2, 2)
	if tr.IsIdentity() {
		t.Error("IsIdentity: scaled transform should not be identity")
	}
}

// Test utility functions
func TestReset(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)
	tr.Translate(5, 7)

	tr.Reset()
	if !tr.IsIdentity() {
		t.Error("Reset: should return to identity matrix")
	}
}

func TestClone(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)
	tr.Translate(5, 7)

	clone := tr.Clone()
	if !tr.Equals(clone) {
		t.Error("Clone: cloned transform should be equal to original")
	}

	// Modify original and ensure clone is unchanged
	tr.Scale(2, 2)
	if tr.Equals(clone) {
		t.Error("Clone: clone should be independent of original")
	}
}

func TestString(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)

	str := tr.String()
	if str == "" {
		t.Error("String: should return non-empty string")
	}
	// Just check it doesn't panic and returns something
}

// Test point operations
func TestApplyToPoint(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)

	p := Point{X: 1, Y: 2}
	result := tr.ApplyToPoint(p)

	if result.X != 2 || result.Y != 6 {
		t.Errorf("ApplyToPoint: expected (2, 6), got (%f, %f)", result.X, result.Y)
	}
}

func TestApplyToPoints(t *testing.T) {
	tr := NewTransform()
	tr.Scale(2, 3)

	points := []Point{{1, 1}, {2, 2}, {3, 3}}
	results := tr.ApplyToPoints(points)

	expected := []Point{{2, 3}, {4, 6}, {6, 9}}
	for i, result := range results {
		if result.X != expected[i].X || result.Y != expected[i].Y {
			t.Errorf("ApplyToPoints[%d]: expected (%f, %f), got (%f, %f)",
				i, expected[i].X, expected[i].Y, result.X, result.Y)
		}
	}
}

// Test advanced transformations
func TestShear(t *testing.T) {
	tr := NewTransform()
	tr.Shear(0.5, 0) // Shear X by 0.5

	x, y := tr.Apply(2, 4)
	expected_x := 2 + 0.5*4 // x + shx*y
	expected_y := 4.0

	if math.Abs(x-expected_x) > 1e-10 || math.Abs(y-expected_y) > 1e-10 {
		t.Errorf("Shear: expected (%f, %f), got (%f, %f)", expected_x, expected_y, x, y)
	}
}

func TestScaleAroundPoint(t *testing.T) {
	tr := NewTransform()
	tr.ScaleAroundPoint(2, 2, 1, 1) // Scale 2x around point (1,1)

	// Point (1,1) should remain unchanged
	x, y := tr.Apply(1, 1)
	if math.Abs(x-1) > 1e-10 || math.Abs(y-1) > 1e-10 {
		t.Errorf("ScaleAroundPoint: center point should be unchanged, got (%f, %f)", x, y)
	}

	// Point (3,3) should become (5,5)
	x, y = tr.Apply(3, 3)
	if math.Abs(x-5) > 1e-10 || math.Abs(y-5) > 1e-10 {
		t.Errorf("ScaleAroundPoint: expected (5, 5), got (%f, %f)", x, y)
	}
}

func TestRotateAroundPoint(t *testing.T) {
	tr := NewTransform()
	tr.RotateAroundPoint(math.Pi/2, 1, 1) // 90° around (1,1)

	// Point (1,1) should remain unchanged
	x, y := tr.Apply(1, 1)
	if math.Abs(x-1) > 1e-10 || math.Abs(y-1) > 1e-10 {
		t.Errorf("RotateAroundPoint: center point should be unchanged, got (%f, %f)", x, y)
	}

	// Point (2,1) should become (1,2)
	x, y = tr.Apply(2, 1)
	if math.Abs(x-1) > 1e-10 || math.Abs(y-2) > 1e-10 {
		t.Errorf("RotateAroundPoint: expected (1, 2), got (%f, %f)", x, y)
	}
}

// Test validation functions
func TestIsNearlyEqual(t *testing.T) {
	tr1 := NewTransform()
	tr1.Scale(2, 3)

	tr2 := NewTransform()
	tr2.Scale(2.0001, 3.0001)

	if !tr1.IsNearlyEqual(tr2, 1e-3) {
		t.Error("IsNearlyEqual: should be nearly equal with large epsilon")
	}

	if tr1.IsNearlyEqual(tr2, 1e-6) {
		t.Error("IsNearlyEqual: should not be nearly equal with small epsilon")
	}
}

func TestIsOrthogonal(t *testing.T) {
	// Rotation matrix should be orthogonal
	tr := NewTransform()
	tr.RotateOrigin(math.Pi / 4)

	if !tr.IsOrthogonal() {
		t.Error("IsOrthogonal: rotation matrix should be orthogonal")
	}

	// Scaled matrix should not be orthogonal (unless scale is ±1)
	tr.Scale(2, 2)
	if tr.IsOrthogonal() {
		t.Error("IsOrthogonal: scaled matrix should not be orthogonal")
	}
}

// Test interpolation
func TestLerp(t *testing.T) {
	tr1 := NewTransform()
	tr1.Scale(1, 1)

	tr2 := NewTransform()
	tr2.Scale(3, 3)

	// Interpolate halfway
	result := tr1.Lerp(tr2, 0.5)
	expected_scale := 2.0

	sx, sy := result.GetScale()
	if math.Abs(sx-expected_scale) > 1e-10 || math.Abs(sy-expected_scale) > 1e-10 {
		t.Errorf("Lerp: expected scale (2, 2), got (%f, %f)", sx, sy)
	}
}

// Test SVG integration
func TestToSVGMatrix(t *testing.T) {
	tr := NewTransform()
	tr.Translate(5, 7) // Apply translate first
	tr.Scale(2, 3)     // Then scale

	svg := tr.ToSVGMatrix()
	if svg == "" {
		t.Error("ToSVGMatrix: should return non-empty string")
	}

	// Should contain the matrix values
	expected := "matrix(2,0,0,3,5,7)"
	if svg != expected {
		t.Errorf("ToSVGMatrix: expected %s, got %s", expected, svg)
	}
}

// Test complex transformation combinations
func TestComplexTransform(t *testing.T) {
	tr := NewTransform()

	// Apply multiple transformations
	tr.Translate(10, 20)
	tr.RotateOrigin(math.Pi / 4)
	tr.Scale(2, 2)
	tr.ReflectX()

	// Test that it's still functional
	x, y := tr.Apply(1, 1)

	// Just verify it doesn't crash and produces some result
	if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
		t.Errorf("ComplexTransform: result should be finite numbers, got (%f, %f)", x, y)
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	// Test zero scale
	tr := NewTransform()
	tr.Scale(0, 1)

	if tr.IsInvertible() {
		t.Error("EdgeCases: zero scale should not be invertible")
	}

	// Test very small numbers
	tr2 := NewTransform()
	tr2.Scale(1e-15, 1e-15)

	det := tr2.Determinant()
	if math.Abs(det) > 1e-10 {
		t.Error("EdgeCases: very small scale should have very small determinant")
	}
}
