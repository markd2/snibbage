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
* VERBS
  - prefix the route pattern with the necessary HTTP method when declaring
  - e.g. `"GET /snorgle/{$}"`
    - stringy API :-(  though I guess it allows custom verbs like SPLUNGE
  - the http methods are UPPER CASE and should be shouted 
  - GET matches GET and HEAD
  - totes OK to delcare nultiple routes that have different verbs
  - the _most specific pattern wins_ rule also applies with route patterns
    that overlap because of an HTTP method
  - no method matches any method, while something like `"POST /toasties"` will
    only match the method POST, so the POST would take preceden e
  - there is no handler nomenclature guidance (at least so far)
    - the book uses a convention of postfixing the names of POST handlers
      - e.g. `func snibbageCreatePost(..)`
* curling iron
  - `curl -i localhost:4000/`
    - GET
  - `curl --head localhost:4000/`
    - HEAD
  - `curl -i -d "" localhost:4000/`
    - POST
    - `-d` flag declares any HTTP POST data to include
  - `curl -iL -d "" http://localhost:4000/snippet/create`
    - `-L` to follow redirects
* Third part rooters
  - the wildcard/method based routing is realtivey new, from Go 1.22 (February
    2024.  WOW)
  - some things not supported
    - sending custom 404 not found and 405 method not allowed
    - using regular expressions in route patterns or wildcards
    - matching multiple HTTP methods in a sigle route declaration
    - automatic support for OPTIONS requests
      - _allow clients to obtain parameters and requirements for
        specific resources and server capabilities without taking
        action on the resource or requesting the resource_
    - routing based on unusual things, like HTTP request headers (headers in the HTTP request)
    - if you need these, get a third-pouty router
    - recommended ones are httprouter, chi, flow, and gorilla/mux.  There's
      a blog post linked in the book with guidance
* Customizing responses
  - default response is 200 OK, Date header, and Content-Length / Content-Type
  - `w.WriteHeader(201)` to return a 201 (Created)
    - can only call it onces per response, get a warning if try again
  - if don't call w.WriteHeader explciitly, the first call to w.Write()
    will send the 200.
  - `net/http` package has constants for HTTP status codes
    - https://pkg.go.dev/net/http#pkg-constants
    - so like `http.StatusCreated` and `http.StatusTeapot`
  - also http.StatusText() for a human-readable version
* headers
  - can edit the _response header map_
  - add one via `w.Header().Add()`
    - e.g. `w.Header.Add("Server", "FORTRAN")
  - make sure add the header before calling w.WriteHeader() or
    w.Write().  By that time it's tooooo late.
  - can also Set(), Del(), Get(), and Values() as well.
  - Set will replace a header
  - Add will append (like for Cache-Control)
  - Del will remove all values for a key
  - Get retrieves the first element
  - Values is a slice of all values for a header
  - When using the set/add/del/get/values, the header name will
    be canonicalized using textproto.CanonicalMIMEHeaderKey
    - converts the first letter and any letter after a dash
      to upper case, all others lowercase
    - beware if you have a case-sensitive header
    - to avoid canonicalization, bash the header map directly
      w.Header()["X-XSS-Protection"] = []string("1; mode=block"}
  - for HTTP/2, Go will always bash the header names and
    values to lowercase per the spec
* Writing response bodies
  - we can w.Write() to blast a string. Nice and simple
  - more common to pass your http.ResponseWriter value to another
    function that writes the response
  - because the http.ResponseWriter value has a Write() method,
    it satisfies the `io.Writer` interface.
    - so can pass the ResponseWriter to anything that takes a Writer
    - so things like io.WriteString, and fmt.Fprint* family.
    - instead of w.Write([]byte("Blorf"))
    - can do
      - `io.WriteString(w, "blah")`
      - `fmt.Frint(w. "BlAh")`
* content sniffing
  - to set the Content-Type header automagically, it
    uses http.DetectCOntentTtype().  If it can't figure
    out, falls back to application/octet-scream
  - cannot distinguish JSON from plain text, so by
    default has text/plain
  - manually set with w.Header().Set("Content-Type", "application/json"))

* **Project Structure**
  - no defined structure from go-land.
  - "don't over-complicate things"
  - there is a popular method:
     - https://go.dev/doc/modules/layout#server-project
     - keep packages in an internal directory
     - keep all go commands in a `cmd` directory
```
root-dir/
  go.mod
  internal/
    snork/
      ...
    greeble/
      ...
  cmd/
    api-clerver/
      main.go
    metrics-spongflongle/
      main.go
    ...
  ... other project directory with non-go code
```
  - run with `go run ./cmd/web`
  - `internal` has a special behavior
    - packages in here can only be imported by code inside the
      _parent_ of the internal directory.  So for snibbage, can only
      be imported by code inside of the snibbage project
    - equivalently, can't be imported by code outside of our project
      - so folks can't creep on it without us being aware

### HTML Templating

* `.tmpl` files don't have any intrinsic special meaning or behavior
* Go's html/template pacakge, has a family of functions for
  safely parsing and rendering HTML templates.
* load the template with `ts := templaet.ParseFiles("path/to/file")` - either
  absolute or relative to the root of the project directory.
* Use with `err = ts.Execute(w, nil)` to actulaly run.
* http.Error sends a lightweight error message and status code back.
* There will be shared / boilerplate / HTML markup to include on every page
  (e.g. Headzor, navigation, meatdata inside the <head> element)
* Template is just regular HTML with some extra `{{actions}}` in double-braces.
* `{{define "name"}} ... {{end}}` defines a distinct nmae template called
  `base`, which contains content want to appear on every page.
* inside that use `{{template "title" .}}` and `{{template "main" .}}` actions
  tp denote that want to invoke other nmed templates (named title
  and main) at a given location in the HTML
  - the dot hiding in there represents any dynamic data that you want to
    pass to the invoked template. (covered later)
* _Partials_ - break out certain bits of HTML that can be re-used in
  different pages or layouts.
* `{{template}}` action invokes one template from another. There's
  also `{{block}} ... {{end}}`. Acts like {{template}} except
  it allows some default content if the template being invoked
  doesn't exist in the current template set.
  - this is useful when want to provide some default content
    (say a sidebar) which individual pages can override on a case-by-case
    basis if they need to. e.g.
```
{{define "base"}}
  <h1>An example template</h1>
  {{block "sidebar" .}}
    <p>Default snidebar content></p>
  {{end}}
{{end}}
```
  - don't _need_ to include default content in the {{block}}/{{end}}
  actions. The invoked template acts like it is optional.  If
   the template exists in the template set, then it will be
   rendered. But if it doesn't, then nothing is displayed
  - nice features of http.FileServer
    - sanitizes all request paths by running them through path.Clean().
      Removes . and .. elements
    - range requests supported (for e.g. large files)
    - Last-Modified and If-Modified-Since transparently supported
      - if haven't changed, get a 304 Not Modified status code
    - Content-Type is automatically set from the file extension via
      mime.TypeByExtension(). Can add your own custom extensions and
      content types by mime.AddExtensionType() if necessary.
  - performance
    - once served once before, the FS cache will be serving from RAM
  - single files in a handler vis http.ServeFile. e.g.
```
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./ui/static/file.zip")
}
```
  - note that http.ServeFile doesn't automaticaly sanitze the file path.
    - if constructing a path from untrused user input, sanitize with
      filepath.Clean() before using.
  - Disabling directory listings
    - easiest is to add a blank index.html. User will get 200 OK. Do
      it for all subdirectories via
      `find ./ui/static -type d -exec touch {}/index.html \;`
    - a better solutuion is to make a custom implementation
      of http.FileSystem and have it reutrn os.ErrNotExist error
      for any directories.
      - https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings

* <b>Static Files</b>
  - example put them into ui/static/css|img|js
  - net/http ships with http.FileServer handler that can use
    to serve files from a specific directory.  Like all
    GET requests that begin with "/static/"

* The http.Handler interface
  - Theory from chapter 2.10
  - "handler" : an object which satisfies the http.Handler interface
    - which is ServeHTTP(ResponseWriter, *Request)
  - so a handler object must have a ServeHTTP method with that exact
    signature
  - so something like

```
type home struct {}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   w.Write(...)
}
```
  - and can register with the servemux using Handle.
    `mux.Handle("/", &home())
  - this is kind of long-winded and confusing
  - more common to write them as a normal function.
    - Um, Actually, these aren't really _handlers_ because it doesn't
      have a ServeHTTP() method
    - But can transform it into a handler with http.HandlerFunc()
      - "adds a ServeHTTP() method to the home function"
      - so when ServeHTTP is called, it turns around and calls
        the code inside of the original function.
      - "a roundabout but convenient way of coercing a normal function
         into satisfying the http.Handler interface
    - HandleFunc() is just some syntactic sugar that transforms a
      function into a handler and registers it in one step.
  - Chaining handlers
    - might have noticed that http.ListenAndServe() takes an http.Handler,
      but we're passing a servemux
    - servemux adopts ServeHTTP, so it can be passed in
    - servmux is just a special kind of handler, but instead of
      providing a resposne itself, passes the request on to a second
      handler. a.k.a. Chaining handlers
      - "very common idiom in Go"
    - our clerver is getting a new HTTP request
      - calls servemux's ServeHTTP
      - looks up the relevant handler based on method/path
      - calls that handler's ServeHTTP
    - can think of a Go web app as a _chain_ of ServeHTTP methods
      being called one after anohter
  - All incoming HTTP requests are served in their own goroutine.
    - so for busy servers, it's very likely that the code in or called
      by your handlers will e running concurrently.
    - beware of races
      - blogpost: https://www.alexedwards.net/blog/understanding-mutexes



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
* Interfaces
  - https://www.alexedwards.net/blog/interfaces-explained
  - interface is like a definiton.  Describes the exact methods that
    some other type must have
  - e.g. the fmt.Stringer interface:
