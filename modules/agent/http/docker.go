package http

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func configDockerRoutes() {
	http.HandleFunc("/docker", func(w http.ResponseWriter, r *http.Request) {

		uuid, _ := g.Hostname()
		RenderDataJson(w, map[string]interface{}{
			"uuid":          uuid,
			"consulAddr":    g.Config().Consul.Addr,
			"consulTimeout": g.Config().Consul.Timeout,
		})
	})

}