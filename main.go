package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	AppSchema  = "app"
	AuthSchema = "auth"
)

// User is a model representing the users table
type User struct {
	ID   uint
	Name string
}

// TableName returns the custom table name for the User model
func (User) TableName() string {
	return AuthSchema + ".user"
}

// Order is a model representing the orders table
type Order struct {
	ID     uint
	Name   string
	UserID uint
}

// TableName returns the custom table name for the Order model
func (Order) TableName() string {
	return AppSchema + ".order"
}

// AuthScope is a custom scope for the auth schema
func AuthScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("auth.user")
	}
}

// AppScope is a custom scope for the app schema
func AppScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("app.order")
	}
}

func main() {
	// Connect to the database
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Create schema if not exist
	db.Exec("CREATE SCHEMA IF NOT EXISTS " + AppSchema)
	db.Exec("CREATE SCHEMA IF NOT EXISTS " + AuthSchema)

	// Auto-migrate the models
	err = db.AutoMigrate(&User{}, &Order{})
	if err != nil {
		panic(err)
	}

	// Create a new user in the auth schema
	user := User{Name: "John"}
	db.Scopes(AuthScope()).Create(&user)

	// Create a new order associated with the user in the app schema
	order := Order{Name: "Order 1", UserID: user.ID}
	db.Scopes(AppScope()).Create(&order)

	// Fetch all users from the auth schema
	var authUsers []User
	db.Scopes(AuthScope()).Find(&authUsers)
	fmt.Println("Auth Users:", authUsers)

	// Fetch all orders from the app schema
	var appOrders []Order
	db.Scopes(AppScope()).Find(&appOrders)
	fmt.Println("App Orders:", appOrders)

	// Fetch all orders from the app schema with the associated user
	var appOrdersWithUser []Order
	db.Scopes(AppScope()).Preload(clause.Associations).Find(&appOrdersWithUser)
	fmt.Println("App Orders with User:", appOrdersWithUser)
}
