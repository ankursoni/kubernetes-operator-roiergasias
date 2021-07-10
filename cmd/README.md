# Getting Started with Roiergasias workflow


## Workflow YAML file syntax
```yaml
version: 0.1

environment: # global environment
  - <variable name>: <variable value>

task:
  - node: <node label> # node is optional and is used to match with kubernetes node having label - node.roiergasias=<node label>
    sequential:
      - <step type 1>:  # currently supported types - print, execute or environment
          - <step 1 argument 1> # can contain {{env:<variable name>}} to resolve environment variables
          - <step 1 argument 2>
          - <step 1 argument 3>
          - ...
      - <step type 2>:
          - <step 2 argument 1>
          - ...
      - ...
  - ...
```
> NOTE:
> 1. Global environment variables are available to all the tasks.
> 2. If "node" is specified in one task, then it must be specified in all the tasks for split workflow to run properly.
> 3. Environment variables defined as a step are available to subsequent steps and tasks, even in split node workflow
> where, they get added to the global environment list for subsequent splits.
> 4. Label the specific kubernetes node to run tasks with "\<node label\>" by running the kubectl command:  
> ```shell
> kubectl label nodes <kubernetes node> node.roiergasias="<node label>" --overwrite
> ```
For details on an example, follow this [README](../examples/hello-world/README.md)

## Command line syntax
```shell
# clone to a local git directory, if not already done so
git clone https://github.com/ankursoni/kubernetes-operator-roiergasias.git

# change to the local git directory
cd kubernetes-operator-roiergasias

# set execute permissions to roiergasias cli
chmod +x cmd/linux/roiergasias cmd/osx/roiergasias


# run the roiergasias cli help
./cmd/linux/roiergasias --help
# or, for mac osx
./cmd/osx/roiergasias --help

<<output
---
Usage:
  roiergasias [OPTIONS] <command>

Application Options:
  -d, --debug  enable debug level for logs

Help Options:
  -h, --help   Show this help message

Available commands:
  run       run a workflow
  split     split a workflow yaml file
  validate  validate a workflow yaml file
  version   display version
---
output


# to get help on a specific command:
./cmd/linux/roiergasias run --help
# or, for mac osx
./cmd/osx/roiergasias run --help

<<output
---
Usage:
  roiergasias [OPTIONS] run [run-OPTIONS]

run a workflow yaml file

Application Options:
  -d, --debug     enable debug level for logs

Help Options:
  -h, --help      Show this help message

[run command options]
      -f, --file= workflow yaml file
---
output
```


## Install Roiergasias operator in Kubernetes
```shell
# install the operator
helm install --repo https://github.com/ankursoni/kubernetes-operator-roiergasias/raw/main/operator/helm/ \
  --version v0.1.1 \
  roiergasias-operator roiergasias-operator

# uninstall the operator
helm uninstall roiergasias-operator
```


## Repository map
```text
┬ 
├── ...
│   📌 --------------------> you are here
├── cmd    ----------------> contains go main starting point for roiergasias workflow cli
│   ├── linux   -----------> contains linux amd64 executable for roiergasias workflow cli
│   └── osx   -------------> contains mac-osx amd64 executable for roiergasias workflow cli
├── ...
```
