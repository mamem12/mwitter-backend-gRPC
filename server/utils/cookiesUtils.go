package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

const (
	grpcgateway_cookie = "grpcgateway-cookie"
)

type User struct {
	ID string
}

func GetCookie(ctx context.Context) string {
	cookies := metautils.ExtractIncoming(ctx).Get(grpcgateway_cookie)

	if cookies == "" {
		return ""
	}

	return cookies
}

func ParseCookie(cookie string, key string) (string, error) {

	cookies := transferCookieSet(cookie)

	var userValue string

	for _, v := range cookies {
		if v.Name == "user" {
			userValue, _ = url.QueryUnescape(v.Value)
			break
		}
	}

	user := make(map[string]interface{})

	err := json.Unmarshal([]byte(userValue), &user)

	if err != nil {
		return "", err
	}

	switch user[key].(type) {
	case string:
		return user[key].(string), nil
	case int:
		return strconv.Itoa(user[key].(int)), nil
	case float64:
		return fmt.Sprintf("%.f", user[key].(float64)), nil
	default:
		return "", fmt.Errorf("not found type!?")
	}

}

func transferCookieSet(cookies string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", cookies)
	req := http.Request{Header: header}

	return req.Cookies()
}
