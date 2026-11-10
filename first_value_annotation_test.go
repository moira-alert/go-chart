package chart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstValueAnnotation(t *testing.T) {
	assert := assert.New(t)

	series := ContinuousSeries{
		XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		YValues: []float64{5.0, 3.0, 3.0, 2.0, 1.0},
	}

	fva := FirstValueAnnotation(series)
	assert.NotEmpty(fva.Annotations)
	fvaa := fva.Annotations[0]
	assert.Equal(1, fvaa.XValue)
	assert.Equal(5, fvaa.YValue)
}

func TestFirstValueAnnotationEmptySeries(t *testing.T) {
	assert := assert.New(t)

	emptySeries := ContinuousSeries{
		XValues: []float64{},
		YValues: []float64{},
	}

	fva := FirstValueAnnotation(emptySeries)
	assert.Empty(fva.Annotations)
}

func TestFirstValueAnnotationNilLikeSeries(t *testing.T) {
	assert := assert.New(t)

	series := ContinuousSeries{}

	fva := FirstValueAnnotation(series)
	assert.Empty(fva.Annotations)
}
