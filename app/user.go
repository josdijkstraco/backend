package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/users/register
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("RegisterUser")
	//enableCors(&w)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", string(body))
		w.WriteHeader(500)
		return
	}

	entry := &UserRegisterRequest{}
	json.Unmarshal(body, entry)
	logrus.WithField("registerdata", entry).Info("Register user data")

	// bcrypt password
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(entry.Password), bcrypt.MinCost)
	if err != nil {
		logrus.WithError(err).Error("BCrypt error")
		w.WriteHeader(400)
		return
	}
	logrus.WithField("password", string(passwordBytes)).Info("Encrypted")

	// store in DB
	connection, err := getDBConnection()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer connection.Close()

	query := fmt.Sprintf("SELECT jpar.register_user('%s', '%s', '%s', %v)",
		entry.License, entry.Email, string(passwordBytes), entry.Notifications)
	row := connection.QueryRow(query)
	var id int
	if err := row.Scan(&id); err != nil {
		sendResponse(w, http.StatusInternalServerError, err.Error())
		logrus.WithError(err).Error("Error registering user")
		return
	}

	logrus.WithField("id", id).Info("User registered")

	// get data
	userData := getUserFromDBByID(id, connection)
	logrus.WithField("data", userData).Info("Returned from DB")

	// create Token
	claims := UserClaims{
		Id:        userData.id,
		FirstName: userData.firstName,
		Lastname:  userData.lastName,
		License:   userData.license,
	}
	token, err := NewAccessToken(claims)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}

	// UserRegisterResponse
	rsp := UserRegisterResponse{
		AccessToken: token,
		ID:          id,
		Name:        userData.firstName + " " + userData.lastName,
		FirstName:   userData.firstName,
		LastName:    userData.lastName,
		Email:       userData.email,
		License:     userData.license,
	}

	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.WithField("data", string(js)).Info("RegisterUser response")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("getUser")

	// vars := mux.Vars(r)
	// id := vars["key"]

	connection, err := getDBConnection()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer connection.Close()

	userData := getUserFromDBByID(1, connection)

	rsp := GetUserResponse{
		ID:        1,
		Name:      userData.firstName + " " + userData.lastName,
		FirstName: userData.firstName,
		LastName:  userData.lastName,
		Email:     userData.email,
		License:   userData.license,
	}

	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.WithField("data", string(js)).Info("GetUser response")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("LoginUser")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", string(body))
		w.WriteHeader(500)
		return
	}

	req := &UserLoginRequest{}
	json.Unmarshal(body, req)
	logrus.WithField("username", req.Username).WithField("password", req.Password).Info("Login request")

	connection, err := getDBConnection()
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	userData := getUserFromDBByEmail(req.Username, connection)
	if userData == nil {
		sendResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userData.encryptedPassword), []byte(req.Password)); err != nil {
		sendResponse(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// create Token
	claims := UserClaims{
		Id:        userData.id,
		FirstName: userData.firstName,
		Lastname:  userData.lastName,
		License:   userData.license,
	}
	token, err := NewAccessToken(claims)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}

	rsp := UserLoginResponse{
		ID:          userData.id,
		Name:        userData.firstName + " " + userData.lastName,
		FirstName:   userData.firstName,
		LastName:    userData.lastName,
		License:     userData.license,
		IsAdmin:     userData.isAdmin,
		AccessToken: token,
	}

	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.WithField("data", string(js)).Info("LoginUser response")
}

// needs a bearer token
func GetProfile(w http.ResponseWriter, r *http.Request) {
	logrus.Info("GetProfile")

	token := GetTokenFromRequest(r)
	if token == "" {
		sendResponse(w, http.StatusBadRequest, "Missing bearer token")
		return
	}

	logrus.WithField("token", token).Info("Token received")

	claims := ParseAccessToken(token)
	logrus.WithField("claims", claims).Info("Claims received")

	connection, _ := getDBConnection()
	userData := getUserFromDBByID(claims.Id, connection)

	// update name, username, email, password
	rsp := GetUserResponse{
		ID:        userData.id,
		Name:      userData.firstName + " " + userData.lastName,
		FirstName: userData.firstName,
		LastName:  userData.lastName,
		Email:     userData.email,
		License:   userData.license,
		IsAdmin:   userData.isAdmin,
	}

	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.WithField("data", string(js)).Info("GetProfile response")
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	logrus.Info("UpdateProfile")

	// create new token

	// update name, username, email, password
	rsp := GetUserResponse{
		ID:        1,
		Name:      "Jos Dijkstra",
		FirstName: "Jos",
		LastName:  "Dijkstra",
		Email:     "josdijkstra@gmail.com",
		License:   "I0023223",
	}

	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.Info("UpdateProfile response")
}

func getUserFromDBByID(id int, conn *sql.DB) *userData {
	query := fmt.Sprintf("SELECT id, first_name, last_name, email, phone, password, is_active, is_admin, license, notifications FROM jpar.user WHERE id = %v", id)
	row := conn.QueryRow(query)
	if row.Err() != nil {
		return nil
	}

	data := &userData{}
	row.Scan(&data.id, &data.firstName, &data.lastName, &data.email, &data.phone,
		&data.encryptedPassword, &data.isActive, &data.isAdmin,
		&data.license, &data.notifications)

	return data
}

func getUserFromDBByEmail(email string, conn *sql.DB) *userData {
	query := fmt.Sprintf("SELECT id, first_name, last_name, email, phone, password, is_active, is_admin, license, notifications FROM jpar.user WHERE email = '%s'", email)
	row := conn.QueryRow(query)
	if row.Err() != nil {
		return nil
	}

	data := &userData{}
	row.Scan(&data.id, &data.firstName, &data.lastName, &data.email, &data.phone,
		&data.encryptedPassword, &data.isActive, &data.isAdmin,
		&data.license, &data.notifications)

	return data
}

func getRealtorFromDBByLicense(license string, conn *sql.DB) *userData {
	query := fmt.Sprintf("SELECT first_name, last_name, phone, license FROM jpar.realtor WHERE license = '%s'", license)
	row := conn.QueryRow(query)
	if row.Err() != nil {
		return nil
	}

	data := &userData{}
	row.Scan(&data.firstName, &data.lastName, &data.phone, &data.license)

	return data
}

func GetUserListings(w http.ResponseWriter, r *http.Request) {
	logrus.Info("GetUserListings")

	token := GetTokenFromRequest(r)
	if token == "" {
		sendResponse(w, http.StatusBadRequest, "Missing bearer token")
		return
	}

	logrus.WithField("token", token).Info("Token received")

	claims := ParseAccessToken(token)
	logrus.WithField("claims", claims).Info("Claims received")

	connection, _ := getDBConnection()
	defer connection.Close()

	query := `
		SELECT l.id, l.realtor_id, oa.full_address, oa.mailing_city, oa.mailing_zip, l.price, l.commission, l.image
		FROM jpar.listing l
		JOIN jpar.owner_address oa ON oa.strap = l.strap
		WHERE l.realtor_id = %v
	`

	rows, err := connection.Query(fmt.Sprintf(query, claims.Id))
	if err != nil {
		logrus.WithError(err).Error("Error retrieving data from DB")
		w.WriteHeader(500)
		return
	}

	result := make([]Listing, 0)
	for rows.Next() {
		var listing Listing
		rows.Scan(&listing.ID, &listing.RealtorID, &listing.Address, &listing.City, &listing.PostalCode,
			&listing.Price, &listing.Commission, &listing.Image)
		result = append(result, listing)
	}
	rows.Close()

	js, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	logrus.WithField("data", string(js)).Info("Response")

}
