package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"samskipnad/internal/auth"
	"samskipnad/internal/config"
	"samskipnad/internal/middleware"
	"samskipnad/internal/models"
	"samskipnad/internal/payments"

	"github.com/gorilla/mux"
)

type Handlers struct {
	db             *sql.DB
	authService    *auth.Service
	paymentService *payments.Service
	templates      *template.Template
}

func New(db *sql.DB, authService *auth.Service, paymentService *payments.Service) *Handlers {
	// Load templates with custom functions
	funcMap := template.FuncMap{
		"divf": func(a, b float64) float64 {
			if b != 0 {
				return a / b
			}
			return 0
		},
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return string(s[0]-32) + s[1:]
		},
		"substr": func(s string, start, length int) string {
			if start < 0 || start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return strings.ToUpper(s[start:end])
		},
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
	}
	
	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/*.html"))

	return &Handlers{
		db:             db,
		authService:    authService,
		paymentService: paymentService,
		templates:      templates,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if h.authService.IsAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	community := config.GetCurrent()
	data := map[string]interface{}{
		"Title":     community.Content.Home.Title,
		"Community": community,
	}

	h.renderTemplate(w, "home-standalone.html", data)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	community := config.GetCurrent()
	
	if r.Method == "GET" {
		data := map[string]interface{}{
			"Title":     "Login",
			"Community": community,
		}
		h.renderTemplate(w, "login-standalone.html", data)
		return
	}

	// POST - handle login
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.authService.Login(email, password)
	if err != nil {
		data := map[string]interface{}{
			"Title":     "Login",
			"Error":     "Invalid email or password",
			"Email":     email,
			"Community": community,
		}
		h.renderTemplate(w, "login-standalone.html", data)
		return
	}

	if err := h.authService.CreateSession(w, r, user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := map[string]interface{}{
			"Title": "Register",
		}
		h.renderTemplate(w, "register-standalone.html", data)
		return
	}

	// POST - handle registration
	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	// For now, all users join the default tenant (1)
	user, err := h.authService.Register(email, password, firstName, lastName, 1)
	if err != nil {
		data := map[string]interface{}{
			"Title":     "Register",
			"Error":     err.Error(),
			"Email":     email,
			"FirstName": firstName,
			"LastName":  lastName,
		}
		h.renderTemplate(w, "register-standalone.html", data)
		return
	}

	if err := h.authService.CreateSession(w, r, user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.authService.DestroySession(w, r); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get upcoming classes
	upcomingClasses, err := h.getUpcomingClasses(user.TenantID)
	if err != nil {
		upcomingClasses = []models.Class{} // Empty slice on error
	}

	// Get user bookings
	userBookings, err := h.getUserBookings(user.ID)
	if err != nil {
		userBookings = []models.Booking{} // Empty slice on error
	}

	community := config.GetCurrent()
	data := map[string]interface{}{
		"Title":           "Dashboard",
		"User":            user,
		"Community":       community,
		"UpcomingClasses": upcomingClasses,
		"UserBookings":    userBookings,
		"UserMembership":  nil, // TODO: Implement membership checking
	}

	h.renderTemplate(w, "dashboard-standalone.html", data)
}

func (h *Handlers) Classes(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	classes, err := h.getUpcomingClasses(user.TenantID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	community := config.GetCurrent()
	data := map[string]interface{}{
		"Title":     "Classes",
		"User":      user,
		"Community": community,
		"Classes":   classes,
	}

	h.renderTemplate(w, "classes.html", data)
}

func (h *Handlers) BookClass(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	classID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Get class details
	class, err := h.getClassByID(classID)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	// Check capacity
	var currentBookings int
	err = h.db.QueryRow("SELECT COUNT(*) FROM bookings WHERE class_id = ? AND status = 'confirmed'", classID).Scan(&currentBookings)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if currentBookings >= class.MaxCapacity {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="booking-error">Class is full</div>`))
		return
	}

	// Check if user already booked
	var existingBooking int
	err = h.db.QueryRow("SELECT COUNT(*) FROM bookings WHERE user_id = ? AND class_id = ?", user.ID, classID).Scan(&existingBooking)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if existingBooking > 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="booking-error">Already booked</div>`))
		return
	}

	// If class is free, book directly
	if class.Price == 0 {
		err = h.bookClassDirectly(user.ID, classID)
		if err != nil {
			http.Error(w, "Failed to book class", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Trigger", "booking-updated")
		w.Write([]byte(`<div class="booking-success">Successfully booked ` + class.Name + `!</div>`))
		return
	}

	// For paid classes, create payment intent
	paymentIntent, err := h.paymentService.CreatePaymentIntent(user.ID, classID, int64(class.Price))
	if err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Return payment form
	data := struct {
		PaymentIntent *struct {
			ID           string
			ClientSecret string
		}
		Class *models.Class
		User  *models.User
	}{
		PaymentIntent: &struct {
			ID           string
			ClientSecret string
		}{
			ID:           paymentIntent.ID,
			ClientSecret: paymentIntent.ClientSecret,
		},
		Class: class,
		User:  user,
	}

	// Render payment form as HTML fragment
	tmpl := `
	<div class="payment-form p-3">
		<h6 class="text-white">Payment Required</h6>
		<div class="payment-summary mb-3">
			<div class="d-flex justify-content-between">
				<span>Class:</span>
				<span class="text-white">{{.Class.Name}}</span>
			</div>
			<div class="d-flex justify-content-between">
				<span>Price:</span>
				<span class="payment-total">${{printf "%.2f" (divf (float64 .Class.Price) 100.0)}}</span>
			</div>
		</div>
		<div id="payment-element-{{.Class.ID}}" class="mb-3"></div>
		<button id="submit-payment-{{.Class.ID}}" class="btn btn-primary w-100">Pay Now</button>
		<div id="payment-messages-{{.Class.ID}}" class="mt-2"></div>
	</div>
	<script src="https://js.stripe.com/v3/"></script>
	<script>
		(function() {
			// Replace with your actual publishable key
			const stripe = Stripe('pk_test_51234567890abcdef'); 
			const elements = stripe.elements({
				clientSecret: '{{.PaymentIntent.ClientSecret}}',
				appearance: {
					theme: 'night',
					variables: {
						colorPrimary: '#ff6b35',
						colorBackground: '#2a2a2a',
						colorText: '#f0f0f0',
						borderRadius: '0px'
					}
				}
			});
			
			const paymentElement = elements.create('payment');
			paymentElement.mount('#payment-element-{{.Class.ID}}');
			
			document.getElementById('submit-payment-{{.Class.ID}}').addEventListener('click', async (e) => {
				e.preventDefault();
				e.target.disabled = true;
				e.target.textContent = 'Processing...';
				
				const {error} = await stripe.confirmPayment({
					elements,
					confirmParams: {
						return_url: window.location.origin + '/payment/success?class_id={{.Class.ID}}'
					}
				});
				
				if (error) {
					document.getElementById('payment-messages-{{.Class.ID}}').innerHTML = 
						'<div class="booking-error">' + error.message + '</div>';
					e.target.disabled = false;
					e.target.textContent = 'Pay Now';
				}
			});
		})();
	</script>`

	t, err := template.New("payment-form").Funcs(template.FuncMap{
		"divf": func(a, b float64) float64 {
			if b != 0 {
				return a / b
			}
			return 0
		},
		"float64": func(i int) float64 {
			return float64(i)
		},
	}).Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) Memberships(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	community := config.GetCurrent()
	data := map[string]interface{}{
		"Title":     "Memberships",
		"User":      user,
		"Community": community,
	}

	h.renderTemplate(w, "memberships.html", data)
}

func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		data := map[string]interface{}{
			"Title": "Profile",
			"User":  user,
		}
		h.renderTemplate(w, "profile.html", data)
		return
	}

	// POST - update profile
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	phone := r.FormValue("phone")

	_, err := h.db.Exec(`
		UPDATE users SET first_name = ?, last_name = ?, phone = ?, updated_at = ?
		WHERE id = ?`, firstName, lastName, phone, time.Now(), user.ID)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Admin handlers
func (h *Handlers) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := map[string]interface{}{
		"Title": "Admin Dashboard",
		"User":  user,
	}

	h.renderTemplate(w, "admin-dashboard-standalone.html", data)
}

func (h *Handlers) AdminClasses(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		classes, err := h.getAllClasses(user.TenantID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Title":   "Manage Classes",
			"User":    user,
			"Classes": classes,
		}

		h.renderTemplate(w, "admin-classes.html", data)
		return
	}

	// POST - create new class
	name := r.FormValue("name")
	description := r.FormValue("description")
	startTime := r.FormValue("start_time")
	endTime := r.FormValue("end_time")
	maxCapacity, _ := strconv.Atoi(r.FormValue("max_capacity"))
	price, _ := strconv.Atoi(r.FormValue("price"))

	startTimeParsed, err := time.Parse("2006-01-02T15:04", startTime)
	if err != nil {
		http.Error(w, "Invalid start time", http.StatusBadRequest)
		return
	}

	endTimeParsed, err := time.Parse("2006-01-02T15:04", endTime)
	if err != nil {
		http.Error(w, "Invalid end time", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec(`
		INSERT INTO classes (tenant_id, name, description, instructor_id, start_time, end_time, max_capacity, price, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.TenantID, name, description, user.ID, startTimeParsed, endTimeParsed, maxCapacity, price, time.Now(), time.Now())
	if err != nil {
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/classes", http.StatusSeeOther)
}

func (h *Handlers) EditClass(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement edit class functionality
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handlers) DeleteClass(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement delete class functionality
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handlers) AdminUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user management
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handlers) AdminRoles(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement role management
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handlers) AdminPayments(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement payment management
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handlers) Calendar(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get month and year from query params, default to current month
	now := time.Now()
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	
	year := now.Year()
	month := now.Month()
	
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}
	
	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = time.Month(m)
		}
	}

	// Get classes for the month
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	
	classes, err := h.getClassesForDateRange(user.TenantID, startOfMonth, endOfMonth)
	if err != nil {
		http.Error(w, "Failed to get classes", http.StatusInternalServerError)
		return
	}

	// Build calendar data
	calendarData := h.buildCalendarData(year, month, classes)

	data := struct {
		Title        string
		User         *models.User
		CalendarData CalendarData
		CurrentMonth time.Time
		PrevMonth    time.Time
		NextMonth    time.Time
	}{
		Title:        "Calendar",
		User:         user,
		CalendarData: calendarData,
		CurrentMonth: time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
		PrevMonth:    time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0),
		NextMonth:    time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0),
	}

	h.renderTemplate(w, "calendar.html", data)
}

type CalendarData struct {
	Year     int
	Month    time.Month
	Days     []CalendarDay
	Classes  []models.Class
}

type CalendarDay struct {
	Day        int
	Date       time.Time
	IsToday    bool
	IsOtherMonth bool
	Classes    []models.Class
}

func (h *Handlers) buildCalendarData(year int, month time.Month, classes []models.Class) CalendarData {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	today := time.Now()
	
	// Start from Monday of the week containing the first day
	startDate := firstDay
	for startDate.Weekday() != time.Monday {
		startDate = startDate.AddDate(0, 0, -1)
	}
	
	// End on Sunday of the week containing the last day
	endDate := lastDay
	for endDate.Weekday() != time.Sunday {
		endDate = endDate.AddDate(0, 0, 1)
	}
	
	var days []CalendarDay
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		var dayClasses []models.Class
		for _, class := range classes {
			if class.StartTime.Day() == d.Day() && class.StartTime.Month() == d.Month() {
				dayClasses = append(dayClasses, class)
			}
		}
		
		days = append(days, CalendarDay{
			Day:          d.Day(),
			Date:         d,
			IsToday:      d.Year() == today.Year() && d.Month() == today.Month() && d.Day() == today.Day(),
			IsOtherMonth: d.Month() != month,
			Classes:      dayClasses,
		})
	}
	
	return CalendarData{
		Year:    year,
		Month:   month,
		Days:    days,
		Classes: classes,
	}
}

func (h *Handlers) getClassesForDateRange(tenantID int, start, end time.Time) ([]models.Class, error) {
	rows, err := h.db.Query(`
		SELECT id, name, description, instructor_id, start_time, end_time, max_capacity, price
		FROM classes
		WHERE tenant_id = ? AND start_time >= ? AND start_time <= ? AND active = true
		ORDER BY start_time ASC`, tenantID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		err := rows.Scan(&class.ID, &class.Name, &class.Description, 
			&class.InstructorID, &class.StartTime, &class.EndTime, 
			&class.MaxCapacity, &class.Price)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (h *Handlers) CalendarDayDetails(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	dateStr := vars["date"]
	
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Get classes for the specific day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	classes, err := h.getClassesForDateRange(user.TenantID, startOfDay, endOfDay)
	if err != nil {
		http.Error(w, "Failed to get classes", http.StatusInternalServerError)
		return
	}

	// Get user's bookings for these classes
	userBookings, err := h.getUserBookings(user.ID)
	if err != nil {
		userBookings = []models.Booking{} // Continue even if we can't get bookings
	}

	// Create a map for quick lookup of booked classes
	bookedClasses := make(map[int]bool)
	for _, booking := range userBookings {
		bookedClasses[booking.ClassID] = true
	}

	data := struct {
		Date          time.Time
		Classes       []models.Class
		BookedClasses map[int]bool
	}{
		Date:          date,
		Classes:       classes,
		BookedClasses: bookedClasses,
	}

	// Render the day details as HTML fragment
	tmpl := `
	<h6>{{.Date.Format "Monday, January 2, 2006"}}</h6>
	{{if .Classes}}
		{{range .Classes}}
		<div class="card mb-3">
			<div class="card-body">
				<h6 class="card-title">{{.Name}}</h6>
				<p class="card-text text-muted">{{.Description}}</p>
				<div class="class-details">
					<div class="detail-item">
						<span class="detail-label">Time:</span>
						<span class="detail-value">{{.StartTime.Format "15:04"}} - {{.EndTime.Format "15:04"}}</span>
					</div>
					<div class="detail-item">
						<span class="detail-label">Duration:</span>
						<span class="detail-value">{{printf "%.0f" .EndTime.Sub(.StartTime).Minutes}} minutes</span>
					</div>
					<div class="detail-item">
						<span class="detail-label">Price:</span>
						<span class="detail-value">
							{{if gt .Price 0}}
								${{printf "%.2f" (divf (float64 .Price) 100.0)}}
							{{else}}
								Free
							{{end}}
						</span>
					</div>
					<div class="detail-item">
						<span class="detail-label">Capacity:</span>
						<span class="detail-value">{{.MaxCapacity}} spots</span>
					</div>
				</div>
				<div class="mt-3">
					{{if index $.BookedClasses .ID}}
					<span class="badge bg-success">Booked</span>
					{{else}}
					<button class="btn btn-primary btn-sm" 
							hx-post="/classes/{{.ID}}/book" 
							hx-target="#booking-result-{{.ID}}"
							hx-swap="innerHTML">
						Book Class
					</button>
					{{end}}
					<div id="booking-result-{{.ID}}" class="mt-2"></div>
				</div>
			</div>
		</div>
		{{end}}
	{{else}}
		<p class="text-muted">No classes scheduled for this day.</p>
	{{end}}`

	t, err := template.New("day-details").Funcs(template.FuncMap{
		"divf": func(a, b float64) float64 {
			if b != 0 {
				return a / b
			}
			return 0
		},
		"float64": func(i int) float64 {
			return float64(i)
		},
	}).Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

// Payment handlers

func (h *Handlers) CreateClassPayment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	classIDStr := vars["id"]
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Get class details to determine price
	class, err := h.getClassByID(classID)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	if class.Price == 0 {
		// Free class, book directly
		err = h.bookClassDirectly(user.ID, classID)
		if err != nil {
			http.Error(w, "Failed to book class", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="booking-success">Class booked successfully!</div>`))
		return
	}

	// Create payment intent for paid class
	paymentIntent, err := h.paymentService.CreatePaymentIntent(user.ID, classID, int64(class.Price))
	if err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Return payment form
	data := struct {
		PaymentIntent *struct {
			ID           string
			ClientSecret string
		}
		Class *models.Class
		User  *models.User
	}{
		PaymentIntent: &struct {
			ID           string
			ClientSecret string
		}{
			ID:           paymentIntent.ID,
			ClientSecret: paymentIntent.ClientSecret,
		},
		Class: class,
		User:  user,
	}

	// Render payment form as HTML fragment
	tmpl := `
	<div class="payment-form">
		<h6>Payment Required</h6>
		<p>Class: {{.Class.Name}}</p>
		<p>Price: ${{printf "%.2f" (divf (float64 .Class.Price) 100.0)}}</p>
		<div id="payment-element"></div>
		<button id="submit-payment" class="btn btn-primary mt-3">Pay Now</button>
		<div id="payment-messages" class="mt-2"></div>
	</div>
	<script src="https://js.stripe.com/v3/"></script>
	<script>
		const stripe = Stripe('pk_test_...'); // Replace with your publishable key
		const elements = stripe.elements({
			clientSecret: '{{.PaymentIntent.ClientSecret}}'
		});
		
		const paymentElement = elements.create('payment');
		paymentElement.mount('#payment-element');
		
		document.getElementById('submit-payment').addEventListener('click', async () => {
			const {error} = await stripe.confirmPayment({
				elements,
				confirmParams: {
					return_url: window.location.origin + '/payment/success?class_id={{.Class.ID}}'
				}
			});
			
			if (error) {
				document.getElementById('payment-messages').innerHTML = 
					'<div class="booking-error">' + error.message + '</div>';
			}
		});
	</script>`

	t, err := template.New("payment-form").Funcs(template.FuncMap{
		"divf": func(a, b float64) float64 {
			if b != 0 {
				return a / b
			}
			return 0
		},
		"float64": func(i int) float64 {
			return float64(i)
		},
	}).Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) PaymentSuccess(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	paymentIntentID := r.URL.Query().Get("payment_intent")
	if paymentIntentID == "" {
		http.Error(w, "Missing payment intent", http.StatusBadRequest)
		return
	}

	// Confirm the payment
	err := h.paymentService.ConfirmPayment(paymentIntentID)
	if err != nil {
		http.Error(w, "Payment confirmation failed", http.StatusInternalServerError)
		return
	}

	// Redirect to dashboard with success message
	http.Redirect(w, r, "/dashboard?payment=success", http.StatusSeeOther)
}

func (h *Handlers) MembershipPayment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		// Show membership options
		data := struct {
			Title string
			User  *models.User
		}{
			Title: "Purchase Membership",
			User:  user,
		}

		h.renderTemplate(w, "membership-payment.html", data)
		return
	}

	// POST - Process membership purchase
	membershipType := r.FormValue("type")
	if membershipType == "" {
		http.Error(w, "Missing membership type", http.StatusBadRequest)
		return
	}

	// Set prices based on type from community config
	community := config.GetCurrent()
	var amount int64
	switch membershipType {
	case "monthly":
		amount = int64(community.Pricing.Monthly * 100) // Convert to cents
	case "yearly":
		amount = int64(community.Pricing.Yearly * 100) // Convert to cents
	default:
		http.Error(w, "Invalid membership type", http.StatusBadRequest)
		return
	}

	paymentIntent, err := h.paymentService.CreateMembershipPaymentIntent(user.ID, membershipType, amount)
	if err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Redirect to payment page
	http.Redirect(w, r, fmt.Sprintf("/payment/membership?payment_intent=%s", paymentIntent.ID), http.StatusSeeOther)
}

// Helper methods for payments

func (h *Handlers) getClassByID(classID int) (*models.Class, error) {
	query := `SELECT id, name, description, instructor_id, start_time, end_time, max_capacity, price, tenant_id
			  FROM classes WHERE id = ? AND active = true`
	
	var class models.Class
	err := h.db.QueryRow(query, classID).Scan(
		&class.ID, &class.Name, &class.Description, &class.InstructorID,
		&class.StartTime, &class.EndTime, &class.MaxCapacity, &class.Price, &class.TenantID,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &class, nil
}

func (h *Handlers) bookClassDirectly(userID, classID int) error {
	query := `INSERT INTO bookings (user_id, class_id, status, created_at, updated_at)
			  VALUES (?, ?, 'confirmed', datetime('now'), datetime('now'))`
	
	_, err := h.db.Exec(query, userID, classID)
	return err
}

// API handlers for HTMX
func (h *Handlers) SearchClasses(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement class search
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handlers) CancelBooking(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement booking cancellation
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Helper methods
func (h *Handlers) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) getUpcomingClasses(tenantID int) ([]models.Class, error) {
	rows, err := h.db.Query(`
		SELECT id, name, description, instructor_id, start_time, end_time, max_capacity, price
		FROM classes
		WHERE tenant_id = ? AND start_time > ? AND active = true
		ORDER BY start_time ASC
		LIMIT 10`, tenantID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		err := rows.Scan(&class.ID, &class.Name, &class.Description, &class.InstructorID,
			&class.StartTime, &class.EndTime, &class.MaxCapacity, &class.Price)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (h *Handlers) getAllClasses(tenantID int) ([]models.Class, error) {
	rows, err := h.db.Query(`
		SELECT id, name, description, instructor_id, start_time, end_time, max_capacity, price
		FROM classes
		WHERE tenant_id = ? AND active = true
		ORDER BY start_time ASC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var class models.Class
		err := rows.Scan(&class.ID, &class.Name, &class.Description, &class.InstructorID,
			&class.StartTime, &class.EndTime, &class.MaxCapacity, &class.Price)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (h *Handlers) getUserBookings(userID int) ([]models.Booking, error) {
	rows, err := h.db.Query(`
		SELECT id, user_id, class_id, status, created_at
		FROM bookings
		WHERE user_id = ? AND status = 'confirmed'
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.ClassID, &booking.Status, &booking.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// DynamicCSS generates CSS based on community configuration
func (h *Handlers) DynamicCSS(w http.ResponseWriter, r *http.Request) {
	community := config.GetCurrent()
	
	w.Header().Set("Content-Type", "text/css")
	
	css := fmt.Sprintf(`/* Dynamic CSS for %s */

