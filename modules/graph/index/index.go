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

package index

import (
	log "github.com/Sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

// 初始化索引功能模块
func Start() {
	InitCache()
	go StartIndexUpdateIncrTask()
	log.Debug("index.Start ok")
}

// index收到一条新上报的监控数据,尝试用于增量更新索引
func ReceiveItem(item *cmodel.GraphItem, md5 string) {
	if item == nil {
		return
	}

	uuid := item.UUID()

	// 已上报过的数据
	if IndexedItemCache.ContainsKey(md5) {
		old := IndexedItemCache.Get(md5).(*IndexCacheItem)
		if uuid == old.UUID { // dsType+step没有发生变化,只更新缓存
			IndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		} else { // dsType+step变化了,当成一个新的增量来处理
			unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		}
		return
	}

	//流程走到这一块，说明数据没有上报过
	// 缓存未命中, 放入增量更新队列
	unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
}

//从graph cache中删除掉某个item, 并删除指定的counter对应的rrd文件
func RemoveItem(item *cmodel.GraphItem) {
	md5 := item.Checksum()
	IndexedItemCache.Remove(md5)
	unIndexedItemCache.Remove(md5)

	//discard data of memory
	checksum := item.Checksum()
	key := g.FormRrdCacheKey(checksum, item.DsType, item.Step)
	poped_items := store.GraphItems.PopAll(key)
	log.Debugf("discard data of item:%v, size:%d", item, len(poped_items))

	rrdFileName := g.RrdFileName(md5, item.DsType, item.Step)
	// TODO: liucong,删除对应的rrd文件.
	log.Debug("alert: should remove rrdfile:", rrdFileName)
}
