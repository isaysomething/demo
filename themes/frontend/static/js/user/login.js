$(document).ready(function() {
    var loginForm = $("#loginForm")

    loginForm.validate({
        rules: {
            email: {
                normalizer: function(value) {
                    return $.trim(value)
                }
            },
            password: {
                minlength: 6,
                normalizer: function(value) {
                    return $.trim(value);
                }
            },
            captcha: {
                required: true,
                normalizer: function(value) {
                    return $.trim(value);
                },
                remote: {
                    url: '/check-captcha',
                    type: 'post',
                    data: {
                        id: function() {
                            return $('input[name="captcha_id"]').val()
                        },
                        captcha: function() {
                            return $('input[name="captcha"]').val()
                        },
                    },
                    dataFilter: function (resp) {
                        var data = JSON.parse(resp)
                        if (data.status == 'error') {
                            return '"' + data.message + '"'
                        }

                        return true
                    }
                }
            }
        },
        errorElement: 'div',
        errorPlacement: function(error, element) {
            error.addClass('form-text text-danger')
            error.appendTo(element.parents('div.form-group'))
        }
    });
    
    loginForm.submit(function(event) {
        event.preventDefault()

        if ($(this).valid()) {
            $.post(loginForm.attr("action"), loginForm.serialize(), function(resp) {
                if (resp.status == 'success') {
                    window.location.href = '/'    
                    return
                }
                alert(resp.message)
            }, 'json')
        }
    })
})