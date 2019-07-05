package models

import (
	"log"
	"os"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

// JwtToken ...
type JwtToken struct {
	jwt.StandardClaims
	UserID uint
}

// Find user by id
func Find(id uint) *User {

	user := &User{}
	DB().Table("users").Where("id = ?", id).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = "" // don't return psw in response
	return user
}

func Login(email, psw string) map[string]interface{} {
	user := &User{}

	err := DB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]interface{}{
				"status":  false,
				"message": "User couldn't be found",
			}
		}

		return map[string]interface{}{
			"status":  false,
			"message": "DB error",
		}
	}

	pswErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(psw))
	if pswErr != nil {
		log.Fatalf("psw compare error: %v", pswErr)
		if pswErr == bcrypt.ErrMismatchedHashAndPassword {
			return map[string]interface{}{
				"status":  false,
				"message": "Passwords do not match!",
			}
		}
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &JwtToken{UserID: user.ID})
	tokenStr, _ := token.SignedString([]byte(os.Getenv("jwt_secret")))
	user.Token = tokenStr

	user.Password = "" // don't return psw in response
	return map[string]interface{}{
		"status":  true,
		"message": "User logged in successfully!",
		"user":    user,
	}
}

// Create creates new user
func (user *User) Create() map[string]interface{} {
	if msg, ok := user.Validate(); !ok {
		return msg
	}

	pswHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(pswHash)

	DB().Create(user)

	if user.ID <= 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "User could not be created",
		}
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &JwtToken{UserID: user.ID})
	tokenStr, _ := token.SignedString([]byte(os.Getenv("jwt_secret")))
	user.Token = tokenStr

	user.Password = "" // don't return psw in response

	return map[string]interface{}{
		"status":  true,
		"message": "User was created successfully",
		"user":    user,
	}
}

// Validate validates user data when attempting to create an account
func (user *User) Validate() (map[string]interface{}, bool) {

	if !validateEmail(user.Email) {
		return map[string]interface{}{
			"status":  false,
			"message": "Email wasn't valid",
		}, false
	}

	if len(user.Password) < 6 {
		return map[string]interface{}{
			"status":  false,
			"message": "Password wasn't valid",
		}, false
	}

	tmpUser := &User{}

	err := DB().Table("users").Where("email = ?", user.Email).First(tmpUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return map[string]interface{}{
			"status":  false,
			"message": "User validate error",
		}, false
	}
	if tmpUser.Email != "" {
		return map[string]interface{}{
			"status":  false,
			"message": "User already exists",
		}, false
	}

	return map[string]interface{}{
		"status":  false,
		"message": "User data was valid",
	}, true
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}
