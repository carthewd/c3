# c3 - CodeCommit CLI
`c3` is a tool which provides AWS CodeCommit functionality from the command line. 

Inspired in part by hub and the GitHub CLI. 

This is still in the very early stages of development 

[![asciicast](https://asciinema.org/a/357519.svg)](https://asciinema.org/a/357519)

## Usage 

- `c3 pr [approve, checkout, create, diff, list, revoke]`
- `c3 link [filepath | pr:123]`

## Installation

Prebuilt binaries (macOS and Linux) from [releases page][]

### Build from source 
0. Prerequisites
* git

1. Clone the repository
```
git clone git@github.com:carthewd/c3.git
cd c3
```

2. Build the project

``` make ```

3. Move the `c3` binary somewhere in your path

``` mv c3 /usr/local/bin/ ```

OR

``` make install ```
