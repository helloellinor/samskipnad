package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"samskipnad/internal/auth"
	"samskipnad/internal/config"
	"samskipnad/internal/database"
	"samskipnad/internal/handlers"
	"samskipnad/internal/middleware"
	"samskipnad/internal/payments"

	"github.com/gorilla/mux"
)

func main() {
	// Load community configuration
	communityName := os.Getenv("COMMUNITY")
	
	// Initialize hot-reload if enabled
	hotReloadEnabled := os.Getenv("HOT_RELOAD_ENABLED") == "true"
	if hotReloadEnabled {
		err := config.InitializeHotReload("config")
		if err != nil {
			log.Printf("Warning: Failed to initialize hot-reload: %v", err)
			log.Println("Falling back to static configuration loading")
		} else {
			log.Println("Hot-reload configuration system initialized")
			
			// Set up reload callback to log configuration changes
			config.SetGlobalReloadCallback(func(name string, cfg *config.Community) {
				log.Printf("🔥 Hot-reload: Configuration '%s' updated - %s", name, cfg.Name)
			})
		}
	}
	
	var community *config.Community
	var err error
	
	if hotReloadEnabled {
		community, err = config.LoadWithHotReload(communityName)
	} else {
		community, err = config.Load(communityName)
	}
	
	if err != nil {
		log.Fatal("Failed to load community config:", err)
	}
	log.Printf("Loaded community config: %s", community.Name)

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

	// Dynamic CSS based on community config
	r.HandleFunc("/css/community.css", handlers.DynamicCSS).Methods("GET")

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
	protected.HandleFunc("/klippekort", handlers.Klippekort).Methods("GET")
	protected.HandleFunc("/klippekort/purchase", handlers.KlippekortPurchase).Methods("POST")
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

	// Klippekort HTMX API routes
	api.HandleFunc("/klippekort/balance", handlers.KlippekortBalance).Methods("GET")
	api.HandleFunc("/klippekort/category", handlers.KlippekortCategory).Methods("GET")
	api.HandleFunc("/klippekort/purchase", handlers.KlippekortPurchaseInstant).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")
	
	// Cleanup hot-reload system if enabled
	if hotReloadEnabled {
		if err := config.ShutdownHotReload(); err != nil {
			log.Printf("Error shutting down hot-reload: %v", err)
		} else {
			log.Println("Hot-reload system shut down gracefully")
		}
	}
	
	log.Println("Server stopped")
}
