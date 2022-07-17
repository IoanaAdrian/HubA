package api

import (
	"context"
	"fmt"
	"github.com/IoanaAdrian/HubA/src/api/handlers"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var routes = []route{
	newRoute("POST", "/", handlers.HomePost),
	newRoute("GET", "/", handlers.HomeGet),
	//newRoute("GET", "/contact", contact),
	//newRoute("GET", "/api/widgets", apiGetWidgets),
	//newRoute("POST", "/api/widgets", apiCreateWidget),
	//newRoute("POST", "/api/widgets/([^/]+)", apiUpdateWidget),
	//newRoute("POST", "/api/widgets/([^/]+)/parts", apiCreateWidgetPart),
	//newRoute("POST", "/api/widgets/([^/]+)/parts/([0-9]+)/update", apiUpdateWidgetPart),
	//newRoute("POST", "/api/widgets/([^/]+)/parts/([0-9]+)/delete", apiDeleteWidgetPart),
	//newRoute("GET", "/([^/]+)", widget),
	//newRoute("GET", "/([^/]+)/admin", widgetAdmin),
	//newRoute("POST", "/([^/]+)/image", widgetImage),
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allowedVerbs []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allowedVerbs = append(allowedVerbs, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			log.Print(fmt.Sprintf("[%s] ", route.method), "Route hit: ", matches)
			return
		}
	}
	if len(allowedVerbs) > 0 {
		w.Header().Set("Allow", strings.Join(allowedVerbs, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}
