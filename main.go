package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// HandlerData is the struct passed to the template.
type HandlerData struct {
	Output string
	Error  string
}

// FormTemplate is the inline HTML template.
const formTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Input Validation Demo</title>
</head>
<body>
    <h1>Input Validation Demo</h1>
    
    <h2>Insecure Input</h2>
    <form method="POST" action="/insecure/input">
        <input type="text" name="input" placeholder="Enter input">
        <button type="submit">Submit Insecure</button>
    </form>

    <h2>Secure Input</h2>
    <form method="POST" action="/secure/input">
        <input type="text" name="input" placeholder="Enter input">
        <button type="submit">Submit Secure</button>
    </form>

    {{if .Error}}
        <div style="color: red;">
            <strong>Error:</strong>
            <p>{{.Error}}</p>
        </div>
    {{end}}

    {{if .Output}}
        <div>
            <strong>Output:</strong>
            <p>{{.Output}}</p>
        </div>
    {{end}}
</body>
</html>
`

// InsecureInputValidate performs minimal validation (insecure).
func InsecureInputValidate(input string) (string, error) {
	if input == "" {
		return "", logError("input cannot be empty")
	}
	return input, nil
}

// SecureInputValidate performs proper validation and sanitization.
func SecureInputValidate(input string) (string, error) {
	if input == "" {
		return "", logError("input cannot be empty")
	}
	if len(input) > 255 {
		return "", logError("input exceeds maximum length of 255 characters")
	}
	// Simple alphanumeric check
	for _, char := range input {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return "", logError("input must be alphanumeric")
		}
	}
	return input, nil
}

// logError creates an error and logs it.
func logError(msg string) error {
	err := errors.New(msg)
	// Log the error message
	log.Println("Error:", err)
	return err
}

// InsecureInputHandler handles the insecure input form.
func InsecureInputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	input := r.FormValue("input")
	output, err := InsecureInputValidate(input)
	data := HandlerData{}
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Output = output
	}
	// tmpl := template.Must(template.New("form").Parse(formTemplate))
	// tmpl.Execute(w, data)
	tmpl := template.Must(template.New("form").Parse(formTemplate))
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// SecureInputHandler handles the secure input form.
func SecureInputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	input := r.FormValue("input")
	output, err := SecureInputValidate(input)
	data := HandlerData{}
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Output = output
	}
	tmpl := template.Must(template.New("form").Parse(formTemplate))
	tmpl.Execute(w, data)
}

func hellohandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	mux := http.NewServeMux()

	// Parse the template once at startup
	tmpl := template.Must(template.New("form").Parse(formTemplate))

	// Root handler to serve the form
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Register handlers for insecure and secure input
	mux.HandleFunc("/insecure/input", InsecureInputHandler)
	mux.HandleFunc("/secure/input", SecureInputHandler)
	mux.HandleFunc("/hello", hellohandler)

	// Start the server
	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
