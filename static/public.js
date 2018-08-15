/**
 * Created by fuzongfei on 2017/9/15.
 */

/**
 * 刷新当前页面
 */
function refresh_page() {
    setTimeout(window.location.reload(), 1000)
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
    var host = $('#s_host').val();
    var database = $('#s_schema').val();
    var operate_type = $('#s_operate').val();
    var contents = myCodeMirror.getValue();
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
                displayPNotify(result.status, result.msg);
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
                        errlevel = '成功';
                    }
                    else if (content.errlevel === 1) {
                        errlevel = '警告';
                    }
                    else if (content.errlevel === 2) {
                        errlevel = '错误';
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
 * 将审核后的工单转换成执行任务
 */
function general_perform_tasks(id, envi_desc) {
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        type: 'post',
        url: "/projects/ol/generate_perform_tasks/",
        dataType: 'json',
        data: {
            'id': id,
            'envi_desc': envi_desc,
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
function work_order_operate_form(id, url, btn_left, btn_right) {
    layui.use('layedit', function () {
        layui.use('layer', function () {
            var layer = layui.layer;
            var csrftoken = $.cookie('csrftoken');

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
                    var addition_info = $('#addition_info').val();
                    var status = btn_left;
                    layer.close(layer.index);
                    ajaxCommit(addition_info, status)
                },

                btn2: function () {
                    var addition_info = $('#addition_info').val();
                    var status = btn_right;
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
var ding_tasks = '';
var modal_tasks_show_selector = $('#modal_tasks_show');

function deploy_tasks_table(tasks) {
    modal_tasks_show_selector.modal('show');
    var $tasks_table = $('#tasks-table');
    // 此处必须destroy table，否则会加载旧数据
    $tasks_table.bootstrapTable('destroy', {silent: true});
    $(function () {
        $tasks_table.bootstrapTable({
            method: 'get',
            dataType: 'json',
            contentType: "application/x-www-form-urlencoded",
            url: '/projects/ol/deploy_tasks/',
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
            uniqueId: "id",
            classes: 'table table-hover table-no-bordered',
            rowStyle: render_row_style,
            queryParams: function (params) {
                ding_tasks = tasks;
                return {
                    tasks: tasks
                }
            },

            columns: [
                {
                    field: 'tasks',
                    title: '任务'
                },
                {
                    field: 'title',
                    title: '标题',
                    formatter: function (value, row) {
                        return '<a target=\'_blank\' href=\'/projects/ol/ol_records/work_order_details/' + row.id + "'>" + value + "</a>"
                    }
                },
                {
                    field: 'proposer',
                    title: '申请人'
                },
                {
                    field: 'test',
                    title: '测试环境',
                    formatter: render_finish_status
                },
                {
                    field: 'staging',
                    title: 'Staging环境',
                    formatter: render_finish_status
                },
                {
                    field: 'product',
                    title: '生产环境',
                    formatter: render_finish_status
                }
            ]
        });
    });

    $('#s_tasks').empty();
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/ol/deploy_tasks/',
        type: 'POST',
        dataType: 'json',
        data: {'tasks': tasks, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: false,
        success: function (result) {
            $('#s_tasks').append(result.num)
        }
    })
}

function render_finish_status(value) {
    var status = parseInt(value);
    var finish_status = [4, 6];
    var unfinish_status = [0, 1, 2, 3, 5];
    var uncommit_status = [-1];
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

/**
 * 渲染行的样式
 * 已完成的自动标记为：绿色(success)
 * 关闭的自动标记为：红色(danger)
 */
function render_row_style(row, index) {
    var finish_status = [4, 6];
    var test = finish_status.indexOf(parseInt(row.test));
    var staging = finish_status.indexOf(parseInt(row.staging));
    var product = finish_status.indexOf(parseInt(row.product));
    if (test >= 0 && staging >= 0 && product >= 0) {
        return {classes: 'success'}
    } else if (parseInt(row.test) === 5 || parseInt(row.staging) === 5 || parseInt(row.product) === 5) {
        return {classes: 'danger'}
    }
    return {};
}

/**
 * 通知未完成工单的开发
 */
function dingNotice() {
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/ol/ding_notice/',
        type: 'POST',
        dataType: 'json',
        data: {'tasks': ding_tasks, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: false,
        success: function (result) {
            displayPNotify(result.status, result.msg)
        }
    })
}

/**
 * 获取指定主机表的元数据信息，tab补全使用
 */
function get_table_meta_info(tab_schema) {
    var csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/projects/get_table_meta_info/',
        type: 'POST',
        dataType: 'json',
        data: {'schema': tab_schema, 'csrfmiddlewaretoken': csrftoken},
        timeout: 30000,
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
 * 执行MySQL查询语句
 */
$('#SqlQueryForm').on('submit', function (e) {
    // 清空输出的结果
    $('#li_append').empty();
    $('#table_append').empty();

    // 获取选中的内容，否则为全部内容
    var contents = '';
    if (myCodeMirror.somethingSelected()) {
        contents = myCodeMirror.getSelection()
    } else {
        contents = myCodeMirror.getValue();
    }
    // 判断输入的内容是否为空
    if (contents.length < 1) {
        displayPNotify(2, '内容不能为空');
        return false;
    }

    if (!query_schema) {
        displayPNotify(2, '请点击选中左侧库或表，在执行查询');
        return false;
    }
    $(this).ajaxSubmit({
        data: {'contents': contents, 'schema': query_schema},
        dataType: 'json',
        beforeSubmit: showLoadingScreen($('body'), '数据查询中，请稍后...'),
        success: function (result) {
            hideLoadingScreen($('body'));
            if (result.status === 0) {
                var data = result.data;
                document.getElementById('typediv1').style.visibility = "visible";

                var li_html = '';
                var table_html = '';
                for (var i in data) {
                    li_html += "<li><a href=\"#tab_" + i + "\" data-toggle=\"tab\">" + "结果" + i + "</a></li>";
                    table_html += "<div class=\"tab-pane\" id=\"tab_" + i + "\">\n" +
                        "<table id=\"table" + i + "\"></table>\n" +
                        "</div>"
                }

                $('#table_append').append(table_html);
                $('#li_append').append(li_html);
                $('.nav-tabs>li>a').first().trigger('click');


                for (var key in data) {
                    var d = data[key];
                    var $table = $("#table" + key);

                    $table.bootstrapTable({
                        columns: d.columnDefinition,
                        data: d.data,
                        search: true,
                        showColumns: true,
                        showRefresh: true,
                        showToggle: true,
                        showExport: true,
                        pageNumber: 1,
                        pageSize: 20,
                        locale: 'zh-CN',
                        pageList: [20, 30, 50],
                        sidePagination: "client",
                        pagination: true,
                        singleSelect: true,
                        minimumCountColumns: 2,
                        matchBrackets: true,
                        lineWrapping: true,
                        rowStyle: function rowStyle(row, index) {
                            return {
                                classes: 'text-nowrap another-class',
                                css: {"font-size": "12px"}
                            };
                        },
                        classes: 'table table-hover'
                    });
                }
            }
            else {
                document.getElementById('typediv1').style.visibility = "hidden";
                displayPNotify(result.status, result.msg)
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status);
                hideLoadingScreen();
            }
        }
    });
    return false;
});