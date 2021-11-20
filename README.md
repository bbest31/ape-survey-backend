# ApeSurvey Backend
Monolithic backend server for the ApeSurvey application. This server exposes an HTTPS API for the ApeSurvey React client to interact with in order to function. 

## Setup & How to Run

Be sure to have Go installed on your machine and have set your GOPATH

### Install Dependencies

In the root folder run:

`go get .`


1. Build the service

```
go build
```
2. Run the executable

```
./ape-survey-backend
```
## Testing

Installing the Go extension for VS Code allows for a built in GUI in the test files to run specific functions and files.

Use the go cli in order to see other options like adding test coverage analysis and benchmarking.

Can also add `-cover` flag for coverage measurement
Use `go test -coverprofile <filename>.out` to output cover analysis. If the test coverage is not high enough you may want to run `go tool cover -html=<filename>.out` to view visual representation of code coverage and see which portions of the code you need to write tests for.


To run test files in a specific directory

```
go test
```

To run a specific test file

```
go test fileName
```

## Environment variables (local setup)

Below is the list of environment variables that need to be set in order for the service to function.
.env file
```
PORT=
AUTH0_API=<protocol://domain>
AUTH0_APP_DOMAIN=<protocol://domain/>
CLIENT_URL=
GOOGLE_APPLICATION_CREDENTIALS=<relative_path>
```

In environment set the following variable
`GOOGLE_APPLICATION_CREDENTIALS=path to key`