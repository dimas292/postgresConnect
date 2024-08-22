package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"env/database" // Ensure this path is correct for your project
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// Config structure for YAML configuration
type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

type User struct {
	Id        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// Load configuration from YAML file
func loadConfigYaml(filename string) (conf Config, err error) {
	// Read the file
	f, err := os.ReadFile(filename)
	if err != nil {
		return conf, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal YAML into the Config struct
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		return conf, fmt.Errorf("error unmarshaling config file: %w", err)
	}
	return conf, nil
}

func main() {

	// Load configuration once
	cfg, err := loadConfigYaml("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// db connection
	db, err := database.ConnectPostgres(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Name,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the database is connected
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	// migrate data
	err = Migrate(db)
	if err != nil {
		log.Print("fail migrate with error ", err.Error())
	}
	// insert data
	err = Insert(db, User{
		Name:  "xxxxx",
		Email: "xxxxx22@gmail.com",
	})
	if err != nil {
		log.Println("fail insert data with error", err.Error())
	}
	// get all data
	users, err := GetAll(db)
	if err != nil {
		log.Print("fail to getAll with error :", err.Error())
		return
	}
	fmt.Println(strings.Repeat("=", 10), "GET ALL", strings.Repeat("=", 10))

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}

	log.Println(strings.Repeat("=", 30), "\n")
	// get one by id

	user, err := GetOneById(db, 1)
	if err != nil {
		log.Print("fail to getAll with error :", err.Error())
		return
	}
	fmt.Println(strings.Repeat("=", 10), "GET ONE", strings.Repeat("=", 10))
	fmt.Printf("%+v\n", user)
	fmt.Println(strings.Repeat("=", 30), "\n")

	// get one end

	// Define the HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Respond with the configuration as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cfg); err != nil {
			http.Error(w, "Failed to encode config", http.StatusInternalServerError)
		}
	})

	// Start the HTTP server
	log.Printf("Server running at port : %s", cfg.App.Port)
	http.ListenAndServe(cfg.App.Port, nil)
}

func Migrate(db *sql.DB) (err error) {
	query := `
	CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name varchar(100),
    email varchar(100),
    created_at timestamptz DEFAULT now()
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		return
	}

	return
}

func Insert(db *sql.DB, user User) (err error) {
	log.Println("try to insert data")

	query := `
	INSERT INTO users (
	name, email) VALUES ($1, $2)
	`

	_, err = db.Exec(query, user.Name, user.Email)
	log.Println("insert data done")
	return
}

func GetAll(db *sql.DB) (users []User, err error) {

	query := `
	SELECT id, name, email, created_at FROM users`

	rows, err := db.Query(query)
	if err != nil {
		return
	}

	for rows.Next() {
		var user = User{}

		err = rows.Scan(
			&user.Id, &user.Name, &user.Email, &user.CreatedAt,
		)

		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}

func GetOneById(db *sql.DB, id int)(user User, err error){

	query := `
	SELECT id, name, email, created_at
	FROM users
	WHERE id=$1`

	rows := db.QueryRow(query, id)
	err = rows.Scan(
		&user.Id, &user.Name, &user.Email, &user.CreatedAt,
	)

	return


}