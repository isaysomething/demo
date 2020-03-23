package core

import (
	"fmt"
	"time"

	"github.com/clevergo/captchas"
	"github.com/clevergo/captchas/drivers"
	"github.com/clevergo/captchas/stores/redisstore"
	"github.com/go-redis/redis/v7"
)

type CaptchaConfig struct {
	Driver string `koanf:"driver"`
	String struct {
		Height int `koanf:"height"`
		Width  int `koanf:"width"`
		Length int `koanf:"length"`
	} `koanf:"string"`
	Math struct {
		Height     int `koanf:"height"`
		Width      int `koanf:"width"`
		NoiseCount int `koanf:"noise_count"`
	} `koanf:"math"`
	Chinese struct {
		Height int `koanf:"height"`
		Width  int `koanf:"width"`
		Length int `koanf:"length"`
	} `koanf:"chinese"`
	Digit struct {
		Height int `koanf:"height"`
		Width  int `koanf:"width"`
		Length int `koanf:"length"`
	} `koanf:"digit"`
	Audio struct {
		Length int `koanf:"length"`
	} `koanf:"audio"`
}

func NewCaptchaManager(cfg CaptchaConfig, store captchas.Store) *captchas.Manager {
	switch cfg.Driver {
	case "string":
		return captchas.New(store, drivers.NewString(

			drivers.StringHeight(cfg.String.Height),
			drivers.StringWidth(cfg.String.Width),
			drivers.StringLength(cfg.String.Length),
		))
	case "chinese":
		return captchas.New(store, drivers.NewChinese(
			drivers.ChineseHeight(cfg.Chinese.Height),
			drivers.ChineseWidth(cfg.Chinese.Width),
			drivers.ChineseLength(cfg.Chinese.Length),
		))
	case "audio":
		return captchas.New(store, drivers.NewAudio(
			drivers.AudioLength(cfg.Audio.Length),
		))
	case "math":
		return captchas.New(store, drivers.NewMath(
			drivers.MathHeight(cfg.Math.Height),
			drivers.MathWidth(cfg.Math.Width),
			drivers.MathNoiseCount(cfg.Math.NoiseCount),
		))
	default:
		return captchas.New(store, drivers.NewDigit(
			drivers.DigitHeight(cfg.Digit.Height),
			drivers.DigitWidth(cfg.Digit.Width),
			drivers.DigitLength(cfg.Digit.Length),
		))
	}
}

func NewCaptchaStore(cfg RedisConfig) captchas.Store {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
	return redisstore.New(
		client,
		redisstore.Expiration(10*time.Minute),
		redisstore.Prefix("captchas"),
	)
}
