/*
@Time    :   2023/04/26 10:12:19
@Author  :   zongfei.fu
@Desc    :   检查TiDB异常的连接是否被释放，如果没有，触发kill操作。
			 原因是TiDB不支持DB接口上下文超时自动结束会话
			 已测试不支持版本：TiDB5.x/TiDB6.x
*/

package tasks

import (
	"context"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/das/dao"
	"goInsight/internal/apps/das/models"
	"strconv"
	"sync"
	"time"
)

type tidbQueryRecord struct {
	Username  string
	RequestID string
	Hostname  string
	Port      int
}

type KillTiDBQuery struct{}

func (k *KillTiDBQuery) getDASRecords() *[]tidbQueryRecord {
	// 获取global.App.Config.Config.MaxExecutionTime/1000时间内tidb异常的查询记录
	var results []tidbQueryRecord
	global.App.DB.Model(&models.InsightDASRecords{}).
		Select("insight_das_records.username, insight_das_records.request_id, insight_db_config.hostname, insight_db_config.port").
		Joins("join insight_db_config on insight_das_records.instance_id = insight_db_config.instance_id").
		Where("insight_das_records.is_finish=1 and insight_das_records.is_kill=0").
		Where("insight_das_records.created_at>= date_sub(now(), interval ? second)", global.App.Config.Das.MaxExecutionTime/1000).
		Where("insight_db_config.db_type = 'TiDB'").
		Scan(&results)
	return &results
}

func (k *KillTiDBQuery) match(row tidbQueryRecord) (*[]string, *[]map[string]interface{}, error) {
	// 根据RequestID和查询时请求用户获取异常查询的信息
	query := fmt.Sprintf(
		`
		select 
			*
		from 
			information_schema.processlist 
		where 
			User="%s" 
			and DB!="information_schema"
			and INFO like "%%%s%%"
	`, global.App.Config.RemoteDB.UserName, row.RequestID)
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	db := dao.DB{
		User:     global.App.Config.RemoteDB.UserName,
		Password: global.App.Config.RemoteDB.Password,
		Host:     row.Hostname,
		Port:     row.Port,
		Database: "information_schema",
		Ctx:      ctx,
	}
	return db.Query(query)
}

func (k *KillTiDBQuery) update(request_id string) {
	// 更新记录
	global.App.DB.Model(&models.InsightDASRecords{}).Where("request_id=?", request_id).Update("is_kill", true)
}

func (k *KillTiDBQuery) kill(row tidbQueryRecord) {
	_, data, err := k.match(row)
	if err != nil {
		return
	}
	for _, d := range *data {
		var killUser string = d["USER"].(string)
		if killUser == global.App.Config.RemoteDB.UserName {
			queryID, _ := strconv.Atoi(d["ID"].(string))
			query := fmt.Sprintf("kill tidb query %d", queryID)
			ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
			defer cancel()
			db := dao.DB{
				User:     global.App.Config.RemoteDB.UserName,
				Password: global.App.Config.RemoteDB.Password,
				Host:     row.Hostname,
				Port:     row.Port,
				Database: "information_schema",
				Ctx:      ctx,
			}
			if err := db.Execute(query); err != nil {
				global.App.Log.Errorf("`kill tidb query %d` Faild, RequestID: %s, Error:%s", queryID, row.RequestID, err.Error())
			} else {
				global.App.Log.Infof("`kill tidb query %d` Success, RequestID: %s,", queryID, row.RequestID)
				k.update(row.RequestID)
			}
		}
	}
}

func (k *KillTiDBQuery) Run() {
	// 开启5个并发
	var wg sync.WaitGroup
	ch := make(chan struct{}, 5)
	for _, i := range *k.getDASRecords() {
		ch <- struct{}{}
		wg.Add(1)
		go func(i tidbQueryRecord) {
			defer wg.Done()
			k.kill(i)
			<-ch
		}(i)
	}
	wg.Wait()
}
