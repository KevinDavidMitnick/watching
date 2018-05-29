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

package db

import (
	"log"
)

// 一个策略组ID对应了多个策略ID
func QueryGroupStrategies() (map[int][]int, error) {
	ret := make(map[int][]int)
	rows, err := DB.Query("select a.id, b.id from strategy_group as a inner join strategy as b on a.id=b.sgrp_id")
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		var sgid, sid int

		err = rows.Scan(&sgid, &sid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if _, ok := ret[sgid]; ok {
			ret[sgid] = append(ret[sgid], sid)
		} else {
			ret[sgid] = []int{sid}
		}
	}

	return ret, nil
}
