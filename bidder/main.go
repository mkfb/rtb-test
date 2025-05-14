package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Define the OpenRTB structs. We only need a subset for this dummy response.
// It's good practice to define structs for request/response handling. Use
// descriptive names and follow Go naming conventions. These structs
// represent the structure of the JSON data we'll be working with.

// BidRequest represents the OpenRTB bid request.  We'll use this to
// unmarshal the incoming request and validate it.  For simplicity,
// we're only including the fields we need for validation.  In a
// real-world scenario, this struct would be much more comprehensive.
type BidRequest struct {
	ID     string    `json:"id"`
	Imp    []Imp     `json:"imp"`
	Device *Device   `json:"device"` // Device is a pointer because it might be nil
	// Source  *Source    `json:"source,omitempty"`
	// Regs    *Regs      `json:"regs,omitempty"`
}

// Imp represents an impression in the bid request.
type Imp struct {
	ID string `json:"id"`
	// Banner  *Banner  `json:"banner,omitempty"`
	// Video   *Video   `json:"video,omitempty"`
	// Audio   *Audio   `json:"audio,omitempty"`
	// Native  *Native  `json:"native,omitempty"`
	// Pmp     *Pmp     `json:"pmp,omitempty"`
	// Ext     any      `json:"ext,omitempty"`
}

// Device represents the device in the bid request.
type Device struct {
	Ua string `json:"ua"`
	// Ip  string `json:"ip"`
	// Geo *Geo   `json:"geo,omitempty"`
	// Os  string `json:"os"`
	// Model string `json:"model"`
	// Ext any    `json:"ext,omitempty"`
}

// BidResponse represents the OpenRTB bid response.
type BidResponse struct {
	ID      string    `json:"id"`
	SeatBid []SeatBid `json:"seatbid"`
	// NoBidReason int       `json:"nobidreason,omitempty"` // omitempty: don't include if value is the default
	Version string    `json:"version"` //Add the version
}

// SeatBid represents a seat's bid(s) within the bid response.
type SeatBid struct {
	Bid  []Bid  `json:"bid"`
	Seat string `json:"seat"`
}

// Bid represents an individual bid.
type Bid struct {
	ID    string  `json:"id"`
	ImpID string  `json:"impid"`
	Price float64 `json:"price"`
	AdM   string  `json:"adm"`
	CrID  string  `json:"crid"`
	// NUrl  string  `json:"nurl,omitempty"` // Add this for win notification
	W int `json:"w"` //Add width
	H int `json:"h"` //Add height
}

// Define the handler function. This function will be called when
// the server receives an HTTP request on the specified path.
//
// w: The http.ResponseWriter is used to construct and send the HTTP response.
// r: The http.Request contains the data from the client's request. We'll
// use this to read the request body and unmarshal it into our
// BidRequest struct.
func bidHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST.  OpenRTB requests should
	// always use POST.  If the method is not POST, we return a
	// Method Not Allowed error.
	if r.Method != http.MethodPost {
		log.Printf("Error: invalid request method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return // Important: return after handling the error.
	}

	// Read the entire request body.  We need to do this because
	// the request body is a stream, and we need to read it all
	// at once to parse it as JSON.  We use io.ReadAll to ensure
	// we read until EOF.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Handle errors!  If reading the body fails, log the error and
		// send a Bad Request error to the client.  We use log.Printf
		// for non-fatal errors that we want to record.
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return // Important: return after handling the error.
	}

	// Unmarshal the JSON data into a BidRequest struct.  This is
	// where we convert the JSON data from the request body into
	// a Go struct that we can work with.  We use the json.Unmarshal
	// function for this.
	var bidRequest BidRequest
	err = json.Unmarshal(body, &bidRequest)
	if err != nil {
		// Handle errors!  If unmarshalling fails, it means the JSON
		// in the request was invalid.  Log the error and send a
		// Bad Request error to the client.
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return // Important: return after handling the error.
	}

	// Validate the required fields.  This is where we check if the
	// request contains all the necessary data.  In a real-world
	// application, you would have more complex validation logic.
	if bidRequest.ID == "" {
		log.Println("Error: request ID is empty")
		http.Error(w, "Request ID is required", http.StatusBadRequest)
		return
	}
	if len(bidRequest.Imp) == 0 {
		log.Println("Error: at least one impression is required")
		http.Error(w, "At least one impression is required", http.StatusBadRequest)
		return
	}
	for _, imp := range bidRequest.Imp {
		if imp.ID == "" {
			log.Println("Error: impression ID is empty")
			http.Error(w, "Impression ID is required", http.StatusBadRequest)
			return
		}
	}
	if bidRequest.Device == nil {
		log.Println("Error: device object is required")
		http.Error(w, "Device object is required", http.StatusBadRequest)
		return
	}
	if bidRequest.Device != nil && bidRequest.Device.Ua == "" {
		log.Println("Error: device UserAgent is required")
		http.Error(w, "Device UserAgent is required", http.StatusBadRequest)
		return
	}

	// Create a dummy BidResponse object. This is the data we will
	// serialize to JSON and send back to the client. The values
	// here are examples; in a real application, these would be
	// populated based on the bid request and your bidding logic.
	response := BidResponse{
		ID:      bidRequest.ID, // Use the ID from the request
		Version: "2.5",       // Set OpenRTB version
		SeatBid: []SeatBid{
			{
				Seat: "seat-id-456",
				Bid: []Bid{
					{
						ID:    "bid-id-789",
						ImpID: bidRequest.Imp[0].ID, // Use the ImpID from the first impression
						Price: 2.50,                 // Use a float for price
						AdM:   "Example Ad", // Basic HTML ad markup
						CrID:  "creative-id-abc",
						W:     300, // width
						H:     250, // height
						// NUrl: "http://example.com/win-notification?bidid=${AUCTION_ID}&price=${AUCTION_PRICE}", // win notice URL,
					},
				},
			},
		},
	}

	// Set the Content-Type header. This tells the client (e.g., the ad exchange)
	// that the response is JSON. It's crucial to set the correct content type.
	w.Header().Set("Content-Type", "application/json")

	// Encode the response struct to JSON. json.NewEncoder creates an encoder
	// that writes to the provided io.Writer (in this case, the http.ResponseWriter).
	// The Encode method then serializes the struct to JSON and writes it to the stream.
	encoder := json.NewEncoder(w)
	err = encoder.Encode(response)
	if err != nil {
		// Handle errors! If the JSON encoding fails, log the error and send a
		// generic server error response. Proper error handling is essential
		// for a robust server. We use log.Printf for non-fatal errors.
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return // Important: return after handling the error.
	}

	// Log a successful response. This is helpful for debugging and monitoring.
	log.Printf("Successfully sent bid response at %s", time.Now().Format(time.RFC3339))
}

// main is the entry point of the program.
func main() {
	// Handle requests to the "/bid" path with the bidHandler function.
	// This line tells the http package to route all requests with the path
	// "/bid" to the bidHandler function.
	http.HandleFunc("/bid", bidHandler)

	// Start the server. Listen on port 8080 and handle incoming requests.
	// The ListenAndServe function blocks, meaning the program will stay
	// running here until the server is stopped. We use a log.Fatal
	// here because if the server fails to start, the program cannot function.
	fmt.Println("Server listening on :8080") //Inform that the server has started
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}