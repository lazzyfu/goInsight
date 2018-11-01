/**
 * Created by fuzongfei on 2017/9/15.
 */

/**
 * 刷新当前页面
 */
function refresh_page() {
    setTimeout("window.location.reload()", 3000)
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
    let opts = {
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


/**
 * 获取指定主机schema的tables信息
 */
function get_table_info(schema) {
    // 格式：host,port,schema
    // 即：10.10.10.1,3306,tbl_name
    if (!schema) {
        schema = $("#s_schema").val();
    }
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/sqlorders/get_tables/',
        type: 'POST',
        dataType: 'json',
        data: {'schema': schema, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: true,
        success: function (result) {
            if (result.status === 0) {
                myCodeMirror.setOption("hintOptions", {tables: result.data.tables});
            }
            else {
                displayPNotify(result.status, result.msg);
            }
        }
    });
}

/**
 * 美化SQL
 */
function beautifySQL() {
    let contents = myCodeMirror.getValue();
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/sqlorders/beautify_sql/',
        type: 'POST',
        dataType: 'json',
        data: {'contents': contents, 'csrfmiddlewaretoken': csrftoken},
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
    let host = $('#s_schema').val();
    let sql_type = $('#s_sql_type').val();
    let contents = myCodeMirror.getValue();
    let csrftoken = $.cookie('csrftoken');

    // 判断输入的SQL内容是否存在
    if (!contents) {
        displayPNotify(2, '请输入要检测的SQL语句');
        return false
    }

    $.ajax({
        url: '/sqlorders/syntax_check/',
        type: 'POST',
        dataType: 'json',
        data: {
            'contents': contents,
            'host': host,
            'sql_type': sql_type,
            'csrfmiddlewaretoken': csrftoken
        },
        timeout: 5000,
        cache: false,
        success: function (result) {
            if (result.status === 2) {
                displayPNotify(result.status, result.msg);
            }
            else {
                $('#td_append').empty();
                let html = '';
                document.getElementById('typediv1').style.visibility = "visible";
                $.each(result.data, function (index, content) {
                    let SQL = content.SQL;
                    let ID = content.ID;
                    let stage = content.stage;
                    let errlevel = '';
                    if (content.errlevel === 0) {
                        errlevel = '成功';
                    }
                    else if (content.errlevel === 1) {
                        errlevel = '警告';
                    }
                    else if (content.errlevel === 2) {
                        errlevel = '错误';
                    }
                    let stagestatus = content.stagestatus;
                    let errormessage = content.errormessage;
                    let Affected_rows = content.Affected_rows;
                    let execute_time = content.execute_time;
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
 * 将审核后的工单转换成执行任务
 */
function general_perform_tasks(id, envi_id) {
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        type: 'post',
        url: "/sqlorders/generate_perform_tasks/",
        dataType: 'json',
        data: {
            'id': id,
            'envi_id': envi_id,
            'csrfmiddlewaretoken': csrftoken
        },
        success: function (result) {
            if (result.status === 0) {
                window.location.href = result.jump_url
            }
            else {
                displayPNotify(result.status, result.msg)
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status)
            }
        }
    });
}

/**
 * 工单操作，审批、反馈、关闭
 */
function sql_orders_audit(id, url, btn_left, btn_right) {
    layui.use('layedit', function () {
        layui.use('layer', function () {
            let layer = layui.layer;
            let csrftoken = $.cookie('csrftoken');

            function ajaxCommit(addition_info, status) {
                $.ajax({
                    type: 'post',
                    url: url,
                    dataType: 'json',
                    data: {
                        'id': id,
                        'status': status,
                        'addition_info': addition_info,
                        'csrfmiddlewaretoken': csrftoken
                    },
                    success: function (result) {
                        displayPNotify(result.status, result.msg);
                    },
                    error: function (jqXHR) {
                        if (jqXHR.status === 403) {
                            displayPNotify(jqXHR.status)
                        }
                    }
                });
            }

            layer.open({
                title: '附加信息：',
                type: 0,
                resize: false,
                content: '<input type="text" id="addition_info" style="margin: 0px; width: 250px;">',
                btn: [btn_left, btn_right],
                btnAlign: 'c',

                yes: function () {
                    let addition_info = $('#addition_info').val();
                    let status = btn_left;
                    layer.close(layer.index);
                    ajaxCommit(addition_info, status)
                },

                btn2: function () {
                    let addition_info = $('#addition_info').val();
                    let status = btn_right;
                    ajaxCommit(addition_info, status)
                },
                // 右上角关闭回调
                cancel: function () {
                }
            })
            ;
        });
    })
}


/**
 * jquery loading
 */
function showLoadingScreen(tag, msg) {
    tag.loading({message: msg});
}

function hideLoadingScreen(tag) {
    tag.loading('stop');
}


/**
 * 渲染任务表格
 */
let ding_tasks = '';
let modal_tasks_show_selector = $('#modal_tasks_show');

function deploy_tasks_table(tasks) {
    modal_tasks_show_selector.modal('show');
    let $tasks_table = $('#tasks-table');
    // 此处必须destroy table，否则会加载旧数据
    $tasks_table.bootstrapTable('destroy', {silent: true});

    $.ajax({
        url: '/sqlorders/get_version_orders_list',
        type: 'GET',
        dataType: 'json',
        data: {'tasks': tasks},
        timeout: 5000,
        cache: false
    }).done(function (d) {
        $tasks_table.bootstrapTable({
            classes: 'table table-hover table-no-bordered',
            cache: false,
            showColumns: true,
            pagination: true,
            search: true,
            showRefresh: true,
            minimumCountColumns: 2,
            pageNumber: 1,
            pageSize: 20,
            locale: 'zh-CN',
            pageList: [10, 20],
            columns: d.columns,
            data: d.data
        })
    });
}

function render_finish_status(value) {
    let status = parseInt(value);
    let finish_status = [4, 6];
    let unfinish_status = [0, 1, 2, 3, 5];
    let uncommit_status = [-1];
    if (finish_status.indexOf(status) >= 0) {
        return "<i class='fa fa-check' style='color:green'></i>"
    } else if (unfinish_status.indexOf(status) >= 0) {
        if (status === 0) {
            return "<a href='#'>待批准</a>"
        } else if (status === 1) {
            return "<a href='#'>未批准</a>"
        } else if (status === 2) {
            return "<a href='#'>已批准</a>"
        } else if (status === 3) {
            return "<a href='#'>处理中</a>"
        } else if (status === 5) {
            return "<a href='#'>已关闭</a>"
        }
    } else if (uncommit_status.indexOf(status) >= 0) {
        return "<i class='fa fa-close' style='color:red'></i>"
    }
}

