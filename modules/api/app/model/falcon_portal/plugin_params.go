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

package falcon_portal

// +-------------+------------------+------+-----+-------------------+----------------+
// | Field       | Type             | Null | Key | Default           | Extra          |
// +-------------+------------------+------+-----+-------------------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL              | auto_increment |
// | grp_id      | int(10) unsigned | NO   | MUL | NULL              |                |
// | dir         | varchar(255)     | NO   |     | NULL              |                |
// | execute_script | varchar(255)  | YES  |     |      ""           |                |
// | execute_interval | int(10) unsigned   | NO  |     |   60        |                |
// | execute_param  | varchar(255)         | YES   |     | ""        |                |
// +-------------+------------------+------+-----+-------------------+----------------+

type PluginParams struct {
	ID              int64  `json:"id" gorm:"column:id"`
	GrpId           int64  `json:"grp_id" gorm:"column:grp_id"`
	Dir             string `json:"dir" gorm:"column:dir"`
	ExecuteScript   string `json:"execute_script" gorm:"column:execute_script"`
	ExecuteInterval int64  `json:"execute_interval" gorm:"column:execute_interval"`
	ExecuteParam    string `json:"execute_param" gorm:"column:execute_param"`
}

func (this PluginParams) TableName() string {
	return "plugin_params"
}
