# Sumo Logic Server TCP

## Description

The goal of this project is the creation of a TCP server in Golang.

The server should receive calls and execute commands based on the information received from the client.

Golang was used in the creation of this project, taking advantage of the language's features for creating goroutines/channels/contexts to better manage and control the different opened threads.

I hope you like the solutions I applied, and I would be very happy to receive any feedback :)

## Requested Features in the Project File

- [x] TCP Server listening on port 3000.
- [x] TCP Server receiving requests.
- [x] Receiving messages in the format of Command and Timeout.
- [x] Adding a Timeout to cancel requests if the executed command exceeds the stipulated limit.
- [x] 0 Timeout for requests without a stipulated timeout.
- [x] TCP Server returning values (duration, command, executed_at, etc.) as requested.
- [x] Capturing STDOUT of the command and adding it to the response field Output.
- [x] Requests being processed in parallel.
- [x] The solution needs to have unit tests.

## Features / Improvements Added by Me

- [x] Added a Client to facilitate evaluator testing.
- [x] Added an option to create a TCP server on different ports using the -p flag.
- [x] Create a build script to generate binary files of the application.
- [x] Create a script to run integration_tests to test paralelism and real world scenarios.
- [x] Create a diagram flow of the application to facilitate presentation and understanding for the evaluators.
- [x] Create a Docker environment to eliminate any chance of incompatibility on the evaluators' computers.
- [x] Create a configurable limit to accept a predefined limit of concurrent requests.
- [x] Add documentation to explain the division of folders and explanations of modules.
- [ ] Create a list of "next improvements" to provide evaluators with a direction on how the project could evolve.

## Layout

![alt text](https://iili.io/2IvxhHg.png)

```tree
├── README.md
├── .gitignore
├── docker-compose.yaml
├── Dockerfile
├── main.go
├── build
│   ├── sumologic_server
│   └── sumologic_server.exe
├── cmd
│   ├── server.go
│   ├── client.go
│   ├── await.go
│   └── cmd.go
├── common
│   └── common.go
├── scripts
│   ├── build.sh
│   └── integration_test.sh.sh
└── server
    ├── models
    │   ├── taskRequest.go
    │   └── taskResult.go
    ├── listener.go
    ├── network.go
    └── server.go
```

A brief description of the layout:

* `README.md` is a detailed description of the project.
* `docker-compose.yaml` a file to assist in creating the local environment and reduce compatibility issues.
* `Dockerfile` template to create Docker images and facilitate the containerization of the application.
* `build` is to hold build outputs.
* `cmd` is where the files that manage the application's commands are located.
* `common` is where auxiliary functions or language functions are wrapped in interfaces to facilitate unit testing
* `scripts` contains scripts to build and test the project.
* `server` where all the TCP server logic is present.

## Tests

Both unit and integration tests were conducted.

The unit tests were primarily carried out within the "Server" component, which currently has over 95% coverage. You can check the current code coverage by running this command in your terminal:
 ```bash
   go test ./... -cover
```

The other components, such as CMD and Common, only implement the functionalities of the Server, and for this reason, they are only covered in the integration tests.
A script was created to facilitate unit testing (particularly to test concurrency and timeout). You can run the script with the following command:
 ```bash
    cd ./scripts/
   ./integration_tests.sh
```
This is an example of all integration tests successfully completed at the moment:
![alt text](https://iili.io/2IvcPoX.png)

## Technologies Used

- **Language**: Go 1.22.6

## Installation

To install and configure the project, follow the steps below:

1. Clone the repository to your local machine:
 ```bash
   git clone https://github.com/hriqueXimenes/sumo_logic_server.git
```
2. Install all dependencies:
```bash
  go get
```
3. Run Server TCP
```bash
go run main.go server -p 3000 -a localhost
