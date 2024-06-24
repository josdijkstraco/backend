package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_ "golang.org/x/crypto/bcrypt"
)

func RunServer() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/products", GetProducts).Methods("GET")
	myRouter.HandleFunc("/api/product/{id}", GetProduct).Methods("GET")

	myRouter.HandleFunc("/api/listings", GetListings).Methods("GET")
	myRouter.HandleFunc("/api/listing/{id}", GetListing).Methods("GET")
	myRouter.HandleFunc("/api/address/{addr}", GetAddress).Methods("GET")
	myRouter.HandleFunc("/api/listing", CreateListing).Methods("POST")

	myRouter.HandleFunc("/api/users/register", RegisterUser).Methods("POST")
	myRouter.HandleFunc("/api/users/login", LoginUser).Methods("POST")
	myRouter.HandleFunc("/api/users/profile/update", UpdateProfile).Methods("PUT")
	myRouter.HandleFunc("/api/users/profile", GetProfile).Methods("GET")
	myRouter.HandleFunc("/api/users/listings", GetUserListings).Methods("GET")

	myRouter.HandleFunc("/api/realtor/license/{license}", GetRealtorLicense).Methods("GET")

	myRouter.Use(loggingMiddleware)

	url := ":8000"

	//os.Setenv("TOKEN_SECRET", "Thisisasecrettoken")

	getDBConnection()

	go run(myRouter)

	logrus.WithField("url", url).Info("Started server")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	logrus.Info("Shutdown request received")
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Request received")
	enableCors(&w)

	connection, _ := getDBConnection()
	rows, err := connection.Query("SELECT * FROM jpar.product")
	if err != nil {
		logrus.WithError(err).Error("Error retrieving data from DB")
		w.WriteHeader(500)
		return
	}

	result := make([]Product, 0)
	for rows.Next() {
		var product Product
		rows.Scan(&product.ID, &product.UserID, &product.Name, &product.Image, &product.Brand, &product.Category, &product.Description, &product.Rating,
			&product.NumReviews, &product.Price, &product.CountInStock, &product.CreatedAt)
		result = append(result, product)
	}
	rows.Close()
	connection.Close()

	js, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	logrus.WithField("data", string(js)).Info("Response")
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	logrus.WithField("id", id).Info("Request received")
	enableCors(&w)

	connection, _ := getDBConnection()
	row := connection.QueryRow(fmt.Sprintf("SELECT * FROM jpar.product WHERE id = %s", id))

	product := Product{}
	row.Scan(&product.ID, &product.UserID, &product.Name, &product.Image, &product.Brand, &product.Category, &product.Description, &product.Rating,
		&product.NumReviews, &product.Price, &product.CountInStock, &product.CreatedAt)
	connection.Close()

	js, _ := json.Marshal(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	logrus.WithField("data", string(js)).Info("Response")
}

func run(router *mux.Router) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)

	if err := http.ListenAndServe("127.0.0.1:8000", handler); err != nil {
		if err != http.ErrServerClosed {
			logrus.WithError(err).Panic("Error starting server")
		} else {
			logrus.Info("Server shutdown")
		}
	}
}

func shutdown(server *http.Server) {
	ctx := context.Background()
	server.Shutdown(ctx)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetListings(w http.ResponseWriter, r *http.Request) {
	logrus.Info("GetListings request received")
	enableCors(&w)

	connection, _ := getDBConnection()

	query := `
		SELECT l.id, l.realtor_id, oa.full_address, oa.mailing_city, oa.mailing_zip, l.price, l.commission, l.image
		FROM jpar.listing l
		JOIN jpar.owner_address oa ON oa.strap = l.strap
	`

	rows, err := connection.Query(query)
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
	logrus.WithField("data", string(js)).Info("GetListings Response")
}

func GetListing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	logrus.WithField("id", id).Info("GetListing request received")
	enableCors(&w)

	connection, _ := getDBConnection()

	query := `
		SELECT l.id, l.realtor_id, j.first_name, j.last_name, oa.full_address, oa.mailing_city, oa.mailing_zip,
			l.price, l.commission, l.image,
			b.number_bedrooms, b.number_full_baths, b.number_3_qtr_baths, b.number_half_baths,
			b.total_finished_sf, oa.owner_name, oa.legal_description, l.strap, oa.folio			
		FROM jpar.listing l
		JOIN jpar.owner_address oa ON oa.strap = l.strap
		JOIN jpar.building b ON b.strap = l.strap
		JOIN jpar.user j ON j.id = l.realtor_id
		WHERE l.id = %v;
`

	row := connection.QueryRow(fmt.Sprintf(query, id))
	var listing Listing
	row.Scan(&listing.ID, &listing.RealtorID, &listing.FirstName, &listing.LastName, &listing.Address, &listing.City, &listing.PostalCode,
		&listing.Price, &listing.Commission, &listing.Image,
		&listing.NumberBedrooms, &listing.NumberFullbacts, &listing.NumberThreeQtrBaths,
		&listing.NumberHalfBaths, &listing.FinishedSquareFt, &listing.Owner, &listing.LegalDescription, &listing.Strap,
		&listing.ParcelNumber)

	js, _ := json.Marshal(listing)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	logrus.WithField("data", string(js)).Info("Response")
}

type SelectResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["addr"]

	logrus.WithField("addr", addr).Info("GetAddress request received")
	enableCors(&w)

	connection, _ := getDBConnection()
	query := fmt.Sprintf("SELECT full_address, strap FROM jpar.owner_address WHERE full_address ILIKE '%s%%'", addr)
	fmt.Println(query)
	rows, err := connection.Query(query)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}

	querySet := []SelectResponse{}
	for rows.Next() {
		var address string
		var strap string

		rows.Scan(&address, &strap)
		rsp := SelectResponse{
			Value: strap,
			Label: address,
		}
		querySet = append(querySet, rsp)
		fmt.Println(address)
	}
	rows.Close()
	connection.Close()

	js, _ := json.Marshal(querySet)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	logrus.WithField("data", string(js)).Info("Response")
}

type CreateListingRequest struct {
	Strap      string `json:"strap"`
	Price      string `json:"price"`
	Commission string `json:"commission"`
	ExpiryDate string `json:"expiryDate"`
	Image      string `json:"image"`
	RealtorID  int    `json:"realtorId"`
}

func CreateListing(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", string(body))
		w.WriteHeader(500)
		return
	}

	listing := &CreateListingRequest{}
	json.Unmarshal(body, listing)

	logrus.WithField("data", listing).Info("CreateListing data")

	connection, _ := getDBConnection()
	query := `
		INSERT INTO jpar.listing (strap, realtor_id, price, commission, expiry_date, image)
		VALUES ('%s', %v, '%s', '%s', '%s', '%s' )
		`
	_, err = connection.Exec(fmt.Sprintf(query, listing.Strap, listing.RealtorID, listing.Price,
		listing.Commission, listing.ExpiryDate, listing.Image))
	if err != nil {
		fmt.Printf("Error inserting listing: %v\n", err)
		w.WriteHeader(500)
		return
	}

	logrus.WithField("listing", listing).WithField("body", string(body)).Info("Creating listing")
	w.WriteHeader(200)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	logrus.Info("CreateToken")

	claims := UserClaims{}
	token, err := NewAccessToken(claims)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return
	}

	fmt.Printf("Token: %s\n", token)

	rsp := UserLoginResponse{}
	js, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)

	logrus.WithField("data", string(js)).Info("CreateToken response")
}
