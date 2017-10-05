## anonbot

[![GoDoc](https://godoc.org/github.com/ReformedDevs/anonbot?status.svg)](https://godoc.org/github.com/ReformedDevs/anonbot)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

This web application simplifies the task of managing a group of anon Twitter accounts. Here is a brief overview of the features that are provided:

- Users can sign up and suggest tweets
- Staff can approve tweets and schedule them
- Tweets are automatically tweeted on schedule

### Building

Assuming you have Docker and GNU Make installed, building the entire application is as simple as running:

    make

The makefile uses the `golang:latest` container to compile the application, so there is no need to have the Go toolchain installed locally. Once compilation is complete, you can run the application and view usage instructions with:

    dist/anonbot --help

### Using with Docker

Once built, the application can easily be run in a Docker container. Simply open a terminal and run:

    docker build -t ReformedDevs/anonbot .

[TODO: describe container env. variables]
