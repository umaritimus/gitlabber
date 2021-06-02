/*
Copyright © 2021 Andy Dorfman <github.com/umaritimus>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controller

import (
	"context"
	"fmt"
	"gitlabber/api"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/render"
	"github.com/spf13/cobra"
)

var (
	secret string
	port   int
	token  string
)

func Router() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/api/{version}", func(r chi.Router) {
		r.Use(ApiVersionCtx())
		r.Mount("/", api.ApiRouter())
	})

	return router
}

func ApiVersionCtx() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			version := chi.URLParam(r, "version")
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			r = r.WithContext(context.WithValue(r.Context(), "api.secret", secret))
			r = r.WithContext(context.WithValue(r.Context(), "api.token", token))
			next.ServeHTTP(w, r)
		})
	}
}

func InitConfig(cmd *cobra.Command) string {
	port, _ = cmd.Flags().GetInt("port")
	secret, _ = cmd.Flags().GetString("secret")
	token, _ = cmd.Flags().GetString("token")

	log.Info(fmt.Sprintf("Listening on port '%d'\n", port))
	log.Info(fmt.Sprintf("Using secret '%s'\n", secret))

	return fmt.Sprintf(":%d", port)
}