```
type Stringer interface {
     String() string
}
```
  - something "satisfies" this interface / "implements" this
    interface if it has a method with that exact same String() string.
  - e.g. Blorf implements / satisfies this interface
```
type Blorf struct {
    Food string
    Greeble string
}
func (b Book) String() string {
    return fmt.Sprintf("arhghghghg %s - %s", b.Food, b.Greeble)
}
```
  - so looks like no explicit conformance, just "hey you adopt the
    proper stuff"
  - why useful?
    - reduce duplication and boilerplate
    - easier to mock instead of using real objects in tests
    - enforce decoupling
  - there isn't an explict declaration of conformance
  - you kind of have to know that something conforms (like file and
    Buffer each have the Write (Writer interface) method
    - useful interfaces: https://www.alexedwards.net/blog/interfaces-explained#useful-interface-types and https://gist.github.com/asukakenji/ac8a05644a2e98f1d5ea8c299541fce9
* Map has Add(), Set(), Del(), Get(), Values()

### configuration and error handladge

* Managing configuration settings
  - currently we have network address and static files hard-coded
  - kinda annoying if need different settings for dev/test/prod)
  - command line flags are common
    - easiset is do something like
    - `addr := flag.String("addr", ":4000", "HTTP network address")`
    - call flag.Parse() to do the thing
    - and then use later via `*addr`
      - dereference the `addr` pointer and get to the underlying string
  - suggests using development-convenient default values
  - Can specify the expected type of a flag (flag.Int(), Bool(), Float64(), Sring(), Duration())
  - `go run ./cmd/web -help` to print help
