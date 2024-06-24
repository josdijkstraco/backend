package app

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func sendResponse(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
	logrus.WithField("message", message).Error("Sending error response")
}

func GetTokenFromRequest(r *http.Request) string {
	header := r.Header.Get("Authorization")
	result := strings.Split(header, " ")
	if len(result) != 2 {
		logrus.WithField("header", header).Error("No token found")
		return ""
	}
	return result[1]
}
