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
	log "github.com/Sirupsen/logrus"
	"strings"
	"time"

	"crypto/tls"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/api"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
)

func HandleCallback(event *model.Event, action *api.Action) {

	teams := action.Uic
	phones := []string{}
	mails := []string{}
	ims := []string{}

	if teams != "" {
		//获取phone mail im 列表
		phones, mails, ims = api.ParseTeams(teams)
		smsContent := GenerateSmsContent(event)
		mailContent := GenerateMailContent(event)
		imContent := GenerateIMContent(event)
		if action.BeforeCallbackSms == 1 {
			redi.WriteSms(phones, smsContent)
			redi.WriteIM(ims, imContent)
		}

		if action.BeforeCallbackMail == 1 {
			redi.WriteMail(mails, smsContent, mailContent)
		}
	}

	message := Callback(event, action)

	if teams != "" {
		if action.AfterCallbackSms == 1 {
			redi.WriteSms(phones, message)
			redi.WriteIM(ims, message)
		}

		if action.AfterCallbackMail == 1 {
			redi.WriteMail(mails, message, message)
		}
	}

}

func Callback(event *model.Event, action *api.Action) string {
	if action.Url == "" {
		return "callback url is blank"
	}

	L := make([]string, 0)
	if len(event.PushedTags) > 0 {
		for k, v := range event.PushedTags {
			L = append(L, fmt.Sprintf("%s:%s", k, v))
		}
	}

	tags := ""
	if len(L) > 0 {
		tags = strings.Join(L, ",")
	}

	data := make(map[string]interface{})
	data["endpoint"] = event.Endpoint
	data["metric"] = event.Metric()
	data["status"] = event.Status
	data["step"] = fmt.Sprintf("%d", event.CurrentStep)
	data["priority"] = fmt.Sprintf("%d", event.Priority())
	data["time"] = event.FormattedTime()
	data["tpl_id"] = fmt.Sprintf("%d", event.TplId())
	data["exp_id"] = fmt.Sprintf("%d", event.ExpressionId())
	data["stra_id"] = fmt.Sprintf("%d", event.StrategyId())
	data["tags"] = tags
	req := httplib.Post(action.Url).SetTimeout(3*time.Second, 20*time.Second)
	req.Header("Content-type", "application/json")
	req.Header("St2-Api-Key", "YTRlMWI2ZmI5N2Y4YjU5ZGNjNTc1ZWM3ODQyNWY0MzFlN2NjZTkwM2MzNDk4MzI5OWJiOWZiNjg2OGRlOTcyNQ")

	//跳过证书验证
	resp := make(map[string]interface{})
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	e := req.ToJson(&resp)

	success := "success"
	if e != nil {
		log.Errorf("callback fail, action:%v, event:%s, error:%s", action, event.String(), e.Error())
		success = fmt.Sprintf("fail:%s", e.Error())
	}
	message := fmt.Sprintf("curl %s %s. resp: %v", action.Url, success, resp)
	log.Debugf("callback to url:%s, event:%s, resp:%v", action.Url, event.String(), resp)

	return message
}
