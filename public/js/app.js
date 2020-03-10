$.fn.serializeObject = function() {
    var parts;
    var obj = [];
    var fields = this.serialize().split('&');
    for (var i in fields) {
        parts = fields[i].split('=');
        obj[parts[0]] = decodeURIComponent(parts[1]);
    }
    return obj;
};

$.fn.serializeArray = function() {
    var parts;
    var res = [];
    var args = this.serialize().split('&');
    for (var i in args) {
        parts = args[i].split('=');
        res.push({ 'name': parts[0], 'value': parts[1] });
    }
    return res;
};

function lockButton(btn) {
    btn.addClass('is-loading');
}

function unlockButton(btn) {
    btn.removeClass('is-loading');
}

document.addEventListener('DOMContentLoaded', () => {
    (document.querySelectorAll('.notification .delete') || []).forEach(($delete) => {
        $notification = $delete.parentNode;
        $delete.addEventListener('click', () => {
            $notification.parentNode.removeChild($notification);
        });
    });

    (document.querySelectorAll('#signOut') || []).forEach(($button) => {
        $form = $button.parentNode;
        $button.addEventListener('click', () => {
            $form.submit();
        });
    });
});

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