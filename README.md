# rails-configd

[![GoDoc](https://godoc.org/github.com/rubenfonseca/rails-configd?status.png)](https://godoc.org/github.com/rubenfonseca/rails-configd) [![Build Status](https://secure.travis-ci.org/rubenfonseca/rails-configd.png?branch=master)](http://travis-ci.org/rubenfonseca/rails-configd)

by Ruben Fonseca (rubenfonseca)

Online rails configuration generator using etcd as source data.

![Demo](https://dl.dropboxusercontent.com/u/110525/rails-configd.gif)

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

    $ go install github.com/rubenfonseca/rails-configd
    $ rails-configd -h

## Example usage

First you have to set the data on your etcd cluster. Let's try to configure the database on our production rails app.

    $ etcdctl set rails_app01/database/production/host lan.db01.example.com
    $ etcdctl set rails_app01/database/production/adapter pg
    $ etcdctl set rails_app01/database/production/database db01
    $ etcdctl set rails_app01/database/production/username dbuser

Then, next to each Rails app, you should run a `rails-configd` daemon:

    $ rails-configd --etcd http://localhost:4001 --etcd-dir /rails_app01 --renderer yaml --yaml-file config/database.yml --reloader touch

This will read from etcd under `/rails_app01` directory and build the `config/database.yml` file. After this, any
change to a key under the `etcd-dir` directory will trigger the generation of a new `database.yml` file, and reload
the rails server by touching `tmp/restart.txt`.

## FAQ

### Why another daemon to do this?

I believe this daemon does only one thing and does it right. I don't want to add another responsability to Rails.
Also, since this is written in Go, you can distribute a binary to all your platforms with no runtime dependencies (aka:
bundle install...)

### But, what if I have more than one config file on my Rails app?

I also do :) Then the sollution is to run one `rails-configd` daemon for each file. Each daemon will read form a different
etcd directory, and you can control the output file with the param `yaml-file`.
Don't worry about resources, go daemons like this use very little memory :-)

## But I'm running this super awesome application server that doesn't support reloading by touching `tmp/restart.txt`!

Then pachtes are welcome :) There's a reloader interface on the code that you can implement with your own reloading method!
