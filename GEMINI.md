# Project: birdstats

## Project Overview

`birdstats` is a command-line tool written in Go for analyzing eBird data from CSV files. It processes a user-provided eBird data file, calculates various statistics, and prints a summary to the console.

The tool calculates and displays:
- Per-species statistics, including common name, number of submissions, total count, and the number of multimedia assets.
- Overall statistics, such as the total number of species, total submissions, total distance traveled, total observation duration, and averages for distance and duration.

The project is organized into two main parts:
- The `main` package (`birdstats.go`), which contains the primary application logic for calculating and printing stats.
- The `ebird` package (`ebird/ebird.go`), which provides the functionality for reading and parsing the eBird CSV data file.

## Building and Running

### Prerequisites

- Go 1.22.0 or later.
- An eBird data file in CSV format. You can download your own data from the eBird website.

### Building the Project

To build the `birdstats` executable, run the following command from the project root directory:

```sh
go build
```

This will create a `birdstats` executable in the project directory.

### Running the Project

To run the `birdstats` tool, you need to provide the path to your eBird CSV file as a command-line argument.

Using `go run`:
```sh
go run birdstats.go path/to/your/ebird_data.csv
```

Using the compiled executable:
```sh
./birdstats path/to/your/ebird_data.csv
```

## Development Conventions

- **Language:** The project is written in Go.
- **Formatting:** The code follows standard Go formatting conventions.
- **Structure:** The project is divided into a `main` application and a reusable `ebird` package.
  - `birdstats.go`: The entry point and main logic of the application.
  - `ebird/ebird.go`: Handles the parsing of eBird CSV data.
- **Dependencies:** The project has no external dependencies outside of the Go standard library.