/* CSS Variables from Community Config */
:root {
    --primary-color: %s;
    --secondary-color: %s;
    --accent-color: %s;
    --success-color: %s;
    --warning-color: %s;
    --danger-color: %s;
    --background-color: %s;
    --surface-color: %s;
    --text-color: %s;
    --muted-color: %s;
    --font-primary: '%s', system-ui, -apple-system, sans-serif;
    --font-secondary: '%s', 'Monaco', 'Menlo', monospace;
    --font-size-base: %s;
}

/* Scandinavian-inspired base styles */
body {
    font-family: var(--font-primary);
    background-color: var(--background-color);
    color: var(--text-color);
    line-height: 1.6;
    font-size: var(--font-size-base);
    margin: 0;
    padding: 0;
}

/* Navigation */
.navbar {
    background-color: var(--surface-color) !important;
    border-bottom: 1px solid var(--primary-color);
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    padding: 1rem 0;
}

.navbar-brand {
    font-weight: 600;
    font-size: 1.25rem;
    color: var(--primary-color) !important;
    text-decoration: none;
}

.navbar-nav .nav-link {
    color: var(--text-color) !important;
    font-weight: 500;
    font-size: 0.9rem;
    transition: color 0.2s ease;
    margin: 0 0.5rem;
}

.navbar-nav .nav-link:hover {
    color: var(--accent-color) !important;
}

