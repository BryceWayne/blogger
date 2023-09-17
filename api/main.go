package main

import (
    "context"
    "log"

    "cloud.google.com/go/firestore"
    firebase "firebase.google.com/go"
    "github.com/BryceWayne/blogger/routes"
    "github.com/BryceWayne/blogger/utils"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
    "google.golang.org/api/option"
)

var client *firestore.Client
var ctx = context.Background()
var jwtSecret []byte
var app *firebase.App
var config *utils.Config // Assume you've a ConfigType in utils

func init() {
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v\n", err)
    }

    config, err = utils.LoadConfig()
    if err != nil {
        log.Fatalf("Could not load config: %v", err)
    }

    jwtSecret = []byte(config.JwtSecret)
    opt := option.WithCredentialsFile(config.GCPCreds)

    // Initialize Firebase App
    app, err = firebase.NewApp(ctx, nil, opt)
    if err != nil {
        log.Fatalf("Error initializing Firebase App: %v\n", err)
    }

    // Initialize Firestore
    client, err = app.Firestore(ctx)
    if err != nil {
        log.Fatalf("Could not initialize Firestore client: %v\n", err)
    }
}

func main() {
    // Create a new Fiber app
    app := fiber.New()
    app.Use(logger.New())

    // Limit each IP to 10 requests per minute
    // app.Use(limiter.New(limiter.Config{
    //     Max: 10,
    // }))

    app.Get("/post", func(c *fiber.Ctx) error {
        return routes.CreatePost(c, client)
    })

    app.Get("/posts", func(c *fiber.Ctx) error {
        return routes.GetPosts(c, client)
    })

    // Start the Fiber app
    log.Fatal(app.Listen(":8080"))
}
