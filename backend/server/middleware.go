package server

import (
	"context"
	"eqweqr/bdkurach/internals/jwttoken"
	"log"
	"net/http"
)

func (server *Server) RoleMiddleware(next http.Handler, needRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearer := r.Header.Get("Authorization")
		jwtToken := jwttoken.GetToken(bearer)
		log.Println(jwtToken)
		token, err := jwttoken.ParseToken(jwtToken, server.Secret)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		role := jwttoken.GetRoles(token)
		id := jwttoken.GetId(token)

		if (role == "worker" || role == "client") && role != needRole {
			log.Println("different role")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
