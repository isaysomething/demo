package controllers

import (
	"net/http"

	"github.com/clevergo/captchas"
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

// NewUser returns a user controller.
func RegisterUser(router clevergo.IRouter, app *frontend.Application, captchaManager *captchas.Manager) {
	u := &user{app, captchaManager}
	router.Get("/user", u.index, clevergo.RouteName("user"))
	router.Get("/login", u.login, clevergo.RouteName("login"))
	router.Post("/login", u.login)
	router.Post("/logout", u.logout, clevergo.RouteName("logout"))
	router.Get("/signup", u.signup, clevergo.RouteName("signup"))
	router.Post("/signup", u.signup)
	router.Post("/user/check-email", u.checkEmail)
	router.Post("/user/check-username", u.checkUsername)
	router.Get("/user/request-reset-password", u.requestResetPassword, clevergo.RouteName("request-reset-password"))
	router.Post("/user/request-reset-password", u.requestResetPassword)
	router.Get("/user/reset-password", u.resetPassword, clevergo.RouteName("reset-password"))
	router.Post("/user/reset-password", u.resetPassword)
	router.Get("/user/verify-email", u.verifyEmail, clevergo.RouteName("verify-email"))
	router.Get("/user/resend-verification-email", u.resendVerificationEmail, clevergo.RouteName("resend-verification-email"))
	router.Post("/user/resend-verification-email", u.resendVerificationEmail)
	router.Post("/user/change-password", u.changePassword, clevergo.RouteName("change-password"))
}

type user struct {
	*frontend.Application
	captchaManager *captchas.Manager
}

func (u *user) index(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if user.IsGuest() {
		ctx.Redirect("/login", http.StatusFound)
		return nil
	}

	return ctx.Render(http.StatusOK, "user/index.tmpl", nil)
}

// Login displays login page and handle login request.
func (u *user) login(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/", http.StatusFound)
		return nil
	}

	if ctx.IsPost() {
		form := forms.NewLogin(u.DB(), user, u.captchaManager)
		if _, err := form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}

		return jsend.Success(ctx.Response, nil)
	}

	return ctx.Render(http.StatusOK, "user/login.tmpl", nil)
}

func (u *user) checkUsername(ctx *clevergo.Context) error {
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

func (u *user) checkEmail(ctx *clevergo.Context) error {
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

func (u *user) resendVerificationEmail(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewResendVerificationEmail(u.DB(), u.Mailer(), u.captchaManager)
		if err := form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}
		u.AddFlash(ctx, bootstrap.NewSuccessAlert("Sent successfully, please check your mailbox."))
		return jsend.Success(ctx.Response, nil)
	}

	return ctx.Render(http.StatusOK, "user/resend-verification-email.tmpl", nil)
}

func (u *user) logout(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		if err := user.Logout(ctx.Request, ctx.Response); err != nil {
			return err
		}
	}

	ctx.Redirect("/login", http.StatusFound)
	return nil
}

func (u *user) signup(ctx *clevergo.Context) error {
	user, _ := u.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/", http.StatusFound)
		return nil
	}

	form := forms.NewSignup(u.DB(), user, u.captchaManager)
	var err error
	form.RegisterOnAfterSignup(listeners.SendVerificationEmail(u.Mailer()))
	if ctx.IsPost() {
		if _, err = form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}

		return jsend.Success(ctx.Response, nil)
	}

	return ctx.Render(http.StatusOK, "user/signup.tmpl", core.ViewData{
		"form":  form,
		"error": err,
	})
}

func (u *user) verifyEmail(ctx *clevergo.Context) (err error) {
	token := ctx.Request.URL.Query().Get("token")
	if token == "" {
		ctx.Redirect("/user/resend-verification-email.tmpl", http.StatusFound)
		return
	}

	form := forms.NewVerifyEmail(u.DB())
	form.Token = token
	if err = form.Handle(ctx); err != nil {
		u.AddFlash(ctx, bootstrap.NewDangerAlert(err.Error()))
		ctx.Redirect("/user/resend-verification-email.tmpl", http.StatusFound)
		return nil
	}

	u.AddFlash(ctx, bootstrap.NewSuccessAlert("Email has been verified successfully."))
	ctx.Redirect("/login", http.StatusFound)
	return
}

func (u *user) requestResetPassword(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewRequestResetPassword(u.DB(), u.Mailer(), u.captchaManager)
		if err := form.Handle(ctx); err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}

		return jsend.Success(ctx.Response, nil)
	}

	return ctx.Render(http.StatusOK, "user/request-reset-password.tmpl", nil)
}

func (u *user) resetPassword(ctx *clevergo.Context) error {
	if ctx.IsPost() {
		form := forms.NewResetPassword(u.DB())
		err := form.Handle(ctx)
		if err != nil {
			return jsend.Error(ctx.Response, err.Error())
		}
		u.AddFlash(ctx, bootstrap.NewSuccessAlert("Password reset successfully."))
		return jsend.Success(ctx.Response, nil)
	}

	return ctx.Render(http.StatusOK, "user/reset-password.tmpl", core.ViewData{
		"token": ctx.Request.URL.Query().Get("token"),
	})
}

func (u *user) changePassword(ctx *clevergo.Context) (err error) {
	user, _ := u.User(ctx)
	if user.IsGuest() {
		ctx.Redirect("/login", http.StatusFound)
		return nil
	}

	identity, _ := user.GetIdentity().(*models.User)
	form := forms.NewChangePassword(u.DB(), identity)
	if err := form.Handle(ctx); err != nil {
		return jsend.Error(ctx.Response, err.Error())
	}

	u.AddFlash(ctx, bootstrap.NewSuccessAlert("Password has been updated."))
	return jsend.Success(ctx.Response, nil)
}
