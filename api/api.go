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
package api

import (
	"gitlabber/api/info"
	"gitlabber/api/mr"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ApiRouter() http.Handler {
	router := chi.NewRouter()
	router.Route("/version", info.VersionRouter)
	router.Route("/status", mr.StatusRouter)
	return router
}
