# log-splitter

Go-based, concurrent log splitter: reads an input log file and splits error logs from standard logs into separate stderr/stdout files.

*Example*

input `file.log`:

```console
INFO   log line one
WARN   log line two
ERROR  log line three
INFO   log line four
```

output:

* `file.log.stderr`:
   ```console
   ERROR  log line three
   ```
* `file.log.stdout`:
   ```console
   INFO   log line one
   WARN   log line two
   INFO   log line four
   ```

## Build and Run

From project root:

```bash
make binary && ./build/main <inputfile>
```

*Example*:

```bash
./build/main testfiles/test2.log
```
