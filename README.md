# rails-configd

[![GoDoc](https://godoc.org/github.com/rubenfonseca/rails-configd?status.png)](https://godoc.org/github.com/rubenfonseca/rails-configd) [![Build Status](https://secure.travis-ci.org/rubenfonseca/rails-configd.png?branch=master)](http://travis-ci.org/rubenfonseca/rails-configd)

by Ruben Fonseca (rubenfonseca)

Online rails configuration generator using etcd as source data.

[![Demo](https://dl.dropboxusercontent.com/u/110525/rails-configd.gif)]

## Features

* Bring the power of [etcd](https://github.com/coreos/etcd) live changes to your Rails app
* Everytime there's a change on an etcd dir, your Rails config files are updated and processes are restarted
* Written in Go so no runtime dependency on your production servers
* Extendable system to support more rendererers (turns etcd data into files) and reloaders (reloads Rails processes)
* Currently supported renderers:
    * YAML - renderes the etcd data to a .yml file
* Currently supported reloaders:
    * Touch - touches `tmp/restart.txt` for passenger compatible servers.

## Installing

TODO

## FAQ

TODO
