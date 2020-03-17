function Captcha(container) {
    var input = container.find('input')
    var media = container.find('img')
    media.click(function() {
        $.post(container.attr('data-url'), {}, function(resp) {
            input.val(resp.data.id)
            media.attr("src", resp.data.data)
        }, 'json')
    })
}

$(document).ready(function(){
    var captchaContainer = $('#captchaContainer')
    if (captchaContainer) {
        new Captcha(captchaContainer)
    }
})