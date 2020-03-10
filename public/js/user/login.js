var loginfFormConstraints = {
    email: {
        presence: true,
        email: true,
    },
    password: {
        presence: true
    }
};

$(function() {
    var loginForm = $('#loginForm')
    var loginUrl = loginForm.attr('action')
    var loginButton = $('#loginButton')

    loginForm.on('submit', function(event) {
        event.preventDefault();

        errors = validate($(this).serializeObject(), loginfFormConstraints)
        if (errors) {
            var field;
            var input;
            var help;
            for (var key in errors) {
                input = loginForm.find('input[name="' + key + '"]');
                field = input.parents('.field');
                help = '';
                for (var i in errors[key]) {
                    help += '<p class="help is-danger">' + errors[key][i] + '</p>'
                }
                field.find('p.help').remove();
                input.addClass('is-danger');
                field.append(help)
            }
            return
        }

        lockButton(loginButton)
        axios.post(loginUrl, $(this).serialize())
            .then(function(response) {
                console.log(response);
                unlockButton(loginButton);

                data = response.data
                if (data.status == 'success') {
                    message('success', 'Login Suceess', 'success');
                } else {
                    message(data.status, data.message ? data.message : 'Internal Error', 'danger');
                }
            })
            .catch(function(error) {
                console.log(error);
                unlockButton(loginButton);
            });

    })
})