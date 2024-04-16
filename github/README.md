# Github Daggerverse Module

## Operations
First of all you can chain functions when calling them with the CLI </br>
this builds the container with inputs passed by setter methods: `with-repo` and `with-branch` 
`dagger call with-repo --repo="git@github.com:your-org/your-priv-repo.git" with-branch --branch="main" container`

once the container is built, you can call any native dagger methods on the container
