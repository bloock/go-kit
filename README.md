# kit

Shared kernel of microservices

## How to works with libraries

Set the enviroment variable GOPRIVATE, to indicate the source of code (alternative)

```Bash
export GOPRIVATE=github.com/bloock
```

If in the "go get" process the command throw any error related with the credential, yo should launch:

```Bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

## How to versioning the library

The workflow should be

```Bash
 --- Normal flow ---
 git add . 
 git commit -m "Commmit msg"
 git push
 --- Closing version ---
 git tag v0.1.0
 git push origin v0.1.1
```

This command should create a new tag into remote github repository to be able use it as a dependency

⚠️  the tag must be unique if it not will throw an error

## How to use it

```Bash
go get github.com/bloock/go-kit
```