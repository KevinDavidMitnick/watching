package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"strings"
)

func NewMetricValue(TS int64, endpoint, metric string, val interface{}, dataType string, tags ...string) *model.MetricValue {
	sec := int64(g.Config().Transfer.Interval)

	mv := model.MetricValue{
		Metric:    metric,
		Value:     val,
		Type:      dataType,
		Endpoint:  endpoint,
		Step:      sec,
		Timestamp: TS,
	}
	if len(tags) > 0 {
		mv.Tags = strings.Join(tags, ",")
	}
	return &mv
}
