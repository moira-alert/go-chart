package chart

import (
	"testing"

	"github.com/moira-alert/go-chart/matrix"
	"github.com/stretchr/testify/assert"
)

func TestPolynomialRegression(t *testing.T) {
	assert := assert.New(t)

	var xv []float64
	var yv []float64

	for i := 0; i < 100; i++ {
		xv = append(xv, float64(i))
		yv = append(yv, float64(i*i))
	}

	values := ContinuousSeries{
		XValues: xv,
		YValues: yv,
	}

	poly := &PolynomialRegressionSeries{
		InnerSeries: values,
		Degree:      2,
	}

	for i := 0; i < 100; i++ {
		_, y := poly.GetValues(i)
		assert.InDelta(float64(i*i), y, matrix.DefaultEpsilon)
	}
}

func TestPolynomialRegressionNilInnerSeries(t *testing.T) {
	assert := assert.New(t)

	poly := &PolynomialRegressionSeries{
		Degree: 2,
	}

	_, y := poly.GetValues(0)
	assert.Equal(0.0, y)

	x, y := poly.GetFirstValues()
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)

	x, y = poly.GetLastValues()
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)
}

func TestPolynomialRegressionEmptyInnerSeries(t *testing.T) {
	assert := assert.New(t)

	emptySeries := ContinuousSeries{
		XValues: []float64{},
		YValues: []float64{},
	}

	poly := &PolynomialRegressionSeries{
		InnerSeries: emptySeries,
		Degree:      2,
	}

	_, y := poly.GetValues(0)
	assert.Equal(0.0, y)

	x, y := poly.GetFirstValues()
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)
}
