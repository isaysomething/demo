$(document).ready(function() {
    var form = $("#requestResetPasswordForm")

    form.validate({
        rules: {
            email: {
                required: true,
                normalizer: function(value) {
                    return $.trim(value)
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
    
    form.submit(function(event) {
        event.preventDefault()

        if ($(this).valid()) {
            $.post(form.attr("action"), form.serialize(), function(resp) {
                if (resp.status == 'success') {
                    window.location.href = '/'    
                    return
                }
                alert(resp.message)
                $('#captchaContainer>img').trigger('click')
                form.find('input[name="captcha"]').val("")
            }, 'json')
        }
    })
})