Canoe
=====

Canoe is a Web Application serving platform for front-end microservice architectures allowing development teams to create small "fragments" of ui that can then be rendered client-side via a pipelining system that uses HTTP2 Server Push

![](canoe.png)

## Basic Server Setup

```go
package main

import (
	"log"
	"net/http"

	"github.com/chrismckenzie/canoe"
)

func main() {
	c := canoe.NewMemoryCache()
	r := canoe.NewHTTPFragmentResolver(http.DefaultClient, c)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/fragments/", canoe.NewFragmentHandler(canoe.WithHTTPResolver(r)))
	http.Handle("/", canoe.NewHandler("/fragments/", http.Dir("examples"), c))

	err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Creating An App Frame

An app frame is an html file that contains `canoe-fragment` elements in places when canoe will render a fragment.

```html
<html>
  <head>
    <title> testing canoe </title>
  </head>

  <body>
    <h2> hello World </h2>

    <canoe-fragment href="https://my-nav-header-service:8080/"></canoe-fragment>
  </body>
</html>
```

The `canoe-fragment` tag is an html custom element that will be in charge of loading and managing the state of the fragment.