/* Cards */
.card {
    background-color: var(--surface-color);
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    transition: all 0.2s ease;
    position: relative;
    overflow: hidden;
}

.card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%%;
    height: 3px;
    background: linear-gradient(90deg, var(--accent-color), var(--secondary-color));
    transform: scaleX(0);
    transition: transform 0.2s ease;
}

.card:hover::before {
    transform: scaleX(1);
}

.card-header {
    background: var(--surface-color);
    border-bottom: 1px solid #e5e7eb;
    font-weight: 600;
    color: var(--primary-color);
    font-size: 0.9rem;
}

.card-body {
    color: var(--text-color);
}

.card-title {
    color: var(--primary-color);
    font-weight: 600;
    margin-bottom: 0.5rem;
}

/* Buttons */
.btn {
    border-radius: 6px;
    font-weight: 500;
    font-size: 0.9rem;
    transition: all 0.2s ease;
    border: 1px solid;
    padding: 0.625rem 1.25rem;
    text-decoration: none;
    display: inline-block;
}

.btn-primary {
    background-color: var(--primary-color);
    border-color: var(--primary-color);
    color: white;
}

.btn-primary:hover {
    background-color: var(--secondary-color);
    border-color: var(--secondary-color);
    color: white;
    transform: translateY(-1px);
}

