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
	"samskipnad/internal/middleware"
	"samskipnad/internal/models"

	"github.com/gorilla/mux"
)

type Handlers struct {
	db          *sql.DB
	authService *auth.Service
	templates   *template.Template
}

func New(db *sql.DB, authService *auth.Service) *Handlers {
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
	}
	
	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/*.html"))

	return &Handlers{
		db:          db,
		authService: authService,
		templates:   templates,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if h.authService.IsAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title": "Samskipnad - Yoga Community Platform",
	}

	h.renderTemplate(w, "home-standalone.html", data)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := map[string]interface{}{
			"Title": "Login",
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
			"Title": "Login",
			"Error": "Invalid email or password",
			"Email": email,
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

	data := map[string]interface{}{
		"Title": "Dashboard",
		"User":  user,
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

	data := map[string]interface{}{
		"Title":   "Classes",
		"User":    user,
		"Classes": classes,
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

	// Check if class exists and has capacity
	var class models.Class
	err = h.db.QueryRow(`
		SELECT id, name, start_time, max_capacity, price, requires_ticket, requires_membership
		FROM classes WHERE id = ? AND active = true`, classID).Scan(
		&class.ID, &class.Name, &class.StartTime, &class.MaxCapacity,
		&class.Price, &class.RequiresTicket, &class.RequiresMembership)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	// Check current bookings
	var currentBookings int
	err = h.db.QueryRow("SELECT COUNT(*) FROM bookings WHERE class_id = ? AND status = 'confirmed'", classID).Scan(&currentBookings)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if currentBookings >= class.MaxCapacity {
		http.Error(w, "Class is full", http.StatusBadRequest)
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
		http.Error(w, "Already booked", http.StatusBadRequest)
		return
	}

	// Create booking
	_, err = h.db.Exec(`
		INSERT INTO bookings (user_id, class_id, status, created_at, updated_at)
		VALUES (?, ?, 'confirmed', ?, ?)`,
		user.ID, classID, time.Now(), time.Now())
	if err != nil {
		http.Error(w, "Failed to book class", http.StatusInternalServerError)
		return
	}

	// Return success for HTMX
	w.Header().Set("HX-Trigger", "booking-updated")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<div class="alert alert-success">Successfully booked %s!</div>`, class.Name)
}

func (h *Handlers) Memberships(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title": "Memberships",
		"User":  user,
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