// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func configConsulHttpRoutes(r *mux.Router) {
	r.HandleFunc("/v1/kv/.*", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		var Data string
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&Data)
		if err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		//r := httplib.Put(url)
		//r.Body(data)
		//ret, err := r.String()
		//fmt.Println(ret, err)

		RenderDataJson(w, Data)
	})
}
