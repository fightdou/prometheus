package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
	"log"
	"strings"
	//blog "github.com"
)
var (
	_SQL_DB *sql.DB
	_ERR error
)
var Expr string
//规则结构定义
type RuleItem struct {
	Name        string            //规则所属组得名称
	Fn          string            //类别
	Interval    int               //规则计算间隔
	Alert       string            //告警名称
	Rules		string			  //告警策略
	Math		string			  //告警判断符
	Values	 	string			  //告警阈值
	For         string            //持续时间
	Resend		string			  //重复告警时间周期
	Labels      map[string]string //规则维度信息
	Annotations map[string]string //规则描述信息
}
// 初始化数据库连接
func Initialzation(dbUrl string){
	_SQL_DB, _ERR = sql.Open("mysql", dbUrl)
	if _ERR != nil{
		log.Fatalf("Open database error: %s\n", _ERR)
		return
	}
	_ERR = _SQL_DB.Ping()
	if _ERR != nil{
		log.Fatal(_ERR)
	}
}
//规则查询,用于将rules表转化为RuleItem结构。
func QueryRuleString()(*list.List, error){
	var(
		rule_labels,rule_annotations string
	)
	l := list.New()
	// 查询规则表
	rows, err := _SQL_DB.Query("select rule_name,rule_fn,rule_interval,rule_alert,rule_rules,rule_math,rule_values,rule_for,rule_resend,rule_labels,rule_annotations from rules;")
	if err != nil{
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next(){
		var item RuleItem
		item.Labels=make(map[string]string)
		item.Annotations=make(map[string]string)
		err := rows.Scan(&item.Name,&item.Fn,&item.Interval,&item.Alert,&item.Rules,&item.Math,&item.Values,&item.For,&item.Resend,&rule_labels,&rule_annotations)
		if err != nil{
			log.Fatal(err)
		}
		//label 数据格式转换
		labels := strings.Split(rule_labels, ",")
		lablen := len(labels)
		for i:=0;i<lablen;i++ {
			pars := strings.Split(labels[i], "=")
			plen := len(pars)
			for j:=0;j<plen;j+=2 {
				item.Labels[pars[j]]=pars[j+1]
			}
		}
		//annotations数据格式转换
		annotations := strings.Split(rule_annotations, ",")
		annlen := len(annotations)
		for k:=0;k<annlen;k++{
			pars := strings.Split(annotations[k], "=")
			plen := len(pars)
			for j:=0;j<plen;j+=2{
				item.Annotations[pars[j]]=pars[j+1]
			}
		}
		l.PushBack(item)
	}
	return l,err
}

