package user

import (
	"errors"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/casbin/casbin/v2"
	"github.com/clevergo/demo/pkg/rest/pagination"

	"github.com/clevergo/captchas"
	"github.com/clevergo/clevergo"
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/forms"
	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/internal/rbac"
	"github.com/clevergo/jsend"
)

type Resource struct {
	*api.Application
	captchaManager *captchas.Manager
	jwtManager     *core.JWTManager
	enforcer       *casbin.Enforcer
	service        Service
}

func New(
	app *api.Application,
	captchaManager *captchas.Manager,
	jwtManager *core.JWTManager,
	enforcer *casbin.Enforcer,
	service Service,
) *Resource {
	return &Resource{
		Application:    app,
		captchaManager: captchaManager,
		jwtManager:     jwtManager,
		enforcer:       enforcer,
		service:        service,
	}
}

func (r *Resource) RegisterRoutes(router clevergo.IRouter) {
	router.Post("/user/login", r.login)
	router.Get("/user/info", r.info)
	router.Post("/user/logout", r.logout)
	router.Get("/users", r.query)
	router.Get("/users/:id", r.get)
	router.Delete("/users/:id", r.delete)
}

func (r *Resource) query(ctx *clevergo.Context) (err error) {
	p := pagination.NewFromContext(ctx)
	qps := new(QueryParams)
	if err = api.DecodeQueryParams(qps, ctx); err != nil {
		return
	}
	p.Items, err = r.service.Query(p.Limit, p.Offset(), qps)
	if err != nil {
		return err
	}
	p.Total, err = r.service.Count()
	if err != nil {
		return err
	}
	return ctx.JSON(200, jsend.New(p))
}

func (r *Resource) get(ctx *clevergo.Context) error {
	return nil
}

func (r *Resource) delete(ctx *clevergo.Context) error {
	id := ctx.Params.String("id")
	user, _ := r.User(ctx)
	if user.GetIdentity().GetID() == id {
		return errors.New("you cannot delete your account")
	}

	sql, _, err := squirrel.Delete("users").Where("id=?", id).ToSql()
	if err != nil {
		return err
	}
	if _, err := r.DB().Exec(sql); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jsend.New(nil))
}

func (r *Resource) login(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	if !user.IsGuest() {
		ctx.Redirect("/backend/", http.StatusFound)
		return nil
	}
	form := forms.NewLogin(r.DB(), user, r.captchaManager)
	v, err := form.Handle(ctx)
	if err != nil {
		return err
	}
	token, err := r.jwtManager.New(v.GetID())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"access_token": token,
	}))
}

func (r *Resource) info(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	roles, err := r.enforcer.GetRolesForUser("user_" + user.GetIdentity().GetID())
	if err != nil {
		return err
	}
	query := squirrel.Select("id").From("auth_items").Where(squirrel.Eq{"item_type": rbac.TypePermission})
	policies, err := r.enforcer.GetImplicitPermissionsForUser("user_" + user.GetIdentity().GetID())
	if err != nil {
		return err
	}
	or := squirrel.Or{}
	for _, v := range policies {
		or = append(or, squirrel.Eq{
			"obj": v[1],
			"act": v[2],
		})
	}
	query = query.Where(or)
	sql, args, err := query.ToSql()
	permissions := []string{}
	rows, err := r.DB().Queryx(sql, args...)
	id := ""
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return err
		}
		permissions = append(permissions, id)
	}

	identity, _ := user.GetIdentity().(*models.User)
	return ctx.JSON(http.StatusOK, jsend.New(core.Map{
		"id":          identity.ID,
		"username":    identity.Username,
		"email":       identity.Email,
		"roles":       roles,
		"permissions": permissions,
	}))
}

func (r *Resource) logout(ctx *clevergo.Context) error {
	user, _ := r.User(ctx)
	if err := user.Logout(ctx.Request, ctx.Response); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jsend.New(nil))
}
