# Prometheus COVID-19 Exporter

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
- `covid_vaccination_vaccinated{state}`: The number of people that already have been vaccinated.
- `covid_vaccination_quote{state}`: The quote of people that already have been vaccinated.

To get nation-wide data, it is possible to pass `state = "Germany"` to any of the above metrics.

## Command-Line Options

The exporter support a few command-line arguments to customize its behaviour to an extent:

- `-port <uint>` (defaults to `8080`): The port at which the exporter should listen on.
- `-delay <duration>` (defaults to `5m`): The delay between data fetching round trips.
