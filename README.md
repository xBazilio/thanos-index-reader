# thanos-index-reader
Tool for reading index from thanos buckets

```bash
usage: indexreader --objstore.config-file=OBJSTORE.CONFIG-FILE --data-dir=DATA-DIR --bucket-ulid=BUCKET-ULID [<flags>]

Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
      --objstore.config-file=OBJSTORE.CONFIG-FILE  
                                 Path to YAML file that contains object store configuration. See format details: https://thanos.io/tip/thanos/storage.md/#configuration
      --data-dir=DATA-DIR        Data dir for storing downloaded from storage data
  -b, --bucket-ulid=BUCKET-ULID  Store bucket ULID as a string
  -l, --label-name=LABEL-NAME    If provided, all values for given label name will be printed, otherwise will print all label names
  -s, --show-stat                Gather stat such as labels value count and size of all values in bytes
      --version                  Show application version.
```
