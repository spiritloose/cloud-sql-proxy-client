# cloud-sql-proxy-client [![Build Status](https://travis-ci.org/spiritloose/cloud-sql-proxy-client.svg)](https://travis-ci.org/spiritloose/cloud-sql-proxy-client)

MySQL CLI launcher with Cloud SQL Proxy.

## Requirements

* MySQL CLI
* CloudSQL Proxy CLI

## Installation

```
$ go get github.com/spiritloose/cloud-sql-proxy-client
```

## Usage

```
$ cloud-sql-proxy-client project-name:us-west1:instance-name -u username -p database
```

## Environment Variables

* `CLOUD_SQL_PROXY_CLIENT_MYSQL`
  * eg: `mycli`

## Author

Jiro Nishiguchi <<jiro@cpan.org>>
