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
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"sync"
)

// 一个HostGroup可以绑定多个Plugin
type SafeGroupPlugins struct {
	sync.RWMutex
	M map[int][]string
	P map[int][]g.PluginParam
}

var GroupPlugins = &SafeGroupPlugins{M: make(map[int][]string), P: make(map[int][]g.PluginParam)}

func (this *SafeGroupPlugins) GetPlugins(gid int) ([]string, bool) {
	this.RLock()
	defer this.RUnlock()
	plugins, exists := this.M[gid]
	return plugins, exists
}

func (this *SafeGroupPlugins) GetPluginParams(gid int) ([]g.PluginParam, bool) {
	this.RLock()
	defer this.RUnlock()
	params, exists := this.P[gid]
	return params, exists
}

func (this *SafeGroupPlugins) Init() {
	m, err := db.QueryPlugins()
	if err != nil {
		return
	}

	p, e := db.QueryPluginParams()
	if e != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
	this.P = p
}

// 根据hostname获取关联的插件
func GetPlugins(hostname string) map[string]g.PluginParam {
	// 因为机器关联了多个Group，每个Group可能关联多个Plugin，故而一个机器关联的Plugin可能重复
	// 所以，如果同一台机器关联多个相同的插件，则随机取一个执行。
	pluginDirs := make(map[string]g.PluginParam)
	hid, exists := HostMap.GetID(hostname)
	if !exists {
		return pluginDirs
	}

	gids, exists := HostGroupsMap.GetGroupIds(hid)
	if !exists {
		return pluginDirs
	}

	for _, gid := range gids {
		plugins, exists := GroupPlugins.GetPlugins(gid)
		if !exists {
			continue
		}
		params, exists := GroupPlugins.GetPluginParams(gid)
		if !exists {
			continue
		}

		for _, plugin := range plugins {
			pluginDirs[plugin] = g.PluginParam{}
			for _, param := range params {
				if param.Dir == plugin {
					pluginDirs[plugin] = param
					break
				}
			}
		}
	}
	return pluginDirs
}
