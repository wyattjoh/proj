# proj

Ever in the situation where you're always jumping to the same folder over and
over again? Well think no more, and drastically reduce your typing by the
creating of "projects" that you can jump to with ease.

## Installation

You must have `go` installed, then run:

```
go get github.com/wyattjoh/proj
```

Which will install the `proj` binary to your `$GOPATH/bin` directory.

Add the following to your shell file to add shortcuts to jump to projects.

```
function p() {
  dir=`proj g $1`
  [[ $? == 0 ]] && cd $dir
}
function p.() { proj a $1 $2; }
alias pl='proj l'
```

## Usage

### Save a project

When you are in a directory that you want to save a project for, do:

```
p. <PROJECT_NAME>
```

### Change directory to a project

```
p <PROJECT_NAME>
```

### List all projects and their directories

```
pl
```

### Advanced

These are all aliases to the `proj` command, run `proj` for more advanced
features. All the project data is additionally stored in `$HOME/.projects.db` as
a Bolt DB.

## Future

- Shell completions
- Delete a project
