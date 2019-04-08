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
 * jquery loading
 */
function showLoadingScreen(tag, msg) {
    tag.loading({message: msg});
}

function hideLoadingScreen(tag) {
    tag.loading('stop');
}


/**
 * 显示通知
 */

function displayPNotify(code, msg) {
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
            in_class: 'slideInDown',
            out_class: 'slideOutUp'
        }
    };
    switch (code) {
        case 0:
            opts.title = "成功";
            opts.type = "success";
            break;
        case 1:
            opts.title = "警告";
            opts.type = "info";
            break;
        case 2:
            opts.title = "错误";
            opts.type = "error";
            break;
    }
    new PNotify(opts);
}

/**
 * 格式化SQL
 */
function beautifySQL() {
    let contents = myCodeMirror.getValue();
    console.log(contents)
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/orders/beautify_sql/',
        type: 'POST',
        dataType: 'json',
        data: {'contents': contents, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: false,
        success: function (result) {
            if (result.code === 0) {
                myCodeMirror.setValue(result.data);
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 400) {
                let code = jqXHR.responseJSON['code'];
                let data = jqXHR.responseJSON['data'];
                displayPNotify(code, data);
            }
        }
    })
}


// 获取指定环境的schema
function getSchema() {
    let s_schema = $('#s_schema');
    let envi_id = $('#s_envi').val();
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/orders/get_schemas/',
        type: 'POST',
        dataType: 'json',
        data: {'envi_id': envi_id, 'csrfmiddlewaretoken': csrftoken},
        timeout: 5000,
        cache: false,
        success: function (result) {
            if (result.code === 0) {
                s_schema.empty();
                let html = '';
                $.each(result.data, function (index, row) {
                    let result = [row.host, row.port, row.schema].join(',');
                    let show_result = [row.comment, row.port, row.schema].join('_');
                    html += "<option data-icon='fa fa-database' value=" + result + ">" + show_result + "</option>"
                });
                s_schema.append(html);
                $('.selectpicker').selectpicker('refresh')
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 400) {
                let code = jqXHR.responseJSON['code'];
                let data = jqXHR.responseJSON['data'];
                displayPNotify(code, data);
            }
        }
    })
}

// 获取系统环境
function getSysEnvironment() {
    $.ajax({
        url: '/orders/get_sysenvi/',
        type: 'GET',
        dataType: 'json',
        timeout: 10000,
        cache: false,
        success: function (data) {
            $("#s_envi").empty();
            $.each(data, function (index, row) {
                $("#s_envi").append(
                    "<option data-icon='glyphicon-record' value=" + row.envi_id + ">" + row.envi_name + "</option>"
                );
            });
            $('.selectpicker').selectpicker('refresh')
        },
    });
}

// 获取任务名
function getOnlineVersion() {
    $.ajax({
        url: '/orders/online_version/no_expire/',
        type: 'GET',
        dataType: 'json',
        timeout: 5000,
        cache: false,
        success: function (data) {
            $.each(data, function (index, data) {
                let option_data = '';
                if (data.is_disable === 1) {
                    option_data = "<option data-subtext='已过期' data-icon='glyphicon-flag' value=" + data.id + " disabled>" + data.version + "</option>";
                } else if (data.is_disable === 0) {
                    option_data = "<option data-icon='glyphicon-flag' value=" + data.id + ">" + data.version + "</option>";
                }
                $("#s_task").append(option_data);
                $('.selectpicker').selectpicker('refresh');
            })
        }
    });
}

// 获取抄送的用户和邮箱
function getEmailCc() {
    $.ajax({
        url: '/users/get_email_cc',
        type: 'GET',
        dataType: 'json',
        timeout: 10000,
        cache: true,
        success: function (data) {
            let html = '';
            $.each(data, function (index, row) {
                html += "<option data-icon='glyphicon-user' data-subtext=" + row.email + " value=" + row.username + ">" + row.username + "</option>"
            });
            $("#s_email_cc").html(html);
            $("#s_reviewer").html(html);
            $('.selectpicker').selectpicker('refresh');
        }
    })
}

// 获取有审核权限的用户
function getAuditor(permission) {
    let allowed_permissions = ['can_audit'];
    if (allowed_permissions.indexOf(permission) < 0) {
        displayPNotify(2, '传入的权限参数错误');
        return false
    }
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/users/get_auditor/',
        type: 'POST',
        dataType: 'json',
        data: {'permission': permission, 'csrfmiddlewaretoken': csrftoken},
        timeout: 10000,
        cache: false,
        success: function (data) {
            let $this = $("#s_auditor");
            $this.empty();
            $.each(data, function (index, row) {
                $this.append(
                    "<option data-icon='glyphicon-user' data-subtext=" + row.displayname + " value=" + row.username + ">" + row.username + "</option>"
                );
            });
            $('.selectpicker').selectpicker('refresh');
        }
    })
}

