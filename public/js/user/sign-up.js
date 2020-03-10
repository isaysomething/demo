$(function() {
    var signUpForm = $('#signUpForm')
    var signUpUrl = signUpForm.attr('action')
    var signUpButton = $('#signUpButton')

    signUpForm.on('submit', function(event) {
        event.preventDefault();

        lockButton(signUpButton)
        axios.post(signUpUrl, $(this).serialize())
            .then(function(response) {
                console.log(response);
                unlockButton(signUpButton);

                data = response.data
                if (data.status == 'success') {
                    message('success', 'Login Suceess', 'success');
                } else {
                    message(data.status, data.message ? data.message : 'Internal Error', 'danger');
                }
            })
            .catch(function(error) {
                console.log(error);
                unlockButton(signUpButton);
            });
    })
})