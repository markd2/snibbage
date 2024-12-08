# snibbage

Learning Go

Working through [Let's Go](https://lets-go.alexedwards.net/) by Alex Edwards.

## Notes

* A module path is a like a canonical name or _identifier_ for the project.
  - almost any string - but focus on uniqueness
  - if intending for others to download and use, make your module path equal to the location the code can be downloaded from
    - e.g. `https://github.com/oop/ack" then module should be "github.com/oop/ack"
* `go mod init snibbage.borkware.com`
* when there's a valid go.mod at the root of the project irectory, the project IS a module
  - makes it easier to manage 3rd party dependencies
  - avoid supply-chain attacks
  - ensure reproducible builds

Herro wold:

```golang
package main

import "fmt"

func main() {
	fmt.Println("Smorgle")
}
```

### Web Stuff

* need a _handler_ - executing application logic and writing HTTP response headers and bodies
* need a _router_ - (a.k.a. _servemux_) Stores a mapping from URL routing patterns and the corresponding handlers.
  - usually one servemux for application containing all your routes (are belong to us)
* Web server - Can use the go application itself, so don't have to have nginx, Apache, or AOLserver
* can use a symbolic port, like `":http"` or `":http-alt``. Looks in e.g. /etc/services when starting the server.
* servemux has different matching rules if route ends with a trailing slash
  - when doesn't have a trailing slash, will only be matched when the request
    URL path exactly matches the pattern in full
  - when ends with a slash, like `"/"` or "/snorgleblorg/", it is a
    _subtree path pattern_.  They're matched (handler is called) when the start
    of the request URL path matches the subtree path.  e.g. like
    `"/**"` or `"/snorgleblorg/**"`
  - this is why "/" actls like a catch-all
  - to prevent the wildcarding, use `{$}` at the end of of hte pattern.
    e.g. `"/{$}"` and `"/snorgleblorg/{$}"`
* Some more servemux stuff
  - request url paths are automatically sanitized (any . / .. / repeated 
    slashes will redirect (301 - permnanet redirect) to a clean URL). 
  - if a subtree path has been registered and request comes in without
    the slash, will redirect (301) to the subtree with the slash.  So
    `/snorgleblorg` will get redirected to `/snorgleblorg/` automagically
  - can include hostnames in the route patterns.  Say for redirecting
    all HTTP requests to a canonical url, or it's the backend for multiple
    sites.
    - e.g. `mux.HanldeFun("fishbattery.com/", fishbatteryHandler)`
  - hostpatterns are matched first.
  - there exists a default servemux, accessed via http.Handle() and
    http.handleFunc(). (stored in http.DefaultServeMux global).  Pass nil
    to your ListenAndServe call to use it.
    - contraindicated b/c less explicit and more _magic_
    - defaultServeMux is a global var in the standard library, it means
      _any_ Go code can register a route. That can lead to mayhem, either
      via local teams adding to the global mux and now we can't figure out
      the canonical list of routes. Plus a compromised package could register
      and do something nefarious.
* Wildcard route patternage
  - can define routes that contain wildcard segments.  Allows more flexible
    routing riles, and pass variables via a request URL
  - `{}` is a _wildcard indentifier` (which explains the use of `{$}` - hat
    before cash
  - e.g. `mux.HandleFunc("/oop/ack/{snorgle}/item/{blorf}", handler)`
    - pattern contains two wildcard seggies.  First will have the identifier
      `snorgle`, and the second `blorf`
  - matching rules still in effect, plus the request path can contain
    _any_ non-empty value for the wildcard segments
    - so theis matches: `/oop/ack/splunge/item/greeble%20bork`
  - the thing between slashes that matches (e.g. snorgle and blorf) must be
    the entirety between slashes.  `/oop/b_{lorf}/` is right out.  As is
    `/{flongwaffle}.html`
  - inside the handler, can get the corresponding value using its
    identifier and the `r.PathValue()` method.
    - `blorf := r.PathValue("blorf")`
    - always returns a string, and can be any alue, so validate and sanity 
      check before doing anything useful with it.
  - precedence - e.g. `/post/flong` and `/post/{id}`.  /post/flong matches
    both of them. 
    - the most specific route pattern wins
    - where _specific_ is one matches only a subset of requests the other does.
    - so `/post/flong` only matches with exactly that, while `/post/{id}`
      casts a wider net
  - nice side effect of the precedence rule is that order of route declaration
    does not matter (yay!)
  - they can still conflict
    - `/post/new/{id}` and `/post/{author}/latest` overlap.  Who handles
      `/post/new/latest`?
    - will cause a runtime panic when initializing routes
  - in general, don't use overlapping routes.
* subtree path patterns with wildcard
  - prior rules still hold, so when pattern ends with `/{$}`, it is a 
    _subtree path pattern_ and only requires that the start of a
    request URL path to match
  - so `/bork/{ookie}/` will match `/bork/1`, `/bork/queens/to/queens/level/three`. To suppress that, use `/bork/{ookie}/{$}`
* Remainder wildcards
  - wildcards normally match a single, non-empty, segment of a request path, ave a special case
  - `...` - something like `/greeble/{greeble...}` will match like 
  `/greeble/1`, `/greeble/a/b/c`, etc  BUT can access the entire wildcard
  part via the `r.PathValue()` call.



### Syntax

 * `go run` : shortcut that compiles the code, creates a binary in /tmp, and runs it.
   - can give it a space-separated list of .go files, path to a package, `.` for current directory, or full module path
   - all are equivalent:
     - `go run .`
     - `go run main.go`
     - `go run snibbage.borkware.com`
 * `func home(w http.ResponseWriter, r *http.Request) {`
  - function that returns nothing.
  - takes a response writer
  - takes a pointer to a struct


### dig in to

- byte slice
  - `[]bytes("blah")` syntax
- are there anonymous functions?
- string interpolation?