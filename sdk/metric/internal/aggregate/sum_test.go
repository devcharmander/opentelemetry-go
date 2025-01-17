// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestSum(t *testing.T) {
	t.Cleanup(mockTime(now))

	t.Run("Int64/DeltaSum", testDeltaSum[int64]())
	t.Run("Float64/DeltaSum", testDeltaSum[float64]())

	t.Run("Int64/CumulativeSum", testCumulativeSum[int64]())
	t.Run("Float64/CumulativeSum", testCumulativeSum[float64]())

	t.Run("Int64/DeltaPrecomputedSum", testDeltaPrecomputedSum[int64]())
	t.Run("Float64/DeltaPrecomputedSum", testDeltaPrecomputedSum[float64]())

	t.Run("Int64/CumulativePrecomputedSum", testCumulativePrecomputedSum[int64]())
	t.Run("Float64/CumulativePrecomputedSum", testCumulativePrecomputedSum[float64]())
}

func testDeltaSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality: metricdata.DeltaTemporality,
		Filter:      attrFltr,
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, alice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      3,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Delta sums are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
	})
}

func testCumulativeSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality: metricdata.CumulativeTemporality,
		Filter:      attrFltr,
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, alice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      14,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      -8,
						},
					},
				},
			},
		},
	})
}

func testDeltaPrecomputedSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality: metricdata.DeltaTemporality,
		Filter:      attrFltr,
	}.PrecomputedSum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice},
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      7,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      14,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Precomputed sums are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
	})
}

func testCumulativePrecomputedSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality: metricdata.CumulativeTemporality,
		Filter:      attrFltr,
	}.PrecomputedSum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice},
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      11,
						},
						{
							Attributes: fltrBob,
							StartTime:  staticTime,
							Time:       staticTime,
							Value:      3,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Precomputed sums are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
	})
}

func BenchmarkSum(b *testing.B) {
	// The monotonic argument is only used to annotate the Sum returned from
	// the Aggregation method. It should not have an effect on operational
	// performance, therefore, only monotonic=false is benchmarked here.
	b.Run("Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.Sum(false)
	}))
	b.Run("Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.Sum(false)
	}))
	b.Run("Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.Sum(false)
	}))
	b.Run("Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.Sum(false)
	}))

	b.Run("Precomputed/Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.PrecomputedSum(false)
	}))
	b.Run("Precomputed/Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.PrecomputedSum(false)
	}))
	b.Run("Precomputed/Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.PrecomputedSum(false)
	}))
	b.Run("Precomputed/Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.PrecomputedSum(false)
	}))
}
