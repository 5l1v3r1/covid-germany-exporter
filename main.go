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
	diseaseCasePreviousDayChange = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "case_previous_day_change",
		Help:      "The number of new cases since the previous day",
	}, []string{"state"})
	diseaseCasesPerHundredThousand = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "cases_per_hundred_thousand",
		Help:      "The number of cases per 10,000 inhabitants",
	}, []string{"state"})
	diseaseSevenDayCasesPerHundredThousand = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "disease",
		Name:      "seven_day_cases_per_hundred_thousand",
		Help:      "The number of new cases in the last week per 10,000 inhabitants",
	}, []string{"state"})

	vaccinationTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "total",
		Help:      "The number of people that need to get vaccinated in total",
	}, []string{"state"})
	vaccinationQuote = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "quote",
		Help:      "The quote of people that have been vaccinated",
	}, []string{"state"})
	vaccinationPer1000Inhabitants = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "per_1000_inhabitants",
		Help:      "The number of vaccinations per 1000 citizens",
	}, []string{"state"})
	vaccinationVaccinated = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "vaccinated",
		Help:      "The number of people that are already vaccinated",
	}, []string{"state"})
	vaccinationDifferenceToPreviousDay = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "difference_to_previous_day",
		Help:      "The number of vaccinations performed during the last 24h period",
	}, []string{"state"})
	vaccinationVaccinatedSecondWave = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "second_wave_vaccinated",
		Help:      "The number of people that are already vaccinated during the 2nd wave",
	}, []string{"state"})
	vaccinationDifferenceToPreviousDaySecondWave = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "covid",
		Subsystem: "vaccination",
		Name:      "second_wave_difference_to_previous_day",
		Help:      "The number of vaccinations performed during the last 24h period during the 2nd wave",
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
		diseaseCasePreviousDayChange.WithLabelValues(province).Set(float64(data.CasePreviousDayChange))
		diseaseCasesPerHundredThousand.WithLabelValues(province).Set(float64(data.CasesPerHundredThousand))
		diseaseSevenDayCasesPerHundredThousand.WithLabelValues(province).Set(float64(data.SevenDayCasesPerHundredThousand))
	}
}

func refreshVaccinationData() {
	vaccinationData, err := covid.GetVaccinationData()
	if err != nil {
		panic(err)
	}

	vaccinationTotal.WithLabelValues("Germany").Set(float64(vaccinationData.Total))
	vaccinationQuote.WithLabelValues("Germany").Set(vaccinationData.Quote)
	vaccinationPer1000Inhabitants.WithLabelValues("Germany").Set(vaccinationData.VaccinationsPer1000Inhabitants)

	vaccinationVaccinated.WithLabelValues("Germany").Set(float64(vaccinationData.Vaccinated))
	vaccinationDifferenceToPreviousDay.WithLabelValues("Germany").Set(float64(vaccinationData.DifferenceToThePreviousDay))
	vaccinationVaccinatedSecondWave.WithLabelValues("Germany").Set(float64(vaccinationData.SecondVaccination.Vaccinated))
	vaccinationDifferenceToPreviousDaySecondWave.WithLabelValues("Germany").Set(float64(vaccinationData.SecondVaccination.DifferenceToThePreviousDay))

	for state, data := range vaccinationData.States {
		vaccinationTotal.WithLabelValues(state).Set(float64(data.Total))
		vaccinationQuote.WithLabelValues(state).Set(data.Quote)
		vaccinationPer1000Inhabitants.WithLabelValues(state).Set(data.VaccinationsPer1000Inhabitants)

		vaccinationVaccinated.WithLabelValues(state).Set(float64(data.Vaccinated))
		vaccinationDifferenceToPreviousDay.WithLabelValues(state).Set(float64(data.DifferenceToThePreviousDay))
		vaccinationVaccinatedSecondWave.WithLabelValues(state).Set(float64(data.SecondVaccination.Vaccinated))
		vaccinationDifferenceToPreviousDaySecondWave.WithLabelValues(state).Set(float64(data.SecondVaccination.DifferenceToThePreviousDay))
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
		diseaseCasePreviousDayChange,
		diseaseCasesPerHundredThousand,
		diseaseSevenDayCasesPerHundredThousand,
		vaccinationTotal,
		vaccinationVaccinated,
		vaccinationQuote,
	)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
