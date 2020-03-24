package controllers

import (
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/internal/listeners"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/bootstrap"
	"github.com/clevergo/form"
	"github.com/clevergo/jsend"
)

type User struct {
	*frontend.Application
}

func NewUser(app *frontend.Application) *User {
	return &User{app}
}

func (u *User) Index(ctx *clevergo.Context) error {
	return u.Render(ctx, "user/index", nil)
}

func (u *User) Login(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/", http.StatusFound)
		return nil
	}

	if ctx.IsPost() {
		form := forms.NewLogin(u.DB(), user, u.CaptcpaManager())
		if _, err := form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}

		return jsend.Success(ctx.Response, nil)
	}

	return u.Render(ctx, "user/login", nil)
}

func (u *User) CheckUsername(ctx *clevergo.Context) error {
	f := forms.NewCheckUsername(u.DB())
	err := form.Decode(ctx.Request, f)
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}
	if err = f.Validate(); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	return jsend.Success(ctx.Response, nil)
}

func (u *User) CheckEmail(ctx *clevergo.Context) error {
	f := forms.NewCheckUserEmail(u.DB())
	err := form.Decode(ctx.Request, f)
	if err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}
	if err = f.Validate(); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	return jsend.Success(ctx.Response, nil)
}

func (u *User) ResendVerificationEmail(ctx *clevergo.Context) error {
	captcha, err := u.CaptcpaManager().Generate()
	if err != nil {
		return err
	}

	form := forms.NewResendVerificationEmail(u.DB(), u.Mailer(), u.CaptcpaManager())
	if ctx.IsPost() {
		if err = form.Handle(ctx); err == nil {
			u.AddFlash(ctx, bootstrap.NewSuccessAlert("Sent successfully, please check your mailbox."))
		} else {
			u.AddFlash(ctx, bootstrap.NewDangerAlert(err.Error()))
			u.Logger().Error(err)
		}
	}
	return u.Render(ctx, "user/resend-verification-email", core.ViewData{
		"form":    form,
		"error":   err,
		"captcha": captcha,
	})
}

func (u *User) Logout(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		if err := user.Logout(ctx.Request, ctx.Response); err != nil {
			return err
		}
	}

	ctx.Redirect("/login", http.StatusFound)
	return nil
}

func (u *User) Signup(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/", http.StatusFound)
		return nil
	}

	form := forms.NewSignup(u.DB(), user, u.CaptcpaManager())
	var err error
	form.RegisterOnAfterSignup(listeners.SendVerificationEmail(u.Mailer()))
	if ctx.IsPost() {
		if _, err = form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}

		return jsend.Success(ctx.Response, nil)
	}

	return u.Render(ctx, "user/signup", core.ViewData{
		"form":  form,
		"error": err,
	})
}

func (u *User) VerifyEmail(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/", http.StatusFound)
		return nil
	}

	token := ctx.Request.URL.Query().Get("verification_token")
	if token != "" {
		form := forms.NewVerifyEmail(u.DB(), user)
		form.Token = token
		err := form.Handle(ctx)
		if err == nil {
			ctx.Redirect("/", http.StatusFound)
			return nil
		}
		u.AddFlash(ctx, bootstrap.NewDangerAlert(err.Error()))
	}

	return u.Render(ctx, "user/verify-email", nil)
}

func (u *User) Setting(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if user.IsGuest() {
		ctx.Redirect("/login", http.StatusFound)
		return nil
	}

	token := ctx.Request.URL.Query().Get("verification_token")
	if token != "" {
		form := forms.NewVerifyEmail(u.DB(), user)
		form.Token = token
		err := form.Handle(ctx)
		if err == nil {
			ctx.Redirect("/", http.StatusFound)
			return nil
		}
		u.AddFlash(ctx, bootstrap.NewDangerAlert(err.Error()))
	}

	return u.Render(ctx, "user/setting", nil)
}

func (u *User) RequestResetPassword(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewRequestResetPassword(u.DB(), u.Mailer(), u.CaptcpaManager())
		if err := form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}

		return jsend.Success(ctx.Response, nil)
	}

	return u.Render(ctx, "user/request-reset-password", nil)
}

func (u *User) ResetPassword(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewResetPassword(u.DB())
		err := form.Handle(ctx)
		if err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}
		u.AddFlash(ctx, bootstrap.NewSuccessAlert("Password reset successfully."))
		return jsend.Success(ctx.Response, nil)
	}

	return u.Render(ctx, "user/reset-password", core.ViewData{
		"token": ctx.Request.URL.Query().Get("token"),
	})
}

func (u *User) ChangePassword(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if user.IsGuest() {
		ctx.Redirect("/login", http.StatusFound)
		return nil
	}

	identity, _ := user.GetIdentity().(*models.User)
	form := forms.NewChangePassword(u.DB(), identity)

	if ctx.IsPost() {
		err := form.Handle(ctx)
		if err == nil {
			u.AddFlash(ctx, bootstrap.NewSuccessAlert("Password has been updated."))
			ctx.Redirect("/", http.StatusFound)
			return nil
		}
		u.AddFlash(ctx, bootstrap.NewDangerAlert(err.Error()))
	}

	return u.Render(ctx, "user/change-password", core.ViewData{
		"form": form,
	})
}
