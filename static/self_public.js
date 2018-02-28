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
function displayPNotify(status, msg, init_fun) {
    var type = '';
    var title = '';

    // 接收后端返回的状态值
    if (status == '200') {
        type = 'success';
        title = 'SUCCESS'
    } else if (status == '400') {
        type = 'error';
        title = 'ERROR'
    } else if (status == '403') {
        type = 'info';
        title = 'WARNING'
    }

    var set_null = function () {
        return null
    };


    // 使用bootstrap样式
    PNotify.prototype.options.styling = "bootstrap3";
    new PNotify({
        title: title,
        text: msg,
        type: type,
        delay: 2000,
        nonblock: {
            nonblock: true
        },
        after_close: function () {
            // 如果传入的变量不存在，则set_null
            if (!init_fun) {
                set_null()
            } else {
                init_fun()
            }
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
    for (i = 0; i < len; i++) {
        pwd += $chars.charAt(Math.floor(Math.random() * maxPos));
    }
    return pwd
}