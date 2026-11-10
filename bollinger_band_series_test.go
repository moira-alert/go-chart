package chart

import (
	"fmt"
	"math"
	"testing"

	"github.com/moira-alert/go-chart/seq"
	"github.com/stretchr/testify/assert"
)

func TestBollingerBandSeries(t *testing.T) {
	s1 := mockValuesProvider{
		X: seq.Range(1.0, 100.0),
		Y: seq.RandomValuesWithMax(100, 1024),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	xvalues := make([]float64, 100)
	y1values := make([]float64, 100)
	y2values := make([]float64, 100)

	for x := 0; x < 100; x++ {
		xvalues[x], y1values[x], y2values[x] = bbs.GetBoundedValues(x)
	}

	for x := bbs.GetPeriod(); x < 100; x++ {
		assert.True(t, y1values[x] > y2values[x], fmt.Sprintf("%v vs. %v", y1values[x], y2values[x]))
	}
}

func TestBollingerBandLastValue(t *testing.T) {
	s1 := mockValuesProvider{
		X: seq.Range(1.0, 100.0),
		Y: seq.Range(1.0, 100.0),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	x, y1, y2 := bbs.GetBoundedLastValues()
	assert.Equal(t, 100.0, x)
	assert.Equal(t, 101, math.Floor(y1))
	assert.Equal(t, 83, math.Floor(y2))
}

func TestBollingerBandsSeriesConcurrentAccess(t *testing.T) {
	assert := assert.New(t)

	s1 := mockValuesProvider{
		X: seq.Range(1.0, 100.0),
		Y: seq.RandomValuesWithMax(100, 1024),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	done := make(chan bool, 10)
	for g := 0; g < 10; g++ {
		go func() {
			for x := 0; x < 100; x++ {
				bbs.GetBoundedValues(x)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	assert.NotNil(bbs.valueBuffer)
}

func TestBollingerBandsSeriesNilInnerSeries(t *testing.T) {
	assert := assert.New(t)

	bbs := &BollingerBandsSeries{}

	x, y1, y2 := bbs.GetBoundedValues(0)
	assert.Equal(t, 0.0, x)
	assert.Equal(t, 0.0, y1)
	assert.Equal(t, 0.0, y2)

	x, y1, y2 = bbs.GetBoundedLastValues()
	assert.Equal(t, 0.0, x)
	assert.Equal(t, 0.0, y1)
	assert.Equal(t, 0.0, y2)
}
