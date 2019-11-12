# Browserbeat

Welcome to Browserbeat.

Browserbeat is a lightweight passive solution for web traffic monitoring. Browserbeat gives sysadmins the ability to monitor web traffic on managed computers without an HTTP proxy, packet capture, or DNS server logs.

## Browserbeat compared to:

### DNS server logs

* Data obtained by Browserbeat is much less noisy than DNS traffic since traffic created by system services and other protocols will show up in DNS server logs.
* You will know if the browser has visited an IP address directly.
* You'll know the user who made the request not just the client IP
* You'll know the web browser used to make the request
* You'll know the operating system used to make the request

### HTTP proxy

* No need to distribute custom certs to monitor HTTPS traffic
* You'll know the user who made the request not just the client IP
* If Browserbeat or the output fails, the user's browsing is not interrupted like if a proxy server goes down
* Less complexity on your network
 
### Packet Sniffing/Capture

* You'll know the user who made the request not just the client IP
* Less complexity on your network
* Depending on the method used for packet sniffing, the user's browsing is not interrupted if a component fails

## Features

* Know the user who made the request
* Know the IP of the computer
* Know the hostname of the computer
* Know the requested hostname
* Know the URL requested
* Know the title of the website
* Know the date & time the request was made
* Know the host OS
* Know the browser that made the request
* Cross-platform support: Windows, macOS, and Linux
* Cross-browser support: see list below
* Output data to all of the standard Elastic Beat outputs

### Browser Support

* :white_check_mark: Chrome - Done
* :white_check_mark: Firefox - Done
* :white_check_mark: Safari (macOS) - Done
* :white_check_mark: Chromium - Done
* :white_check_mark: Chrome Canary - Done
    * :white_check_mark: Chrome Beta (linux) - Done
    * :white_check_mark: Chrome Dev (linux) - Done
* :white_check_mark: Vivaldi - Done
* :white_check_mark: Opera - Done
* :yellow_heart: K-Meleon - Haven't investigated yet
* :yellow_heart: Brave - Haven't investigated yet
* :yellow_heart: Epic - Haven't investigated yet
* :sos: Edge - M$ is too cool and uses an ESE database anyone know of an ESE DB library for go?
* :sos: IE 11 - M$ is too cool and uses an ESE database anyone know of an ESE DB library for go?

Feel free to suggest more browsers.

# Development

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/MelonSmasher/browserbeat`

## Getting Started with Browserbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Browserbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Browserbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/MelonSmasher/browserbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Browserbeat run the command below. This will generate a binary
in the same directory with the name browserbeat.

```
make
```


### Run

To run Browserbeat with debugging output enabled, run:

```
./browserbeat -c browserbeat.yml -e -d "*"
```


### Test

To test Browserbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Browserbeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Browserbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/MelonSmasher/browserbeat
git clone https://github.com/MelonSmasher/browserbeat ${GOPATH}/src/github.com/MelonSmasher/browserbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