.btn-outline-secondary {
    background-color: transparent;
    border-color: var(--muted-color);
    color: var(--muted-color);
}

.btn-outline-secondary:hover {
    background-color: var(--muted-color);
    color: white;
}

.btn-success {
    background-color: var(--success-color);
    border-color: var(--success-color);
    color: white;
}

.btn-warning {
    background-color: var(--warning-color);
    border-color: var(--warning-color);
    color: var(--primary-color);
}

.btn-danger {
    background-color: var(--danger-color);
    border-color: var(--danger-color);
    color: white;
}

/* Hero Section */
.hero-section {
    background: linear-gradient(135deg, var(--surface-color) 0%%, var(--background-color) 100%%);
    padding: 4rem 2rem;
    text-align: center;
    border-radius: 12px;
    margin: 2rem 0;
}

.hero-section h1 {
    color: var(--primary-color);
    font-weight: 700;
    font-size: 2.5rem;
    margin-bottom: 1rem;
}

.hero-section .lead {
    color: var(--muted-color);
    font-size: 1.1rem;
    max-width: 600px;
    margin: 0 auto 2rem;
    line-height: 1.6;
}

/* Feature boxes */
.feature-box {
    background-color: var(--surface-color);
    border-radius: 8px;
    padding: 2rem;
    text-align: center;
    transition: all 0.2s ease;
    border: 1px solid #e5e7eb;
    height: 100%%;
}

