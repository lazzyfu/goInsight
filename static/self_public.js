/**
 * Created by fuzongfei on 2017/9/15.
 */

/**
 * 刷新当前页面
 */
var refresh_page = function () {
    window.location.reload()
};

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
    var type = '';
    var title = '';

    // 接收后端返回的状态值
    if (status === 0) {
        type = 'success';
        title = 'SUCCESS'
    } else if (status === 1) {
        type = 'notice';
        title = 'NOTICE'
    } else if (status === 2) {
        type = 'error';
        title = 'ERROR'
    }

    PNotify.prototype.options.styling = "bootstrap3";

    new PNotify({
        title: title,
        text: msg,
        type: type,
        delay: 1500,
        nonblock: {
            nonblock: true
        },
        animate: {
            animate: true,
            in_class: 'rotateInDownLeft',
            out_class: 'rotateOutUpRight'
        }
    });
}


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
    var sql_content = myCodeMirror.getValue();
    var database = $("#select_db").val();
    var host = $("#select_env").val();
    var op_action = $("#select_op").val();
    var csrftoken = $.cookie('csrftoken');

    // 判断输入的SQL内容是否存在
    if (!sql_content) {
        displayPNotify(2, '请输入要检测的SQL语句');
        return false
    }

    $.ajax({
        url: '/projects/syntax_check/',
        type: 'POST',
        dataType: 'json',
        data: {
            'sql_content': sql_content,
            'database': database,
            'host': host,
            'op_action': op_action,
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
 * 生成工单
 */

<!-- 有效性验证 -->
$('#auditCommitForm').validator().on('submit', function (e) {
    var sql_content = myCodeMirror.getValue();
    if (sql_content.length < 10) {
        displayPNotify(2, '审核内容不能为空或小于10个字符');
        return false;
    }
    if (e.isDefaultPrevented()) {
        // 验证不通过
        displayPNotify(2, '表单无效，请完成填写');
    } else {
        // 验证通过
        $('#auditCommitForm').ajaxSubmit({
            data: {'sql_content': sql_content},
            dataType: 'json',
            success: function (result) {
                if (result.status === 0) {
                    window.parent.location.href = result.jump_url
                }
                else {
                    displayPNotify(result.status, result.msg)
                }
            }
        });
        return false;
    }
});

/**
 * 获取指定主机的数据库库名列表
 */
function getDatabaseList() {
    $("#select_db").empty();
    var host = $("#select_env").val();
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/db_list/',
        type: 'POST',
        dataType: 'json',
        data: {'host': host, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: false,
        success: function (result) {
            if (result.status === 0) {
                var html = '';
                $.each(result.data, function (index, db) {
                    html += "<option value=" + db + ">" + db + "</option>"
                });
                $('#select_db').append(html);
                $('.selectpicker').selectpicker('refresh')
            }
            else {
                displayPNotify(result.status, result.msg);
            }
        }
    })
}