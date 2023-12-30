package utils

import (
	"log"
	"net/url"
)

func AddQueryToUrl(link string, query map[string]string) string {
	u, err := url.Parse(link)

	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()

	for k, v := range query {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