* can also do environment varialbes
  - link to 12-factor app. https://12factor.net/config
  - `addr := os.Getenv("SNIBBAGE_ADDR")
    - can't specify a default (get an empty string if the env var doesn't exist)
    - don't get -help
    - return value is always a string
  - so, pass the env var as a command-line flag!
    - `go run ./cmd/web -addr=$SNIBBAGE_ADDR
* Boolean flags
  - omitting the value (but providing the flag) is the same as writing -flag=true
  - so `go run ./blah -flag=true` is the same as `go run ./blah -flag`
  - have to use -flag=false to set it to false
* pre-existing variables
  - can parse int addresses of pre-existing variables
    - e.g. flag.StringVar(), IntVar(), BoolVar(), etc
```
type confing struct {
     addr string
     staticDir string
}

var cfg config
flag.StringVar(&cfg.addr", "addr", ":4000", " seem to be a spoon
flag.Parse()
```

* Structured Logging
  - log.Printf and log.Fatal are easy to use, using the Go's standard logger
    from the `log` package.
    - outpus with the local date and time and writes to stderr
  - many applications that's Good Enough
  - if you do :alot: of logging, might want them easier to filter and work with
    - like severities
    - or a consisten structure so it's easy for tools to parse
  - `log/slog` can create _structured loggers_ that output log entries in
    a set format.  including:
    - timestamp (ms precision)
    - severity (Debug / Info(e) / Warn / Error)
    - log message (string value)
    - any nubmer of attributes (key/value pairs) with additional info
  - structured loggers have a _structured logging handler_ assoicated with
    them. (distinct from HTTP handlers)
    - this handler controls how log entries are formatted and where they
      ge written to.
  - create with
