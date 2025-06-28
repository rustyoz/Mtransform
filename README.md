# Mtransform [![Go](https://github.com/rustyoz/Mtransform/actions/workflows/go.yml/badge.svg)](https://github.com/rustyoz/Mtransform/actions/workflows/go.yml)

A simple and efficient 2D matrix transformation library for Go.

## Overview

Mtransform provides a clean API for performing 2D transformations using 3x3 transformation matrices. It supports common operations like translation, scaling, rotation, and skewing, making it ideal for graphics programming, game development, and geometric calculations.

## Installation

```bash
go get github.com/rustyoz/Mtransform
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/rustyoz/Mtransform"
)

func main() {
    // Create a new transformation matrix
    t := mtransform.NewTransform()
    
    // Apply transformations
    t.Translate(10, 20)        // Move by (10, 20)
    t.Scale(2, 2)              // Scale by 2x
    t.RotateOrigin(math.Pi/4)  // Rotate 45 degrees around origin
    
    // Apply transformation to a point
    x, y := t.Apply(5, 5)
    fmt.Printf("Transformed point: (%.2f, %.2f)\n", x, y)
}
```

## API Reference

### Types

#### `Transform`
```go
type Transform [3][3]float64
```
A 3x3 transformation matrix for 2D transformations.

### Functions

#### `Identity() Transform`
Creates and returns an identity transformation matrix.

#### `NewTransform() *Transform`
Creates a new transformation matrix initialized to the identity matrix.

#### `MultiplyTransforms(a Transform, b Transform) Transform`
Multiplies two transformation matrices and returns the result.

### Methods

#### `Apply(x float64, y float64) (float64, float64)`
Applies the transformation to the given coordinates and returns the transformed point.

#### `MultiplyWith(b Transform)`
Multiplies the current transformation matrix with another transformation matrix in-place.

#### `Scale(x float64, y float64)`
Applies scaling transformation to the matrix.

#### `Translate(x float64, y float64)`
Applies translation transformation to the matrix.

#### `RotateOrigin(angle float64)`
Applies rotation around the origin. Angle is in radians.

#### `RotatePoint(angle float64, x float64, y float64)`
Applies rotation around a specific point. Angle is in radians.

#### `SkewX(angle float64)`
Applies skew transformation along the X-axis. Angle is in radians.

#### `SkewY(angle float64)`
Applies skew transformation along the Y-axis. Angle is in radians.

#### `Equals(t2 *Transform) bool`
Compares two transformation matrices for equality.

## Examples

### Combining Transformations

```go
t := mtransform.NewTransform()

// Chain multiple transformations
t.Translate(100, 50)       // Move to (100, 50)
t.Scale(1.5, 1.5)          // Scale up by 1.5x
t.RotateOrigin(math.Pi/6)  // Rotate 30 degrees

// Apply to multiple points
points := [][]float64{{0, 0}, {10, 10}, {20, 0}}
for _, point := range points {
    x, y := t.Apply(point[0], point[1])
    fmt.Printf("(%.1f, %.1f) -> (%.1f, %.1f)\n", point[0], point[1], x, y)
}
```

### Matrix Multiplication

```go
// Create two separate transformations
t1 := mtransform.NewTransform()
t1.Scale(2, 2)

t2 := mtransform.NewTransform()
t2.Translate(10, 10)

// Combine them
combined := mtransform.MultiplyTransforms(*t1, *t2)

// Or multiply in-place
t1.MultiplyWith(*t2)
```

### Rotation Around a Point

```go
t := mtransform.NewTransform()

// Rotate 90 degrees around point (50, 50)
t.RotatePoint(math.Pi/2, 50, 50)

// Apply to a point
x, y := t.Apply(60, 50)  // Point 10 units to the right of rotation center
fmt.Printf("Rotated point: (%.1f, %.1f)\n", x, y)  // Should be (50, 60)
```

## License

This project is licensed under the terms specified in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
