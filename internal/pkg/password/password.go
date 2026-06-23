package password

import "golang.org/x/crypto/bcrypt"

const cost = 12

// Hash bcrypt-hashes a plain-text password
func Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare checks whether a plain-text password matches the hash
func Compare(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