// 语法规则检查
function incepSyntaxCheck(schema, sql_type, contents) {
    let csrftoken = $.cookie('csrftoken');

    // 判断输入的SQL内容是否存在
    if (!contents) {
        displayPNotify(2, '请输入要检测的SQL语句');
        return false
    }

    $.ajax({
        url: '/orders/syntax_check/',
        type: 'POST',
        dataType: 'json',
        data: {
            'schema': schema,
            'sql_type': sql_type,
            'contents': contents,
            'csrfmiddlewaretoken': csrftoken
        },
        timeout: 10000,
        cache: false,
        success: function (result) {
            console.log(result);
            if (result.code === 2) {
                displayPNotify(result.code, result.data);
            }
            if (result.code === 0) {
                let html = '';
                document.getElementById('typediv1').style.visibility = "visible";
                $.each(result.data, function (index, content) {
                    let SQL = content.SQL;
                    let ID = content.ID;
                    let stage = content.stage;
                    let errlevel = '';
                    if (content.errlevel === 0) {
                        errlevel = '成功';
                    } else if (content.errlevel === 1) {
                        errlevel = '警告';
                    } else if (content.errlevel === 2) {
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
                $('#td_append').html(html);
            }
        }
    })
}

// 提交工单
$('#OrdersCommitForm').validator().on('submit', function (e) {
    let contents = myCodeMirror.getValue();

    if (contents.length < 10) {
        displayPNotify(2, '审核内容不能为空或小于10个字符');
        return false;
    }
    if (e.isDefaultPrevented()) {
        // 验证不通过
        displayPNotify(2, '表单无效，请完成填写');
    } else {
        // 验证通过
        $(this).ajaxSubmit({
            data: {'contents': contents},
            dataType: 'json',
            success: function (result) {
                if (result.code === 0) {
                    window.parent.location.href = result.data
                } else {
                    displayPNotify(result.code, result.data)
                }
            },
            error: function (jqXHR) {
                if (jqXHR.status === 400) {
                    let code = jqXHR.responseJSON['code'];
                    let data = jqXHR.responseJSON['data'];
                    displayPNotify(code, data);
                }
                if (jqXHR.status === 403) {
                    let data = jqXHR.responseJSON['detail'];
                    displayPNotify(1, data);
                }
            }
        });
        return false;
    }
});


// 用户操作，如：审核、执行、反馈等，需要触发的动作
function ordersOp(id, url, left_btn, right_btn) {
    let csrftoken = $.cookie('csrftoken');
    swal({
        buttons: {
            left_btn: {
                text: left_btn,
                value: left_btn,
                visible: true
            },
            right_btn: {
                text: right_btn,
                value: right_btn,
                visible: true
            },
        },
        text: '输入附加的信息',
        content: {
            element: "input",
            attributes: {
                placeholder: "",
                type: "text",
                require: false
            },
        },
    })
        .then(value => {
            // value为null时，关闭窗口
            if (!value) {
                return false
            }
            // 获取用户输入的信息
            let input_val = $('.swal-content__input').val();

            $.ajax({
                url: url,
                type: 'POST',
                dataType: 'json',
                data: {
                    'id': id,
                    'status': value,
                    'msg': input_val,
                    'csrfmiddlewaretoken': csrftoken
                },
                timeout: 10000,
                cache: false,
                success: function (result) {
                    displayPNotify(result.code, result.data)
                },
                error: function (jqXHR) {
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
            })
        })
}

// 将工单生成子任务
function generate_subtasks(id, envi_id) {
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        type: 'post',
        url: "/orders/generate_subtasks/",
        dataType: 'json',
        beforeSubmit: showLoadingScreen($('body'), '数据处理中，请稍后...'),
        data: {
            'id': id,
            'envi_id': envi_id,
            'csrfmiddlewaretoken': csrftoken
        },
        success: function (result) {
            hideLoadingScreen($('body'));
            if (result.code === 0) {
                window.location.href = result.data
            } else {
                displayPNotify(result.code, result.data)
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 400) {
                let code = jqXHR.responseJSON['code'];
                let data = jqXHR.responseJSON['data'];
                displayPNotify(code, data);
                setTimeout("hideLoadingScreen($('body'))", 2000);
            }
            // 处理403状态，一般为权限验证的错误
            if (jqXHR.status === 403) {
                let data = jqXHR.responseJSON['detail'];
                displayPNotify(1, data);
                setTimeout("hideLoadingScreen($('body'))", 2000);
            }
        }
    });
}


/**
 * 渲染上线任务列表
 */
let modal_OnlineVersion_selector = $('#modal_tasks_show');

function load_online_version_list(version) {
    modal_OnlineVersion_selector.modal('show');
    let $tasks_table = $('#tasks-table');
    let csrftoken = $.cookie('csrftoken');
    $.ajax({
        url: '/orders/online_version/detail/',
        type: 'POST',
        dataType: 'json',
        data: {'version': version, 'csrfmiddlewaretoken': csrftoken},
        timeout: 10000,
        cache: false
    }).done(function (d) {
        $tasks_table.bootstrapTable('destroy').bootstrapTable({
            classes: 'table table-hover table-no-bordered',
            cache: false,
            showColumns: true,
            pagination: true,
            search: true,
            showRefresh: true,
            minimumCountColumns: 2,
            pageNumber: 1,
            pageSize: 10,
            locale: 'zh-CN',
            pageList: [10, 20],
            columns: d.columns,
            data: d.data
        })
    });
}

function render_onlinetasks_status(value) {
    let status = parseInt(value);
    let uncommit_status = [-1];
    if (uncommit_status.indexOf(status) >= 0) {
        return "<i class='fa fa-close' style='color:red'></i>"
    }
    if (status === 0) {
        return "<a href='#'>待批准</a>"
    }
    if (status === 1) {
        return "<a href='#' class='text-red'>未批准</a>"
    }
    if (status === 2) {
        return "<a href='#'>已批准</a>"
    }
    if (status === 3) {
        return "<a href='#'>处理中</a>"
    }
    if (status === 4) {
        return "<a href='#' class='text-green'>已完成</a>"
    }
    if (status === 5) {
        return "<a href='#' class='text-red'>已关闭</a>"
    }
    if (status === 6) {
        return "<a href='#' class='text-green'>已复核</a>"
    }
    if (status === 7) {
        return "<a href='#' class='text-green'>已勾住</a>"
    }
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
        url: '/query/get_tables/',
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
