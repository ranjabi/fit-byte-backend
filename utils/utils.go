package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"fit-byte/constants"
	"fit-byte/models"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.SALT_ROUND)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func CreateClaims(user *models.User) (string, error) {
	tokenAuth := jwtauth.New(constants.HASH_ALG, []byte(constants.JWT_SECRET), nil)
	claims := map[string]any{
		"userId":    user.Id,
		"userEmail": user.Email,
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetJsonResponse(w http.ResponseWriter, statusCode int, response any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
}

func AppHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if err, ok := err.(*models.AppError); ok {
				if err.Code != 0 {
					http.Error(w, err.Error(), err.Code)
					return
				}
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AllowContentType(contentTypes ...string) func(http.Handler) http.Handler {
	allowedContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, ctype := range contentTypes {
		allowedContentTypes[strings.TrimSpace(strings.ToLower(ctype))] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				// Skip check for empty content body
				next.ServeHTTP(w, r)
				return
			}

			s := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))

			if _, ok := allowedContentTypes[s]; ok {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
		})
	}
}

func GetJSONTagName(field reflect.StructField) string {
	tag := field.Tag.Get("db")

	parts := strings.Split(tag, ",")

	return parts[0]
}

func BuildPartialUpdateQuery(tableName, idField, idValue string, data interface{}) (string, pgx.NamedArgs, error) {
	val := reflect.ValueOf(data).Elem()
	typ := reflect.TypeOf(data).Elem()

	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	args := pgx.NamedArgs{}
	var setClauses []string
	index := 1

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldName := GetJSONTagName(typ.Field(i))

		if fieldName == "" || fieldName == "-" {
			continue
		}

		if !fieldValue.IsNil() {
			if fieldName == idField {
				setClauses = append(setClauses, fmt.Sprintf("%s = @%sNew", fieldName, fieldName))
				args[fieldName+"New"] = fieldValue.Elem().Interface() // Dereference the pointer
			} else {
				setClauses = append(setClauses, fmt.Sprintf("%s = @%s", fieldName, fieldName))
				args[fieldName] = fieldValue.Elem().Interface()
			}
			index++
		}
	}

	if len(setClauses) == 0 {
		query = fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE %s = @%s
		`, tableName, idField, idField)
		args = pgx.NamedArgs{
			idField: idValue,
		}
		return query, args, nil
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE %s = @%s", idField, idField)
	args[idField] = idValue
	query += " RETURNING *"

	return query, args, nil
}

func IsValidHost(host string) bool {
	return strings.Count(host, ".") > 0
}

func IsValidURI(urlString string) bool {
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}

	return IsValidHost(parsedURL.Hostname())
}

func CalculateCaloriesBurned(activity string, minutes int) int {
	activityCaloriesBurnedPerMinute := map[string]int{
		"Walking":    4,
		"Yoga":       4,
		"Stretching": 4,
		"Cycling":    8,
		"Swimming":   8,
		"Dancing":    8,
		"Hiking":     10,
		"Running":    10,
		"HIIT":       10,
		"JumpRope":   10,
	}

	return activityCaloriesBurnedPerMinute[activity] * minutes
}

func GenerateS3FileURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("S3_BUCKET_NAME"), os.Getenv("AWS_REGION"), key)
}