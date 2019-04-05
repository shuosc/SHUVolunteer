package token

import (
	"SHUVolunteer/model/student"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func GenerateJWT(student student.Student) string {
	result, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentId": student.Id,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return result
}
