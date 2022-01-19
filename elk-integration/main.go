package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

type Logger struct {
	log *zerolog.Logger
}

func newLogger(debug bool) *Logger {
	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &Logger{log: &logger}
}

func processHandler(w http.ResponseWriter, r *http.Request) {

	value := r.URL.Query().Get("age")
	age, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Println("can't convert string to integer")
		return
	}
	
	writeResponse(w, http.StatusOK, "You entered age as "+strconv.Itoa(int(age)))
}


func Serve(l *Logger) error {

	http.HandleFunc("/data", sillyLogger(processHandler, l)) // /data?age=
	l.log.Debug().Msgf("Getting into first server...")
	return http.ListenAndServe("0.0.0.0:8080", nil)
}

func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}

// some silly middleware
func sillyLogger(next http.HandlerFunc, l *Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		l.log.Info().Msgf("Got new request from user")
		next(rw, r)
	}
}

func main() {
	debugMode := flag.Bool("debug", false, "show debug messages")
	flag.Parse()
	logger := newLogger(*debugMode)
	
	logger.log.Info().Msgf("Server is starting...")

	go func() {
		logger.log.Debug().Msgf("Getting into second server...")
		http.HandleFunc("/year", sillyLogger(func(w http.ResponseWriter, r *http.Request) {
			value := r.URL.Query().Get("year")

			w.Write([]byte("Replay from second HTTP server\n"))
			w.Write([]byte("You entered year: " + value))
		}, logger))
		http.ListenAndServe(":8081", nil)
	}()

	
	if err := Serve(logger); err != nil {
		logger.log.Error().Msgf("The app is closing due to some error: %v", err)
	}

	logger.log.Debug().Msgf("The application is closing")
}
