package cron

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"log"
	"time"
)

func SyncStrategyGroup() {
	duration := time.Duration(g.Config().Hbs.Interval) * time.Second
	for {
		syncStrategyGroup()
		time.Sleep(duration)
	}
}

func syncStrategyGroup() {
	var strategyGroupResponse model.GroupStrigesResponse
	err := g.HbsClient.Call("Hbs.GetGroupStrategies", model.NullRpcRequest{}, &strategyGroupResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetGroupStrategies:", err)
		return
	}

	rebuildStrategyGroupMap(&strategyGroupResponse)
}

func rebuildStrategyGroupMap(strategyGroupResponse *model.GroupStrigesResponse) {
	m := strategyGroupResponse.GroupStriges
	g.StrategyGroupMap.ReInit(m)
}
