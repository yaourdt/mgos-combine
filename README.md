# mgos-combine

A tool to combine all parts of a Mongoose OS firmware ZIP-file into a single binary.

## Install

### Local install

Download the latest binary from the [release page](https://github.com/yaourdt/mgos-combine/releases).
There are two options, the Linux binary (named `mgos-combine-ubuntu` as it is
compiled on ubuntu), and a version for Mac OS named `mgos-combine-macos`.

Move the downloaded binary to a folder within you path variable, make it
executable, and you are good to go:

```
sudo mv mgos-combine-ubuntu /usr/local/bin/mgos-combine
sudo chmod a+x /usr/local/bin/mgos-combine
```

### Install using Docker

The docker image is mainly provided for use in CI/CD pipelines. It is located at
`docker.pkg.github.com/yaourdt/mgos-combine/mgos-combine`, starting from version 0.2.2 upwards,
tag names correspond to release versions.

Run it as

```
docker run -v /path/to/host/fw/dir:/build docker.pkg.github.com/yaourdt/mgos-combine/mgos-combine -o /build/out.bin
```

## Usage

Usage: `mgos-combine [options]`

```
  -i, --input    Path to the firmware zip file (default: ./build/fw.zip)
  -o, --output   Name of the output file (default: output.bin)
  -s, --size     Output file size in KB (default: 4096)
  -f, --force    Force writing to an output file that is too small
  -h, --help     Show this help
  -v, --version  Prints current version and exits. All other options will be ignored
```

## Compile the software yourself

With a working go environment on your machine, just `git clone` this repository,
run `go get` to install missing libraries, and run `go build -o mgos-combine .`
to compile.

## Acknowledgments
Parts of the initial code were developed by [ert](https://github.com/ertugrul-sevgili)
as part of a coding exercise.
