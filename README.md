# cqrs

[![Build Status](https://github.com/screwyprof/cqrs/actions/workflows/test.yml/badge.svg)](https://github.com/screwyprof/cqrs/actions/workflows/test.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/cfd6595574228033fd95/maintainability)](https://codeclimate.com/github/screwyprof/cqrs/maintainability)
[![Codecov](https://codecov.io/gh/screwyprof/cqrs/branch/master/graph/badge.svg)](https://codecov.io/gh/screwyprof/cqrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/screwyprof/payment)](https://goreportcard.com/report/github.com/screwyprof/cqrs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/screwyprof/cqrs.svg)
[![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)

## Description

`CQRS` is a library that provides an implementation of the Command Query Responsibility Segregation (CQRS) pattern
for Go. It aims to help you create clean, modular, and scalable applications by separating the read and write concerns
of your domain.

## Stable and Unstable Components

The library is currently considered unstable as it hasn't reached `v1.0.0` yet. However, it consists of both stable and 
unstable components. The stable components are less likely to have their API changed, while the unstable components are 
still under active development and might have breaking changes in the future. The unstable components can be found under 
the `x/` directory. Keep in mind that the whole library is under development, and it is recommended to always check the 
latest changes and updates before using it in your projects.

## Installation

To install the library, use the following command:

```bash
go get github.com/screwyprof/cqrs
```
## Quick Example

One of the stable components is the `aggregate` package, which provides a way to create event-sourced aggregates. An
aggregate is a domain object that processes commands and produces events as a result. Event sourcing means that the
aggregate's state is derived from its event history.

To use the `aggregate` package, you'll need to define your own identifier and event types. The package provides a
`FromAggregate` function to convert your domain aggregate into an event-sourced aggregate.

A runnable example demonstrating the usage of the `aggregate` package can be found in the `example_test.go` file
within the `aggregate` package. This example showcases how to define your own domain aggregate, commands, and events,
and how to process commands and apply events using the event-sourced aggregate. More examples can be found in the
the `examples` directory.

## Documentation

Full documentation can be found on [GoDoc](https://pkg.go.dev/github.com/screwyprof/cqrs).

## Contributing

Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open a new issue or submit a pull request.

## License

The Interactor Library is released under the [MIT License](https://opensource.org/licenses/MIT).

## Credits

This project was highly inspired by the following projects:

- https://github.com/gregoryyoung/m-r/blob/master/SimpleCQRS/Domain.cs
- https://github.com/MarkNijhof/Fohjin/tree/master/Fohjin.DDD.Example
- https://github.com/edumentab/cqrs-starter-kit/tree/master/sample-app
- https://github.com/jankronquist/rock-paper-scissors-in-java
