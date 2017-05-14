# snap-plugin-publisher-signifai
Publishes snap metrics/events to SignifAI metrics service.

SignifAI is a machine intelligence platform that helps TechOps teams get to answers faster by learning from their expertise, not generic algorithims. SignifAI helps TechOps deliver more uptime by intelligently prioritizing alerts, quickly identifying the root cause of an issue and correlating all the relevant log, events and metric data associated with the issue.


1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements 
* [golang 1.8+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download signifai plugin binary:
You can get the pre-built binaries for your OS and architecture under the plugin's [release](https://github.com/SignifAi/snap-plugin-publisher-signifai/releases) page.  For Snap, check [here](https://github.com/intelsdi-x/snap/releases).


#### To build the plugin binary:
Fork https://github.com/SignifAi/snap-plugin-publisher-signifai

Clone repo into `$GOPATH/src/github.com/SignifAi/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-signifai.git
```


#### Building
The following provides instructions for building the plugin yourself if
you decided to download the source. We assume you already have a $GOPATH
setup for [golang development](https://golang.org/doc/code.html). The
repository utilizes [glide](https://github.com/Masterminds/glide) for
library management.

build:
  ```make```

testing:
  ```make test```

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

#### Load the Plugin
Once the framework is up and running, you can load the plugin.
```
$ snaptel plugin load snap-plugin-publisher-signifai
Plugin loaded
Name: signafai
Version: 1
Type: publisher
Signed: false
Loaded Time: Sat, 18 Mar 2017 13:28:45 PDT
```

#### Task File
You need to create or update a task file to use the signafai publisher
plugin. We have provided an example, __tasks/signifai.yaml_ shown below. In
our example, we utilize the psutil collector so we have some data to
work with. There are three (3) configuration settings you can use.

Setting|Description|Required?|
|-------|-----------|---------|
|token|The Signafai [JWT token](https://docs.signifai.io).|Yes|


```
---
  version: 1
  schedule:
    type: "simple"
    interval: "5s"
  max-failures: 10
  workflow:
    collect:
      config:
      metrics:
        /intel/psutil/load/load1: {} 
        /intel/psutil/load/load15: {}
        /intel/psutil/load/load5: {}
        /intel/psutil/vm/available: {}
        /intel/psutil/vm/free: {}
        /intel/psutil/vm/used: {}
      publish:
        - plugin_name: "signafai"
          config:
            token: "1234ABCD"
```

Once the task file has been created, you can create and watch the task.
```
$ snaptel task create -t tasks/signafai.yaml
Using task manifest to create task
Task created
ID: 72869b36-def6-47c4-9db2-822f93bb9d1f
Name: Task-72869b36-def6-47c4-9db2-822f93bb9d1f
State: Running

$ snaptel task list
ID                                       NAME
STATE     ...
72869b36-def6-47c4-9db2-822f93bb9d1f
Task-72869b36-def6-47c4-9db2-822f93bb9d1f    Running   ...
```

## Documentation

docs.signifai.io

### Roadmap

We keep working on more feature and will update the publisher as needed.

## Community Support

Open an issue and we will respond.

## Contributing We love contributions!

The most immediately helpful way you can benefit this plug-in is by cloning the repository, adding some further examples and submitting a pull request.

## License
Released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@SignifAi](https://github.com/SignifAi/)
* Info: www.signifai.io
