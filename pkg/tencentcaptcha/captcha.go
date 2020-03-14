package tencentcaptcha

import (
	"errors"

	v20190722 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
)

type Captcha struct {
	client      *v20190722.Client
	appID       uint64
	secretKey   string
	captchaType uint64
}

func New(client *v20190722.Client, appID uint64, secretKey string) *Captcha {
	return &Captcha{
		client:      client,
		appID:       appID,
		secretKey:   secretKey,
		captchaType: 9,
	}
}

func (c *Captcha) AppID() uint64 {
	return c.appID
}

func (c *Captcha) Validate(ticket, randstr, ipAddr string) (bool, error) {
	req := v20190722.NewDescribeCaptchaResultRequest()
	req.Ticket = &ticket
	req.Randstr = &randstr
	req.AppSecretKey = &c.secretKey
	req.CaptchaAppId = &c.appID
	req.CaptchaType = &c.captchaType
	req.UserIp = &ipAddr
	resp, err := c.client.DescribeCaptchaResult(req)
	if err != nil {
		return false, err
	}
	if *resp.Response.CaptchaCode != 1 {
		return false, errors.New(*resp.Response.CaptchaMsg)
	}

	return true, nil
}
