# Home


    __   __          _
    \ \ / /_ _ _ __ (_) ___
     \ V / _` | '_ \| |/ __|
      | | (_| | | | | | (__
      |_|\__,_|_| |_|_|\___|
    Yet another node info collector

(previously [respond-collector](https://github.com/FreifunkBremen/respond-collector))

[![Build Status](https://travis-ci.org/FreifunkBremen/yanic.svg?branch=master)](https://travis-ci.org/FreifunkBremen/yanic)
[![Coverage Status](https://coveralls.io/repos/github/FreifunkBremen/yanic/badge.svg?branch=master)](https://coveralls.io/github/FreifunkBremen/yanic?branch=master)
[![codecov](https://codecov.io/gh/FreifunkBremen/yanic/branch/master/graph/badge.svg)](https://codecov.io/gh/FreifunkBremen/yanic)
[![Go Report Card](https://goreportcard.com/badge/chaos.expert/FreifunkBremen/yanic)](https://goreportcard.com/report/chaos.expert/FreifunkBremen/yanic)

`yanic` is a respondd client that fetches, stores and publishes information about a Freifunk network.

## The goals:

* Generating JSON for [Meshviewer](https://github.com/ffrgb/meshviewer)
* Storing statistics in [InfluxDB](https://influxdata.com/) or [Graphite](https://graphiteapp.org/) to be analyzed by [Grafana](http://grafana.org/)
* Provide a little webserver for a standalone installation with a meshviewer
