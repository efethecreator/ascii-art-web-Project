package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	http.HandleFunc("/background-color", backgroundColorHandler) // New handler
	http.HandleFunc("/text-color", textColorHandler)             // New handler
	log.Println("Server started at :8080")
	log.Println("Server started at :8065")
	http.ListenAndServe(":8080", nil)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	asciiArt := generateAsciiArt(text, banner)

	w.Write([]byte(asciiArt))
}

func generateAsciiArt(text, banner string) string {
	bannerFileName := banner + ".txt"
	groups := readFile(bannerFileName)

	characters := make(map[int]string)
	index := 32

	for i := 0; i < len(groups); i++ {
		characters[index] = groups[i]
		index++
	}

	words := strings.Split(text, "\n")

	var asciiArt strings.Builder

	for _, word := range words {
		if word == "" {
			asciiArt.WriteString("\n\n")
		} else {
			for line := 0; line < 9; line++ {
				for i := 0; i < len(word); i++ {
					ascii := characters[int(word[i])]
					asciiLines := strings.Split(ascii, "\n")
					if line < len(asciiLines) {
						asciiArt.WriteString(asciiLines[line])
					}
				}
				if line != 8 {
					asciiArt.WriteString("\n")
				}
			}
			asciiArt.WriteString("\n\n")
		}
	}

	return asciiArt.String()
}

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var lines []string
	var group []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		group = append(group, scanner.Text())
		if len(group) == 9 {
			lines = append(lines, strings.Join(group, "\n"))
			group = nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file: %v", err)
	}

	return lines
}

func backgroundColorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	color := r.FormValue("color")
	w.Write([]byte(color)) // Respond with the selected color
	// You can add logic here to do further processing based on the selected color
}

func textColorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	color := r.FormValue("color")
	w.Write([]byte(color)) // Respond with the selected color
	// You can add logic here to do further processing based on the selected color
}

func man() {
	http.HandleFunc("/", handler)            // "/" yolu için işleyici belirtin
	http.HandleFunc("/404", notFoundHandler) // "/404" yolu için 404 işleyici belirtin
	http.ListenAndServe(":8080", nil)        // 8080 portunda HTTP sunucusunu başlatın
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Merhaba, Dünya!") // Örnek bir sayfa
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Sayfa Bulunamadı")
}
