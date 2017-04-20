# godep-check

`godep-check` helps you to understand how your `Godeps.json` file is in sync with your $GOPATH!

Just run `godep-check` in a repo with a Godeps folder or with a Godeps.json file and it will print out how much the repos in your $GOPATH are in sync telling you if

- you are in sync
- you have a (probably) old version
- you have a newer version (and how much new)
- your work tree is dirty!

![terminal-img](https://cloud.githubusercontent.com/assets/1763949/25254544/621fd8ac-2626-11e7-871e-2f9d52d4ca14.png)

Run it with the verbose option `-v` to check additional details:

![terminal-img-verbose](https://cloud.githubusercontent.com/assets/1763949/25254543/621fa684-2626-11e7-8dc7-c5c236454b59.png)