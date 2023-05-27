User Subscription App
This is a simple web application that allows users to subscribe to receive emails with the current exchange rate of
Bitcoin to Ukrainian Hryvnia (UAH). The application is built using Go programming language and utilizes HTTP requests 
to interact with external APIs and store user data in a JSON file.

Features
User subscription: Users can subscribe to receive email notifications with the current BTC to UAH exchange rate.
Email notifications: The application fetches the exchange rate from a third-party service and sends email notifications to the subscribed users.
User existence check: The application checks if a user already exists before adding them to the subscription list.
Prerequisites
Before running the application, make sure you have the following installed:

Go (version 1.16 or higher)
Git
Getting Started
Clone the repository:

Copy code
git clone https://github.com/maksloposha/user-subscription-app.git
Navigate to the project directory:


Copy code
cd user-subscription-app
Build and run the application:

Copy code
go run main.go
The application will start running on http://localhost:8000.

API Endpoints
The application provides the following API endpoints:

POST /sendEmails: Allows you to send the exchange rate of Hryvnia to Bitcoin at users.json
POST /subscribe: Allows users to subscribe by providing their user ID and email address.
GET /rate: Retrieves the current BTC to UAH exchange rate from a third-party API.
Configuration
The application uses a few configuration parameters that you can modify as needed. Update the following variables in the main function of the main.go file:

smtpServer: The SMTP server address for sending emails.
smtpPort: The SMTP server port.
smtpUsername: Your SMTP username.
smtpPassword: Your SMTP password.
Make sure to update these variables with the appropriate values for your SMTP server.

Testing
The application includes unit tests to ensure the correctness of the core functionalities. To run the tests, use the following command:


Copy code
go test