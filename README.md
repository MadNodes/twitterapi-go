# twitterapi for go sdk

> https://docs.twitterapi.io/

## install
```bash
go get github.com/MadNodes/twitterapi-go
```

## usage
```go
import (
	"github.com/MadNodes/twitterapi-go"
)

func main() {
    x := New(xApiKey, WithProxy(proxy))

	if err := x.Login(username, email, password, nil); err != nil {
		return
	}
	
    // TODO something
}
```

## todo

- [ ] test
- [ ] add full api support
- [ ] add example code
