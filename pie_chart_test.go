package chart

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPieChart(t *testing.T) {
	assert := assert.New(t)

	pie := PieChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 10, Label: "Blue"},
			{Value: 9, Label: "Green"},
			{Value: 8, Label: "Gray"},
			{Value: 7, Label: "Orange"},
			{Value: 6, Label: "HEANG"},
			{Value: 5, Label: "??"},
			{Value: 2, Label: "!!"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	pie.Render(PNG, b)
	assert.NotZero(b.Len())
}

func TestPieChartDropsZeroValues(t *testing.T) {
	assert := assert.New(t)

	pie := PieChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 0, Label: "Gray"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	assert.Nil(err)
}

func TestPieChartAllZeroValues(t *testing.T) {
	assert := assert.New(t)

	pie := PieChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 0, Label: "Blue"},
			{Value: 0, Label: "Green"},
			{Value: 0, Label: "Gray"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	assert.NotNil(err)
}

func TestPieChartGetHeight(t *testing.T) {
	assert := assert.New(t)

	pie := PieChart{}
	assert.Equal(DefaultChartHeight, pie.GetHeight())

	pie2 := PieChart{Height: 600}
	assert.Equal(600, pie2.GetHeight())

	pie3 := PieChart{Height: 0}
	assert.Equal(DefaultChartHeight, pie3.GetHeight())
}

func TestPieChartGetWidth(t *testing.T) {
	assert := assert.New(t)

	pie := PieChart{}
	assert.Equal(DefaultChartWidth, pie.GetWidth())

	pie2 := PieChart{Width: 800}
	assert.Equal(800, pie2.GetWidth())
}
