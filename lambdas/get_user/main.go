package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Health check for the database connection
    // retry if err, maxes out at three tries
    if err := pingDatabaseWithRetry(ctx, db, 3, 2*time.Second); err != nil {
        log.Printf("Error pinging database: %v", err)
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusInternalServerError,
            Body:       "Database connection error.",
        }, nil
    }

    query := "SELECT * FROM Users WHERE UserID = ?"
    userID := request.RequestContext.Identity.CognitoIdentityPoolID

    type User struct {
        UserID                    string
        Username                  string
        Email                     string
        PhoneNumber               string
        DateOfBirth               string
        AccountVerificationStatus string
        CreatedAt                 string
        UpdatedAt                 string
        AccountBalance            string
    }

    var user User
    err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Email, &user.PhoneNumber, &user.DateOfBirth, &user.AccountVerificationStatus, &user.CreatedAt, &user.UpdatedAt, &user.AccountBalance)

    if err != nil {
        if err == sql.ErrNoRows {
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusNotFound,
                Body:       fmt.Sprintf("User with ID %s not found.", userID),
            }, nil
        }
        log.Printf("Error executing query: %v", err)
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusInternalServerError,
            Body:       fmt.Sprintf("Internal server error: %v", err),
        }, nil
    }

    userJSON, err := json.Marshal(user)
    if err != nil {
        log.Printf("Error marshalling user struct to JSON: %v", err)
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusInternalServerError,
            Body:       fmt.Sprintf("Error processing request: %v", err),
        }, nil
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(userJSON),
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    if err := initializeDatabase(); err != nil {
        fmt.Printf("failed to initialize database: %v", err)
    }
    lambda.Start(handler)
}

func initializeDatabase() error {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2"),
    })
    if err != nil {
        log.Printf("Error creating AWS session: %v", err)
        return err
    }

    ssmSvc := ssm.New(sess)
    paramName := "/application/dev/database/credentials"
    withDecryption := true
    param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
        Name:           &paramName,
        WithDecryption: &withDecryption,
    })
    if err != nil {
        log.Printf("Error getting parameter: %v", err)
        return err
    }

    var dbCreds struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Host     string `json:"host"`
        Port     int    `json:"port"`
    }
    err = json.Unmarshal([]byte(*param.Parameter.Value), &dbCreds)
    if err != nil {
        log.Printf("Error parsing JSON: %v", err)
        return err
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/user_management", dbCreds.Username, dbCreds.Password, dbCreds.Host, dbCreds.Port)
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Printf("Error opening database: %v", err)
        return err
    }
    fmt.Println("Connected to the MySQL database successfully!")
    return nil
}

func pingDatabaseWithRetry(ctx context.Context, db *sql.DB, attempts int, delay time.Duration) error {
    for i := 0; i < attempts; i++ {
        if err := db.PingContext(ctx); err != nil {
            if i < (attempts - 1) {
                log.Printf("Database ping failed, retrying in %v... (%d/%d)", delay, i+1, attempts)
                time.Sleep(delay)
                continue
            }
            return fmt.Errorf("database ping attempts failed after %d tries: %w", attempts, err)
        }
        return nil // Ping successful
    }
    return fmt.Errorf("failed to ping database after %d attempts", attempts)
}