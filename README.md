[![Go Test](https://github.com/amunchet/autopuller-go/actions/workflows/go-tests.yml/badge.svg)](https://github.com/amunchet/autopuller-go/actions/workflows/go-tests.yml)
# Autopuller-Go
Rewrite of the Autopuller project in Go.  This moves away from the Docker implementation and handles things simply in Go.

## Overall process
1. Check if there is a new commit to the master branch of a given repo
2. If there is, check if the tests have passed (Github actions)
3. If they have, then git pull
4. If that succeeds, then execute a docker compose rebuild and restart
5. Go back to sleep for 60 seconds
6. Repeat
