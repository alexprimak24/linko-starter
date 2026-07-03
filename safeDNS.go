package main

import "net/url"

func _(dsn string) string {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "invalid dsn"
	}
	_, ok := parsed.User.Password()
	if !ok {
		return parsed.String()
	}
	parsed.User = url.UserPassword(parsed.User.Username(), "***")
	return parsed.String()
}

// logger.Info("connecting to database",
// 	"dsn", safeDSN(dsn),
// )
// 2024-01-15T10:30:45.123Z INFO msg="connecting to database" dsn="postgres://admin:***@db.example.com/appdb"
