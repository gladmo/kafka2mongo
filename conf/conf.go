package conf

import "os"

func Getenv(env, def string) string {
	e := os.Getenv(env)
	if e != "" {
		return e
	}

	return def
}
