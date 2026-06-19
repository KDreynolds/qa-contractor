package main

import (
	"log"
	"net/http"

	"qa-contractor/email"
	"qa-contractor/templates"
	"qa-contractor/templates/components"
)

func main() {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Pages
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		templates.Home().Render(r.Context(), w)
	})
	mux.HandleFunc("GET /services", func(w http.ResponseWriter, r *http.Request) {
		templates.Services().Render(r.Context(), w)
	})
	mux.HandleFunc("GET /portfolio", func(w http.ResponseWriter, r *http.Request) {
		templates.Portfolio().Render(r.Context(), w)
	})
	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		templates.About().Render(r.Context(), w)
	})
	mux.HandleFunc("GET /contact", func(w http.ResponseWriter, r *http.Request) {
		templates.Contact().Render(r.Context(), w)
	})

	// Contact form submission
	mux.HandleFunc("POST /contact", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		emailAddr := r.FormValue("email")
		message := r.FormValue("message")

		// Basic validation
		if name == "" || emailAddr == "" || message == "" {
			// Re-render form with errors (simplified: just return the form)
			components.ContactForm().Render(r.Context(), w)
			return
		}

		log.Printf("Contact form submission: name=%s email=%s message=%s", name, emailAddr, message)

		if err := email.SendContactEmail(name, emailAddr, message); err != nil {
			log.Printf("Failed to send email: %v", err)
			// Still show success to the user — we have their info logged
		}

		components.ContactSuccess().Render(r.Context(), w)
	})

	log.Println("Server starting on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
