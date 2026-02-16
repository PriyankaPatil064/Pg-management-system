package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"pg-management-system/internal/database"
	"pg-management-system/internal/gql"
	"pg-management-system/internal/handlers"
	"pg-management-system/internal/middleware"

	"net/http/pprof"

	"github.com/graphql-go/handler"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// Try loading from root directory if running from cmd/server
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	database.Connect()
	database.InitSchema()

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PG Management System is running"))
	}).Methods("GET")

	// Auth Routes
	r.HandleFunc("/auth/google/login", handlers.GoogleLogin).Methods("GET")
	r.HandleFunc("/auth/google/callback", handlers.GoogleCallback).Methods("GET")

	// Protected API Routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Room Routes
	api.HandleFunc("/rooms", handlers.CreateRoom).Methods("POST")
	api.HandleFunc("/rooms", handlers.GetAllRooms).Methods("GET")
	api.HandleFunc("/rooms/{id}", handlers.GetRoomByID).Methods("GET")
	api.HandleFunc("/rooms/{id}", handlers.UpdateRoom).Methods("PUT")
	api.HandleFunc("/rooms/{id}", handlers.DeleteRoom).Methods("DELETE")

	// Guest Routes
	api.HandleFunc("/guests", handlers.CreateGuest).Methods("POST")
	api.HandleFunc("/guests", handlers.GetAllGuests).Methods("GET")
	api.HandleFunc("/guests/{id}", handlers.GetGuestByID).Methods("GET")
	api.HandleFunc("/guests/{id}", handlers.UpdateGuest).Methods("PUT")
	api.HandleFunc("/guests/{id}", handlers.DeleteGuest).Methods("DELETE")

	// Payment Routes
	api.HandleFunc("/payments", handlers.CreatePayment).Methods("POST")
	api.HandleFunc("/payments", handlers.GetAllPayments).Methods("GET")
	api.HandleFunc("/payments/{id}", handlers.GetPaymentByID).Methods("GET")
	api.HandleFunc("/payments/{id}", handlers.UpdatePayment).Methods("PUT")
	api.HandleFunc("/payments/{id}", handlers.DeletePayment).Methods("DELETE")
	api.HandleFunc("/payments/guest/{id}", handlers.GetPaymentsByGuestID).Methods("GET")

	// GraphQL Route (Protected)
	h := handler.New(&handler.Config{
		Schema:   &gql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})
	api.Handle("/graphql", h)

	// Profiling Routes
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	// Print all registered routes
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, _ := route.GetMethods()
		fmt.Printf("Registered Route: %s %v\n", path, methods)
		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port", port)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	hWithCors := c.Handler(r)
	log.Println("Attempting to listen on :" + port)
	err := http.ListenAndServe(":"+port, hWithCors)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
