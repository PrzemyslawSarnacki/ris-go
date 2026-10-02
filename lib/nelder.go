/*
Package ris provides tools to fit a generalized logistic model to weightlifting data
and compute the RIS (Relative Intensity Score) index.
*/
package ris

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat"
)

// FitRISParamsNelder fits the generalized logistic model to data, returning FitResult.
// normalizer: e.g., 100 for index normalization.
func FitRISParamsNelder(data []DataPoint, normalizer float64) (FitResult, error) {
	if len(data) < 5 {
		return FitResult{}, ErrNotEnoughData
	}

	sortedData := make([]DataPoint, len(data))
	copy(sortedData, data)
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].BodyWeight < sortedData[j].BodyWeight
	})

	x := make([]float64, len(sortedData))
	y := make([]float64, len(sortedData))
	for i, dp := range sortedData {
		x[i] = dp.BodyWeight
		y[i] = dp.Total
	}

	// estimate A: min(y) - (max(y)-min(y))/3
	minY, maxY := floats.Min(y), floats.Max(y)
	minX, maxX := floats.Min(x), floats.Max(x)
	A := minY - (maxY-minY)*(4.0/3.0)*0.25

	// initial parameter guesses: A, K, B, V, Q
	medianX := x[len(x)/2]
	init := []float64{maxY * 1.05, 0.05, medianX, 1.0}

	// Define least-squares objective with soft penalty constraints for bounds
	problem := optimize.Problem{
		Func: func(params []float64) float64 {
			K, B, V, Q := params[0], params[1], params[2], params[3]

			// Penalty checks to prevent degenerate flat curves (e.g. Q -> 0, V drifting out of range)
			penalty := 0.0
			if Q < 0.01 {
				penalty += math.Pow(0.01-Q, 2) * 1e6
			}
			if V < minX || V > maxX {
				if V < minX {
					penalty += math.Pow(minX-V, 2) * 1e4
				} else {
					penalty += math.Pow(V-maxX, 2) * 1e4
				}
			}
			if B < 0.001 || B > 1.0 {
				penalty += 1e5
			}
			if K < maxY*0.8 || K > maxY*1.5 {
				penalty += 1e5
			}

			sum := 0.0
			p := Params{A: A, K: K, B: B, V: V, Q: Q}
			for i := range x {
				d := GeneralizedLogistic(x[i], p) - y[i]
				sum += d * d
			}
			return sum + penalty
		},
	}

	// Perform optimization using Nelder-Mead
	settings := optimize.Settings{GradientThreshold: 1e-8, FuncEvaluations: 20000}
	method := &optimize.NelderMead{}
	result, err := optimize.Minimize(problem, init, &settings, method)
	if err != nil {
		return FitResult{}, err
	}

	opt := result.X
	params := Params{A: A, K: opt[0], B: opt[1], V: opt[2], Q: opt[3]}

	var sse float64
	for i := range x {
		diff := GeneralizedLogistic(x[i], params) - y[i]
		sse += diff * diff
	}
	rmse := math.Sqrt(sse / float64(len(x)))

	// compute RIS*Total for each data point and fit linear model
	scores := make([]float64, len(data))
	for i := range data {
		idx := params.Inverse(data[i].BodyWeight, normalizer)
		scores[i] = idx * data[i].Total
	}

	xOrig := make([]float64, len(data))
	for i, dp := range data {
		xOrig[i] = dp.BodyWeight
	}

	slope, intercept := stat.LinearRegression(xOrig, scores, nil, false)

	return FitResult{
		Params:        params,
		LineSlope:     slope,
		LineIntercept: intercept,
		RMSE:          rmse,
	}, nil
}
