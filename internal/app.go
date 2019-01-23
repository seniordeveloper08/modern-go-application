package internal

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goph/emperror"
	"github.com/goph/logur"
	"github.com/gorilla/mux"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkeradapter"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkerdriver"
	"github.com/sagikazarmark/modern-go-application/internal/httpbin"
)

// NewApp returns a new HTTP application.
func NewApp(logger logur.Logger, publisher message.Publisher, errorHandler emperror.Handler) http.Handler {
	sayHello := greeting.NewHelloService(
		greetingadapter.NewSayHelloEvents(publisher),
		greetingadapter.NewLogger(logger),
		errorHandler,
	)
	helloWorldController := greetingdriver.NewHTTPController(sayHello, errorHandler)

	router := mux.NewRouter()

	router.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write([]byte(template))
	})

	router.Path("/hello").Methods("POST").HandlerFunc(helloWorldController.SayHello)

	router.PathPrefix("/httpbin").Handler(http.StripPrefix("/httpbin", httpbin.New()))

	return router
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.Logger) error {
	sayHelloHandler := greetingworkerdriver.NewSayHelloEventHandler(
		greetingworker.NewSayHelloEventLogger(greetingworkeradapter.NewLogger(logger)),
	)

	err := router.AddNoPublisherHandler(
		"log_said_hello_to",
		"said_hello_to",
		subscriber,
		sayHelloHandler.SaidHelloTo,
	)
	if err != nil {
		return err
	}

	return nil
}
