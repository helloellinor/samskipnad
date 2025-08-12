package main

import (
	"log"
	"net/http"
	"os"

	"samskipnad/internal/auth"
	"samskipnad/internal/database"
	"samskipnad/internal/handlers"
	"samskipnad/internal/middleware"
	"samskipnad/internal/payments"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := database.Init()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize services
	authService := auth.NewService(db)
	paymentService := payments.NewService(db)
	handlers := handlers.New(db, authService, paymentService)

	// Set up routes
	r := mux.NewRouter()
	
	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Public routes
	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.Register).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthRequired(authService))
	
	protected.HandleFunc("/dashboard", handlers.Dashboard).Methods("GET")
	protected.HandleFunc("/classes", handlers.Classes).Methods("GET")
	protected.HandleFunc("/calendar", handlers.Calendar).Methods("GET")
	protected.HandleFunc("/classes/{id}/book", handlers.BookClass).Methods("POST")
	protected.HandleFunc("/memberships", handlers.Memberships).Methods("GET")
	protected.HandleFunc("/profile", handlers.Profile).Methods("GET", "POST")
	
	// Payment routes
	protected.HandleFunc("/payment/membership", handlers.MembershipPayment).Methods("GET", "POST")
	protected.HandleFunc("/payment/success", handlers.PaymentSuccess).Methods("GET")

	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminRequired())
	
	admin.HandleFunc("/", handlers.AdminDashboard).Methods("GET")
	admin.HandleFunc("/classes", handlers.AdminClasses).Methods("GET", "POST")
	admin.HandleFunc("/classes/{id}/edit", handlers.EditClass).Methods("GET", "POST")
	admin.HandleFunc("/classes/{id}/delete", handlers.DeleteClass).Methods("POST")
	admin.HandleFunc("/users", handlers.AdminUsers).Methods("GET")
	admin.HandleFunc("/roles", handlers.AdminRoles).Methods("GET", "POST")
	admin.HandleFunc("/payments", handlers.AdminPayments).Methods("GET")

	// API routes for HTMX
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthRequired(authService))
	api.HandleFunc("/classes/search", handlers.SearchClasses).Methods("GET")
	api.HandleFunc("/bookings/cancel/{id}", handlers.CancelBooking).Methods("DELETE")
	api.HandleFunc("/calendar/day/{date}", handlers.CalendarDayDetails).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}