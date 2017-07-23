## escape-client package

Create a package

### Synopsis


Create a package

```
escape-client package [flags]
```

### Options

```
  -d, --deployment string    Deployment name (default "<release name>")
  -e, --environment string   The logical environment to target (default "dev")
  -f, --force                Overwrite output file if it exists
  -h, --help                 help for package
  -i, --input string         The location of the Escape plan. (default "escape.yml")
  -s, --state string         Location of the Escape state file (default "escape_state.json")
  -u, --uber                 Build an uber package containing all dependencies
```

### Options inherited from parent commands

```
      --collapse-logs    Collapse logs (Default: true) (default true)
  -c, --config string    Global Escape configuration file (default "~/.escape_config")
  -l, --level string     Log level: debug, info, warn, error (default "info")
      --profile string   Configuration profile (default "default")
```

### SEE ALSO
* [escape-client](escape-client.md)	 - Package and deployment manager

###### Auto generated by spf13/cobra on 20-Jul-2017