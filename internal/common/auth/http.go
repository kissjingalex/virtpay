package auth

type ctxKey int

const (
	userContextKey ctxKey = iota
)

type User struct {
	UUID  string
	Email string
	Role  string

	DisplayName string
}
