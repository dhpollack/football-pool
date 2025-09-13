// Package auth provides authentication context keys and types for the football pool application.
package auth

type contextKey string

// EmailKey is the context key used to store and retrieve user email from request context.
const EmailKey contextKey = "email"
