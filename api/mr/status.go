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
package mr

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Context().Value("request.api.version").(string) {
	case "1":
		w.Write([]byte(fmt.Sprintf("api : %d", r.Context().Value("api.version").(int))))
	case "2":
		w.Write([]byte(fmt.Sprintf("secret : %s", r.Context().Value("api.secret").(string))))
	case "3":
		w.Write([]byte(fmt.Sprintf("token : %s" /*r.Context().Value("api.token").(string)*/, "shhhh!")))

		url := fmt.Sprintf("%s/projects/%s", r.Context().Value("api.url").(string), r.Context().Value("api.project").(string))
		var bearer = fmt.Sprintf("Bearer %s", r.Context().Value("api.token").(string))
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", bearer)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		w.Write([]byte(body))

	case "4":
		w.Write([]byte(fmt.Sprintf("status : %s", "ok")))
	case "5":
		w.Write([]byte(fmt.Sprintf("request : %s", r.Context().Value("request.api.version").(string))))
	default:
		http.Error(w, http.StatusText(403), 403)
		return
	}
}

func StatusRouter(r chi.Router) {
	r.Get("/", GetStatus)
}
