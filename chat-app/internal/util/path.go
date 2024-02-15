package util

import (
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

func FilePathPrefix(fileAddr string) string {
	if strings.TrimSpace(fileAddr) == "" {
		// no photo image
		fileAddr = "system/no-photo.jpg"
	}

	address, _ := url.JoinPath(viper.GetString("app.files_download_prefix"), fileAddr)
	return address
}
