package main

import (
	"html/template"
	"net/http"
)

// PageData is a struct to hold the data for your template.
type PageData struct {
	Title   string
	Heading string
	Content string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Define your HTML template.
		const htmlTemplate = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Heading}}</h1>
			<p>{{.Content}}</p>
		</body>
		</html>
		`

		// Parse the template.
		tmpl, err := template.New("mytemplate").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a data instance.
		data := PageData{
			Title:   "My Page",
			Heading: "Welcome to My Page",
			Content: "This is the content of my page.",
		}

		// Execute the template and write the result to the HTTP response.
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the web server.
	http.ListenAndServe(":8888", nil)
}
