$(document).ready(function() {
    var changePasswordForm = $("#changePasswordForm")

    changePasswordForm.validate({
        rules: {
            password: {
                minlength: 6,
                normalizer: function(value) {
                    return $.trim(value)
                }
            },
            new_password: {
                minlength: 6,
                normalizer: function(value) {
                    return $.trim(value)
                }
            },
            confirm_password: {
                normalizer: function(value) {
                    return $.trim(value)
                },
                equalTo: "#newPassword"
            }
        },
        errorElement: 'div',
        errorPlacement: function(error, element) {
            error.addClass('form-text text-danger')
            error.appendTo(element.parents('div.form-group'))
        }
    });
    
    changePasswordForm.submit(function(event) {
        event.preventDefault()

        if ($(this).valid()) {
            $.post(changePasswordForm.attr("action"), changePasswordForm.serialize(), function(resp) {
                if (resp.status == 'success') {
                    window.location.reload()  
                    return
                }
                alert(resp.message)
                $('#captchaContainer>img').trigger('click')
                changePasswordForm.find('input[name="captcha"]').val("")
            }, 'json')
        }
    })
})