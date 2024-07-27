package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NishantBansal2003/Brokerax/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_           = godotenv.Load(".env")
	SECRET_KEY  = os.Getenv("SECRET_KEY")
	MONGODB_URI = os.Getenv("MONGODB_URI")
)

const (
	dbName  = "test"
	colName = "UserInfo"
)

var collection *mongo.Collection

func Connect() {
	// Special method runs first and one time only
	clientOptions := options.Client().ApplyURI(MONGODB_URI).
		SetMaxPoolSize(50).
		SetServerSelectionTimeout(2 * time.Second)

		// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection reference is ready!")
}

// MongoDB helpers
func insertNewUser(newUser model.User) {
	inserted, err := collection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Fatal("Error creating user ", err)
	}
	fmt.Println("Inserted new User", inserted)
}

func findUser(email string) (*model.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	var result model.User
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(userIDp primitive.ObjectID, email string) (string, error) {
	userID := userIDp.Hex()
	claims := Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func Login(c *fiber.Ctx) error {
	var data LoginRequest

	// Parse the request body into the data struct
	if err := c.BodyParser(&data); err != nil {
		// If parsing fails, fallback to query parameters
		data.Email = c.Query("email")
		data.Password = c.Query("password")
	}
	existingUser, err := findUser(data.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Error finding user",
			"detail": err.Error(), // Include the error message
		})
	}

	if existingUser == nil || existingUser.Password != data.Password {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Wrong details Please Check Once",
		})
	}

	token, err := GenerateToken(existingUser.ID, existingUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Issue while generating token, Please Try Again!",
		})
	}
	return c.JSON(fiber.Map{
		"User":  existingUser,
		"Token": token,
	})
}

type SignUpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func ConvertToUser(req SignUpRequest) *model.User {
	return &model.User{
		ID:        primitive.NewObjectID(), // Generate a new ObjectID for the user
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		UserType:  "regular",       // Set default or determine dynamically
		Credits:   1000000,         // Default value
		Stocks:    []model.Stock{}, // Default empty slice
	}
}

func Signup(c *fiber.Ctx) error {
	var data SignUpRequest

	// Parse the request body into the data struct
	if err := c.BodyParser(&data); err != nil {
		// If parsing fails, fallback to query parameters
		data.FirstName = c.Query("firstName")
		data.LastName = c.Query("lastName")
		data.Email = c.Query("email")
		data.Password = c.Query("password")
	}
	_, err := findUser(data.Email)
	if err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User Already Exists",
		})
	}
	newUser:=*ConvertToUser(data)
	insertNewUser(newUser)
	token, err := GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Issue while generating token, Please Try Again!",
		})
	}
	return c.JSON(fiber.Map{
		"User":  newUser,
		"Token": token,
	})
}

func Logout(c *fiber.Ctx) error {
	
	return nil
}

func Portfolio(c *fiber.Ctx) error {
	return nil
}

func AddStock(c *fiber.Ctx) error {
	return nil
}

func RemoveStock(c *fiber.Ctx) error {
	return nil
}