.feature-box:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 25px rgba(0,0,0,0.1);
}

.feature-box h3 {
    color: var(--primary-color);
    font-weight: 600;
    margin-bottom: 1rem;
    font-size: 1.25rem;
}

.feature-box p {
    color: var(--muted-color);
    line-height: 1.6;
    margin: 0;
}

/* Forms */
.form-control {
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    padding: 0.75rem;
    background-color: var(--surface-color);
    color: var(--text-color);
    transition: border-color 0.2s ease;
}

.form-control:focus {
    border-color: var(--accent-color);
    box-shadow: 0 0 0 2px rgba(208, 135, 112, 0.1);
    outline: none;
}

.form-label {
    color: var(--primary-color);
    font-weight: 500;
    margin-bottom: 0.5rem;
}

/* Profile Card */
.profile-card {
    background: linear-gradient(135deg, var(--surface-color) 0%%, var(--background-color) 100%%);
    border-radius: 12px;
    padding: 2rem;
    text-align: center;
    border: 1px solid #e5e7eb;
    transition: all 0.3s ease;
}

.profile-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.profile-avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%%;
    background: linear-gradient(135deg, var(--accent-color), var(--secondary-color));
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 1rem;
    font-size: 1.5rem;
    font-weight: 700;
    color: white;
    text-transform: uppercase;
}

