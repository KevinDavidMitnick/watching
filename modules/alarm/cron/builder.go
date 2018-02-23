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
	"fmt"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"io/ioutil"
	"net/http"
)

type NameStruct struct {
	DisplayName string `json:"display_name"`
	Parent_name string `json:"parent_name"`
	Zone_name   string `json:"zone_name"`
	Database    string `json:"database"`
}

type DataStruct struct {
	Name *NameStruct `json:"data"`
}

func BuildCommonSMSContent(event *model.Event) string {
	var data DataStruct
	endpoint := event.Endpoint
	addr := g.Config().CmdbConfig.Addr + "/" + event.Endpoint
	request, _ := http.NewRequest("GET", addr, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			err := json.Unmarshal(body, &data)
			if err == nil && data.Name != nil && data.Name.DisplayName != "" {
				endpoint = string(data.Name.DisplayName)
			}
		}
	}

	// change by liucong format mail title
	return fmt.Sprintf(
		"[P%d][%s][0%d][%s][%s]",
		event.Priority(),
		event.Status,
		event.CurrentStep,
		endpoint,
		event.Metric(),
	)
}

func BuildCommonIMContent(event *model.Event) string {
	var data DataStruct
	endpoint := event.Endpoint
	addr := g.Config().CmdbConfig.Addr + "/" + event.Endpoint
	request, _ := http.NewRequest("GET", addr, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			err := json.Unmarshal(body, &data)
			if err == nil && data.Name != nil && data.Name.DisplayName != "" {
				endpoint = string(data.Name.DisplayName)
			}
		}
	}
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
	// get hostname from cmdb ,modify by liucong.
	var data DataStruct
	addr := g.Config().CmdbConfig.Addr + "/" + event.Endpoint
	endpoint := "Endpoint(Uuid):" + event.Endpoint

	request, _ := http.NewRequest("GET", addr, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			err := json.Unmarshal(body, &data)
			if err == nil && data.Name != nil && data.Name.DisplayName != "" {
				endpoint = "对象(Object):" + string(data.Name.DisplayName)
			}
		}
	}

	link := g.Link(event)
	// change by liucong format mail
	line := "------------------------------------"
	status := "类型(Type)"
	if event.Status == "OK" {
		status += ":" + "恢复"
	} else {
		status += ":" + "告警"
	}
	level := fmt.Sprintf("级别(Level):P%d", event.Priority())
	timestamp := fmt.Sprintf("时间(Timestamp):%s", event.FormattedTime())
	uuid := "Uuid:" + string(event.Endpoint)
	position := "位置信息(Position):" + string(data.Name.Database) + " " + string(data.Name.Zone_name) + " " + string(data.Name.Parent_name)
	metric := "指标(Metric):" + event.Metric()
	note := "描述(Description):" + event.Note()
	tags := "标签(Tags):" + utils.SortedTags(event.PushedTags)
	meta := "元数据(Meta-data):"
	funcs := "函数(func):" + event.Func() + ":" + utils.ReadableFloat(event.LeftValue) + event.Operator() + utils.ReadableFloat(event.RightValue())
	times := fmt.Sprintf("报警次数(Strategy):最大(max)%d次，当前(current)第%d次", event.MaxStep(), event.CurrentStep)
	tpl := "模板(Tpl):" + link
	return fmt.Sprintf(
		"%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\t%s\r\n\t%s\r\n\t%s\r\n",
		line,
		status,
		level,
		timestamp,
		endpoint,
		uuid,
		position,
		metric,
		note,
		tags,
		meta,
		funcs,
		times,
		tpl,
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateIMContent(event *model.Event) string {
	return BuildCommonIMContent(event)
}
