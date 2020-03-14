$(document).ready(function() {
    var loginForm = $("#loginForm")
    var loginForm = $("#username")
    var loginForm = $("#loginForm")
    loginForm.validate()
    var captchaTicket = $('#captchaTicket')
    var captchaRandstr = $('#captchaRandstr')
    var captcha = new TencentCaptcha(document.getElementById("btnLogin"), "{{ .captchaAppID }}", function(res) {
        captchaTicket.val(res.ticket)
        captchaRandstr.val(res.randstr)
        loginForm.submit()
    }, {})
})