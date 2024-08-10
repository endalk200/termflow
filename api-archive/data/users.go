package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Utility function to load users from file
func loadUsers() ([]User, error) {
	file, err := os.ReadFile("data/users.json")
	if err != nil {
		return nil, err
	}

	var users []User
	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func saveUsers(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("data/users.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetUsers() []User {
	users, err := loadUsers()
	if err != nil {
		panic("Error reading data file")
	}

	return users
}

func FindUserByEmail(email string) (User, error) {
	users := GetUsers()

	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user with [Email] %s not found", email)
}

func FindUserById(id int) (User, error) {
	users := GetUsers()

	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user with ID %d not found", id)
}

func CreateUser(data User) (User, error) {
	users, err := loadUsers()
	if err != nil {
		return User{}, err
	}

	for _, user := range users {
		if data.Email == user.Email {
			log.Printf("There is already a user with email: [%s]", data.Email)
			return User{}, fmt.Errorf("There is already a user with email: [%s]", data.Email)
		}
	}

	var id int
	if len(users) == 0 {
		id = 1
	} else {
		// Sort the slice by Id
		sort.Slice(users, func(i, j int) bool {
			return users[i].Id < users[j].Id
		})

		// Get the highest Id
		latestUser := users[len(users)-1]
		latestId := latestUser.Id
		id = latestId + 1
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Something went wrong while hashing password: %s", err)
		return User{}, fmt.Errorf("Something went wrong while saving user account")
	}

	newUser := User{
		Id:        id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  string(password),
	}

	users = append(users, newUser)

	err = saveUsers(users)
	if err != nil {
		log.Printf("Something went wrong while writing new data to file %s", err)
		return User{}, fmt.Errorf("Something went wrong while creating a user account")
	}

	return newUser, nil
}

func DeleteUser(id int) error {
	users, err := loadUsers()
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			return saveUsers(users)
		}
	}

	return fmt.Errorf("user with ID %d not found", id)
}
