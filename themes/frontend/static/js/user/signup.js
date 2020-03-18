$(document).ready(function() {
    var signupForm = $("#signupForm")
    var username = signupForm.find('input[name="username"]')
    var email = signupForm.find('input[name="email"]')

    signupForm.validate({
        rules: {
            username: {
                minlength: 5,
                normalizer: function(value) {
                    return $.trim(value)
                },
                remote: {
                    url: '/user/check-username',
                    type: 'post',
                    data: {
                        username: function() {
                            return username.val()
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
            },
            email: {
                normalizer: function(value) {
                    return $.trim(value)
                },
                remote: {
                    url: '/user/check-email',
                    type: 'post',
                    data: {
                        email: function() {
                            return email.val()
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
            },
            password: {
                minlength: 6,
                normalizer: function(value) {
                    return $.trim(value);
                }
            },
            captcha: {
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
            error.addClass('invalid-feedback')
            error.appendTo(element.parents('div.form-group'))
        }
    });
    
    signupForm.submit(function(event) {
        event.preventDefault()

        if ($(this).valid()) {
            $.post('/signup', signupForm.serialize(), function(resp) {
                if (resp.status == 'success') {
                    window.location.href = "/login"    
                    return
                }
                alert(resp.message)
                $('#captchaContainer>img').trigger('click')
                signupForm.find('input[name="captcha"]').val("")
            }, 'json')
        }
    })
})