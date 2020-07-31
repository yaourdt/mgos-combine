# mgos-combine

A tool to combine all parts of a Mongoose OS firmware ZIP-file into a single binary.

## Usage

Usage: `mgos-combine [options]`

```
  -m, --manifest Path to the JSON manifest file (default: manifest.json)
  -o, --output   Name of the output file (default: output.bin)
  -s, --size     Output file size in KB (default: 4096)
  -f, --force    Force writing to an output file that is too small
  -h, --help     Show this help
  -v, --version  Prints current version and exits. All other options will be ignored
```

## Acknowledgments
Parts of the initial code were developed by [ert](https://github.com/ertugrul-sevgili)
as part of a coding exercise.
