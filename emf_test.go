// Copyright (c) 2020 Coda Solutions Ltd
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package emf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDimension(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddDimension("FunctionVersion", "$LATEST", "Environment", "Production")

	assert.Len(t, cwm.dimensions, 1)
	assert.Equal(t, "$LATEST", cwm.dimensions[0]["FunctionVersion"])
	assert.Equal(t, "Production", cwm.dimensions[0]["Environment"])
}

func TestDimensionWithoutValueIsDefaultedToNil(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddDimension("FunctionVersion", "$LATEST", "Environment")

	assert.Len(t, cwm.dimensions, 1)
	assert.Equal(t, "$LATEST", cwm.dimensions[0]["FunctionVersion"])
	assert.Equal(t, "", cwm.dimensions[0]["Environment"])
}

func TestAddMultipleDimensions(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddDimension("FunctionVersion", "$LATEST")
	cwm.AddDimension("Region", "Ireland")

	assert.Len(t, cwm.dimensions, 2)
	assert.Equal(t, "$LATEST", cwm.dimensions[0]["FunctionVersion"])
	assert.Equal(t, "Ireland", cwm.dimensions[1]["Region"])
}

func TestMetricSetsTargetMember(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddMetric("ExecutionTime", Seconds, 1)

	assert.Contains(t, cwm.metrics, "ExecutionTime")
	assert.Len(t, cwm.metrics["ExecutionTime"].values, 1)
	assert.Equal(t, float64(1), cwm.metrics["ExecutionTime"].values[0])
	assert.Equal(t, Seconds, cwm.metrics["ExecutionTime"].unit)
}

func TestMetricAppendsToExistingTargetMember(t *testing.T) {
	times := []float64{1, 3, 5, 6, 8}

	cwm := NewMetric("test-lambda-metrics")
	cwm.AddMetric("ExecutionTimes", Seconds, times[:1]...)
	cwm.AddMetric("ExecutionTimes", Seconds, times[1:]...)

	assert.Len(t, cwm.metrics["ExecutionTimes"].values, len(times))
	assert.Equal(t, times, cwm.metrics["ExecutionTimes"].values)
}

func TestMetricUnitCannotBeChangedOnceSet(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddMetric("ExecutionTimes", Seconds, 1)
	cwm.AddMetric("ExecutionTimes", Microseconds, 100)

	assert.Equal(t, Seconds, cwm.metrics["ExecutionTimes"].unit)
}

func TestPropertiesAreSet(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddProperties("requestID", "9d0ff7b8-be31-4f4f-a301-765bc975ad61", "functionVersion", "$LATEST")

	assert.Contains(t, cwm.properties, "requestID")
	assert.Equal(t, "9d0ff7b8-be31-4f4f-a301-765bc975ad61", cwm.properties["requestID"])
	assert.Contains(t, cwm.properties, "functionVersion")
	assert.Equal(t, "$LATEST", cwm.properties["functionVersion"])
}

func TestPropertyWithoutValueIsDefaultedToNil(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddProperties("requestID", "9d0ff7b8-be31-4f4f-a301-765bc975ad61", "functionVersion")

	assert.Contains(t, cwm.properties, "functionVersion")
	assert.Equal(t, "", cwm.properties["functionVersion"])
}

func TestDuplicatePropertyOverridesValue(t *testing.T) {
	cwm := NewMetric("test-lambda-metrics")
	cwm.AddProperties("requestID", "9d0ff7b8-be31-4f4f-a301-765bc975ad61", "functionVersion", "$LATEST")
	cwm.AddProperties("functionVersion", "1")

	assert.Equal(t, "1", cwm.properties["functionVersion"])
}
