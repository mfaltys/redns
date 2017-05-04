# Redns
[![Build Status (Travis)](https://travis-ci.org/mfaltys/redns.svg?branch=master)](https://travis-ci.org/mfaltys/redns)  
Redns is a DNS based indicator of compromise (IOC) written in go.  This tool is designed to be a very low
overhead plug-and-play approach to implimenting an IOC for organizational or
personal use.  If you want to track our day-to-day activities, check out our
[trello board](https://trello.com/b/5KMHrR6L/redns).

## Using redns

### Getting precompiled binaries
This project pushes up a binary on every project commit and tag.
You can find these binaries in the following locations (note we only pre compile
for 64 bit linux architectures):  
[browse binaries](https://cryo.unixvoid.com/bin/redns/)  
[latest redns](https://cryo.unixvoid.com/bin/redns/redns-latest-linux-amd64)  
[latest redns_cli](https://cryo.unixvoid.com/bin/redns/redns_cli-latest-linux-amd64)

### Compiling from source
Redns uses Make for compilation/testing.  Use the following commands to buid redns
from source.
- First make sure you have golang installed and configured
- To pull all dependencies use `make dependencies`
- To build dynamically use `make` or `make redns`
  - The dynamically compiled binary will end up in the `bin/` directory
- To build statically use `make stat`
  - The statically compiled binary will end up in the `bin/` directory
- To build the cli use `make stat_cli`
  - The portable binary will end up in the `bin/` directory
- To clean up all compiled binaries use `make clean`

### Configuration
Currently a work in progress as the configuration is changing often in the early
stages (pre-release) of redns.  
[Configuration docs](https://mfaltys.github.io/redns_docs/configuration/index)  


## Documentation  
You can find our documentation over on [github pages](https://mfaltys.github.io/redns_docs)  
* [Milestone 1](https://mfaltys.github.io/redns_docs/milestone.1/index)  
* [Milestone 2](https://mfaltys.github.io/redns_docs/milestone.2/index)  
* [Milestone 3](https://mfaltys.github.io/redns_docs/milestone.3/index)

### Contributions
Big shoutout to [MLHale](https://github.com/MLHale) for the project name, turns out that was the hardest part.
