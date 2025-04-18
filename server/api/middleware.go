package main

import (
	"chat/utils"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	missingAuthHeaderError = errors.New("Authorization header missing")
	invalidAuthHeaderError = errors.New("Invalid Authorization header")
)

var allowedUnauthorizedRoutes = [3]string{
	"GET /",
	"POST /login",
	"POST /user",
}

func isAuthNeeded(method, path string) bool {
	url := method + " " + path
	for _, route := range allowedUnauthorizedRoutes {
		if route == url {
			return false
		}
	}
	return true
}

func (app *application) authenticateDummy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !isAuthNeeded(request.Method, request.URL.Path) {
			next.ServeHTTP(writer, request)
			return
		}

		header := request.Header.Get("Authorization")
		if header == "" {
			utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, missingAuthHeaderError)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, invalidAuthHeaderError)
			return
		}

		userIdString := headerParts[1]
		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, err)
			return
		}

		req := utils.ContextSetUser(request, userId)
		next.ServeHTTP(writer, req)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !isAuthNeeded(request.Method, request.URL.Path) {
			next.ServeHTTP(writer, request)
			return
		}

        var err error
        authTokenCookie, err := request.Cookie("authToken")
        if err != nil {
            err = errors.New("missing authorization token")
            utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, err)
            return
        }
		userId, err := app.authService.VerifyToken(authTokenCookie.Value)
		if err != nil {
			utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, err)
			return
		}

		req := utils.ContextSetUser(request, userId)
		next.ServeHTTP(writer, req)
	})
}

func (app *application) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		log.Println(request.Method, request.URL.Path, time.Since(start))
	})
}

func (app *application) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
        writer.Header().Set("Access-Control-Allow-Credentials", "true")
        writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours cache

        if request.Method == "OPTIONS" {
            writer.WriteHeader(http.StatusNoContent)
            return
        }

		next.ServeHTTP(writer, request)
	})
}
