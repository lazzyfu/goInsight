# -*- coding:utf-8 -*-
# edit by fuzongfei


# DB类别
rdsCategory = (
    (1, 'mysql'),
    (2, 'tidb'),
    (3, 'clickhouse'),
)

# 提供商
rdsTypeChoice = (
    (0, '华为云RDS'),
    (1, '阿里云RDS'),
    (2, '自建RDS')
)

# 用途
useTypeChoice = (
    (0, 'SQL审核'),
    (1, 'SQL查询'),
)

# 字符集
characterChoice = (
    ('utf8', 'utf8'),
    ('utf8mb4', 'utf8mb4')
)

# 操作类型选择
sqlTypeChoice = (
    ('DDL', u'DDL工单'),
    ('DML', u'DML工单'),
    ('EXPORT', u'导出工单')
)

# 导出工单支持的文件格式
fileFormatChoice = (
    ('xlsx', 'xlsx格式'),
    ('csv', 'csv格式')
)

# 工单备注
orderRemark = (
    ('立即执行', '立即执行'),
    ('窗口执行', '窗口执行'),
    ('上线执行', '上线执行')
)

# SQL工单状态
sqlProgressChoice = (
    (0, u'待审核'),
    (1, u'已驳回'),
    (2, u'已批准'),
    (3, u'处理中'),
    (4, u'已完成'),
    (5, u'已关闭'),
    (6, u'已复核'),
    (7, u'已勾住')
)

# SQL工单任务执行状态
taskProgressChoice = (
    (0, u'未执行'),
    (1, u'已完成'),
    (2, u'处理中'),
    (3, u'失败'),
    (4, u'暂停'),
)
