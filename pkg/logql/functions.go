package logql

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/prometheus/prometheus/promql"
)

func getAggregationOverTimeFn(expr *rangeAggregationExpr) (RangeVectorAggregator, error) {
	switch expr.operation {
	case OpRangeTypeRate:
		return rate(expr.left.interval), nil
	case OpRangeTypeCount:
		return countOverTime, nil
	case OpRangeTypeAvg:
		return avgOverTime, nil
	case OpRangeTypeMax:
		return maxOverTime, nil
	case OpRangeTypeMin:
		return minOverTime, nil
	case OpRangeTypeStddev:
		return stddevOverTime, nil
	case OpRangeTypeStdvar:
		return stdvarOverTime, nil
	case OpRangeTypeSum:
		return sumOverTime, nil
	case OpRangeTypeQuantile:
		return quantileOverTime(expr.param), nil
	default:
		return nil, fmt.Errorf("unsupported aggregation over time: %s", expr.operation)
	}
}

// rate calculate the per-second rate of log lines.
func rate(selRange time.Duration) func(ts int64, samples []promql.Point) float64 {
	return func(ts int64, samples []promql.Point) float64 {
		return float64(len(samples)) / selRange.Seconds()
	}
}

// countOverTime counts the amount of log lines.
func countOverTime(ts int64, samples []promql.Point) float64 {
	return float64(len(samples))
}

func avgOverTime(ts int64, samples []promql.Point) float64 {
	var mean, count float64
	for _, v := range samples {
		count++
		mean += (v.V - mean) / count
	}
	return mean
}

func maxOverTime(ts int64, samples []promql.Point) float64 {
	max := samples[0].V
	for _, v := range samples {
		if v.V > max || math.IsNaN(max) {
			max = v.V
		}
	}
	return max
}

func minOverTime(ts int64, samples []promql.Point) float64 {
	min := samples[0].V
	for _, v := range samples {
		if v.V < min || math.IsNaN(min) {
			min = v.V
		}
	}
	return min
}

func sumOverTime(ts int64, samples []promql.Point) float64 {
	var sum float64
	for _, v := range samples {
		sum += v.V
	}
	return sum
}

func stddevOverTime(ts int64, samples []promql.Point) float64 {
	var aux, count, mean float64
	for _, v := range samples {
		count++
		delta := v.V - mean
		mean += delta / count
		aux += delta * (v.V - mean)
	}
	return math.Sqrt(aux / count)
}

func stdvarOverTime(ts int64, samples []promql.Point) float64 {
	var aux, count, mean float64
	for _, v := range samples {
		count++
		delta := v.V - mean
		mean += delta / count
		aux += delta * (v.V - mean)
	}
	return aux / count
}

func quantileOverTime(q float64) func(ts int64, samples []promql.Point) float64 {
	return func(ts int64, samples []promql.Point) float64 {
		values := make(vectorByValueHeap, 0, len(samples))
		for _, v := range samples {
			values = append(values, promql.Sample{Point: promql.Point{V: v.V}})
		}
		if len(values) == 0 {
			return math.NaN()
		}
		if q < 0 {
			return math.Inf(-1)
		}
		if q > 1 {
			return math.Inf(+1)
		}
		sort.Sort(values)

		n := float64(len(values))
		// When the quantile lies between two samples,
		// we use a weighted average of the two samples.
		rank := q * (n - 1)

		lowerIndex := math.Max(0, math.Floor(rank))
		upperIndex := math.Min(n-1, lowerIndex+1)

		weight := rank - math.Floor(rank)
		return values[int(lowerIndex)].V*(1-weight) + values[int(upperIndex)].V*weight
	}
}
