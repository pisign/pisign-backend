# pisign-backend
[![Build Status](https://travis-ci.org/pisign/pisign-backend.svg?branch=master)](https://travis-ci.org/pisign/pisign-backend)
[![codecov](https://codecov.io/gh/pisign/pisign-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/pisign/pisign-backend)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

### Development

#### Interacting with `main.go`
The main entry point into the program is `main.go` in the root directory. You can either build an executable called
`main` using `go build main.go` or run it directly using `go run main.go`. With both methods, there are currently
two subcommands that can be used:

##### Linux / OSX
1) `./main run [--port=9000]`: This runs the entire backend server. The `-p` or `--port` flag can be
added to specify the port of the server

2) `./main run <name>`: This creates a new api with the specified name. It modifies all necessary files
and creates new skeleton code inside the `api/<name>` folder.

##### Windows
1) If you installed from the package, run the program via powershell with `.\<PROGRAM_NAME>.exe`

*NOTE* there is a bug where you have to press control + c if the program does not find the `.\assets\` folder. Once you get that error message, if you press control + c one time, the program continues like normal. 

### Testing 

When testing a specific package, inside the package's directory, use 

`go test -coverprofile /tmp/cp.out`

to run tests and show test coverage. To view the coverage in a browser, use 

`go tool cover -html=/tmp/cp.out` 

to visualize what parts of the code are not being tested.

To test all files, run 

`go test ./...`

from the root directory to test everything inside the repo. 
