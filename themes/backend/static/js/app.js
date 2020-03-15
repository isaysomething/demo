$.ajaxSetup({
    beforeSend: function (xhr) {
        xhr.setRequestHeader("X-CSRF-Token",  $('input[name="gorilla.csrf.Token"]').val())
    }
});