```
loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{...})
logger := slog.New(loggerHandler)
// or combine
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{...}))
```
  - use NewTextHandler to make a structured logging handler, takes two args
    - the write destination.
    - pointer to a slog.HandlerOptions struct (https://pkg.go.dev/log/slog#HandlerOptions)
    - if happy with the defaults, can pass ni
  - once you have a structured logger, can write an entry with a specific
    severeity.  `logger.Info("blah")` or `logger.Error("out of fuds")`
    - these are variadic methods that can accept an arbitrary number of
      key-value pairs, like
    - `logger.Info("request receivededed", "method", "TUNA", "path", "/splunge")`
      - yields `time=2024-03-18T11:29:23.000+00:00 level=INFO msg="request received" method=GET path=/`
    - keys must be strings, but values of any type.
    - if keys or values contain `"` or `=` or whitespace, will be
      wrapped in double-quoes.
    - ther is no equivalent of log.Fatal()
  - type safety of key value pairs
    - can do the variadic thing
    - or do Any: `logger.Info("blargle", slog.Any("addr", ":4000"))`
    - or a type-pecuilar function, like slog.String(), Int, Bool, Time,
      and Duration
      - `logger.Info("grumble cake", slog.String("addr", ":4000"))`
  - slog.NewTextHandler makes a handler that writes plaintex log entries.
  - can also write as JSON objects, using `slog.NewJSONHandler()`
  - can filter the noise by setting the log level.      
    - by default uses Info
    - use slog.HandlerOptions to override this
```
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
```
  - can also include file name and line number, via `AddSource: true` in the
    HandlerOptions
* Decoupled logging
  - like writing to os.Stdout.  app and the logging are decoupled.  The
    app isn't concerned with the routing or storage of logs
  - like in staging or prod environments, can redirect the stream to a final
    dstination, say to disk or to Splunk
  - e.g. `go run ./cmd/web >> /tmp/web.logge`

* Concurrent logging
  - custom loggers via slog.New) are concurrency-safe. Share and enjoy
  - multiple sructured loggers writing to the same destination have to be
    careful with and ensure the underlying Write() is also safe for
    concurrent use

* Dependency Injection
  - handlers.go uses the old standard logger, not our cool new hotness.
    -  blog post. https://www.alexedwards.net/blog/organising-database-access
  - for applications where all your handlers are in the same package, can
    inject dependencies to put them into a custom `application` struct, then
    define handler functions as methods against `application`
  - Closures can be used.
    - that application thing won't work if handlers are spread across packages
    - so make a stand-o-lone config package that exports an Application struct,
      and have your handler functions close over this to form a closures.
      This gist https://gist.github.com/alexedwards/5cd712192b4831058b21 has
      a more fleshy example. Plus this from the book:
```
// package config
type Applications struct {
    Logger *slog.Logger
}

// package greeblebork

func ExampleHandler(app *config.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ...
        ts, err := template.ParseFiles(files...)
        if err != nil {
            // app captured by the argument passed in
            app.Logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
            http.Error(w, "interfnace error", http.StatusInternalServerError)
            return
        }
        ...
    }
}

// package main
func main() {
    app := &config.Application {
        Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
    }
    ...
    mux.Handle("/", greeblebork.ExampleHandler(app))
}
```

