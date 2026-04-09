package hasher

import "golang.org/x/crypto/bcrypt"

func Hash(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	return string(b), err
}

func Compare(hash, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}
