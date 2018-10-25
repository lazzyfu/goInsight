/**
 * Created by fuzongfei on 2017/9/15.
 */

/**
 * 创建websocket
 * status
 * 1: 输出执行当前SQL语句的processlist信息
 * 2：输出inception osc的执行进度
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

        if (status === 1) {
            $this.empty();
            output_html = renderSqlProcesslistResult(data);
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
 * 渲染执行当前SQL语句的processlist信息
 * {'ID': 5703, 'USER': 'yops', 'HOST': '10.10.1.25:63032', 'DB': 'aa',
   'COMMAND': 'Sleep', 'TIME': 0, 'STATE': '', 'INFO': None, 'TIME_MS': 44,
   'ROWS_SENT': 0, 'ROWS_EXAMINED': 0}
 */

function renderSqlProcesslistResult(data) {
    let html = "<p class=\"text-danger\">该SQL的SHOW PROCESSLIST实时输出：</p>";
    for (let key in data) {
        html += "<p><b>" + key + "</b>: " + data[key] + "</p>"
    }
    return html
}


function renderGhostResult(data) {
    return "<dt></dt><dd>" + data.replace(/\n/g, '\.' + '<br>') + "</dd>";
}