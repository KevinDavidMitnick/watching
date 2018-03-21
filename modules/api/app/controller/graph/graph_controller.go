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

package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/httplib"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
	grh "github.com/open-falcon/falcon-plus/modules/api/graph"
	tcache "github.com/toolkits/cache/localcache/timedcache"
	"net/http"
)

var (
	localStepCache = tcache.New(600*time.Second, 60*time.Second)
)

type APIEndpointObjGetInputs struct {
	Endpoints []string `json:"endpoints" form:"endpoints"`
	Deadline  int64    `json:"deadline" form:"deadline"`
}

func EndpointObjGet(c *gin.Context) {
	inputs := APIEndpointObjGetInputs{
		Deadline: 0,
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if len(inputs.Endpoints) == 0 {
		h.JSONR(c, http.StatusBadRequest, "endpoints missing")
		return
	}

	var result []m.Endpoint = []m.Endpoint{}
	dt := db.Graph.Table("endpoint").
		Where("endpoint in (?) and ts >= ?", inputs.Endpoints, inputs.Deadline).
		Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, dt.Error)
		return
	}

	endpoints := []map[string]interface{}{}
	for _, r := range result {
		endpoints = append(endpoints, map[string]interface{}{"id": r.ID, "endpoint": r.Endpoint, "ts": r.Ts})
	}

	h.JSONR(c, endpoints)
}

type APIEndpointRegexpQueryInputs struct {
	Q     string `json:"q" form:"q"`
	Label string `json:"tags" form:"tags"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

func EndpointRegexpQuery(c *gin.Context) {
	inputs := APIEndpointRegexpQueryInputs{
		//set default is 500
		Limit: 500,
		Page:  1,
	}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if inputs.Q == "" && inputs.Label == "" {
		h.JSONR(c, http.StatusBadRequest, "q and labels are all missing")
		return
	}

	labels := []string{}
	if inputs.Label != "" {
		labels = strings.Split(inputs.Label, ",")
	}
	qs := []string{}
	if inputs.Q != "" {
		qs = strings.Split(inputs.Q, " ")
	}

	var offset int = 0
	if inputs.Page > 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}

	var endpoint []m.Endpoint
	var endpoint_id []int
	var dt *gorm.DB
	if len(labels) != 0 {
		dt = db.Graph.Table("endpoint_counter").Select("distinct endpoint_id")
		for _, trem := range labels {
			dt = dt.Where(" counter like ? ", "%"+strings.TrimSpace(trem)+"%")
		}
		dt = dt.Limit(inputs.Limit).Offset(offset).Pluck("distinct endpoint_id", &endpoint_id)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}
	}
	if len(qs) != 0 {
		dt = db.Graph.Table("endpoint").
			Select("endpoint, id")
		if len(endpoint_id) != 0 {
			dt = dt.Where("id in (?)", endpoint_id)
		}

		for _, trem := range qs {
			dt = dt.Where(" endpoint regexp ? ", strings.TrimSpace(trem))
		}
		dt.Limit(inputs.Limit).Offset(offset).Scan(&endpoint)
	} else if len(endpoint_id) != 0 {
		dt = db.Graph.Table("endpoint").
			Select("endpoint, id").
			Where("id in (?)", endpoint_id).
			Scan(&endpoint)
	}
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, dt.Error)
		return
	}

	endpoints := []map[string]interface{}{}
	for _, e := range endpoint {
		endpoints = append(endpoints, map[string]interface{}{"id": e.ID, "endpoint": e.Endpoint})
	}

	h.JSONR(c, endpoints)
}

type uuidName struct {
	Display_name string `json:"display_name"`
	Uuid         string `json:"uuid"`
}

func getNameRelUuid(url string) (map[string]string, error) {
	req := httplib.Get(url).SetTimeout(3*time.Second, 20*time.Second)
	resp, e := req.String()
	if e != nil {
		log.Warnln(e.Error())
		return nil, e
	}

	type Resp struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    []uuidName `json:"data"`
	}
	var rdata Resp
	json.Unmarshal([]byte(resp), &rdata)

	cache := make(map[string]string, 0)
	for _, data := range rdata.Data {
		uuid := data.Uuid
		display := data.Display_name
		cache[uuid] = display
	}
	return cache, nil

}

func EndpointCounterRegexpQuery(c *gin.Context) {
	eid := c.DefaultQuery("eid", "")
	metricQuery := c.DefaultQuery("metricQuery", ".+")
	limitTmp := c.DefaultQuery("limit", "3000")
	limit, err := strconv.Atoi(limitTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	pageTmp := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageTmp)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	var offset int = 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	if eid == "" {
		h.JSONR(c, http.StatusBadRequest, "eid is missing")
	} else {
		eids := utils.ConverIntStringToList(eid)
		if eids == "" {
			h.JSONR(c, http.StatusBadRequest, "input error, please check your input info.")
			return
		} else {
			eids = fmt.Sprintf("(%s)", eids)
		}

		var counters []m.EndpointCounter
		dt := db.Graph.Table("endpoint_counter").Select("endpoint_id, counter, step, type").Where(fmt.Sprintf("endpoint_id IN %s", eids))
		if metricQuery != "" {
			qs := strings.Split(metricQuery, " ")
			if len(qs) > 0 {
				for _, term := range qs {
					dt = dt.Where("counter regexp ?", strings.TrimSpace(term))
				}
			}
		}
		dt = dt.Limit(limit).Offset(offset).Scan(&counters)
		if dt.Error != nil {
			h.JSONR(c, http.StatusBadRequest, dt.Error)
			return
		}

		/*get counter alias*/
		var counteraliass []m.CounterAlias
		var gdt *gorm.DB
		gdt = db.Graph.Table("counter_alias").Find(&counteraliass)
		if gdt.Error != nil {
			h.JSONR(c, expecstatus, dt.Error)
			return
		}
		cache := make(map[string]string, 0)
		counerIdCache := make(map[string]uint, 0)
		for _, counteralias := range counteraliass {
			counter := counteralias.Counters
			alias := counteralias.Alias
			cache[counter] = alias
			counerIdCache[counter] = counteralias.ID
		}

		countersResp := []interface{}{}
		var current string
		var currid uint = 0
		for _, c := range counters {
			if v, ok := cache[c.Counter]; ok {
				current = v
				currid = counerIdCache[c.Counter]
			} else {
				current = ""
			}

			classification := classify(c.Counter, current)
			countersResp = append(countersResp, map[string]interface{}{
				"endpoint_id":      c.EndpointID,
				"counter":          c.Counter,
				"counter_alias":    current,
				"step":             c.Step,
				"type":             c.Type,
				"counter_alias_id": currid,
				"classification":   classification,
			})
		}
		h.JSONR(c, countersResp)
	}
	return
}

type APIQueryGraphDrawData struct {
	HostNames []string `json:"hostnames" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
	ConsolFun string   `json:"consol_fun" binding:"required"`
	StartTime int64    `json:"start_time" binding:"required"`
	EndTime   int64    `json:"end_time" binding:"required"`
	Step      int      `json:"step"`
}

