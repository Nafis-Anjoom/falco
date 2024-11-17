package utils

import (
	"context"
	"net/http"
)

type contextKey string
const userContextKey contextKey = "user"

func GetUserFromRequest(request *http.Request) int64 {
    userId, ok := request.Context().Value(userContextKey).(int64)
    // should never happen in prod
    if !ok {
        panic("missing user in request")
    }

    return userId
}

func SetUserInRequest(request *http.Request, userId int64) *http.Request {
    context := context.WithValue(request.Context(), userContextKey, userId)
    return request.WithContext(context)
}
