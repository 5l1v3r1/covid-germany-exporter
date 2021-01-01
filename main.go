package main

import (
	"flag"
	"fmt"
	"github.com/NoizeMe/prometheus-covid-exporter/pkg/covid"
	"github.com/NoizeMe/prometheus-covid-exporter/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	diseaseCases = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "cases",
		Help:      "The number of cases in Germany",
	}, []string{"state"})
	diseaseDeaths = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "deaths",
		Help:      "The number of deaths in relation with COVID in Germany",
	}, []string{"state"})

	vaccinationTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "total",
		Help:      "The number of people that need to get vaccinated in total",
	}, []string{"state"})
	vaccinationVaccinated = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "vaccinated",
		Help:      "The number of people that are already vaccinated",
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
			province = "combined"
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

	vaccinationTotal.WithLabelValues("combined").Set(float64(vaccinationData.Total))
	vaccinationVaccinated.WithLabelValues("combined").Set(float64(vaccinationData.Vaccinated))

	for state, data := range vaccinationData.States {
		vaccinationTotal.WithLabelValues(state).Set(float64(data.Total))
		vaccinationVaccinated.WithLabelValues(state).Set(float64(data.Vaccinated))
	}
}

func main() {
	port := *flag.Uint("port", 0, "The port at which the exporter should listen on")
	flag.Parse()

	if port == 0 {
		if envPort, hasEnvPort := os.LookupEnv("COVID_EXPORTER_PORT"); hasEnvPort {
			if parsedEnvPort, err := strconv.ParseUint(envPort, 10, 32); err != nil {
				port = uint(parsedEnvPort)
			}
		}
	}

	if port == 0 {
		port = 8080
	}

	refreshDiseaseData()
	refreshVaccinationData()

	refreshJob := utils.CreatePeriodic(2*time.Minute, func() {
		refreshDiseaseData()
		refreshVaccinationData()
	})
	refreshJob.Start()

	registry := prometheus.NewRegistry()
	registry.MustRegister(diseaseCases)
	registry.MustRegister(diseaseDeaths)
	registry.MustRegister(vaccinationTotal)
	registry.MustRegister(vaccinationVaccinated)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
