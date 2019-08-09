<h1 align="center">UnbelievaBoat API Wrapper</h1><br>
<p align="center">
  <a href="https://unb.pizza/">
    <img alt="Unb" title="Unb" src="https://i.imgur.com/tUiCsY5.jpg" width="400">
  </a>
</p>

<p align="center">
  A Go library for the UnbelievaBoat Discord bot API.<br/>
  By <a href="https://github.com/BaileyJM02/">Bailey</a> and <a href="https://github.com/BaileyJM02/unb-api-go/graphs/contributors">Contributors</a>
</p>

<p align="center">
  <a href="https://unb.pizza/">
    Website
  </a>
|
  <a href="https://discordapp.com/invite/YMJ2dGp">
    Discord
  </a>
</p>

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Feedback](#feedback)
- [Install Process](#Install-process)
- [Acknowledgments](#acknowledgments)

## Introduction

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/badge/Chat_On-Discord-008080.svg?style=flat-square)](https://discordapp.com/invite/YMJ2dGp)

The UnbelievaBoat API library is an importable package to make using the UnbelievaBoat API easier to use within Go. Although updates may be few and far, all new features are added as they are released. As I talk to the developer of the bot, I am made aware of upcoming changes. **PRs, bug reports and feature requests are welcome!**

Version 1 ([`/v1`](https://github.com/BaileyJM02/unb-api-go/tree/master/v1)) uses version 1 of the UnbelievaBoat API and should be imported as `github.com/BaileyJM02/unb-api-go/v1`, more on the install process [here](#install-process). This is allows for the second version of the api to be installed as `github.com/BaileyJM02/unb-api-go/v2` etc. upon release.

> This may change as Go Module support will soon be implemented.

## Features

A few of the things you can do with **unb-api-go**:

* Connect to the UnbelievaBoat API
* See user balance
* See guild leaderboard
* Set custom http.Client
* And more...

## Feedback

Feel free to send me feedback on [Twitter](https://twitter.com/BaileyJM02) or, preferably, [file an issue](https://github.com/baileyjm02/unb-api-go/issues/new). Feature requests are always welcome. If you wish to contribute, please take a quick look at the [guidelines](./CONTRIBUTING.md)!

If there's anything you'd like to chat about, please feel free to join our [Discord Server](https://discordapp.com/invite/YMJ2dGp) and mention me: `@Bailey#0004`!

## Install Process

> I'm guessing you already have an environment setup.
1. Install the project.
	-  `$ go get github.com/BaileyJM02/unb-api-go/v1` to install the wrapper using version one.
2.  Create a new instance to use the library. 
	- ( `token` is your token found [here](https://unb.pizza/api/docs))
```go
import(
	"github.com/BaileyJM02/unb-api-go/v1"
)

function main() {
	// assign and create a new instance to api.
	api := v1.New(token);
}
```

3. Use functions like so: `api.GetBalance(guildID, userID)` . Where `guildID` and `userID` are representatives of their Discord values.

## Acknowledgments
Thank you to all of the contributors who have added and improved the project!
