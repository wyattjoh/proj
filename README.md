# proj

## Shell alias

Add the following to your shell file to add shortcuts to jump to projects.

```
function p() {
  dir=`proj g $1`
  [[ $? == 0 ]] && cd $dir
}
function p.() { proj a $1 $2; }
alias pl='proj l'
```
