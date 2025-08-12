package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to database
	db, err := sql.Open("sqlite3", "./samskipnad.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Add test classes with Norwegian pricing
	classes := []struct {
		name        string
		description string
		startTime   time.Time
		endTime     time.Time
		price       int // in øre (NOK cents)
		maxCapacity int
	}{
		{
			name:        "Morning Yoga Flow",
			description: "Gentle yoga flow to start your day. Perfect for beginners and experienced yogis alike.",
			startTime:   time.Now().Add(24 * time.Hour).Truncate(time.Hour).Add(9 * time.Hour),  // Tomorrow 9 AM
			endTime:     time.Now().Add(24 * time.Hour).Truncate(time.Hour).Add(10 * time.Hour), // Tomorrow 10 AM
			price:       0,                                                                      // Free class
			maxCapacity: 15,
		},
		{
			name:        "Power HIIT",
			description: "High-intensity interval training for busy parents. Burn calories fast!",
			startTime:   time.Now().Add(24 * time.Hour).Truncate(time.Hour).Add(18 * time.Hour), // Tomorrow 6 PM
			endTime:     time.Now().Add(24 * time.Hour).Truncate(time.Hour).Add(19 * time.Hour), // Tomorrow 7 PM
			price:       5900,                                                                   // 59 NOK
			maxCapacity: 12,
		},
		{
			name:        "Mindful Meditation",
			description: "Find your center in chaos. Perfect for overwhelmed parents who need 5 minutes of peace.",
			startTime:   time.Now().Add(48 * time.Hour).Truncate(time.Hour).Add(12 * time.Hour),                // Day after tomorrow 12 PM
			endTime:     time.Now().Add(48 * time.Hour).Truncate(time.Hour).Add(12*time.Hour + 30*time.Minute), // 12:30 PM
			price:       0,                                                                                     // Free class
			maxCapacity: 20,
		},
		{
			name:        "Strength & Core",
			description: "Build functional strength while the kids are at school. No judgment zone.",
			startTime:   time.Now().Add(48 * time.Hour).Truncate(time.Hour).Add(10 * time.Hour), // Day after tomorrow 10 AM
			endTime:     time.Now().Add(48 * time.Hour).Truncate(time.Hour).Add(11 * time.Hour), // 11 AM
			price:       8900,                                                                   // 89 NOK
			maxCapacity: 10,
		},
		{
			name:        "Evening Stretch & Relax",
			description: "Unwind after a chaotic day. Gentle stretches and relaxation techniques.",
			startTime:   time.Now().Add(72 * time.Hour).Truncate(time.Hour).Add(19 * time.Hour), // 3 days from now 7 PM
			endTime:     time.Now().Add(72 * time.Hour).Truncate(time.Hour).Add(20 * time.Hour), // 8 PM
			price:       6900,                                                                   // 69 NOK
			maxCapacity: 15,
		},
		{
			name:        "Weekend Warrior Bootcamp",
			description: "Intense weekend workout for when you actually have childcare. Bring your A-game.",
			startTime:   time.Now().AddDate(0, 0, 7-int(time.Now().Weekday())+6).Truncate(time.Hour).Add(8 * time.Hour),                // Next Saturday 8 AM
			endTime:     time.Now().AddDate(0, 0, 7-int(time.Now().Weekday())+6).Truncate(time.Hour).Add(9*time.Hour + 30*time.Minute), // 9:30 AM
			price:       12900,                                                                                                         // 129 NOK
			maxCapacity: 8,
		},
	}

	for _, class := range classes {
		_, err := db.Exec(`
			INSERT INTO classes (tenant_id, name, description, instructor_id, start_time, end_time, max_capacity, price, requires_ticket, requires_membership, active, created_at, updated_at)
			VALUES (1, ?, ?, 1, ?, ?, ?, ?, false, false, true, datetime('now'), datetime('now'))
		`, class.name, class.description, class.startTime, class.endTime, class.maxCapacity, class.price)

		if err != nil {
			log.Printf("Failed to insert class %s: %v", class.name, err)
		} else {
			log.Printf("Added class: %s", class.name)
		}
	}

	log.Println("Test data added successfully!")
}
