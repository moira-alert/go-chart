package chart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastValueAnnotation(t *testing.T) {
	assert := assert.New(t)

	series := ContinuousSeries{
		XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		YValues: []float64{5.0, 3.0, 3.0, 2.0, 1.0},
	}

	lva := LastValueAnnotation(series)
	assert.NotEmpty(lva.Annotations)
	lvaa := lva.Annotations[0]
	assert.Equal(5, lvaa.XValue)
	assert.Equal(1, lvaa.YValue)
}

func TestLastValueAnnotationEmptySeries(t *testing.T) {
	assert := assert.New(t)

	emptySeries := ContinuousSeries{
		XValues: []float64{},
		YValues: []float64{},
	}

	lva := LastValueAnnotation(emptySeries)
	assert.Empty(lva.Annotations)
}

func TestLastValueAnnotationNilLikeSeries(t *testing.T) {
	assert := assert.New(t)

	series := ContinuousSeries{}

	lva := LastValueAnnotation(series)
	assert.Empty(lva.Annotations)
}
