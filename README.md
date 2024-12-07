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

### Syntax

* `func home(w http.ResponseWriter, r *http.Request) {`
  - function that returns nothing.
  - takes a response writer
  - takes a pointer to a struct


### dig in to

- byte slice
  - `[]bytes("blah")` syntax
- are there anonymous functions?