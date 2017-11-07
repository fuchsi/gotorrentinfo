# gotorrentinfo
gotorrentinfo parses .torrent files and displays information about the torrent and the files it references.

## Installation
```bash
$ dep ensure
$ go install
```

## Usage
```
Usage: gotorrentinfo [-defhnv] <filename>
 -d, --detailed    Show detailed information about the files
 -e, --everything  Print everything about the torrent
 -f, --files       Show files within the torrent
 -h, --help        Show this help message and exit
 -n, --nocolors    No ANSI colour
 -v, --version     Print version and quit
```