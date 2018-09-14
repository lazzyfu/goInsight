
function csrfSafeMethod(method) {
    // these HTTP methods do not require CSRF protection
    return (/^(GET|HEAD|OPTIONS|TRACE)$/.test(method));
}

/*
 全局Ajax中添加请求头X-CSRFToken，用于跨过CSRF验证
 */
$.ajaxSetup({
    beforeSend: function (xhr, settings) {
        if (!csrfSafeMethod(settings.type) && !this.crossDomain) {
            xhr.setRequestHeader("X-CSRFToken", $.cookie('csrftoken'));
        }
    }
});


/*获取有审核权限的用户*/
$(function () {
    $.ajax({
        url: '/opsorder/ops_audit',
        type: 'GET',
        dataType: 'json',
        timeout: 5000,
        cache: false,
        success: function (result) {
            $("#verifiers").empty();
            $.each(result, function (index, row) {
                $("#verifiers").append(
                    "<option data-icon='glyphicon-user' data-subtext=" + row.displayname + " value=" + row.username + ">" + row.username + "</option>");
            });
            $('.selectpicker').selectpicker('refresh');
        }
    })
});


/*添加工单 */
$("#order_add").on("click", function () {
    $('.msg-error').remove();

    var data = {};

    $("form input:text").each(function () {
        data[$(this).attr('name')] = $(this).val().trim();
    });

    $("form select").each(function () {
        data[$(this).attr('name')] = $(this).val().trim();
    });

    data["order_id"] = $(".form-horizontal").attr("data");
    data["content"] = myCodeMirror.getValue();


    $.ajax({
        url: '/opsorder/add',
        type: 'POST',
        data: data,
        traditional: true,
        dataType: 'json',

        success: function (data) {
            if (data.status == 0) {
                displayPNotify(0, "提交成功");
                window.location.href = data.data;
            } else {
                $.each(data.msg, function (k, v) {
                    var tag = '<span class="msg-error">' + v[0].message + '</span>';
                    $("#" + k).parent().append(tag);
                });
                setTimeout(function () {
                    $('.msg-error').remove();
                }, 5000);
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status)
            }
        }
    })
});


/* 上传文件 */
function UploadFile(){
    if ($("#upload_file").hasClass("btn-success")) {
        $('.msg-error').remove();

        var form = new FormData();
        var file_obj = $('#filename')[0].files[0];
        form.append("filename", file_obj);

        $.ajax({
            type: 'POST',
            url: '/opsorder/upload',
            data: form,
            processData: false,
            contentType: false,
            success: function(data){
                if (data.status == 0) {
                    $(".form-horizontal").attr("data", data.data);
                    displayPNotify(0, "文件上传成功");
                    $("#upload_file").removeClass("btn-success");
                    $("#upload_file").addClass("btn-default");
                } else {
                    displayPNotify(2, data.msg);
                }
            },
            error: function (jqXHR) {
                if (jqXHR.status === 403) {
                    displayPNotify(jqXHR.status)
                }
            }
        })
    } else {
        $('.msg-error').remove();
        var tag = '<span class="msg-error">文件已上传</span>';
        $('#upload_file').after(tag)
    }
}


/* 同意工单执行 */
$("button[name='agree_order']").on("click", function () {
    data = {};
    data["id"] = $(this).attr("data");

    $.ajax({
        url: '/opsorder/agree',
        type: 'POST',
        data: data,
        traditional: true,
        dataType: 'json',

        success: function (data) {
            if (data.status == 0) {
                displayPNotify(0, data.msg);
                window.location.reload();
            } else {
                displayPNotify(2, data.msg);
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status)
            }
        }
    })
});

/* 执行工单*/
$("button[name='do_ops']").on("click", function () {
    data = {};
    data["id"] = $(this).attr("data");

    $.ajax({
        url: '/opsorder/doing',
        type: 'POST',
        data: data,
        traditional: true,
        dataType: 'json',

        success: function (data) {
            if (data.status == 0) {
                displayPNotify(0, data.msg);
                window.location.reload();
            } else {
                displayPNotify(2, data.msg);
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status)
            }
        }
    })
});


/*关闭工单*/
$("button[name='close_order_submit']").on("click", function () {
    data = {};
    data["id"] = $(this).attr("data");
    data["close_text"] = $("textarea[name='close_text']").val();

    $.ajax({
        url: '/opsorder/close',
        type: 'POST',
        data: data,
        traditional: true,
        dataType: 'json',

        success: function (data) {
            if (data.status == 0) {
                window.location.reload();
            } else {
                displayPNotify(2, data.msg);
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status)
            }
        }
    })
});

$("button[name='done_order_submit']").on("click", function () {
    data = {};
    data["id"] = $(this).attr("data");
    data["done_text"] = $("textarea[name='done_text']").val();

    $.ajax({
        url: '/opsorder/done',
        type: 'POST',
        data: data,
        traditional: true,
        dataType: 'json',

        success: function (data) {
            if (data.status == 0) {
                window.location.reload();
            } else {
                displayPNotify(2, data.msg);
            }
        },
        error: function (jqXHR) {
            if (jqXHR.status === 403) {
                displayPNotify(jqXHR.status)
            }
        }
    })
});

$("button[name='close_order']").on("click", function () {
    $("#myModal").modal('toggle');
    $("button[name='close_order_submit']").attr("data", $(this).attr("data"))
});


$("button[name='done_order']").on("click", function () {
    $("#myModal1").modal('toggle');
    $("button[name='done_order_submit']").attr("data", $(this).attr("data"))
});


function clear_modal() {
    $("button[name='dis_modal']").on("click", function () {
        $('#myModal').modal('hide');
        $('#myModal1').modal('hide');
    })
}
