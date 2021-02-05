package wrap

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {
	Convey("TestContext", t, func() {
		ctx := NewContext(context.Background(),
			WithCtxDisableMetric(),
			WithCtxDisableMetric(),
			WithCtxMetricCustomLabelValue("hello"),
			WithCtxTraceTag("key1", "val1"),
			WithCtxTraceTag("key2", "val2"),
		)

		options := FromContext(ctx)

		So(options.DisableMetric, ShouldBeTrue)
		So(options.DisableTrace, ShouldBeTrue)
		So(options.MetricCustomLabelValue, ShouldEqual, "hello")
		So(options.TraceTags, ShouldResemble, map[string]string{
			"key1": "val1",
			"key2": "val2",
		})
	})
}
