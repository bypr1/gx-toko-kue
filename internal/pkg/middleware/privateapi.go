package middleware

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"service/internal/pkg/config"
	"service/internal/pkg/encryption"
	error2 "service/internal/pkg/error"
	"time"
)

func AuthenticatePrivateAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := r.Header.Get("Client-ID")
		if clientId == "" {
			error2.ErrXtremePrivateAPIAuthentication("Client ID is missing!")
		}

		clientName := r.Header.Get("Client-Name")
		if clientName == "" {
			error2.ErrXtremePrivateAPIAuthentication("Client Name missing!")
		}

		clientSecret := r.Header.Get("Client-Secret")
		if clientSecret == "" {
			error2.ErrXtremePrivateAPIAuthentication("Client Secret missing!")
		}

		key := fmt.Sprintf("%s-%s", clientId, clientName)
		if _, ok := config.XtremeCache.Get(key); !ok {
			ec := encryption.NewPrivateAPIEncryption(clientId)

			decrypt, err := ec.Decrypt(clientName, clientSecret)
			if err != nil {
				error2.ErrXtremePrivateAPIAuthentication(err.Error())
			}

			err = bcrypt.CompareHashAndPassword([]byte(decrypt), []byte(config.PrivateAPICredential[clientName]["key"].(string)))
			if err != nil {
				error2.ErrXtremePrivateAPIAuthentication("Your client secret is incorrect!")
			}

			config.XtremeCache.Set(key, decrypt, time.Hour)
		}

		next.ServeHTTP(w, r)
	})
}
