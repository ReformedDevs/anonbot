## anonbot

[![Build Status](https://ci.quickmediasolutions.com/buildStatus/icon?job=anonbot)](https://ci.quickmediasolutions.com/job/anonbot)
[![GoDoc](https://godoc.org/github.com/ReformedDevs/anonbot?status.svg)](https://godoc.org/github.com/ReformedDevs/anonbot)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

This web application simplifies the task of managing a group of anon Twitter accounts. Here is a brief overview of the features that are provided:

- Users can sign up and suggest tweets
- Staff can approve tweets and schedule them
- Tweets are automatically tweeted on schedule

### Building

Assuming you have Docker and GNU Make installed, building the entire application is as simple as running:

    make

The makefile uses the `golang:latest` container to compile the application, so there is no need to have the Go toolchain installed locally. Once compilation is complete, you can run the application with:

    dist/anonbot

### Usage

anonbot requires an SQL database (currently limited to PostgreSQL) which must be initialized before first use. Assuming you have PostgreSQL running locally, the command would look like this:

    dist/anonbot \
        --db-driver postgres \
        --db-args 'dbname=postgres user=postgres' \
        migrate

Once the migration completes, you can run the command above without "migration" to launch the web server. By default, the application is accessible on port 8000.

### Using with Docker

The application can easily be run in a Docker container. To build the container, simply open a terminal and run:

    docker build -t reformeddevs/anonbot .
