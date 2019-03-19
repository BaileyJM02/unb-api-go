<h1 align="center">UnbelievaBoat API Wraper</h1><br>
<p align="center">
  <a href="https://open-source.hue.observer/pre-micro/">
    <img alt="Unb" title="Unb" src="https://i.imgur.com/tUiCsY5.jpg" width="400">
  </a>
</p>

<p align="center">
  A Go library for the UnbelievaBoat Discord bot API.
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

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/badge/Chat_On-Discord-008080.svg?style=flat-square)](https://discordapp.com/invite/YMJ2dGp)

*In writing*

## Features

A few of the things you can do with unb-api-go:

* Connect to the UnbelievaBoat API
* See user balance
* See guild leaderboard
* Set custom http.Client
* More in development...

## Feedback

Feel free to send me feedback on [Twitter](https://twitter.com/BaileyJM02) or [file an issue](https://github.com/baileyjm02/unb-api-go/issues/new). Feature requests are always welcome. If you wish to contribute, please take a quick look at the [guidelines](./CONTRIBUTING.md)!

If there's anything you'd like to chat about, please feel free to join our [Discord Server](https://discordapp.com/invite/YMJ2dGp) and mention me: `@Bailey#0004`!

## Install Process

> I'm guessing you already have an environment setup.
- `go get github.com/BaileyJM02/unb-api-go/v1` to install the wrapper (v1)
- `api := v1.New(token)` to create a new instance, `token` is your token found [here](https://unb.pizza/api/docs)
- `api.GetBalance(guildID, userID)` functions follow the `api` var we created above.

## Acknowledgments

Thanks to [Codenvy](https://codenvy.io) for supporting me with an awesome IDE while I don't have a computer. 
