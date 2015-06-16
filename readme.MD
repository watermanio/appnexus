AppNexus REST API Go Client
===========================
*Connect to the AppNexus REST API through your Golang application*

**Important** This is very much beta and subject to breaking changes while under development. Do not use in production just yet!

Of the published [AppNexus API Services](https://wiki.appnexus.com/display/adnexusdocumentation/API+Services) listed, this package covers the following:

* Auth API Service [Docs](https://wiki.appnexus.com/display/adnexusdocumentation/Auth+API+Service)
* Member API Service [Docs](https://wiki.appnexus.com/display/adnexusdocumentation/Member+Service)
* Segment Service [Docs](https://wiki.appnexus.com/display/adnexusdocumentation/Segment+Service)

Support for the remaining services should follow - pull requests welcome :)

Getting started
--------------
```Bash
$ cd $GOPATH
$ go get github.com/adwww/appnexus
$ go get -u all
```

Be sure to run the tests with `go test` and have a look at the `examples/main.go` file for a usage demonstration.
