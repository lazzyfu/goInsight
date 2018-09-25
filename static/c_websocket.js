/**
 * Created by fuzongfei on 2017/9/15.
 */

/**
 * 创建websocket
 * status
 * 0：输出inception的执行结果
 * 1：输出inception osc的执行进度
 */
function CreateWebSocket() {
    let socket = new WebSocket('ws://' + window.location.host + '/ws/');
    socket.onmessage = function (message) {
        document.getElementById('typediv1').style.visibility = "visible";
        let result = JSON.parse(message.data);
        let status = result.status;
        let data = result.data;
        let output_html = '';
        let $this = $('#output_append');

        if (status === 0) {
            $this.empty();
            output_html = renderIncepExecResult(data);
            $this.append(output_html);
        }
        else if (status === 1) {
            $this.empty();
            output_html = renderIncepOSCResult(data);
            $this.append(output_html);
        }
        else if (status === 3) {
            output_html = renderGhostResult(data);
            $this.append(output_html);
            let element = document.getElementById('output_append');
            element.scrollTop = element.scrollHeight;
        }
    };
    socket.onclose = function () {
        setTimeout(function () {
            CreateWebSocket()
        }, 2000);
    };
}

/**
 * 渲染inception执行输出
 */
function renderIncepExecResult(data) {
    let html = '';
    $.each(data, function (index, row) {
        let SQL = row.SQL;
        let ID = row.ID;
        let stage = row.stage;
        let errlevel = '';
        if (row.errlevel === 0) {
            errlevel = '成功'
        }
        else if (row.errlevel === 1) {
            errlevel = '警告'
        }
        else if (row.errlevel === 2) {
            errlevel = '错误'
        }
        let stagestatus = row.stagestatus;
        let errormessage = row.errormessage;
        let Affected_rows = row.Affected_rows;
        let execute_time = row.execute_time;
        html += "<dt> ID：</dt>" + "<dd>" + ID + "</dd>" +
            "<dt> 阶段：</dt>" + "<dd>" + stage + "</dd>" +
            "<dt> 状态：</dt>" + "<dd>" + stagestatus + "</dd>" +
            "<dt> 错误级别：</dt>" + "<dd>" + errlevel + "</dd>" +
            "<dt> 扫描/影响行数：</dt>" + "<dd>" + Affected_rows + "</dd>" +
            "<dt> 耗时：</dt>" + "<dd>" + execute_time + "</dd>" +
            "<dt> 错误信息：</dt>" + "<dd>" + errormessage + "</dd>" +
            "<dt> SQL语句：</dt>" + "<dd>" + SQL + "</dd>" + "<br>"
    });
    return html
}

/**
 * 渲染inception osc输出
 */
function renderIncepOSCResult(data) {
    let html = '';
    $.each(data, function (index, row) {
        let TABLE = row.DBNAME + '.' + row.TABLENAME;
        let REMAINTIME = row.REMAINTIME;
        let SQLSHA1 = row.SQLSHA1;
        let PERCENT = row.PERCENT;
        let INFOMATION = row.INFOMATION.replace(/\n/g, '\.' + '<br>');

        html += "<dt> 正在操作表：</dt>" + "<dd>" + TABLE + "</dd>" +
            "<dt> SQLSHA1：</dt>" + "<dd>" + SQLSHA1 + "</dd>" +
            "<dt> 预估剩余时间：</dt>" + "<dd class='text-red'>" + REMAINTIME + "</dd>" +
            "<dt> OSC执行进度：</dt>" + "<dd class='text-red'>" + PERCENT + "%" + "</dd>" +
            "<dt> OSC输出：</dt>" + "<dd>" + INFOMATION + "</dd>"
    });
    return html
}


function renderGhostResult(data) {
    return "<dt></dt><dd>" + data.replace(/\n/g, '\.' + '<br>') + "</dd>";
}