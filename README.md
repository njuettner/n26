# N26 CLI 🚀

Use your [N26](https://n26.com) account via command-line

Latest Version: v0.4

## Features 🙌

- Get your **latest transactions**
- See your **balance**
- See all of your **N26 accounts**
- Get your **account information**
- Get your **bank statements via PDF**
- See your **N26 savings and investment**
- See your **N26 cards**
- Block/Unblock your **N26 cards**
- List all **N26 categories**

## Requirement

If you never used n26 cli just run `n26 init` to setup the configuration

or

you create a YAML file **n26.yml** in your ~/.config directory:

```yaml
username: your-email@domain.com
password: n26-password
```

## Installation

### Mac

```sh
brew install njuettner/n26/n26
```

### Windows/Linux 

[See latest release](https://github.com/njuettner/n26/releases/latest)

### Go

```bash
go get github.com/njuettner/n26
```

## How to use it 🤔

### Bash/ZSH Shell Completion

Add an additional statement to your bash_profile or zsh_profile:

```bash
eval "$(n26 --completion-script-bash)"
```

or

```bash
eval "$(n26 --completion-script-zsh)"
```

```bash
usage: n26 [<flags>] <command> [<args> ...]

A command-line to interact with your N26 bank account

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  init
    Setup the configuration to use N26 CLI

  categories
    Show N26 categories

  transactions [<amount>]
    Show N26 latest transactions (Number by Default: 5)

  balance
    Show N26 balance

  contacts
    Show N26 contacts

  account info
    Show N26 account information

  account limit
    Show N26 account limit

  account stats
    Show N26 account statistics

  account status
    Show N26 account status

  statement [<statementID>]
    Get N26 statement, will be saved as PDF files

  savings
    Show N26 savings and investments

  cards
    Show N26 cards

  block-card [<cardID>]
    Block N26 Card

  unblock-card [<cardID>]
    Unblock N26 Card
```
