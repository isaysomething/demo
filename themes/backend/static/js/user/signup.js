$(document).ready(function() {
    $.validator.addMethod("isEmailTaken", function(value, element) {
        $.post(apiHost + "/v1/check-user-email", {email: $(this).val()}, function(resp) {
            console.log(resp)
            return resp.status == "status"
        }, 'json')
    }, "* email was taken");

    var signupForm = $("#signupForm")
    var username = signupForm.find('input[name="username"]')
    var email = signupForm.find('input[name="email"]')

    /*username.bind('input propertychange', function() {
        console.log($(this).val())
        $.post(apiHost + "/v1/check-username", {username: $(this).val()}, function(resp) {
            console.log(resp)
            if (resp.status != "status") {

            }
        }, 'json')
    })*/

    /*email.bind('input propertychange', function() {
        console.log($(this).val())
        $.post(apiHost + "/v1/check-user-email", {email: $(this).val()}, function(resp) {
            console.log(resp)
            if (resp.status != "status") {

            }
        }, 'json')
    })*/

    signupForm.validate({
        rules: {
            username: {
                required: true,
                normalizer: function(value) {
                    return $.trim(value);
                }
            },
        },
       // errorElement: 'div',
        // errorClass: 'form-text text-danger',
        /*errorPlacement: function(error, element) {
            console.log(error.status)
            var parent = element.parents('.form-group')
            parent.find('.form-text').remove()
            error.attr('class', "form-text text-danger")
            parent.append(error)
        }*/
    });
    
    var captchaTicket = $('#captchaTicket')
    var captchaRandstr = $('#captchaRandstr')
    signupForm.submit()
    var captcha = new TencentCaptcha(document.getElementById("btnSignup"), "{{ .captchaAppID }}", function(res) {
        captchaTicket.val(res.ticket)
        captchaRandstr.val(res.randstr)
        signupForm.submit()
    }, {})
})