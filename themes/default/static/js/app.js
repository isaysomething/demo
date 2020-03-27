$.ajaxSetup({
    beforeSend: function (xhr) {
        xhr.setRequestHeader("X-CSRF-Token",  $('input[name="gorilla.csrf.Token"]').val())
    }
})

$(function() {
    var lang = Cookies.get("lang")
    var langSelector = $('#langSelector')
    langSelector.find('a[data-id="' + lang + '"]').addClass('active')
    langSelector.find('a').on('click', function() {
        langSelector.find('a').removeClass('active')
        var self = $(this)
        self.addClass('active')
        Cookies.set("lang", self.attr("data-id"))
        window.location.reload()
    })

    if (!navigator.cookieEnabled) {
        message("Cookie", "Cookie was disabled by your browser", "error")
    }
})
