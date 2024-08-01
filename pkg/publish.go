package pkg

import (
	"net/http"
)

func Publish(req *http.Request, res *http.Response) {
	req.Header.Get("Authorization")
}
