package main

import (
	"flag"
	"fmt"
	"github.com/NoizeMe/prometheus-covid-exporter/covid"
	"github.com/jtaczanowski/go-scheduler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	defaultPort = 8080

	defaultInitialDelay = 0 * time.Second
	defaultRefreshDelay = 5 * time.Minute
)

var (
	diseaseCases = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "cases",
		Help:      "The number of cases in Germany",
	}, []string{"state"})
	diseaseDeaths = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "deaths",
		Help:      "The number of deaths in relation with COVID in Germany",
	}, []string{"state"})

	vaccinationTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "total",
		Help:      "The number of people that need to get vaccinated in total",
	}, []string{"state"})
	vaccinationVaccinated = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "vaccinated",
		Help:      "The number of people that are already vaccinated",
	}, []string{"state"})
	vaccinationQuote = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "quote",
		Help:      "The quote of people that have been vaccinated",
	}, []string{"state"})
)

func refreshDiseaseData() {
	governmentData, err := covid.GetGovernmentData("de", false)
	if err != nil {
		panic(err)
	}

	for _, data := range governmentData {
		province := strings.ReplaceAll(data.Province, "\n", "")
		if province == "Total" {
			province = "Germany"
		}

		diseaseCases.WithLabelValues(province).Set(float64(data.Cases))
		diseaseDeaths.WithLabelValues(province).Set(float64(data.Deaths))
	}
}

func refreshVaccinationData() {
	vaccinationData, err := covid.GetVaccinationData()
	if err != nil {
		panic(err)
	}

	vaccinationTotal.WithLabelValues("Germany").Set(float64(vaccinationData.Total))
	vaccinationVaccinated.WithLabelValues("Germany").Set(float64(vaccinationData.Vaccinated))
	vaccinationQuote.WithLabelValues("Germany").Set(vaccinationData.Quote)

	for state, data := range vaccinationData.States {
		vaccinationTotal.WithLabelValues(state).Set(float64(data.Total))
		vaccinationVaccinated.WithLabelValues(state).Set(float64(data.Vaccinated))
		vaccinationQuote.WithLabelValues(state).Set(data.Quote)
	}
}

func main() {
	var (
		port         uint
		refreshDelay time.Duration
	)

	flag.UintVar(&port, "port", defaultPort, "the port at which the exporter should listen on.")
	flag.DurationVar(&refreshDelay, "delay", defaultRefreshDelay, "the delay between data fetching round trips.")
	flag.Parse()

	log.Println("Starting prometheus COVID exporter")
	log.Printf("  - Port: %d\n", port)
	log.Printf("  - Delay: %s\n", refreshDelay.String())

	refreshDiseaseData()
	scheduler.RunTaskAtInterval(refreshDiseaseData, refreshDelay, defaultInitialDelay)

	refreshVaccinationData()
	scheduler.RunTaskAtInterval(refreshVaccinationData, refreshDelay, defaultInitialDelay)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		diseaseCases,
		diseaseDeaths,
		vaccinationTotal,
		vaccinationVaccinated,
	)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
