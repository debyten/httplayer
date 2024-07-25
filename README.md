# httplayer

> Disclaimer: this package is NOT another go HTTP router.

## What is httplayer
Httplayer is a package designed to enhance the route building phase atop the standard Go HTTP router.
## How is structured
You can access to the following functionalities:
- Builder: A struct facilitating the construction of simple routes with middleware.
- RoutingDefinition: A struct for building a collection of routes with shared middleware.
- ServiceBuilder: A struct for assembling multiple services with distinct middlewares.
- Mounting: Utilities functions to mount the built routes/services

### Builder
The Builder struct offers a fluent api for constructing individual routes:

```go
package main

import (
	"github.com/debyten/httplayer"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	corsMw := cors.New(cors.Options{})
	myRoute := httplayer.NewBuilder(http.MethodGet, http.MethodPost).
		Path("/api/user").
		Middleware(corsMw).
		Handler(func(w http.ResponseWriter, r *http.Request) {
        // call stack is corsMw > handler
	}).
		Build()
	mux := http.NewServeMux()
	httplayer.MountRoute(mux, myRoute)
	http.ListenAndServe(":8080", mux)
}
```

### RoutingDefinition
This includes more advanced middleware capabilities. You can define routes like so:

```go
package main

import (
	"github.com/debyten/httplayer"
	"github.com/rs/cors"
	"net/http"
)

func userHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	corsMw := cors.New(cors.Options{})
	def := httplayer.NewDefinition(corsMw).
		Add(http.MethodGet, "/api/user", userHandler)
	    // Add other routes ...
	mux := httplayer.Mount(def)
	http.ListenAndServe(":8080", mux)
}
```

### ServiceBuilder
The ServiceBuilder is useful for combining multiple services with different middlewares:
```go
package main

import (
	"github.com/debyten/httplayer"
	"github.com/rs/cors"
	"net/http"
)

func NewUserApi(svc UserService) httplayer.Routing {
	return userApi{svc: svc}
}

type userApi struct {
	svc UserService
}

func (u userApi) Routes(l *httplayer.RoutingDefinition) []httplayer.Route {
	return l.
		Add(http.MethodGet, "/api/users/{id}", u.findByIDApi).
		Add(http.MethodPost, "/api/users", u.createApi).
		Done()
}

func NewLoginApi(svc LoginService) httplayer.Routing {
	return loginApi{svc: svc}
}

type loginApi struct {
	svc LoginService
}

func (u loginApi) Routes(l *httplayer.RoutingDefinition) []httplayer.Route {
	return l.
		Add(http.MethodPost, "/api/login", u.loginApi, u.loginRateLimit).
		Done()
}

// Middleware
func (u loginApi) loginRateLimit(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // rate limit middleware impl
		h(w, r)
	}
}


func main() {
	corsMw := cors.New(cors.Options{})
	authMiddleware := ...
	userApi := NewUserApi(...)
	loginApi := NewLoginApi(...)
	protectedApis := httplayer.NewServiceBuilder(authMiddleware).Add(userApi)
	publicApis := httplayer.NewServiceBuilder().Add(loginApi)
	mux := httplayer.MountServices(protectedApis, publicApis)
	// we can add cors on top
	h := corsMw.Handler(mux)
	http.ListenAndServe(":8080", h)
}

```

#### Detach and Merge
These two functions allows to group features, for example, if we define an rbac middleware like so:
```go

package main

func RBAC(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// implementation
		}
	}
}
```

Then we can use the middleware while implementing the Service:

```go
package main

type statsApi struct {
	svc StatsService
}

func (u statsApi) Routes(l *httplayer.RoutingDefinition) []httplayer.Route {
	adminRoutes := l.Detach(RBAC("admin")).
		Add(http.MethodPost, "/api/stats/purge", u.purgeApi)
	    Add(http.MethodPost, "/api/stats", u.pushStatApi)
	
	userRoutes := l.Detach(RBAC("user", "admin")).
		Add(http.MethodGet, "/api/stats", u.viewStatsApi).
	    Add(http.MethodPatch, "/api/stats/{id}", u.pinStatApi)
	return httplayer.Merge(adminRoutes, userRoutes)
}
```

In this way we share the middlewares from `l` (which is the main routing definition)

### Mounting
To finalize the routing setup, you can employ the following functions to mount the routes onto a `http.ServeMux`:
- `MountRoute`: This function installs the specified route onto the provided http.ServeMux.
- `Mount`: Use this function to install a collection of RoutingDefinitions onto a new `http.ServeMux`.
- `DynamicMount`: This function installs the given `RoutingDefinition` onto a new `http.ServeMux` and utilizes the input channel to dynamically register new Route instances in a separate goroutine.

> For the last use case for example, you can think to an api gateway which automatically registers routes as new services are discovered.
