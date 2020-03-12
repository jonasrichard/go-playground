# go-playground
Go proof of concept projects

### Build for debug

`go build -gcflags='all=-N -l' -tags nopkcs11 -ldflags='-linkmode internal'`

Then you can use delve for debugging.

### Vim convenience

You need to install `gotags` in order to generate Go tag files.

`brew install gotags`
