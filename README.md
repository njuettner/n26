# N26 CLI

Use your [N26](https://n26.com) account via command-line

Latest Version: v0.1

## Requirement

You only need to create a YAML file **n26.yml** in your ~/.config directory.

```
username: your-email@domain.com
password: n26-password
```

## Installation

You could either use `go get github.com/njuettner/n26` or just download the binary release ([Windows](https://github.com/njuettner/n26/releases/download/v0.1/n26_windows_amd64.exe)/[Linux](https://github.com/njuettner/n26/releases/download/v0.1/n26_linux_amd64)/[Mac OS](https://github.com/njuettner/n26/releases/download/v0.1/n26_darwin_amd64))

## How to use it

```
usage: n26 [<flags>] <command> [<args> ...]

A command-line to interact with your N26 bank account

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  transactions [<amount>]
    N26 latest transactions (Number by Default: 5)

  balance
    N26 balance

  contacts
    N26 contacts

  account info
    N26 account information

  account limit
    N26 account limit

  account stats
    N26 account statistics

  account status
    N26 account status

  statements [<statementID>]
    N26 statements, will be saved as PDF files

  cards
    N26 cards
```