# energy-prometheus-exporter

A Prometheus exporter for home energy metrics. It supports:
* [Homewizard](https://www.homewizard.com/) devices, based on the [API v2](https://api-documentation.homewizard.com/docs/category/api-v2).
* [Solaredge](https://www.solaredge.com/) inverter with/without battery.

## Building

* The `make` command (e.g. [GNU make](https://www.gnu.org/software/make/manual/make.html)).
* The [Golang toolchain](https://golang.org/doc/install) (version 1.24 or later).

In a shell, execute: `make` (or `make build`).

The build artifacts can be cleaned by using: `make clean`.

## Configuration

The configuration can be provided via command flags or environment variables.
Use the `--help` command line flag to see the usage.
