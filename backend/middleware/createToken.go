package middleware

// func RegisterNewUser(next http.Handler, secret string) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token, err := jwttoken.CreateToken(secret, "root", "user")
// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
// 		w.Write([]byte(token))
// 	})
// }
