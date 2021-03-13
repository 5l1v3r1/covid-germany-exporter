# Prometheus COVID-19 Exporter for Germany

[![Build Status](https://img.shields.io/github/workflow/status/jangraefen/covid-germany-exporter/Build?logo=GitHub)](https://github.com/jangraefen/covid-germany-exporter/actions?query=workflow:Build)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/jangraefen/covid-germany-exporter)](https://pkg.go.dev/mod/github.com/jangraefen/covid-germany-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/jangraefen/covid-germany-exporter)](https://goreportcard.com/report/github.com/jangraefen/covid-germany-exporter)

A metric exporter for Prometheus that scraps [disease.sh](https://disease.sh) and
the [RKI Vaccination Data](https://rki-vaccination-data.vercel.app) APIs for new data. By default, data is retrieved
every five minutes. Usually data is only updated once a day, so the fetching interval can be overwritten by passing a
command-line argument.

## Exported metrics

- `covid_disease_cases{state}`: The number of total cases in the given state.
- `covid_disease_deaths{state}`: The number of total deaths in the given state.
- `covid_disease_case_previous_day_change{state}`: The number of new cases since the previous day.
- `covid_disease_cases_per_hundred_thousand{state}`: The number of cases per 10,000 inhabitants.
- `covid_disease_seven_day_cases_per_hundred_thousand{state}`: The number of new cases in the last week per 10,000
  inhabitants.
- `covid_vaccination_total{state}`: The number of people that need to get vaccinated. This is basically the number of
  inhabitants for the given state.
- `covid_vaccination_quote{state}`: The quote of people that already have been vaccinated.
- `covid_vaccination_per_1000_inhabitants{state}`: The number of vaccinations per 1000 citizens.
- `covid_vaccination_vaccinated{state}`: The number of people that already have been vaccinated.
- `covid_vaccination_difference_to_previous_day{state}`: The number of vaccinations performed during the last 24h
  period.
- `covid_vaccination_second_wave_vaccinated{state}`: The number of people that are already vaccinated during the 2nd
  wave.
- `covid_vaccination_second_wave_difference_to_previous_day{state}`: The number of vaccinations performed during the
  last 24h period during the 2nd wave.

To get nation-wide data, it is possible to pass `state = "Germany"` to any of the above metrics.

## Command-Line Options

The exporter support a few command-line arguments to customize its behaviour to an extent:

- `-port <uint>` (defaults to `8080`): The port at which the exporter should listen on.
- `-delay <duration>` (defaults to `5m`): The delay between data fetching round trips.
