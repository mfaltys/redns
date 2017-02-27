# Doic
[![Build Status (Travis)](https://travis-ci.org/mfaltys/doic.svg?branch=master)](https://travis-ci.org/mfaltys/doic)  
Doic is a DNS based IOC written in go.  This tool is designed to be a very low
overhead plug-and-play approach to implimenting an IOC for organizational or
personal use.  If you want to track our day-to-day activities, check out our
[trello board](https://trello.com/b/5KMHrR6L/project-requirement-clicitation-cybr4580).

## Executive Summary
Many industry malware detection tools (used in the workplace or for personal use)
are either platform specific, intrusive, or high overhead and need to be
maintained. Doic is a malware detection tool which plans to overcome these
problems in a small, easy to deploy package.  

### Goals and Objectives
Doic plans to combat the problems outined in the executive summary in the
following ways:

- **Easy to impliment** If its run as the main DNS server it can detect
  malicious content without having to run any client-side software. This also
  means it can be deployed in a large corperate environment just as easily as 
  on a personal computer.
- **Low overhead**. Doic is written in golang and the entire system can be deployed
  on most architectures in under 10MB. The tool is also statically compiled
  meaning no extensive dependencies, it is shipped with everything it needs to
  run.
- **Non intrusive**. Being non-client side, it does not run in the background
  therefore is no-overhead on your workstations. Doic has no annoying popops and
  will only redirect traffic if told to do so.

Doic is and always will remain open-source. After version 1.0 is released we
encourage critique and peer reviews. If you have any questions or feedback feel
free to contact us or better yet submit a PR.  
Doic will continue to be fast and reliable, DNS requests passing through the
tool add a negligible amount of time to the request and its probably faster than
your current firewall implimented DNS server. Doic is tested for errors with
unit tests and function test being constantly ran against the code over on our
projects [travis page](https://travis-ci.org/mfaltys/doic).

![granttpro is not free](https://s3.amazonaws.com/unixvoid-blog/trickedbyabusiness.png)

## Project Risks

|Risk name | Impact | Likelihood | Description |
|----------|--------|------------|-------------|
|Domain list availability (16) | 8 | 2 | If open doimain lists are hard to get a hold of proving malicious intent becomes difficult |
|Alerting system (9) | 3 | 3 | Finding a way to alert users could take a fair amount of time with one developer |
|Procrastination (12) | 3 | 4 | I always pull through with enough coffee |
|Unknown DNS issues (12) | 3 | 4 | If clients try an unsupported DNS call it will not resolve but instead return a DNS SERVFAIL and timeout if not handled correctly. I dont know what I dont know |
|Unit testing (8) | 2 | 4 | Its difficut to unit test DNS calls in go, feature testing is much easier. It could have unknown negative side effects|
