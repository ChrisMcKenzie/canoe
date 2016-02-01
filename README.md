Canoe
=====

Canoe is streaming template language experiment aim at providing a better experience
than current front-end templating engines (PHP, handlebars, go-template). It 
will accomplish this by providing a auto optimizing interpreter, as well as a toolchain for 
testing, and formatting your templates. 

Canoe will offer both a CGI and HTTP interface for handling requests.

**Early Example:** __Subject to change__
```
<=
  import (
    "http"
  )

  func test() {
    res, err := http.get("http://google.com/")
    if err != nil {
      fatal("[ERR]:", err)
    }

    return res
  }

  log(test())

  title := "hello world"
=>

<h1> {{title}} </h1>
```
