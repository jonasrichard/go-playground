# go-playground
Go proof of concept projects

### Build for debug

`go build -gcflags='all=-N -l' -tags nopkcs11 -ldflags='-linkmode internal'`

Then you can use delve for debugging.
