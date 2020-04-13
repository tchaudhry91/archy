ZSH Archaeologist
=================

Simple tool to operate on ZSH_HISTORY files. Grab the CLI for your system from the releases.

```
Usage:
  archy [command]

Available Commands:
  help        Help about any command
  import      import allows you to add external entries to your base zsh history

Flags:
      --baseHistoryFile string   base history file to operate on (default "/Users/tchaudhr/.zsh_history")
      --config string            config file (default is $HOME/.archy.yaml)
  -h, --help                     help for archy
      --hostname string          Override the hostname value for entries from this machine (default "Tanmays-MacBook-Pro.local")
      --token string             the token to communicate with the remote service

Use "archy [command] --help" for more information about a command.
```

### Import history from another file

`archy import --maxEntries 6000 local f1 f2 f3 f4`

This will import entries from files `f1 f2 f3 f4` into your base history file, keeping only 6000 of the newest entries after merging.