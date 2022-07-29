package docs

import "net/http"

func SwaggerServefile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger.yaml")
}