* Centralized Error Handling
  - move error handling code into helper methods
  - help separate our concerens (blog https://deviq.com/principles/separation-of-concerns) 

* extra info(e)
  - can use debug.Stack() to get a stack trace for the _current goroutine_
  - can isolate the application routes, moving the setup out of main
  - is reducing main() to
    - parsing the runtime configuration settings for the application
    - establishing the dependencies for the handlers
    - runnign the http server

### Databasey things

* brew install mysql
* brew services start mysql
* brew services stop mysql
* mysql -u root -
* user 'web'@'localhost', password 'snork'
* `mysql -D snippetbox -u web -p`
* get the driver via `go get`
  - comprehensive list: https://go.dev/wiki/SQLDrivers
  - but using the go-sql-driver/mysql
  - go get will recursively download any dependencies that the package has
  - uses semantic versioning (so attach @v1 to get latest v1, like v1.666.52)
  - omit the suffix to get the latest/greatest

```
% cd $PROJECT_DIR
% go get github.com/go-sql-driver/mysql@v1
```

* go.mod grows a `require` section with the actual version numbers of the
  packages it is using.  Makes it easy to have multiple projects on the
  same machien use different versions of the same project.
  - though that makes my head hurt thinking about it...
  - `// indirect` indicates a package doesn't directly appera in any import
    statement
  - the go.sum file has the checksums of the packages
    - commit it
    - https://go.dev/wiki/Modules#should-i-commit-my-gosum-file-as-well-as-my-gomod-file
  - `go mod verify` will verify the checksums
  - `go mod download` will download all the dependencies of the project
  - when running / testing / building, the exact package versions are used
  - handy for creating reproducible builds
  - once package is in go.mod, package and version are fixed.
  - `go get -u github.com/....` to update to latest minor / patch release
  - `go get -u github.com/...@v2.0.0` - upgrade to a specific version
  - `go get github.com/...@none` - to forget the package
    - or if removed all references in code, can `go mod tidy`

* Database programming
  - connection pool
  - use sql.Open.  First arg is driver name, second is data soucre name, 
    a.k.a. Connection String or DSN
    - which is database-peculiar. docs https://github.com/go-sql-driver/mysql#dsn-data-source-name
    - parseTime=true is a driver-specific parameter. This one converts
      sql time and date fields to go time.Time objects
    - returns a sql.DB object, *not* a connection. It's a pool
      - go manages the opening and closing of connections automagically
      - safe for concurrent access, so can use from web handlers safely
      - pool is intended to be long-lived.  like make in `main()` and then
        pass it to the handlers.
      - calling sql.Open() in an http handler is grounds for immediate 
        dismissal and community taunting
      - sql.Open() doesn't actually create any connections, just initialzes
        the pool for future use.  Actual connections are lazy\
      - db.Ping() verifies things are set up correctly to create a connection
        and check for errors. 
```
db,err := sql.Open("mysql", "web:pass@/snippetbox?parseTime=true")
if err != nil {
    // cry
}
```
  - not really database, but covered here.  The import for mysql is like `_ "github.com/go-sql-driver/mysql"`
    - main.go doesn't actually use anything in there, so go will complain.
    - we need the driver's `init()` to run so it can register itself
    - so work around to alias the package name to the blank identifier
      - standard practice for most sql drivers

* data modelling
  - a.k.a. Service Layer or Data Access Layer
  - encapslate the code for working with the databsae in a packge to
    the rest of of the application

* modules
  - subdir of internal/
  - pull in via `"snibbage.borkware.com/internal/models"` (or whatever is in your go.mod)

* sql queries
  - three different methods:
    - DB.QueryRow() - single row
    - DB.Query() - multiple rows
    - DB.Exec() - for non-select statemetns (e.g. insert, delete)
      - returns a sql.Result, with two methods
      - LastInsertId() - an integer generated by the db in response to the command.
        - typically from an autoincrement column
        - not supported by postgresql (https://github.com/lib/pq/issues/24)
          - use QueryRow with Returning
      - RowsAffected() - number of rows (int64) affected by the statement
      - it's common to ignore the sql.Result (`_, err := m.DB.Exec...`)
      - behind the scenes, it creates a new prpared statement, passes the parameter values. The db ewxecutes it. when done, it cloes/deallocates the prepared statement
    - statement is a string, with ? as placeholders. It is database dependent. Mysql is ?, pgsql is $N ($1 $2, etc)
```
stmt := `insert into snippets (title, content, created, expires)
     values(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
```

* Single record queries
  - use a select statement to return a single record (like using an id
    for a primary key and doing an equality check)
  - after QueryRow, do a row.Scan to poke values into the structure
    - errors are deferred ntil Scan, so can chain the QueryRow and the Scan
  - the driver automatically converts the raw sql db output to the required
    go types. Things should Just Work. Like char/varchar/text -> string, boolean -> bool, etc
    - due to a quirk of mysql, need to do that `parseTime=true` on the DSN
    - otherwise it returns them as []byte objects. (because, of course, that's 
      so useful as a default behavior...)
  - the statement will return an sql.ErrNoRows directly.  Author suggests
    using a custom error type instead (so if it's ErrNoRows, return 
    ErrNoRecordOr8Track) to encapsulate the model from its expression as 
    a database
    - recommended way is using Is to compare errors.  `if error.Is(err, models.ErrNoRecord)` rather than older-go use of the equality operator
      - because you can wrap errors: https://go.dev/blog/go1.13-errors#wrapping-errors-with-w
      - and Is will unwrap errors as necessary checking for a match
      - there's also an .As(), to check if a (maybe wrapped) error has a
        specifc type

* Multiple record queries
  - return multiple rows. queries like
```
select id, title, content, created, expres from snippets
where xpires > UTC_TIMESTAMP() order by id desc limit 10
```
  - iterate over rows.Next
  - calling `defer rows.Close()` is critical - make sure the resultset is
    closed so the underlying database connection is recycled

* Miscellaney
  - Go doesn't do well is handling NULL values
    - e.g. can't convert NULL into a string when doing a `Scan()` to stuff
      a structure
    - roughly the fix is to change the field scanning into from a `string`
      to a `sqlNullStrong`.  See this gist: https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9
    - but in general, avoid null values altogether
  - transactions
    - any calls to Exec/Query/QueryRow will opportunisticaly use any
      connection from the database pool
    - but say you need to balance a lock tables with an unlock tables, which
      need be done same on the same connection
    - wrap in a transaction
      `tx, err := m.DB.Begin()`, `defer tx.Rollback()`, do the work, `tx.Commit()`
    - _must_ call Rollback() or Commit() before leaving function
  - prepared statements
    - Exec/Query/QueryRow all use prepared statements behind the scenes.
    - could use DB.Prepare()
    - (code snibbage in the book)
    - prepared statements exist on db connections, and because we have a pool,
      sql.Stmnt tries to get to the same pool. If it's in use, you gotta
      wait
    - also possible that a large number of prepared statements will be
      created on multiple connections
    - so, wait for manually preparing statements to prove it's a problem


### dig in to

- byte slice
  - `[]bytes("blah")` syntax
- are there anonymous functions?
- string interpolation?
- strings with backticks vs double-quotes (multi-line?)
- Header Map.
- slices in general
- method signature?
  - `func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {`
- "%+v" Fprintf format specifier
- why
   - `mux.Handle("/", &home())`
- how does method/function resolution work when using interfaces
  polymorphically
- dig deeper into templating: https://templ.guide/
- what is `Snippet{}` vs `Snippet()`?
- also `var s Snippet` // initialize a new zeroed snippet struct
- piece-wise initialization of structs? e.g.
```
slog.HandlerOptions{
    Level: slog.LevelDebug,
}
```
- what is this bulk asignment thing?
```
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
        )
```

### Emacs fun

* https://github.com/dominikh/go-mode.el
* available from melpa. (might need to refresh the package library)
* M-x gofmt to trigger go-format
* can also set a save hook, but I found that annoying last time I tried
* go mode hook -  sets 4-space tabs. also preserves tabs. Uncomment the line to use spaces.

```
(add-hook 'go-mode-hook
  (lambda ()
    (setq-default)
    (setq tab-width 4)
    (setq standard-indent 4)
;;    (setq indent-tabs-mode nil)
))
```

