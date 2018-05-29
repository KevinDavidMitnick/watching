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

////////////////////////////////////////////////////////////////////////////
// | Field       | Type             | Null | Key | Default | Extra          |
// +-------------+------------------+------+-----+---------+----------------+
// | id          | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
// | name        | varchar(255)     | NO   |     |         |                |
// | tpl_id      | int(10) unsigned | NO   | MUL | 0       |                |
// +-------------+------------------+------+-----+---------+----------------+
////////////////////////////////////////////////////////////////////////////

type StrategyGroup struct {
	ID    int64  `json:"id" gorm:"column:id"`
	Name  string `json:"name" gorm:"column:name"`
	TplId int64  `json:"tpl_id" gorm:"column:tpl_id"`
}

func (this StrategyGroup) TableName() string {
	return "strategy_group"
}
