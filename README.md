# What is Webhookr?

Quickly test [webhook](http://en.wikipedia.org/wiki/Webhook) callbacks.

Start a new Webhookr and begin receiving seeing HTTP requests right in your browser. Data is not retained, so refreshing will clear out a Webhookr.
Webhookrs are completely ephemeral and your history is only stored in your cookies.

Feel free to branch code and host as necessary. Free service at [http://webhookr.com](http://webhookr.com)

## Usage

Go to a webhookr, such as [http://webhookr.com/2c93b147-6431-45ec-9672-b3a522881185](http://webhookr.com/2c93b147-6431-45ec-9672-b3a522881185)
and then in the command line try sending a request over a la:

    curl -X POST http://webhookr.com/2c93b147-6431-45ec-9672-b3a522881185 -d 'SomeParam=1&B=hi' -H "X-Special-Header: test"

