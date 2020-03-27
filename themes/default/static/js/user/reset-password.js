$(document).ready(function() {
    var form = $("#resetPasswordForm")

    form.validate({
        rules: {
            password: {
                required: true,
                minlength: 6,
                normalizer: function(value) {
                    return $.trim(value)
                }
            },
            confirm_password: {
                required: true,
                normalizer: function(value) {
                    return $.trim(value)
                },
                equalTo: "#password"
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
                    window.location.href = '/login'    
                    return
                }
                alert(resp.message)
            }, 'json')
        }
    })
})