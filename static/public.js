/**
 * Created by fuzongfei on 2017/9/15.
 */

/**
 * 刷新当前页面
 */
function refresh_page() {
    window.location.reload()
}

/**
 * 移除初始化的通知
 */
$(document).ready(function () {
    $('.ui-pnotify').remove();
});


/**
 * 显示通知
 */

function displayPNotify(status, msg) {
    var opts = {
        title: "Over Here",
        text: msg,
        delay: 2000,
        shadow: true,
        nonblock: {
            nonblock: true
        },
        animate: {
            animate: true,
            in_class: 'zoomInLeft',
            out_class: 'zoomOutRight'
        }
    };
    switch (status) {
        case 0:
            opts.title = "成功";
            opts.type = "success";
            break;
        case 1:
            opts.title = "通知";
            opts.type = "notice";
            break;
        case 2:
            opts.title = "错误";
            opts.type = "error";
            break;
        case 403:
            opts.title = "403";
            opts.type = "error";
            opts.text = "权限拒绝";
            break;
    }
    new PNotify(opts);
}

// function displayPNotify(status, msg) {
//     var type = '';
//     var title = '';
//
//     // 接收后端返回的状态值
//     if (status === 0) {
//         type = 'success';
//         title = '成功'
//     } else if (status === 1) {
//         type = 'notice';
//         title = '通知'
//     } else if (status === 2) {
//         type = 'error';
//         title = '错误'
//     } else if (status === 403) {
//         type = 'info';
//         title = '403';
//         msg = '权限拒绝，您没有权限操作'
//     }
//
//     new PNotify({
//         title: title,
//         text: msg,
//         type: type,
//         delay: 2000,
//         shadow: true,
//         // styling: 'bootstrap3',
//         nonblock: {
//             nonblock: true
//         },
//         animate: {
//             animate: true,
//             in_class: 'zoomInLeft',
//             out_class: 'zoomOutRight'
//         }
//     });
// }


/**
 * 生成随机字符串
 */
function random_str(len) {
    len = len || 1;
    var $chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!@#$%^&*';
    var maxPos = $chars.length;
    var pwd = '';
    for (var i = 0; i < len; i++) {
        pwd += $chars.charAt(Math.floor(Math.random() * maxPos));
    }
    return pwd
}


/**
 * 美化SQL
 */
function beautifySQL() {
    var sql_content = myCodeMirror.getValue();
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/beautify_sql/',
        type: 'POST',
        dataType: 'json',
        data: {'sql_content': sql_content, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: false,
        success: function (result) {
            if (result.status === 2) {
                displayPNotify(result.status, result.msg);
            }
            else {
                $.each(result, function (index, sql) {
                    myCodeMirror.setValue(sql);
                });
            }
        }
    })
}


/**
 * 执行inception语法检查
 */

function incepSyntaxCheckForm() {
    var contents = myCodeMirror.getValue();
    var database = $("#s_database").val();
    var host = $("#s_host").val();
    var operate_type = $("#s_operate").val();
    var csrftoken = $.cookie('csrftoken');

    // 判断输入的SQL内容是否存在
    if (!contents) {
        displayPNotify(2, '请输入要检测的SQL语句');
        return false
    }

    $.ajax({
        url: '/projects/syntax_check/',
        type: 'POST',
        dataType: 'json',
        data: {
            'contents': contents,
            'database': database,
            'host': host,
            'operate_type': operate_type,
            'csrfmiddlewaretoken': csrftoken
        },
        timeout: 5000,
        cache: false,
        success: function (result) {
            if (result.status === 2) {
                displayPNotify(result.status, result.msg)
            }
            else {
                $('#td_append').empty();
                var html = '';
                document.getElementById('typediv1').style.visibility = "visible";
                $.each(result.data, function (index, content) {
                    var SQL = content.SQL;
                    var ID = content.ID;
                    var stage = content.stage;
                    var errlevel = '';
                    if (content.errlevel === 0) {
                        errlevel = '成功'
                    }
                    else if (content.errlevel === 1) {
                        errlevel = '警告'
                    }
                    else if (content.errlevel === 2) {
                        errlevel = '错误'
                    }
                    var stagestatus = content.stagestatus;
                    var errormessage = content.errormessage;
                    var Affected_rows = content.Affected_rows;
                    var execute_time = content.execute_time;
                    html += "<tr>" +
                        "<td>" + ID + "</td>" +
                        "<td>" + stage + "</td>" +
                        "<td>" + stagestatus + "</td>" +
                        "<td>" + errlevel + "</td>" +
                        "<td>" + errormessage + "</td>" +
                        "<td>" + SQL + "</td>" +
                        "<td>" + Affected_rows + "</td>" +
                        "<td>" + execute_time + "</td>" +
                        "</tr>";
                });
                $('#td_append').append(html);
            }
        }
    })
}


/**
 * 获取指定主机的schema信息和tables信息
 */
function getTablesList() {
    $('#s_database').empty();
    var host = $("#s_host").val();
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/schema_info/',
        type: 'POST',
        dataType: 'json',
        data: {'host': host, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: true,
        success: function (result) {
            if (result.status === 0) {
                var html = '';
                $.each(result.data.schema, function (index, db) {
                    html += "<option value=" + db + ">" + db + "</option>"
                });
                $('#s_database').append(html);
                $('.selectpicker').selectpicker('refresh');

                myCodeMirror.setOption("hintOptions", {tables: result.data.tables});
            }
            else {
                // 此处必须刷新表格
                $('.selectpicker').selectpicker('refresh');
                displayPNotify(result.status, result.msg);
            }
        }
    });
}

/**
 * jquery loading
 */
function showLoadingScreen() {
    $('body').loading()
}

function hideLoadingScreen() {
    $('body').loading('stop');
}