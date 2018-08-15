status : 返回状态
msg : 返回消息
data: 返回的数据
jump_url：返回的目标url

status:
0 : success
1 : notice
2 : error
403: permission deny

格式：
context = {'status': 2, 'msg': 'SQL语法检查未通过', 'data': {'name': 'zs'}}