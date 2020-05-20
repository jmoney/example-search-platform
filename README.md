# example-search-platform

## Building

```bash
make -B
```

This will build each module individually and create the necessary docker images out of them.  The key is to make sure you run with `-B`.  This tells `make` to ignore `freshness` and build each step anyways.