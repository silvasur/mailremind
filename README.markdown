# mailremind

mailremind is a web-based service that sends you mails with a delay-

## Why?

I often send myself an email to remind me of something. If the event I want to be remindered of is in the not-so-close future this method does not work so well, since the mail is then not new any more. Also mailremind allows you to send mails repetetive based on a schedule, so it can be used for reoccurring events, like birthdays.

## Installation

Get the sources and build mailremind with `go build`.

Or simply run `go get github.com/silvasur/mailremind`. This will place the compiled binary in your `$GOPATH/bin` directory.

## Configuration

All config stuff is in the mailremind.ini file.

## Running

Simply run `./mailremind -config mailremind.ini` and mailremind will listen on the configured address (net.laddr) and handle http requests.

## Public mailremind installations

* [mailremind.silvasur.net](http://mailremind.silvasur.net)
