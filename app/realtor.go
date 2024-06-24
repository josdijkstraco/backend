package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func GetRealtorLicense(w http.ResponseWriter, r *http.Request) {
	logrus.Info("GetRealtorLicense")

	vars := mux.Vars(r)
	license := vars["license"]
	logrus.WithField("license", license).Info("License check")

	connection, _ := getDBConnection()
	defer connection.Close()

	userData := getRealtorFromDBByLicense(license, connection)
	if userData == nil {
		sendResponse(w, http.StatusNotFound, "License not found")
		return
	}

	rsp := UserLoginResponse{
		ID:        userData.id,
		Name:      userData.firstName + " " + userData.lastName,
		FirstName: userData.firstName,
		LastName:  userData.lastName,
		IsAdmin:   userData.isAdmin,
	}

	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.WithField("data", string(js)).Info("GetRealtorLicense response")
}