func QueryGraphDrawData(c *gin.Context) {
	endpointcache, _ := getNameRelUuid(config.CMDB_ADDR)

	var inputs APIQueryGraphDrawData
	var err error
	if err = c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	/*get counter alias*/
	var counteraliass []m.CounterAlias
	var gdt *gorm.DB
	gdt = db.Graph.Table("counter_alias").Find(&counteraliass)
	if gdt.Error != nil {
		h.JSONR(c, expecstatus, gdt.Error)
		return
	}
	cache := make(map[string]string, 0)
	for _, counteralias := range counteraliass {
		counter := counteralias.Counters
		alias := counteralias.Alias
		cache[counter] = alias
	}

	respData := []*cmodel.GraphQueryResponse{}
	for _, host := range inputs.HostNames {
		for _, counter := range inputs.Counters {
			var step int
			if inputs.Step > 0 {
				step = inputs.Step
			} else {
				step, err = getCounterStep(host, counter)
				if err != nil {
					continue
				}
			}
			data, _ := fetchData(host, counter, inputs.ConsolFun, inputs.StartTime, inputs.EndTime, step)

			if v, ok := cache[data.Counter]; ok {
				data.Counter = v
			}

			if v, ok := endpointcache[data.Endpoint]; ok {
				data.Endpoint = v
			}
			respData = append(respData, data)
		}
	}
	h.JSONR(c, respData)
}

