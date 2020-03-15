$(document).ready(function() {
    var signupForm = $("#signupForm")
    var username = signupForm.find('input[name="username"]')
    var email = signupForm.find('input[name="email"]')
    var password = signupForm.find('input[name="password"]')

    signupForm.validate({
        rules: {
            username: {
                required: true,
                minlength: 5,
                normalizer: function(value) {
                    return $.trim(value)
                },
                remote: {
                    url: '/backend/user/check-username',
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
                required: true,
                normalizer: function(value) {
                    return $.trim(value)
                },
                remote: {
                    url: '/backend/user/check-email',
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
                required: true,
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
                    url: '/backend/check-captcha',
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
    
    signupForm.submit(function(event) {
        event.preventDefault()

        if ($(this).valid()) {
            $.post('/backend/signup', signupForm.serialize(), function(resp) {
                if (resp.status == 'success') {
                    window.location.href = "/backend/login"    
                    return
                }
                alert(resp.message)
            }, 'json')
        }
    })
})