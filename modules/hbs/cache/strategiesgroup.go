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

package cache

import (
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"sync"
)

// 做一个map缓存策略组与策略的对应关系
type SafeGroupStrategiesMap struct {
	sync.RWMutex
	M map[int][]int
}

var GroupsStrategiesMap = &SafeGroupStrategiesMap{M: make(map[int][]int)}

func (this *SafeGroupStrategiesMap) GetGroupStrategies(sgid int) ([]int, bool) {
	this.RLock()
	defer this.RUnlock()
	strategiesIds, exists := this.M[sgid]
	return strategiesIds, exists
}

func (this *SafeGroupStrategiesMap) GetMap() map[int][]int {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeGroupStrategiesMap) Init() {
	m, err := db.QueryGroupStrategies()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
