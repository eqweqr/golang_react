package middleware

import (
	"eqweqr/bdkurach/internals/jwttoken"
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.Handler, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		jwtToken := jwttoken.GetToken(bearer)
		_, err := jwttoken.ParseToken(jwtToken, secretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		// if !jwttoken.IsValid(*parsedToken) {
		// w.WriteHeader(http.StatusUnauthorized)
		// fmt.Println("unauthorized user")
		// return
		// }

		next.ServeHTTP(w, r)
	})
}
