# Sumo Logic Server TCP

## Description

This project is part of the "take-home" phase for SumoLogic. In my attempt to join the software engineering team, I was tasked with creating a TCP server in Golang. I aimed to create a simple, testable, and easily readable project for the evaluators.

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
- [ ] Create a design for the solution to facilitate presentation and understanding for the evaluators.
- [x] Create a Docker environment to eliminate any chance of incompatibility on the evaluators' computers.
- [ ] Create a configurable limit to accept a predefined limit of concurrent requests.
- [ ] Add documentation to explain the division of folders and explanations of modules.
- [ ] Create a list of "next improvements" to provide evaluators with a direction on how the project could evolve.

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
