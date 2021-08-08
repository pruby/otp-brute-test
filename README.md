# otp-brute-test

A test server for brute-forcing One-Time Passwords. This maintains a lookup table of valid codes at any point in time
so that it can respond extremely quickly.

This generates a new, random OTP key on each start-up. Note that while this simulates the time-rolling properties of TOTP,
the underlying algorithm is actually HOTP with a counter. This simplifies the verification process while behaving in exactly
the same way from the point of view of attempting to guess codes.

## Installation and Running

This needs to be compiled with a recent version of the Go runtime.

1. Download and install the Go runtime from https://golang.org/
2. Check out this repository using `git clone https://github.com/pruby/otp-brute-test`.
3. Enter the checked out directory, and install dependencies by running `go mod download`.
4. Build the server by running `go build server.go`. This will create a binary called "server", or "server.exe" on Windows.
5. Run the binary using `./server`

## Options

You can control various options using environment variables.

* `OTP_DIGITS` - number of digits to use in the code.
* `OTP_PERIOD` - number of seconds between rotating a code.
* `OTP_WINDOW` - number of codes valid at any given time.
* `OTP_SERVER_ADDR` - address for server to bind to (`host:port`)

## Guessing codes

The server binds by default to port 3000, and can be invoked with a URL like the following:

`http://127.0.0.1:3000/check?code=123456`

It will return code 400 with a body of "BAD" for invalid requests to this endpoint (i.e. missing a code), code 401 with body "NO" for an incorrect guess and
code 200 with body "OK" for a correct one. It otherwise does nothing with the code.

