# fitren :: Fitness Tracking Engine

## Introduction

Some businesses offer a bonus or stipend to employees who exercise regularly.
However, entering the data needed for auditing purposes is often redundant and error-prone.
The Fitness Tracking Engine (fitren) is a daemon intended to make it easier for you to audit gym, dojo, mo gwoon, etc. attendance.
Conbined with supporting software on a mobile device, Fitren reduces the administravia to a single click to verify attendance.

## Dependencies

Installation requires that you have Go version 1.1 or later installed.
You can find installation instructions for Go at [the Go language download page.](http://golang.org/doc/install)

Since Fitren sends e-mails, [you will need an account through Mailgun](http://www.mailgun.com).  For low-volume accounts, this will cost you nothing.

## Installation

You may proceed by executing the following commands at any Bash or compatible prompt:

    git clone git@github.com:sam-falvo/fitren
    cd fitren
    ./setup.sh
    . ./env.sh
    go install fitrend

## Configuring Fitren

### How To
You can edit the program directly with the following command:

    vim src/fitrend/main.go

Substitute your editor of choice for `vim` above.  Once you complete your configuration changes, you can rebuild the program with:

    go install fitrend

### Configuration Reference

The following settings exist inside a `const` declaration:

    mgDomain      = "..."
    mgApiKey      = "key-..."
    mgFromUser    = "FCE <fce@...>"
    mgSubject     = "How was class/gym tonight?"
    mgText        = "You will need HTML e-mail support to use this application.\n"
    webhookDomain = "localhost:8081"

`mgDomain` sets your Mailgun domain through which e-mail will be sent.
You'll get a sandbox domain when you sign up to Mailgun.
You can always create a new domain later on.

`mgApiKey` specifies your Mailgun API key.
You can retrieve this via your Mailgun account control panel.

`mgFromUser` specifies the identity of Fitren as seen by recipients of its emails.

`mgSubject` specifies the subject heading for the query sent to users.

`mgText` specifies the plain-text message to appear in an HTML e-mail.
This should probably be left alone.

`webhookDomain` specifies the HTTP server's listening address.

## Running Fitren

At present, configuration for Fitren is internal to the source code of the daemon.
Thus, all that's needed to launch it is:

    bin/fitrend

## Using Fitren

The following commands may be used to inform Fitren of the gym(s), dojo(s), et. al. you attend.  For example, I train at both Aikido West and Ving Tsun Sito.  Thus, instead of `gym-name`, I might specify the names `aw` or `vts` in its place.

    curl -X GET http://localhost:8081/gyms
    curl -X POST http://localhost:8081/gyms/gym-name
    curl -X DELETE http://localhost:8081/gyms/gym-name

The following commands may be used to inform Fitren of the users it's supposed to track.  Users have first and last names, along with an e-mail address to contact them with.

    curl -X GET http://localhost:8081/users
    curl -X POST http://localhost:8081/users/user-id/first-name/last-name/email
    curl -X DELETE http://localhost:8081/users/user-id

The following commands are used to retrieve the attendance records of all users registered with this Fitren instance, and to record an attendance between a user and a gym.

The `timestamp` is formatted as YYYYMMDD, with no dashes, dots, slashes, etc.

    curl -X GET http://localhost:8081/attendance
    curl -X POST http://localhost:8081/attendance/user-id/gym-name/timestamp

Generally speaking, the normal mode of operation for Fitren is to serve as a webhook responder for clients residing on the users' mobile device(s).  For instance, a mobile client might use a geolocation proximity trigger to `POST` to the `/webhooks/ask` endpoint.  This will cause an email to be sent to the user with a link embedded in the message to acknowledge attendance.  The link will refer to the `webhooks/ack` endpoint below.
    
    curl -X POST http://localhost:8081/webhooks/ask/user-id/gym-name/timestamp
    curl -X GET http://localhost:8081/webhooks/ack/ask-id
