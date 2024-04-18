# Github Daggerverse Module

## Operations
First of all you can chain functions when calling them with the CLI </br>
this builds the container with inputs passed by setter methods: `with-url` and `with-branch` 
```
dagger call -m github \
    with-url --addr="https://github.com/some-org/some-repo.git" \
    with-branch --branch="main" \
    container --repo-path="/src" --token=env:GITHUB_PAT
```


once the container is built, you can call any native dagger methods on the container </br>
for example: </br>
`directory --path "/src" entries --path "some/sub/path"`