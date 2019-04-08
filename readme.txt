返回状态：
data = {'code': code, 'data': msg|data}
code:
0：success
1：warning
2: error

return Response(data=data, status=status.HTTP_200_OK)


AJAX处理返回状态：
// 处理200状态
success: function (result) {
    if (result.code === 0) {
        displayPNotify(result.code, result.data)
    }
},
// 处理异常
error: function (jqXHR) {
    // 处理400状态，一般为serializer验证的错误
    if (jqXHR.status === 400) {
        let code = jqXHR.responseJSON['code'];
        let data = jqXHR.responseJSON['data'];
        displayPNotify(code, data);
    }
    // 处理403状态，一般为权限验证的错误
    if (jqXHR.status === 403) {
        let data = jqXHR.responseJSON['detail'];
        displayPNotify(1, data);
    }
}