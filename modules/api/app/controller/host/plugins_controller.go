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

package host

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
)

type APICreatePluginInput struct {
	GrpId   int64  `json:"hostgroup_id" binding:"required"`
	DirPaht string `json:"dir_path" binding:"required"`
}

type APICreatePluginParamsInput struct {
	GrpId           int64  `json:"hostgroup_id" binding:"required"`
	DirPath         string `json:"dir_path" binding:"required"`
	ExecuteScript   string `json:"execute_script" binding:"required"`
	ExecuteInterval int64  `json:"execute_interval" binding:"required"`
	ExecuteParam    string `json:"execute_param" binding:"required"`
}

func CreatePlugin(c *gin.Context) {
	var inputs APICreatePluginInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	//user, _ := h.GetUser(c)
	//if !user.IsAdmin() {
	//	hostgroup := f.HostGroup{ID: inputs.GrpId}
	//	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
	//		h.JSONR(c, expecstatus, dt.Error)
	//		return
	//	}
	//	if hostgroup.CreateUser != user.Name {
	//		h.JSONR(c, badstatus, "You don't have permission!")
	//		return
	//	}
	//}
	plugin := f.Plugin{Dir: inputs.DirPaht, GrpId: inputs.GrpId, CreateUser: "root"}
	if dt := db.Falcon.Save(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, plugin)
	return
}

func CreatePluginParams(c *gin.Context) {
	var inputs APICreatePluginParamsInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	plugin := f.PluginParams{Dir: inputs.DirPath, GrpId: inputs.GrpId, ExecuteScript: inputs.ExecuteScript, ExecuteInterval: inputs.ExecuteInterval, ExecuteParam: inputs.ExecuteParam}
	if dt := db.Falcon.Save(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, plugin)
	return
}

func GetPluginOfGrp(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	plugins := []f.Plugin{}
	if dt := db.Falcon.Where("grp_id = ?", grpID).Find(&plugins); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, plugins)
	return
}

func GetPluginParamsOfGrp(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	plugin_params := []f.PluginParams{}
	if dt := db.Falcon.Where("grp_id = ?", grpID).Find(&plugin_params); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, plugin_params)
	return
}

func DeletePlugin(c *gin.Context) {
	pluginIDtmp := c.Params.ByName("id")
	if pluginIDtmp == "" {
		h.JSONR(c, badstatus, "plugin id is missing")
		return
	}
	pluginID, err := strconv.Atoi(pluginIDtmp)
	if err != nil {
		log.Debugf("pluginIDtmp: %v", pluginIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	plugin := f.Plugin{ID: int64(pluginID)}
	if dt := db.Falcon.Find(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	//user, _ := h.GetUser(c)
	//if !user.IsAdmin() {
	//	hostgroup := f.HostGroup{ID: plugin.GrpId}
	//	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
	//		h.JSONR(c, expecstatus, dt.Error)
	//		return
	//	}
	//	if hostgroup.CreateUser != user.Name && plugin.CreateUser != user.Name {
	//		h.JSONR(c, badstatus, "You don't have permission!")
	//		return
	//	}
	//}

	if dt := db.Falcon.Delete(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("plugin:%v has been deleted", pluginID))
	return
}

func DeletePluginParams(c *gin.Context) {
	pluginIDtmp := c.Params.ByName("id")
	if pluginIDtmp == "" {
		h.JSONR(c, badstatus, "plugin id is missing")
		return
	}
	pluginID, err := strconv.Atoi(pluginIDtmp)
	if err != nil {
		log.Debugf("pluginIDtmp: %v", pluginIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	plugin := f.PluginParams{ID: int64(pluginID)}
	if dt := db.Falcon.Find(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	//user, _ := h.GetUser(c)
	//if !user.IsAdmin() {
	//	hostgroup := f.HostGroup{ID: plugin.GrpId}
	//	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
	//		h.JSONR(c, expecstatus, dt.Error)
	//		return
	//	}
	//	if hostgroup.CreateUser != user.Name && plugin.CreateUser != user.Name {
	//		h.JSONR(c, badstatus, "You don't have permission!")
	//		return
	//	}
	//}

	if dt := db.Falcon.Delete(&plugin); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("plugin:%v has been deleted", pluginID))
	return
}
