/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-DB管理系统(BlueKing-BK-DBM) available.
 * Copyright (C) 2017-2023 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at https://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package model

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm/clause"

	"dbm-services/common/go-pubpkg/logger"
)

const (
	// StringItem string
	StringItem = "string"
	// ArryItem arry
	ArryItem = "arry"
	// IntItem int
	IntItem = "int"
	// BoolItem bool
	BoolItem = "bool"
)

// TbSyntaxRule [...]
type TbSyntaxRule struct {
	ID int `gorm:"primaryKey;column:id;type:int(11);not null" json:"-"`
	// 规则组名称
	GroupName string `gorm:"uniqueIndex:group;column:group_name;type:varchar(64);not null" json:"group_name"`
	// 子规则项，一个规则可能包括过个子规则
	RuleName string          `gorm:"uniqueIndex:group;column:rule_name;type:varchar(64);not null" json:"rule_name"`
	Item     json.RawMessage `gorm:"column:item;type:varchar(1024);not null" json:"item"`
	ItemType string          `gorm:"column:item_type;type:varchar(128);not null" json:"item_type"`
	// 规则表达式
	Expr string `gorm:"column:expr;type:varchar(128);not null" json:"expr"`
	// 规则提示信息
	Desc string `gorm:"column:desc;type:varchar(512);not null" json:"desc"`
	// 0:作为普通检查项,1:禁用命中该规则的行为
	WarnLevel int16 `gorm:"column:warn_level;type:smallint(2);not null" json:"warn_level"`
	// 1：启用，0:禁用
	Status bool `gorm:"column:status;type:tinyint(1);not null" json:"status"`
}

// GetTableName get sql table name.获取数据库名字
func (obj *TbSyntaxRule) GetTableName() string {
	return "tb_syntax_rules"
}

func init() {
	if err := InitRule(); err != nil {
		logger.Fatal("init syntax rule failed %s", err.Error())
		return
	}
}

// InitRule init rules
func InitRule() (err error) {
	initRules := []TbSyntaxRule{}
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "CommandRule",
		RuleName:  "HighRiskCommandRule",
		Expr:      "Val in Item",
		ItemType:  ArryItem,
		Item: []byte(
			`["drop_table", "drop_index", "lock_tables", "drop_db", "analyze","rename_table", 
			"drop_procedure", "drop_view","drop_trigger","drop_function", "drop_server", 
			"drop_event", "drop_compression_dictionary","optimize", "alter_tablespace"]`),
		Desc:      "高危命令",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "CommandRule",
		RuleName:  "BanCommandRule",
		Expr:      "Val in Item",
		ItemType:  ArryItem,
		Item: []byte(`["truncate", "revoke", "kill", "reset", "drop_user", "grant",
					"create_user", "revoke_all", "shutdown", "lock_tables_for_backup",
					"reset", "purge", "lock_binlog_for_backup","lock_tables_for_backup",
					"install_plugin", "uninstall_plugin","alter_user"]`),
		Desc:      "高危变更类型",
		WarnLevel: 1,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "CreateTableRule",
		RuleName:  "SuggestBlobColumCount",
		Expr:      "Val >= Item ",
		ItemType:  IntItem,
		Item:      []byte(`10`),
		Desc:      "建议单表Blob字段不要过多",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "CreateTableRule",
		RuleName:  "SuggestEngine",
		Expr:      "not (Val contains Item) and ( len(Val) != 0 )",
		ItemType:  StringItem,
		Item:      []byte(`"innodb"`),
		Desc:      "建议使用Innodb表",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "CreateTableRule",
		RuleName:  "NeedPrimaryKey",
		Expr:      "Val == Item",
		ItemType:  IntItem,
		Item:      []byte(`1`),
		Desc:      "建议包含主键",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "CreateTableRule",
		RuleName:  "DefinerRule",
		Expr:      "Val not in Item ",
		ItemType:  ArryItem,
		Item:      []byte(`["ADMIN@localhost"]`),
		Desc:      "必须指定definer",
		WarnLevel: 0,
		Status:    true,
	})

	initRules = append(initRules, TbSyntaxRule{
		GroupName: "AlterTableRule",
		RuleName:  "HighRiskType",
		Expr:      "Val in Item",
		ItemType:  ArryItem,
		Item:      []byte(`["drop_column"]`),
		Desc:      "高危变更类型",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "AlterTableRule",
		RuleName:  "HighRiskPkAlterType",
		Expr:      "Val in Item",
		ItemType:  ArryItem,
		Item:      []byte(`["add_column", "add_key", "change_column"]`),
		Desc:      "主键高危变更类型",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "AlterTableRule",
		RuleName:  "AlterUseAfter",
		Expr:      "Val != Item",
		ItemType:  StringItem,
		Item:      []byte(`""`),
		Desc:      "变更表时使用了after",
		WarnLevel: 0,
		Status:    true,
	})
	initRules = append(initRules, TbSyntaxRule{
		GroupName: "AlterTableRule",
		RuleName:  "AddColumnMixed",
		Expr:      "( Item in Val ) && ( len(Val) > 1 )",
		ItemType:  StringItem,
		Item:      []byte(`"add_column"`),
		Desc:      "加字段和其它alter table 类型混用，可能导致非在线加字段",
		WarnLevel: 0,
		Status:    true,
	})

	initRules = append(initRules, TbSyntaxRule{
		GroupName: "DmlRule",
		RuleName:  "DmlNotHasWhere",
		Expr:      " Val != Item ",
		ItemType:  BoolItem,
		Item:      []byte(`true`),
		Desc:      "没有使用WHERE或者LIMIT,可能会导致全表数据更改",
		WarnLevel: 0,
		Status:    true,
	})

	for i := range initRules {
		if err = CreateRule(&initRules[i]); err != nil {
			logger.Error("初始化规则失败%s", err.Error())
			return err
		}
	}
	return err
}

// CreateRule create rule
func CreateRule(m *TbSyntaxRule) (err error) {
	return DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(m).Error
}

// GetAllRule get all rules
func GetAllRule() (rs []TbSyntaxRule, err error) {
	err = DB.Find(&rs).Error
	return
}

// GetRuleByName get rules group by group name
func GetRuleByName(group, rulename string) (rs TbSyntaxRule, err error) {
	err = DB.Where("group_name = ? and rule_name = ? ", group, rulename).First(&rs).Error
	return
}

// GetItemVal get item val
func GetItemVal(rule TbSyntaxRule) (val interface{}, err error) {
	switch rule.ItemType {
	case ArryItem:
		var d []string
		if err = json.Unmarshal(rule.Item, &d); err != nil {
			logger.Error("umarshal failed %s", err.Error())
			return nil, err
		}
		val = d
	case StringItem:
		var d string
		if err = json.Unmarshal(rule.Item, &d); err != nil {
			logger.Error("umarshal failed %s", err.Error())
			return nil, err
		}
		val = d
	case IntItem:
		var d int
		if err = json.Unmarshal(rule.Item, &d); err != nil {
			logger.Error("umarshal failed %s", err.Error())
			return nil, err
		}
		val = d
	case BoolItem:
		var d bool
		if err = json.Unmarshal(rule.Item, &d); err != nil {
			logger.Error("umarshal failed %s", err.Error())
			return nil, err
		}
		val = d
	default:
		return nil, fmt.Errorf("unrecognizable type:%s", rule.ItemType)
	}
	return val, err
}
