package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 Rewrite 方法
func TestRewrite(t *testing.T) {
	tests := []struct {
		inputSQL    string
		expectedSQL string
	}{
		// 基本的Insert语句
		{
			inputSQL:    "INSERT INTO test1 (id, order_no, created_time, updated_time) VALUES (14, 'Order014', '2024-06-25 12:00:00', '2024-06-25 12:00:00');",
			expectedSQL: "select 1 from DUAL",
		},
		// insert into ... select语句
		{
			inputSQL:    "insert into t2 select * from t1 where id >0;",
			expectedSQL: "select * from t1 where id > 0",
		},
		// insert into ... select ... limit ...语句
		{
			inputSQL:    "insert into t2 select * from t1 where id >0 limit 10;",
			expectedSQL: "select * from t1 where id > 0 limit 10",
		},
		// insert into ... select ... order by ... limit ...语句
		{
			inputSQL:    "insert into t2 select * from t1 where id >0 order by c1 desc limit 10;",
			expectedSQL: "select * from t1 where id > 0 order by c1 desc limit 10",
		},
		// insert into ... select distinct ... order by ... limit ...语句
		{
			inputSQL:    "insert into t2 select distinct id from t1 where id >0 order by c1 desc limit 10;",
			expectedSQL: "select distinct id from t1 where id > 0 order by c1 desc limit 10",
		},
		// insert into ... select ... group by ... limit ...语句
		{
			inputSQL:    "insert into t2 select c1,count(1) as cnt from t1 where id >0 group by c1 limit 10;",
			expectedSQL: "select c1, count(1) as cnt from t1 where id > 0 group by c1 limit 10",
		},
		// 基本的delete语句
		{
			inputSQL:    "delete from test1 where id=1;",
			expectedSQL: "select * from test1 where id = 1",
		},
		{
			inputSQL:    "delete from test1 where order_no='order005';",
			expectedSQL: "select * from test1 where order_no = 'order005'",
		},
		// delete ... where range
		{
			inputSQL:    "delete from test1 where created_time < '2024-03-20 17:57:36';",
			expectedSQL: "select * from test1 where created_time < '2024-03-20 17:57:36'",
		},
		{
			inputSQL:    "delete from t1 where id >0;",
			expectedSQL: "select * from t1 where id > 0",
		},
		// delete ... where ... limit ...
		{
			inputSQL:    "delete from t1 where id >0 limit 10;",
			expectedSQL: "select * from t1 where id > 0 limit 10",
		},
		// delete ... where ... order by ... limit ...
		{
			inputSQL:    "delete from t1 where id >0 order by c1 desc limit 10;",
			expectedSQL: "select * from t1 where id > 0 order by c1 desc limit 10",
		},
		// delete ... with subquery
		{
			inputSQL:    "DELETE FROM test2 WHERE order_no = (SELECT order_no FROM (SELECT order_no FROM test2 WHERE id = 4) AS temp);",
			expectedSQL: "select * from test2 where order_no = (select order_no from (select order_no from test2 where id = 4) as temp)",
		},
		// update simple
		{
			inputSQL:    "update t1 set c1=1 where id =10;",
			expectedSQL: "select * from t1 where id = 10",
		},
		// update with range
		{
			inputSQL:    "update t1 set c1=1 where id between 10 and 20;",
			expectedSQL: "select * from t1 where id between 10 and 20",
		},
		{
			inputSQL:    "update t1 set c1=1 where id >0;",
			expectedSQL: "select * from t1 where id > 0",
		},
		// update ... limit ...
		{
			inputSQL:    "update t1 set c1=1 where id >0 limit 10;",
			expectedSQL: "select * from t1 where id > 0 limit 10",
		},
		// update ... order by ... limit ...
		{
			inputSQL:    "update t1 set c1=1 where id >0 order by c1 desc limit 10;",
			expectedSQL: "select * from t1 where id > 0 order by c1 desc limit 10",
		},
		// update with multi tables join
		{
			inputSQL:    "update t1 inner join t2 on t1.id=t2.id2  set t1.c1=t2.c1 where c11=1;",
			expectedSQL: "select * from t1 join t2 on t1.id = t2.id2 where c11 = 1",
		},
		{
			inputSQL:    "update t1,t2 set t1.c1=t2.c1 where t1.id=t2.id2 and c11=1;",
			expectedSQL: "select * from t1, t2 where t1.id = t2.id2 and c11 = 1",
		},
		{
			inputSQL:    "update t1,t2 set t1.c1=t2.c1 where t1.id=t2.id2 and c11=1 limit 10;",
			expectedSQL: "select * from t1, t2 where t1.id = t2.id2 and c11 = 1 limit 10",
		},
		{
			inputSQL:    "REPLACE INTO test2 (id, order_no, created_time, updated_time) VALUES (1, 'replaced_order_no_1', '2022-01-01 10:00:00', '2022-01-01 10:00:00');",
			expectedSQL: "select * from t1 where id > 0 limit 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.inputSQL, func(t *testing.T) {
			rw, err := NewRewrite(tt.inputSQL)
			if err != nil {
				t.Fatalf("failed to create Rewrite: %v", err)
			}

			err = rw.RewriteDML2Select()
			if err != nil {
				t.Fatalf("Rewrite failed: %v", err)
			}

			assert.Equal(t, rw.SQL, tt.expectedSQL)
		})
	}
}
