// @Title
// @Description
// @Author
// @Update
package config

import "os"

func Config(envvar string) string {
	return os.Getenv(envvar)
}
