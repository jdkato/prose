// Package mathutil holds the handful of statistics prose's readability
// metrics need.
//
// These are small enough that depending on a statistics library for them costs
// more than it saves — a library this size should not pull in a transitive
// dependency to compute a mean.
package mathutil

import "math"

// Mean returns the arithmetic mean of xs, or NaN if xs is empty.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	var sum float64
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// StdDev returns the population standard deviation of xs, or NaN if xs is
// empty.
//
// Population rather than sample: readability formulas treat the text in hand
// as the whole population, not a draw from a larger one.
func StdDev(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	m := Mean(xs)
	var variance float64
	for _, x := range xs {
		variance += (x - m) * (x - m)
	}
	return math.Sqrt(variance / float64(len(xs)))
}

// Round rounds x to the given number of decimal places, with halves rounded
// away from zero.
//
// Note that this is not math.Round's behaviour scaled up: math.Round works to
// whole numbers only, and naively scaling it changes which way exact halves
// go for negative input.
func Round(x float64, places int) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x = -x
	}

	precision := math.Pow(10, float64(places))
	digit := x * precision

	_, decimal := math.Modf(digit)

	var rounded float64
	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	return rounded / precision * sign
}
