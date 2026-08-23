package auth

import "pkg/user"

<<<<<<< HEAD
func ValidateToken(t string) (*user.Profile, error) {
	return user.Find(t)
}
=======
func ValidateToken(t string) (*user.Profile, error) {
	return user.FindByToken(t)
}
>>>>>>> branch
