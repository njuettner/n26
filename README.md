# N26 CLI 🚀

Use your [N26](https://n26.com) account via command-line

Latest Version: v0.3

## Features 🙌

- Get your **latest transactions**
- See your **balance**
- See all of your **N26 accounts**
- Get your **account information**
- Get your **bank statements via PDF**
- See your **N26 savings and investment**
- See your **N26 cards**

## Requirement

If you never used n26 cli before you can run n26 init to setup the configuration

or

you can create the YAML file **n26.yml** in your ~/.config directory by yourself:

```yaml
username: your-email@domain.com
password: n26-password
```

## Installation

You could either use `go get github.com/njuettner/n26` or just download the [latest release](https://github.com/njuettner/n26/releases/latest).

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
    Setting up the configuration to use N26 CLI

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

  savings
    N26 savings and investments

  cards
    N26 cards
```
