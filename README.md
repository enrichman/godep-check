# godep-check

`godep-check` helps you to understand how your `Godeps.json` file is in sync with your $GOPATH!

Just run `godep-check` in a repo with a Godeps folder or with a Godeps.json file and it will print out how much the repos in your $GOPATH are in sync telling you if

- you are in sync
- you have a (probably) old version
- you have a newer version (and how much new)
- your work tree is dirty!

Run it with the verbose option `-v` to check additional details