function Captcha(container) {
    var input = container.find('input')
    var media = container.find('img')
    this.reload = function() {
        $.post('/captcha', {}, function(resp) {
            input.val(resp.data.id)
            media.attr("src", resp.data.data)
        }, 'json')
    }
    media.click(this.reload)
}

$(document).ready(function(){
    var captchaContainer = $('#captchaContainer')
    if (captchaContainer) {
        let captcha = new Captcha(captchaContainer)
        captcha.reload()
    }
})
