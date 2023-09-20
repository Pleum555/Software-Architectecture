package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

func GetCurrentLocation(w http.ResponseWriter, r *http.Request) {
	userValue := context.Get(r, "user")
	user, _ := userValue.(*CustomClaims)

	fmt.Println("Username:", user.Username)
	fmt.Println("Password:", user.Password)

	return
}
