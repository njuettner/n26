# N26 CLI

Documentation will be improved soon.

Create a YAML in your ~/.config a file n26.yml

```
username: your-email@domain.com
password: n26-password
```

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

Current Version: v0.1
