ZSH Archaeologist
=================

Simple tool to operate on ZSH_HISTORY files. Grab the CLI for your system from the releases.

```
A zsh history manipulator

Usage:
  archy [command]

Available Commands:
  export      exports your command history
  help        Help about any command
  import      import allows you to add external entries to your base zsh history
  login       login to the remote service
  register    register to the remote service

Flags:
      --baseHistoryFile string   base history file to operate on (default "$HOME/.zsh_history")
      --config string            config file (default is $HOME/.archy.yaml)
  -h, --help                     help for archy
      --hostname string          Override the hostname value for entries from this machine (default "QEAirArch")
      --remote string            Address of the remote service to contact (default "https://archy.tux-sudo.com")
      --token string             the token to communicate with the remote service

Use "archy [command] --help" for more information about a command.
```

### Import history from another file

`archy import --maxEntries 6000 local f1 f2 f3 f4`

This will import entries from files `f1 f2 f3 f4` into your base history file, keeping only 6000 of the newest entries after merging.

### Remote service

This project comes with a remote service that can hold your ZSH entries from multiple machines. It is still a work in progress, the basic operations work as follows:

- Register (see `archy register -h`)
- Login (see `archy login -h`)
- Export Remote (see `archy export remote -h`)

This will be expanded as the service develops. 