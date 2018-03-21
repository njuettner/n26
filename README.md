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

You could either use `go get github.com/njuettner/n26` or just download the binary release (Windows/Linux/Mac OS)

## How to use it

Usage:

```
usage: n26 [<flags>] <command> [<args> ...]

A command-line to interact with N26

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  transactions [<amount>]
    N26 Transactions (Number by Default: 5)

  balance
    N26 Balance

  contacts
    N26 Contacts

  account info
    Info

  account limit
    Limit

  account stats
    Statistics

  account status
    Status

  statements [<statementID>]
    N26 Bank Statements

  cards
    Cards
```

