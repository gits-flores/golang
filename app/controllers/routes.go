package controllers

import (
	"github.com/gits-flores/golang/app/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.Register)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	s.Router.HandleFunc("/articles", middlewares.SetMiddlewareJSON(s.CreateArticle)).Methods("POST")
	s.Router.HandleFunc("/articles", middlewares.SetMiddlewareJSON(s.GetArticles)).Methods("GET")
	s.Router.HandleFunc("/articles/{id}", middlewares.SetMiddlewareJSON(s.GetArticle)).Methods("GET")
	s.Router.HandleFunc("/articles/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateArticle))).Methods("PUT")
	s.Router.HandleFunc("/articles/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteArticle)).Methods("DELETE")
}