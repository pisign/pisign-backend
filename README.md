# pisign-backend

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![codecov](https://codecov.io/gh/pisign/pisign-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/pisign/pisign-backend)
### Development

Copy the script `post-commit` into your git hooks folder, 

`./.git/hooks/`

This makes it so that the specfile in `spec` gets automatically updates whenever we update the 
external API types defined in the `types/` folder.

### Testing 

When testing a specific package, inside the package's directory, use 

`go test -coverprofile /tmp/cp.out`

to run tests and show test coverage. To view the coverage in a browser, use 

`go tool cover -html=/tmp/cp.out` 

to visualize what parts of the code are not being tested.

To test all files, run 

`go test ./...`

from the root directory to test everything inside the repo. 