.profile-stats {
    display: flex;
    justify-content: space-around;
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid #e5e7eb;
}

.profile-stat {
    text-align: center;
}

.profile-stat-number {
    display: block;
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--accent-color);
    line-height: 1;
}

.profile-stat-label {
    display: block;
    font-size: 0.8rem;
    color: var(--muted-color);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-top: 0.25rem;
}

/* Attribution */
.attribution {
    text-align: center;
    padding: 1rem;
    color: var(--muted-color);
    font-size: 0.8rem;
    border-top: 1px solid #e5e7eb;
    margin-top: 2rem;
}

.attribution a {
    color: var(--accent-color);
    text-decoration: none;
}

.attribution a:hover {
    text-decoration: underline;
}

/* Responsive */
@media (max-width: 768px) {
    .hero-section h1 {
        font-size: 2rem;
    }
    
    .feature-box {
        margin-bottom: 1rem;
    }
    
    .profile-stats {
        flex-direction: column;
        gap: 1rem;
    }
}`,
		community.Name,
		community.Colors.Primary,
		community.Colors.Secondary,
		community.Colors.Accent,
		community.Colors.Success,
		community.Colors.Warning,
		community.Colors.Danger,
		community.Colors.Background,
		community.Colors.Surface,
		community.Colors.Text,
		community.Colors.Muted,
		community.Fonts.Primary,
		community.Fonts.Secondary,
		community.Fonts.SizeBase,
	)
	
	fmt.Fprint(w, css)
}