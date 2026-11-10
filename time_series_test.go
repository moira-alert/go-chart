package chart

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeSeriesGetValue(t *testing.T) {
	assert := assert.New(t)

	ts := TimeSeries{
		Name: "Test",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}

	x0, y0 := ts.GetValues(0)
	assert.NotZero(x0)
	assert.Equal(1.0, y0)
}

func TestTimeSeriesValidate(t *testing.T) {
	assert := assert.New(t)

	cs := TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	assert.Nil(cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
	}
	assert.NotNil(cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	assert.NotNil(cs.Validate())
}

func TestTimeSeriesGetValuesOutOfBounds(t *testing.T) {
	assert := assert.New(t)

	ts := TimeSeries{
		Name: "Test",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
		},
		YValues: []float64{1.0, 2.0, 3.0},
	}

	x, y := ts.GetValues(-1)
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)

	x, y = ts.GetValues(3)
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)

	x, y = ts.GetValues(100)
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)
}

func TestTimeSeriesGetFirstValuesEmpty(t *testing.T) {
	assert := assert.New(t)

	ts := TimeSeries{
		Name:    "Test",
		XValues: []time.Time{},
		YValues: []float64{},
	}

	x, y := ts.GetFirstValues()
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)
}

func TestTimeSeriesGetLastValuesEmpty(t *testing.T) {
	assert := assert.New(t)

	ts := TimeSeries{
		Name:    "Test",
		XValues: []time.Time{},
		YValues: []float64{},
	}

	x, y := ts.GetLastValues()
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)
}

func TestTimeSeriesGetLastValuesEmptyY(t *testing.T) {
	assert := assert.New(t)

	ts := TimeSeries{
		Name:    "Test",
		XValues: []time.Time{time.Now(), time.Now()},
		YValues: []float64{},
	}

	x, y := ts.GetLastValues()
	assert.Equal(0.0, x)
	assert.Equal(0.0, y)
}