func QueryGraphLastPoint(c *gin.Context) {
	var inputs []cmodel.GraphLastParam
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	respData := []*cmodel.GraphLastResp{}

	for _, param := range inputs {
		one_resp, err := grh.Last(param)
		if err != nil {
			log.Warn("query last point from graph fail:", err)
		} else {
			respData = append(respData, one_resp)
		}
	}

	h.JSONR(c, respData)
}

func DeleteGraphEndpoint(c *gin.Context) {
	var inputs []string = []string{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	type DBRows struct {
		Endpoint  string
		CounterId int
		Counter   string
		Type      string
		Step      int
	}
	rows := []DBRows{}
	dt := db.Graph.Raw(
		`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a, endpoint_counter as b
		where b.endpoint_id = a.id
		AND a.endpoint in (?)`, inputs).Scan(&rows)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	var affected_counter int64 = 0
	var affected_endpoint int64 = 0

	if len(rows) > 0 {
		var params []*cmodel.GraphDeleteParam = []*cmodel.GraphDeleteParam{}
		for _, row := range rows {
			param := &cmodel.GraphDeleteParam{
				Endpoint: row.Endpoint,
				DsType:   row.Type,
				Step:     row.Step,
			}
			fields := strings.SplitN(row.Counter, "/", 2)
			if len(fields) == 1 {
				param.Metric = fields[0]
			} else if len(fields) == 2 {
				param.Metric = fields[0]
				param.Tags = fields[1]
			} else {
				log.Error("invalid counter", row.Counter)
				continue
			}
			params = append(params, param)
		}
		grh.Delete(params)
	}

	tx := db.Graph.Begin()

	if len(rows) > 0 {
		var cids []int = make([]int, len(rows))
		for i, row := range rows {
			cids[i] = row.CounterId
		}

		dt = tx.Table("endpoint_counter").Where("id in (?)", cids).Delete(&m.EndpointCounter{})
		if dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			tx.Rollback()
			return
		}
		affected_counter = dt.RowsAffected

		dt = tx.Exec(`delete from tag_endpoint where endpoint_id in 
			(select id from endpoint where endpoint in (?))`, inputs)
		if dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			tx.Rollback()
			return
		}
	}

	dt = tx.Table("endpoint").Where("endpoint in (?)", inputs).Delete(&m.Endpoint{})
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		tx.Rollback()
		return
	}
	affected_endpoint = dt.RowsAffected
	tx.Commit()

	h.JSONR(c, map[string]int64{
		"affected_endpoint": affected_endpoint,
		"affected_counter":  affected_counter,
	})
}

type APIGraphDeleteCounterInputs struct {
	Endpoints []string `json:"endpoints" binding:"required"`
	Counters  []string `json:"counters" binding:"required"`
}

func DeleteGraphCounter(c *gin.Context) {
	var inputs APIGraphDeleteCounterInputs = APIGraphDeleteCounterInputs{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	type DBRows struct {
		Endpoint  string
		CounterId int
		Counter   string
		Type      string
		Step      int
	}
	rows := []DBRows{}
	dt := db.Graph.Raw(`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a,
		endpoint_counter as b
		where b.endpoint_id = a.id 
		AND a.endpoint in (?)
		AND b.counter in (?)`, inputs.Endpoints, inputs.Counters).Scan(&rows)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	if len(rows) == 0 {
		h.JSONR(c, map[string]int64{
			"affected_counter": 0,
		})
		return
	}

	var params []*cmodel.GraphDeleteParam = []*cmodel.GraphDeleteParam{}
	for _, row := range rows {
		param := &cmodel.GraphDeleteParam{
			Endpoint: row.Endpoint,
			DsType:   row.Type,
			Step:     row.Step,
		}
		fields := strings.SplitN(row.Counter, "/", 2)
		if len(fields) == 1 {
			param.Metric = fields[0]
		} else if len(fields) == 2 {
			param.Metric = fields[0]
			param.Tags = fields[1]
		} else {
			log.Error("invalid counter", row.Counter)
			continue
		}
		params = append(params, param)
	}
	grh.Delete(params)

	tx := db.Graph.Begin()
	var cids []int = make([]int, len(rows))
	for i, row := range rows {
		cids[i] = row.CounterId
	}

	dt = tx.Table("endpoint_counter").Where("id in (?)", cids).Delete(&m.EndpointCounter{})
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		tx.Rollback()
		return
	}
	affected_counter := dt.RowsAffected
	tx.Commit()

	h.JSONR(c, map[string]int64{
		"affected_counter": affected_counter,
	})
}

