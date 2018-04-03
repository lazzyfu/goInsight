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

// /**
//  * 显示通知
//  */
// function displayPNotify(status, msg, jump_url) {
//     var type = '';
//     var title = '';
//
//     // 接收后端返回的状态值
//     if (status == '200') {
//         type = 'success';
//         title = 'SUCCESS'
//     } else if (status == '201') {
//         type = 'notice';
//         title = 'NOTICE'
//     } else if (status == '400') {
//         type = 'error';
//         title = 'ERROR'
//     } else if (status == '403') {
//         type = 'info';
//         title = 'WARNING'
//     }
//
//     var set_null = function () {
//         return null
//     };
//
//
//     // 使用bootstrap样式
//     PNotify.prototype.options.styling = "bootstrap3";
//     new PNotify({
//         title: title,
//         text: msg,
//         type: type,
//         delay: 1500,
//         nonblock: {
//             nonblock: true
//         },
//         after_close: function () {
//             // 如果传入的变量不存在，则set_null
//             if (!init_fun) {
//                 set_null()
//             } else {
//                 init_fun()
//             }
//         }
//     });
// }

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
 * 获取当前用户所属的项目组
 */
$(function () {
    $.ajax({
        url: '/projects/group_info/',
        type: 'GET',
        dataType: 'json',
        timeout: 5000,
        cache: false,
        success: function (data) {
            $.each(data, function (index, data) {
                $("#select_group").append(  //此处向select中循环绑定数据
                    "<option data-icon='fa fa-th-large' value=" + data.group_id + ">" + data.group_name + "</option>");
            });
            $('.selectpicker').selectpicker('refresh');
        }
    })
});

/**
 *获取被选中项目的的DBA和Leader
 */
$(function () {
    $("#select_group").change(function () {
        var group_id = $(this).val();
        var csrftoken = $.cookie('csrftoken');
        $.ajax({
            url: '/projects/audit_user/',
            type: 'POST',
            dataType: 'json',
            timeout: 5000,
            data: {
                'group_id': group_id,
                'csrfmiddlewaretoken': csrftoken
            },
            cache: false,
            success: function (data) {
                $("#select_dba").empty();
                $("#select_leader").empty();
                $.each(data, function (index, data) {
                    if (data.user_role == 'DBA') {
                        $("#select_dba").append(  //此处向select中循环绑定数据
                            "<option data-icon='fa fa-user' data-subtext=" + data.email + " value=" + data.username + ">" + data.username + "</option>");
                    }
                    if (data.user_role == 'Leader') {
                        $("#select_verifier").append(  //此处向select中循环绑定数据
                            "<option data-icon='fa fa-user' data-subtext=" + data.email + " value=" + data.username + ">" + data.username + "</option>");
                    }
                });
                $('.selectpicker').selectpicker('refresh');
            }
        })
    });
    $("select_group").validator('update');
});

/**
 * 获取用户选择的项目的联系人
 */
$(function () {
    $("#select_group").change(function () {
        var group_id = $(this).val();
        var csrftoken = $.cookie('csrftoken');
        $.ajax({
            url: '/projects/contacts_info/',
            type: 'POST',
            dataType: 'json',
            timeout: 5000,
            data: {
                'group_id': group_id,
                'csrfmiddlewaretoken': csrftoken
            },
            cache: false,
            success: function (data) {
                $("#select_contact").empty();
                $.each(data, function (index, row) {
                    var html = '';
                    for (var i = 0; i < row.split(",").length; i++) {
                        var result = row.split(",")[i];
                        var username = result.split(":")[0];
                        var contact_id = result.split(":")[1];
                        var email = result.split(":")[2];
                        html += "<option data-icon='fa fa-user' data-subtext=" + email + " value=" + contact_id + ">" + username + "</option>";
                    }
                    $("#select_contact").append(html)
                });
                $('.selectpicker').selectpicker('refresh');
            }
        })
    })
});

/**
 * 获取备注信息
 */
$(function () {
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/remark_info/',
        type: 'POST',
        dataType: 'json',
        timeout: 5000,
        data: {
            'csrfmiddlewaretoken': csrftoken
        },
        cache: false,
        success: function (data) {
            $.each(data, function (index, data) {
                $("#select_remark").append(  //此处向select中循环绑定数据
                    "<option data-icon='fa fa-tag' value=" + data.remark + ">" + data.remark + "</option>");
            });
            $('.selectpicker').selectpicker('refresh');
        }
    })
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
        success: function (data) {
            if (data.status === 0) {
                var html = '';
                $.each(data.msg, function (index, db) {
                    html += "<option value=" + db + ">" + db + "</option>"
                });
                $('#select_db').append(html);
                $('.selectpicker').selectpicker('refresh')
            }
            else {
                displayPNotify(data.status, data.msg);
            }
        }
    })
}