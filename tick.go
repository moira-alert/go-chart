package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/moira-alert/go-chart/util"
)

const prettyTicksTolerance = 1e-10

// TicksProvider is a type that provides ticks.
type TicksProvider interface {
	GetTicks(r Renderer, defaults Style, vf ValueFormatter) []Tick
}

// Tick represents a label on an axis.
type Tick struct {
	Value float64
	Label string
}

// Ticks is an array of ticks.
type Ticks []Tick

// Len returns the length of the ticks set.
func (t Ticks) Len() int {
	return len(t)
}

// Swap swaps two elements.
func (t Ticks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less returns if i's value is less than j's value.
func (t Ticks) Less(i, j int) bool {
	return t[i].Value < t[j].Value
}

// String returns a string representation of the set of ticks.
func (t Ticks) String() string {
	var values []string
	for i, tick := range t {
		values = append(values, fmt.Sprintf("[%d: %s]", i, tick.Label))
	}
	return strings.Join(values, ", ")
}

// GenerateContinuousTicks generates a set of ticks.
func GenerateContinuousTicks(r Renderer, ra Range, isVertical bool, style Style, vf ValueFormatter) []Tick {
	if vf == nil {
		vf = FloatValueFormatter
	}

	var ticks []Tick
	min, max := ra.GetMin(), ra.GetMax()

	if ra.IsDescending() {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	} else {
		ticks = append(ticks, Tick{
			Value: min,
			Label: vf(min),
		})
	}

	minLabel := vf(min)
	style.GetTextOptions().WriteToRenderer(r)
	labelBox := r.MeasureText(minLabel)

	var tickSize float64
	if isVertical {
		tickSize = float64(labelBox.Height() + DefaultMinimumTickVerticalSpacing)
	} else {
		tickSize = float64(labelBox.Width() + DefaultMinimumTickHorizontalSpacing)
	}

	domain := float64(ra.GetDomain())
	domainRemainder := domain - (tickSize * 2)
	intermediateTickCount := int(math.Floor(float64(domainRemainder) / float64(tickSize)))

	rangeDelta := math.Abs(max - min)
	tickStep := rangeDelta / float64(intermediateTickCount)

	roundTo := util.Math.GetRoundToForDelta(rangeDelta) / 10
	intermediateTickCount = util.Math.MinInt(intermediateTickCount, DefaultTickCountSanityCheck)

	for x := 1; x < intermediateTickCount; x++ {
		var tickValue float64
		if ra.IsDescending() {
			tickValue = max - util.Math.RoundUp(tickStep*float64(x), roundTo)
		} else {
			tickValue = min + util.Math.RoundUp(tickStep*float64(x), roundTo)
		}
		ticks = append(ticks, Tick{
			Value: tickValue,
			Label: vf(tickValue),
		})
	}

	if ra.IsDescending() {
		ticks = append(ticks, Tick{
			Value: min,
			Label: vf(min),
		})
	} else {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	}

	return ticks
}

// PrettyTicker is an interface that defines a method for checking whether pretty ticks should be enabled.
type PrettyTicker interface {
	GetEnablePrettyTicks() bool
}

// allowGeneratePrettyContiniousTicks is a method that determines whether the GeneratePrettyContiniousTicks
// function can be called, which does a lot of maths transformations that don't work on large floats.
func allowGeneratePrettyContiniousTicks(enablePrettyTicks bool, ra Range) bool {
	return enablePrettyTicks && math.Abs(ra.GetMax()-ra.GetMin()) > prettyTicksTolerance
}

// GeneratePrettyContinuousTicks generates a set of ticks at visually pleasing intervals.
// Based on http://vis.stanford.edu/files/2010-TickLabels-InfoVis.pdf by Justin Talbot et. al.
func GeneratePrettyContinuousTicks(r Renderer, ra Range, isVertical bool, style Style, vf ValueFormatter) []Tick {
	if vf == nil {
		vf = FloatValueFormatter
	}

	prettyStepsPriorityList := []float64{1, 5, 2, 2.5, 4, 3}

	const (
		simplicityParam = 0.2
		coverageParam   = 0.25
		densityParam    = 0.5
		legibilityParam = 0.05
	)

	rangeMin, rangeMax := ra.GetMin(), ra.GetMax()

	if rangeMin >= rangeMax || ra.GetDomain() == 0 {
		return []Tick{}
	}

	renderedLabelExample := vf(rangeMin)
	style.GetTextOptions().WriteToRenderer(r)
	renderedLabelSizePx := r.MeasureText(renderedLabelExample)

	var actualLabelSizePx, desiredPaddedLabelSizePx float64
	if isVertical {
		actualLabelSizePx = math.Max(float64(renderedLabelSizePx.Height()), 1)
		desiredPaddedLabelSizePx = actualLabelSizePx + DefaultMinimumTickVerticalSpacing
	} else {
		actualLabelSizePx = math.Max(float64(renderedLabelSizePx.Width()), 1)
		desiredPaddedLabelSizePx = actualLabelSizePx + DefaultMinimumTickHorizontalSpacing
	}
	availableSpacePx := float64(ra.GetDomain())
	desiredTicksCount := math.Min(
		math.Max(math.Floor(availableSpacePx/desiredPaddedLabelSizePx), 2), // less than 2 leads to incorrect density calculation
		DefaultTickCountSanityCheck)

	prettyStepsCount := float64(len(prettyStepsPriorityList))

	var bestTickMin, bestTickMax, bestTickStep float64
	bestScore := -2.0
	stepsToSkip := 1.0

OUTER:
	for {
		for prettyStepIndex, prettyStep := range prettyStepsPriorityList {
			simplicityMax := calculateSimplicityMax(float64(prettyStepIndex), prettyStepsCount, stepsToSkip)

			if simplicityParam*simplicityMax+
				coverageParam+
				densityParam+
				legibilityParam < bestScore {
				break OUTER
			}

			ticksCount := 2.0
			for {
				densityMax := calculateDensityMax(ticksCount, desiredTicksCount)

				if simplicityParam*simplicityMax+
					coverageParam+
					densityParam*densityMax+
					legibilityParam < bestScore {
					break
				}

				delta := (rangeMax - rangeMin) / (ticksCount + 1) / stepsToSkip / prettyStep
				stepSizeMultiplierLog := math.Ceil(math.Log10(delta))

				for {
					tickStep := stepsToSkip * prettyStep * math.Pow(10, stepSizeMultiplierLog)
					coverageMax := calculateCoverageMax(rangeMin, rangeMax, tickStep*(ticksCount-1))

					if simplicityParam*simplicityMax+
						coverageParam*coverageMax+
						densityParam*densityMax+
						legibilityParam < bestScore {
						break
					}

					minStart := math.Floor(rangeMax/tickStep)*stepsToSkip - (ticksCount-1)*stepsToSkip
					maxStart := math.Ceil(rangeMin/tickStep) * stepsToSkip

					if minStart > maxStart {
						stepSizeMultiplierLog += 1
						continue
					}

					for start := minStart; start <= maxStart; start++ {
						tickMin := start * (tickStep / stepsToSkip)
						tickMax := tickMin + tickStep*(ticksCount-1)

						coverage := calculateCoverage(rangeMin, rangeMax, tickMin, tickMax)
						simplicity := calculateSimplicity(prettyStepsCount, float64(prettyStepIndex), stepsToSkip, tickMin, tickMax, tickStep)
						density := calculateDensity(ticksCount, desiredTicksCount, rangeMin, rangeMax, tickMin, tickMax)
						legibility := 1.0 // format is out of our control (provided by ValueFormatter)
						// font size is out of our control (provided by Style)
						// orientation is out of our control
						if actualLabelSizePx*ticksCount > availableSpacePx {
							legibility = math.Inf(-1) // overlap is unacceptable
						}

						score := simplicityParam*simplicity +
							coverageParam*coverage +
							densityParam*density +
							legibilityParam*legibility

						// original algorithm allows ticks outside value range, but it breaks rendering in this library
						if score > bestScore && tickMin >= rangeMin && tickMax <= rangeMax {
							bestTickMin = tickMin
							bestTickMax = tickMax
							bestTickStep = tickStep
							bestScore = score
						}
					}
					stepSizeMultiplierLog++
				}
				ticksCount++
			}
		}
		stepsToSkip++
	}

	var ticks []Tick
	if bestTickStep == 0 {
		return ticks
	}

	if ra.IsDescending() {
		for tickValue := bestTickMax; tickValue > bestTickMin-bestTickStep/2; tickValue -= bestTickStep {
			ticks = append(ticks, Tick{
				Value: tickValue,
				Label: vf(tickValue),
			})
		}
	} else {
		for tickValue := bestTickMin; tickValue < bestTickMax+bestTickStep/2; tickValue += bestTickStep {
			ticks = append(ticks, Tick{
				Value: tickValue,
				Label: vf(tickValue),
			})
		}
	}

	return ticks
}

func calculateSimplicity(prettyStepsCount, prettyStepIndex, stepsToSkip, tickMin, tickMax, tickStep float64) float64 {
	var hasZeroTick float64
	if tickMin <= 0 && tickMax >= 0 && math.Mod(tickMin, tickStep) < 10e-10 {
		hasZeroTick = 1
	}

	return 1 - prettyStepIndex/(prettyStepsCount-1) - stepsToSkip + hasZeroTick
}

func calculateSimplicityMax(prettyStepIndex, prettyStepsCount, stepsToSkip float64) float64 {
	return 2 - prettyStepIndex/(prettyStepsCount-1) - stepsToSkip
}

func calculateCoverage(rangeMin, rangeMax, tickMin, tickMax float64) float64 {
	return 1 - 0.5*(util.Math.SqrFloat(rangeMax-tickMax)+util.Math.SqrFloat(rangeMin-tickMin))/util.Math.SqrFloat(0.1*(rangeMax-rangeMin))
}

func calculateCoverageMax(rangeMin, rangeMax, span float64) float64 {
	if span <= rangeMax-rangeMin {
		return 1
	}

	return 1 - util.Math.SqrFloat((rangeMax-rangeMin)/2)/util.Math.SqrFloat(0.1*(rangeMax-rangeMin))
}

func calculateDensity(ticksCount, desiredTicksCount, rangeMin, rangeMax, tickMin, tickMax float64) float64 {
	ticksDensity := (ticksCount - 1) / (tickMax - tickMin)
	desiredTicksDensity := (desiredTicksCount - 1) / (math.Max(tickMax, rangeMax) - math.Min(tickMin, rangeMin))

	return 2 - math.Max(ticksDensity/desiredTicksDensity, desiredTicksDensity/ticksDensity)
}

func calculateDensityMax(ticksCount, desiredTicksCount float64) float64 {
	if ticksCount >= desiredTicksCount {
		return 2 - (ticksCount-1)/(desiredTicksCount-1)
	}

	return 1
}
