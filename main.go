package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const GEMINI_MODEL = "gemini-2.5-flash"

func main() {

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load API Key dari environment variable
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY tidak ditemukan di environment variable")
	}

	// Inisialisasi Gemini AI client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Gagal inisialisasi Gemini client: %v", err)
	}
	defer client.Close()

	// Setup router Chi
	r := chi.NewRouter()

	// Middleware CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // 5 menit
	}))

	// Endpoint test
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server ready, Gemini AI integrated 🚀"))
	})

	// Endpoint POST /api/chat
	r.Post("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var requestBody struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request format: %v", err), http.StatusBadRequest)
			return
		}

		// Check if messages array exists and is not empty
		if len(requestBody.Messages) == 0 {
			http.Error(w, "No messages provided", http.StatusBadRequest)
			return
		}

		// Prepare chat history
		ctx := context.Background()
		model := client.GenerativeModel(GEMINI_MODEL)

		// Convert messages to Gemini format
		var parts []genai.Part
		for _, msg := range requestBody.Messages {
			// Use the role in the content if it's not empty
			content := msg.Content
			if msg.Role != "" {
				content = fmt.Sprintf("%s: %s", msg.Role, msg.Content)
			}
			parts = append(parts, genai.Text(content))
		}

		// Generate response
		resp, err := model.GenerateContent(ctx, parts...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating response: %v", err), http.StatusInternalServerError)
			return
		}

		// Extract and return the response
		response := ""
		for _, cand := range resp.Candidates {
			if len(cand.Content.Parts) > 0 {
				response = fmt.Sprintf("%s", cand.Content.Parts[0])
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"response": response,
		})
	})

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server ready on http://localhost:%s", port)
	http.ListenAndServe(":"+port, r)
}

func extractText(resp interface{}) (string, error) {
	var text string

	// Attempt to parse the response
	data, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error extracting text: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("error extracting text: %v", err)
	}

	// Safely extract text from response structure
	if candidates, _ := result["candidates"].([]interface{}); len(candidates) > 0 {
		if candidate, _ := candidates[0].(map[string]interface{}); candidate != nil {
			if content, _ := candidate["content"].(map[string]interface{}); content != nil {
				if parts, _ := content["parts"].([]interface{}); len(parts) > 0 {
					if part, _ := parts[0].(map[string]interface{}); part != nil {
						text, _ = part["text"].(string)
					}
				}
			}
		}
	}

	// Fallback to null coalescing
	if text == "" {
		text = "JSON.stringify(resp, null, 2)"
	}

	return text, nil
}
