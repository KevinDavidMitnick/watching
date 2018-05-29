// Copyright 2017 Xiaomi, Inc.
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

package cron

import (
	"encoding/json"
	"time"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func ReadHighEvent() {
	queues := g.Config().Redis.HighQueues
	if len(queues) == 0 {
		return
	}

	for {
		event, err := popEvent(queues)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		consume(event, true)
	}
}

func ReadLowEvent() {
	queues := g.Config().Redis.LowQueues
	if len(queues) == 0 {
		return
	}

	for {
		event, err := popEvent(queues)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		consume(event, false)
	}
}

func popEvent(queues []string) (*cmodel.Event, error) {

	count := len(queues)

	params := make([]interface{}, count+1)
	for i := 0; i < count; i++ {
		params[i] = queues[i]
	}
	// set timeout 0
	params[count] = 0

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	reply, err := redis.Strings(rc.Do("BRPOP", params...))
	if err != nil {
		log.Errorf("get alarm event from redis fail: %v", err)
		return nil, err
	}

	var event cmodel.Event
	err = json.Unmarshal([]byte(reply[1]), &event)
	if err != nil {
		log.Errorf("parse alarm event fail: %v", err)
		return nil, err
	}

	log.Debugf("pop event: %s", event.String())

	//insert event into database
	eventmodel.InsertEvent(&event)
	// events no longer saved in memory

	return &event, nil
}

func ReadStrategyGroupEvent() {
	for {
		consumeStrategyGroupEvent()
		time.Sleep(time.Second)
	}
}

func getReisKeyArray(strategyKey string) []string {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	var keys []string

	value, err := rc.Do("HKEYS", strategyKey)
	if err == nil && value != nil {
		for _, v := range value.([]interface{}) {
			data := v.([]byte)
			keys = append(keys, string(data))
		}
	}

	return keys
}

func consumeStrategyGroupEvent() {
	strategyGroup := g.StrategyGroupMap.Get()
	log.Printf("1.[DEBUG] start consume strategy group event,%v.........", strategyGroup)
	if strategyGroup != nil && len(strategyGroup) > 0 {
		rc := g.RedisConnPool.Get()
		defer rc.Close()
		for k, v := range strategyGroup {
			strategyKey := fmt.Sprintf("StrategyGroup_%d", k)
			keyArray := getReisKeyArray(strategyKey)
			log.Printf("2.[DEBUG] strategy group event key array is : %v", keyArray)
			for _, key := range keyArray {
				lastValue, err := rc.Do("HGET", strategyKey, key)
				if (err == nil) && (lastValue != nil) {
					if strValue, ok := lastValue.([]byte); ok {
						value := make(map[int]*cmodel.Event)
						if err := json.Unmarshal(strValue, &value); (err == nil) && (len(v) == len(value)) {
							pushMap := make(map[string][]string)
							cnt := 0
							var eventTime int64
							for _, id := range v {
								if tmpEvent, flag := value[id]; flag {
									if eventTime == 0 {
										eventTime = tmpEvent.EventTime
									}
									diffTime := tmpEvent.EventTime - eventTime
									if eventStr, err := json.Marshal(tmpEvent); diffTime < 60 && err == nil {
										cnt = cnt + 1
										redisKey := fmt.Sprintf("event:p%v", tmpEvent.Priority())
										pushMap[redisKey] = append(pushMap[redisKey], string(eventStr))
									}
								}
							}
							if cnt == len(v) {
								for k, v := range pushMap {
									for _, e := range v {
										rc.Do("LPUSH", k, e)
									}
								}
								log.Printf("3.[DEBUG] push map is  : %v", pushMap)
							}
							rc.Do("HDEL", strategyKey, key)
							log.Printf("4.[DEBUG] del strategy group %s,%s", strategyKey, key)
						}
					}
				}
			}
		}
	}
}
