package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adeo/iwc-gcp-firewall-api/handlers"
	"github.com/adeo/iwc-gcp-firewall-api/helpers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Enable http access log on testing
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"request_uri": r.RequestURI,
			"user_agent":  r.UserAgent(),
		}).Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Define JSON as default returned content type
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Init logger to be Stackdriver compliant
	helpers.InitLogger()

	// Set port to listen to
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter().StrictSlash(true)
	if os.Getenv("CI") == "" {
		r.Use(loggingMiddleware)
	}

	// Define routers
	projectRouter := r.PathPrefix("/project/{project}").Subrouter()
	serviceProjectRouter := projectRouter.PathPrefix("/service_project/{service_project}").Subrouter()
	applicationRouter := serviceProjectRouter.PathPrefix("/application/{application}").Subrouter()
	ruleRouter := applicationRouter.PathPrefix("/firewall_rule/{rule}").Subrouter()

	// Manage sets of rules routes
	applicationRouter.Path("").Methods(http.MethodGet).HandlerFunc(handlers.ListFirewallRuleHandler)

	// Manage a specific rule
	ruleRouter.Path("").Methods(http.MethodPost).HandlerFunc(handlers.CreateFirewallRuleHandler)
	ruleRouter.Path("").Methods(http.MethodGet).HandlerFunc(handlers.GetFirewallRuleHandler)
	ruleRouter.Path("").Methods(http.MethodDelete).HandlerFunc(handlers.DeleteFirewallRuleHandler)

	// Other endpoints routes
	r.Path("/_health").Methods(http.MethodGet).HandlerFunc(handlers.HealthCheckHandler)

	// Override default error handlers
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	r.Use(contentTypeMiddleware)

	srv := http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Printf("Listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)
	logrus.Println("Shutting down server")
	os.Exit(0)
}
