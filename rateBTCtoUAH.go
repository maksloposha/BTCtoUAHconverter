package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
type ExchangeRate struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse JSON request body into User struct
	var newUser User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists in the file
	exists, err := userExists(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Open the file in append mode
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the existing file content
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the file is empty, write the first user object
	if fileInfo.Size() == 0 {
		_, err = file.WriteString("[\n")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Remove the closing square bracket from the existing content
		newContent := fileContent[:len(fileContent)-2]

		// Write the updated content back to the file
		err = ioutil.WriteFile("users.json", newContent, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Append a comma separator
		_, err = file.WriteString(",\n")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Convert user to JSON
	jsonData, err := json.MarshalIndent(newUser, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write user JSON to file
	_, err = file.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the closing square bracket
	_, err = file.WriteString("\n]\n")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

func userExists(newUser User) (bool, error) {
	// Open the file in read mode
	file, err := os.Open("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	// Read the file content
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return false, err
	}

	// Unmarshal the file content into a slice of User structs
	var users []User
	err = json.Unmarshal(fileContent, &users)
	if err != nil {
		return false, err
	}

	// Check if user already exists in the slice
	for _, user := range users {
		if user.ID == newUser.ID {
			return true, nil
		}
	}

	return false, nil
}

func getBTCUANRate(w http.ResponseWriter, r *http.Request) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(exchangeRate.Bitcoin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func sendEmails(w http.ResponseWriter, r *http.Request) {
	// Get the current BTC to UAH exchange rate from a third-party service
	rate, err := http.Get("http://localhost:8000/rate")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the subscribed email addresses from your database or storage
	subscribedEmails, err := readUsersFromFile("users.json")

	// Compose the email message with the exchange rate
	subject := "BTC to UAH Exchange Rate"
	body := fmt.Sprintf("The current BTC to UAH exchange rate is %.2f", rate)

	// Send the email to each subscribed email address
	for _, email := range subscribedEmails {
		err := sendEmail(email.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send email to %s: %s", email, err.Error())
		}
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Emails sent successfully"))
}
func sendEmail(email, subject, body string) error {
	// Configure the SMTP server settings
	smtpServer := "smtp.example.com"
	smtpPort := 587
	smtpUsername := "your-smtp-username"
	smtpPassword := "your-smtp-password"

	// Compose the email message
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
func readUsersFromFile(filePath string) ([]User, error) {
	// Read the JSON file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Declare a slice to store the users
	var users []User

	// Unmarshal the JSON data into the users slice
	err = json.Unmarshal(fileContent, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func main() {
	http.HandleFunc("/sendEmails", sendEmails)
	http.HandleFunc("/rate", getBTCUANRate)
	http.HandleFunc("/subscribe", addUser)
	fmt.Println("Server listening on port 8000")

	log.Fatal(http.ListenAndServe(":8000", nil))

}
