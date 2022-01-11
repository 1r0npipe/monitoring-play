package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Namespace   = "eugenespace"
	LabelMethod = "method"
	LabelStatus = "status"
)

type app struct {
	ageHistogram         *prometheus.HistogramVec
	requestsCounter      prometheus.Counter
	lastInsertedAgeGauge prometheus.Gauge
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {

	value := r.URL.Query().Get("age")
	age, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Println("can't convert string to integer")
		return
	}

	defer func() {
		a.ageHistogram.With(prometheus.Labels{LabelMethod: r.Method}).Observe(float64(age))
		a.requestsCounter.Inc()
		a.lastInsertedAgeGauge.Set(float64(age))
	}()

	writeResponse(w, http.StatusOK, "You entered age as "+strconv.Itoa(int(age)))
}

func (a *app) Init() error {
	// prometheus type: histogram
	a.ageHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "Age",
		Help:      "The distribution of the input age",
		Buckets:   []float64{10, 20, 30, 50, 70, 100},
	}, []string{LabelMethod})

	// prometheus type: counter
	a.requestsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "number_requests",
		Help:      "The number of connected requests",
	})

	// prometheus type: gauge
	a.lastInsertedAgeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "last_input_age",
		Help:      "Last inserted age",
	})

	prometheus.MustRegister(a.ageHistogram)
	prometheus.MustRegister(a.requestsCounter)
	prometheus.MustRegister(a.lastInsertedAgeGauge)

	return nil
}

func (a *app) Serve() error {

	http.HandleFunc("/data", sillyLogger(a.processHandler)) // /data?age=
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	return http.ListenAndServe("0.0.0.0:8080", nil)
}

func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}

// some silly middleware
func sillyLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Got request from external user")
		next(rw, r)
	}
}

func main() {

	myApp := app{}

	go func() {
		
		yearGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "last_input_year_birth",
			Help:      "Last inserted year of birth",
		})
		prometheus.MustRegister(yearGauge)

		http.HandleFunc("/year", sillyLogger(func(w http.ResponseWriter, r *http.Request) {
			value := r.URL.Query().Get("year")
			year, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				log.Println("can't convert string to integer")
				return
			}
			w.Write([]byte("Replay from second HTTP server\n"))
			w.Write([]byte("You entered year: " + value))
			defer func() {
				yearGauge.Set(float64(year))
				myApp.requestsCounter.Inc()
			}()
		}))
		http.ListenAndServe(":8081", nil)
	}()

	if err := myApp.Init(); err != nil {
		log.Fatal(err)
	}

	if err := myApp.Serve(); err != nil {
		log.Fatal(err)
	}
}
