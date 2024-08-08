package chart

import (
	"math"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func FuzzGeneratePrettyContinuousTicks(f *testing.F) {
	const (
		pngSize                    = 1024
		generatePrettyTicksTimeout = 5 * time.Second
	)

	font, err := GetDefaultFont()
	if err != nil {
		f.Errorf("failed to get default font: %s", err.Error())
	}

	r, err := PNG(pngSize, pngSize)
	if err != nil {
		f.Errorf("failed to create a new PNG: %s", err.Error())
	}

	r.SetFont(font)

	vf := FloatValueFormatter

	f.Fuzz(func(t *testing.T, min, max float64, domain int) {
		if min > max {
			min, max = max, min
		}

		ra := &ContinuousRange{
			Min:    min,
			Max:    max,
			Domain: domain,
		}

		enablePrettyTicks := true

		timer := time.NewTimer(generatePrettyTicksTimeout)
		defer timer.Stop()

		done := make(chan struct{})

		go func() {
			if allowGeneratePrettyContiniousTicks(enablePrettyTicks, ra) {
				GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
			}
			close(done)
		}()

		select {
		case <-timer.C:
			t.Errorf("Timeout on %#v\n", ra)
		case <-done:
		}
	})
}

func TestGenerateContinuousTicks(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    0.0,
		Max:    10.0,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
}

func TestGenerateContinuousTicksForVerySmallRange(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    0,
		Max:    1e-15,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
}

func TestGenerateContinuousTicksDescending(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:        0.0,
		Max:        10.0,
		Domain:     256,
		Descending: true,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
	assert.Len(ticks, 11)
	assert.Equal(10.0, ticks[0].Value)
	assert.Equal(9.0, ticks[1].Value)
	assert.Equal(1.0, ticks[len(ticks)-2].Value)
	assert.Equal(0.0, ticks[len(ticks)-1].Value)
}

func TestGenerateContinuousPrettyTicks(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    37.5,
		Max:    60.1,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
	assert.Equal(ticks, []Tick{
		{Label: "38.00", Value: 38},
		{Label: "40.00", Value: 40},
		{Label: "42.00", Value: 42},
		{Label: "44.00", Value: 44},
		{Label: "46.00", Value: 46},
		{Label: "48.00", Value: 48},
		{Label: "50.00", Value: 50},
		{Label: "52.00", Value: 52},
		{Label: "54.00", Value: 54},
		{Label: "56.00", Value: 56},
		{Label: "58.00", Value: 58},
		{Label: "60.00", Value: 60}})
}

func TestGenerateContinuousPrettyTicksForEmptyRange(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    1.0,
		Max:    1.0,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.Empty(ticks)

	ra = &ContinuousRange{
		Min:    1.0,
		Max:    2.0,
		Domain: 0,
	}

	ticks = GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.Empty(ticks)
}

func TestGeneratePrettyTicksForVerySmallRange(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    1e-100,
		Max:    1e-99,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
	assert.Len(ticks, 9)
}

func TestGeneratePrettyTicksForVeryLargeRange(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    1e-100,
		Max:    1e+100,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
	assert.Len(ticks, 10)
}

func TestGeneratePrettyTicksForVerySmallDomain(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    0.0,
		Max:    10.0,
		Domain: 1,
	}

	vf := FloatValueFormatter

	ticks := GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.Empty(ticks)
}

func TestGeneratePrettyTicksForVeryLargeDomain(t *testing.T) {
	assert := assert.New(t)

	f, err := GetDefaultFont()
	assert.Nil(err)

	r, err := PNG(1024, 1024)
	assert.Nil(err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    0.0,
		Max:    10.0,
		Domain: math.MaxInt32,
	}

	vf := FloatValueFormatter

	ticks := GeneratePrettyContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(ticks)
	assert.Len(ticks, 1001)
}

func TestAllowGeneratePrettyContiniousTicks(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	type args struct {
		enablePrettyTicks bool
		ra                Range
	}

	testcases := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Allow generate pretty continious ticks with correct parameters",
			args: args{
				enablePrettyTicks: true,
				ra: &ContinuousRange{
					Min: 0,
					Max: 100,
				},
			},
			expected: true,
		},
		{
			name: "Don't allow generate pretty continious ticks with disabled EnablePrettyTicks",
			args: args{
				enablePrettyTicks: false,
				ra: &ContinuousRange{
					Min: 0,
					Max: 100,
				},
			},
			expected: false,
		},
		{
			name: "Don't allow generate pretty continious ticks with too small a difference between max and min",
			args: args{
				enablePrettyTicks: true,
				ra: &ContinuousRange{
					Min: 0,
					Max: 1e-12,
				},
			},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actual := allowGeneratePrettyContiniousTicks(testcase.args.enablePrettyTicks, testcase.args.ra)
			assert.Equal(testcase.expected, actual)
		})
	}
}
