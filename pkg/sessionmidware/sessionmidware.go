package sessionmidware

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/clevergo/clevergo"
)

func New(manager *scs.SessionManager) clevergo.MiddlewareFunc {
	return func(next clevergo.Handle) clevergo.Handle {
		return func(ctx *clevergo.Context) error {
			var token string
			cookie, err := ctx.Request.Cookie(manager.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			context, err := manager.Load(ctx.Request.Context(), token)
			if err != nil {
				manager.ErrorFunc(ctx.Response, ctx.Request, err)
				return err
			}

			ctx.Request = ctx.Request.WithContext(context)
			w := ctx.Response
			resp := &response{ResponseWriter: ctx.Response}
			ctx.Response = resp
			err = next(ctx)
			ctx.Response = w

			switch manager.Status(context) {
			case scs.Modified:
				token, expiry, err := manager.Commit(context)
				if err != nil {
					manager.ErrorFunc(resp.ResponseWriter, ctx.Request, err)
					return err
				}
				manager.WriteSessionCookie(resp.ResponseWriter, token, expiry)
			case scs.Destroyed:
				manager.WriteSessionCookie(resp.ResponseWriter, "", time.Time{})
			}

			if resp.code != 0 {
				resp.ResponseWriter.WriteHeader(resp.code)
			}
			resp.buf.WriteTo(resp.ResponseWriter)
			return err
		}
	}
}

type response struct {
	http.ResponseWriter
	buf         bytes.Buffer
	code        int
	wroteHeader bool
}

func (r *response) Write(b []byte) (int, error) {
	return r.buf.Write(b)
}

func (r *response) WriteHeader(code int) {
	if !r.wroteHeader {
		r.code = code
		r.wroteHeader = true
	}
}

func (r *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := r.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (r *response) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
