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

// Unit represents a unit of measure and is associated with a metric value
type Unit string

// A unit of measure to be associated with a metric value
const (
	None            Unit = "None"
	Seconds         Unit = "Seconds"
	Microseconds    Unit = "Microseconds"
	Bytes           Unit = "Bytes"
	Kilobytes       Unit = "Kilobytes"
	Megabytes       Unit = "Megabytes"
	Gigabytes       Unit = "Gigabyte"
	Terabits        Unit = "Terabits"
	Percent         Unit = "Percent"
	Count           Unit = "Count"
	BytesSecond     Unit = "Bytes/Second"
	KilobytesSecond Unit = "Kilobytes/Second"
	MegabytesSecond Unit = "Megabytes/Second"
	GigabytesSecond Unit = "Gigabytes/Second"
	TerabytesSecond Unit = "Terabytes/Second"
	BitsSecond      Unit = "Bits/Second"
	KilobitsSecond  Unit = "Kilobits/Second"
	MegabitsSecond  Unit = "Megabits/Second"
	GigabitsSecond  Unit = "Gigabits/Second"
	TerabitsSecond  Unit = "Terabits/Second"
	CountSecond     Unit = "Count/Second"
)

// CloudWatchMetric defines the CloudWatch embedded metric format used to instruct
// CloudWatch Logs to automatically extract metric values embedded in structured log events.
// The embedded metric format allows you can use CloudWatch to graph and create alarms on
// the extracted metric values for real-time incident detection.
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch_Embedded_Metric_Format_Specification.html
type CloudWatchMetric struct {
	namespace  string
	dimensions []map[string]string
	metrics    map[string]metric
	properties map[string]string
}

type metric struct {
	key    string
	values []float64
	unit   Unit
}

// NewMetric creates a new namespaced blank CloudWatch Embedded Format metric
// ready for capturing metrics
func NewMetric(namespace string) CloudWatchMetric {
	return CloudWatchMetric{
		namespace:  namespace,
		dimensions: []map[string]string{},
		metrics:    map[string]metric{},
		properties: map[string]string{},
	}
}

// AddDimension will associate a new dimension set with all of the current metric
// values. It is important to note, that CloudWatch treats a unique key/value
// combination as a separate metric. If the cardinality of a particular dimension
// value is expected to be high, it is advisable to use a property instead.
//
// A dimension set is restricted to 9 keys in length, and will automatically
// be truncated when the metric is serialised. If a dimension key is provided
// without a value it be defaulted to nil
func (em *CloudWatchMetric) AddDimension(kv ...string) {
	set := map[string]string{}

	var l int
	if l = len(kv); l%2 != 0 {
		l--
		// A dimension has been provided with no value. Default it to nil
		set[kv[l]] = ""
	}

	for i := 0; i < l; i += 2 {
		set[kv[i]] = kv[i+1]
	}

	em.dimensions = append(em.dimensions, set)
}

// AddMetric will either add or update an existing metric. A metric is uniquely
// defined by its key and a corresponding unit should be provided describing its
// data type. If the keyed metric already exists, the values will simple be appended.
// The CloudWatch Embedded Metric Format supports a maximum of 150 metrics, each
// with a maximum of 100 values.
func (em *CloudWatchMetric) AddMetric(key string, unit Unit, values ...float64) {
	if m, ok := em.metrics[key]; ok {
		m.values = append(m.values, values...)
		em.metrics[key] = m
	} else {
		em.metrics[key] = metric{
			key:    key,
			unit:   unit,
			values: values,
		}
	}
}

// AddProperties allows additional properties to be appended to the root of the
// CloudWatch metric JSON object, but these will not be submitted to CloudWatch Metrics.
// However these additional properties will be searchable by CloudWatch Logs Insights.
// This is useful for contextual and potentially high-cardinality data that is
// not appropriate for CloudWatch Metrics dimensions.
//
// All properties are composed of a key value pair. And any number of properties
// can be set at once. If a property with the same key, is defined more than once,
// it will always be overwritten with the latest value. If a property is provided
// without a corresponding value, its value will be defaulted to nil. This can have
// the side effect of properties being omitted when marshalling to JSON.
func (em *CloudWatchMetric) AddProperties(props ...string) {
	var l int
	if l = len(props); l%2 != 0 {
		l--
		// A property has been provided with no value. Default it to nil
		em.properties[props[l]] = ""
	}

	for i := 0; i < l; i += 2 {
		em.properties[props[i]] = props[i+1]
	}
}
