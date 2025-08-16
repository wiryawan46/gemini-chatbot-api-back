# Gemini Chatbot API

A simple HTTP server that provides an API for interacting with Google's Gemini AI model.

## Prerequisites

- Go 1.16 or higher
- A Google Gemini API key

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/chatbot-gemini.git
   cd chatbot-gemini
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the project root and add your Gemini API key:
   ```env
   API_KEY=your_gemini_api_key_here
   PORT=3000  # Optional: Change the default port if needed
   ```

## Running the Server

1. Start the server:
   ```bash
   go run main.go
   ```

2. The server will start on `http://localhost:3000` by default.

## API Endpoints

### POST /api/chat

Send messages to the Gemini AI model.

**Request:**
```json
{
  "messages": [
    {
      "role": "user",
      "content": "Hello, who are you?"
    }
  ]
}
```

**Response:**
```json
{
  "response": "I am an AI assistant created by Google."
}
```

## Environment Variables

- `API_KEY`: Your Google Gemini API key (required)
- `PORT`: Port to run the server on (default: 3000)

## Testing

You can test the API using `curl`:

```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Hello, who are you?"
      }
    ]
  }'
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
