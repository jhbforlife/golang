# translate cli
a cli program that interacts with my [translate server](https://github.com/golang/tree/main/cmd/translate/server) (currently not live)

## Usage
```
translate [options] original

-s, --source
    source language

-h, --help
    print this message

-t, --to *required
    language to translate to
```

## Examples
```
translate -s fr -t en bonjour, comment ça va?
```
```
translate --to french what is your name?
```