func fetchData(hostname string, counter string, consolFun string, startTime int64, endTime int64, step int) (resp *cmodel.GraphQueryResponse, err error) {
	qparm := grh.GenQParam(hostname, counter, consolFun, startTime, endTime, step)
	// log.Debugf("qparm: %v", qparm)
	resp, err = grh.QueryOne(qparm)
	if err != nil {
		log.Debugf("query graph got error: %s", err.Error())
	}
	return
}

func getCounterStep(endpoint, counter string) (step int, err error) {
	cache_key := fmt.Sprintf("step:%s/%s", endpoint, counter)
	s, found := localStepCache.Get(cache_key)
	if found && s != nil {
		step = s.(int)
		return
	}

	var rows []int
	dt := db.Graph.Raw(`select a.step from endpoint_counter as a, endpoint as b
		where b.endpoint = ? and a.endpoint_id = b.id and a.counter = ? limit 1`, endpoint, counter).Scan(&rows)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	if len(rows) == 0 {
		err = errors.New("empty result")
		return
	}
	step = rows[0]
	localStepCache.Set(cache_key, step, tcache.DefaultExpiration)

	return
}

type APICounterAliasInputs struct {
	Alias    string `json:"alias"`
	Counters string `json:"counters"`
	Type     uint   `json:"type"`
}

func CreateCounterAlias(c *gin.Context) {
	var inputs APICounterAliasInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	counteralias := m.CounterAlias{
		Counters: inputs.Counters,
		Alias:    inputs.Alias,
		Htype:    inputs.Type,
	}

	if dt := db.Graph.Create(&counteralias); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, counteralias)
	return
}

type APIUpdateCounterAlias struct {
	New_alias string `json:"new_alias" binding:"required"`
	Id        int64  `json:"id" binding:"required"`
}

func UpdateCounterAlias(c *gin.Context) {
	var inputs APIUpdateCounterAlias
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	//user, _ := h.GetUser(c)
	counteralias := m.CounterAlias{
		Alias: inputs.New_alias,
	}
	if dt := db.Graph.Table("counter_alias").Where("id=?", inputs.Id).Update(&counteralias); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}

	h.JSONR(c, map[string]string{
		"message": "update success",
		"id":      fmt.Sprintf("%d", inputs.Id),
	})
	return
}

func GetCounterAlias(c *gin.Context) {
	q := c.DefaultQuery("q", ".+")

	var counteraliass []m.CounterAlias
	var dt *gorm.DB
	dt = db.Graph.Table("counter_alias").Where("counters regexp ?", q).Find(&counteraliass)
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, counteraliass)
	return
}

func DeleteCounterAlias(c *gin.Context) {
	idTmp := c.Params.ByName("nid")
	if idTmp == "" {
		h.JSONR(c, badstatus, "counter_alias id is missing")
		return
	}
	counterID, err := strconv.Atoi(idTmp)
	if err != nil {
		log.Debugf("counterIDtmp: %v", counterID)
		h.JSONR(c, badstatus, err)
		return
	}
	plugin := m.CounterAlias{ID: uint(counterID)}
	if dt := db.Graph.Find(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}

	if dt := db.Graph.Delete(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("counter alias:%v has been deleted", counterID))
	return
}

func classify(counter, alias string) map[string]interface{} {
	counterList := strings.SplitN(counter, "/", 2)
	metrics := counterList[0]
	var tags string
	if len(counterList) > 1 {
		tags = counterList[1]
	} else {
		tags = ""
	}
	metricsSlice := strings.Split(metrics, ".")
	fMetric := metricsSlice[0]
	sMetric := metricsSlice[1:]

	class := make(map[string]interface{})
	child := ""
	if alias != "" {
		class["name"] = fMetric
		class["child"] = map[string]string{"name": alias, "child": child}
	} else {
		class["name"] = fMetric
		if tags != "" {
			child = tags
		}
		sMetricStr := strings.Join(sMetric, ".")
		class["child"] = map[string]string{"name": sMetricStr, "child": child}
	}
	return class
}
