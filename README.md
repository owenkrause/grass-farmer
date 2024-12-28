# [Get Grass](https://www.getgrass.io/) Air Drop Bot


<img src="cli.png" width="400px">


Decided to open source this after the airdrop. Request based script to simulate the get grass browser extension on multiple accounts.

## Features

Generate accounts
- automatically sign up and set a proxy

Farm
- simulates the browser extension

Check accounts
- check the points earned on each account

## Setup

1. paste emails into ``config/accounts.txt``
2. paste proxies into ``config/proxies.txt``
> **Note:** maintain 1:1 email to proxy ratio and must use sticky residential ips
3. run the script
  ```bash
  go run *.go
  ```
