# Forage

## What is Forage?
An Internet-of-Things (IoT) application written in [Go][go] using [MongoDB][mongo] that adds smart functionality to my kitchen. Goal is to utilize a barcode scanner and scale on my Raspberry Pi to automatically track the groceries I have stored to quickly tell me what recipes I can cook and when items will expire.

Its goals are to:
1. Quickly identify what recipes I can make from the food I have stocked.
2. Reduce/prevent food waste by alerting (via [Twilio][twilio] SMS) and creating a shopping list (as a [Trello][trello] card) when items are nearing their expiration dates. Food items are defined as [JSON][json] documents in a [MongoDB][mongo] database.

## Features
Presently, *Forage* can handle the following:
- Ingredient expiration: know which of the ingredients in your kitchen will expire and how soon.
- Recipe availability: know which recipes you can cook tonight from what you have available.
- Shopping list curation: see which ingredients need replacing, without risk of forgetting.
- SMS alerting: be reminded of when its time to go grocery shopping.

## Depenencies
- [adlio/trello][packageTrello]: Trello API client
- [twilio/twilio-go][packageTwilio]: Twilio API client
- [mongo-driver][packageMongo]: MongoDB driver
- [sirupsen/logrus][packageLogrus]: Logging
- [gorilla/mux][packageMux]: HTTP router/multiplexer

See `go.mod` for details.

## Usage
### Environment Variables
- `FORAGE_CONTEXT_TIMEOUT`:
- `FORAGE_INTERVAL` (CURRENTLY UNUSED):
- `FORAGE_LOOKAHEAD`: the amount of time the Expiration job [checks ahead][checksAhead] for expiring items.
- `FORAGE_TIME`: the time of day at which the Expiration job is [scheduled to execute][checkExpirationsScheduled].
- `FORAGE_TIMEZONE`: the timezone in which this instance is hosted.
- `LISTEN_SOCKET`: the socket upon which to listen for incoming connections.
- `LOGRUS_LEVEL`: the log granularity threshold (e.g. `DEBUG`, `INFO`, `WARN`, `ERROR`).
- `MONGO_URI`: the database connection string.
- `TRELLO_MEMBER`: 
- `TRELLO_BOARD`: the Trello board upon which to place the shopping list card.
- `TRELLO_LIST`: the Trello list upon which to place the shopping list card.
- `TRELLO_LABELS`: the Trello labels to be added to the shopping list card.
- `TRELLO_API_KEY`:
- `TRELLO_API_TOKEN`: 
- `TWILIO_ACCOUNT_SID`: 
- `TWILIO_AUTH_TOKEN`: 
- `TWILIO_PHONE_FROM`: the Twilio phone number assigned to this instance of Forage from which to send SMS messages.
- `TWILIO_PHONE_TO`: the recipent phone number for SMS messages from Twilio.

#### Example
```
export FORAGE_CONTEXT_TIMEOUT="5s"
export FORAGE_INTERVAL="24h"
export FORAGE_LOOKAHEAD="48h"
export FORAGE_TIME="19:00"
export FORAGE_TIMEZONE="America/New_York"
export LISTEN_SOCKET="127.0.0.1:8001"
export LOGRUS_LEVEL="DEBUG"
export MONGO_URI="mongodb://127.0.0.1:27017"
```

### Steps
1. `<start mongodb instance>`
2. `go build`
3. `source .env && ./forage`

Easy as 1,2,3.

## Experience

## Design Goals
- Comprehensive testing
- Clear readability
- Detailed logging

## Repository Hierarchy
- `api/`: Code that defines the REST API and background job(s).
  - `api_test.go`: Primary test file, containing unit tests for every case in HTTP handlers and their supporting functions.
  - `api.go`: Primary API handler and background job functions.
  - `listen.go`: Code for the listening/serving of HTTP requests.
  - `mocks.go`: Defines various mock functions and test variables for use in `api_test.go`.
- `clients/`: Code for interfacing with 3rd-party APIs.
  - `mongo.go`: MongoDB interface.
  - `trello.go`: Trello interface.
  - `twilio.go`: Twilio interface.
- `config/`: Definition of global configuration (including client interface definitions) as well as an abstract MongoDB interface to enable mocking.
- `operations/`: Various scripts intended to be run manually.
  - `populate.js`: Script to populate MongoDB with documents describing each ingredient & recipe being tracked.
- `tests/`: Unit tests and mocks.
- `utils/`: Miscellaneous helper code.
- `vendor/`: Vendored dependencies.

## Contributing
1) Write your code, following the current style
2) Write your tests & ensure complete coverage
3) Submit a pull request and I'll review it :)

[go]: https://golang.org/
[json]: https://www.json.org/json-en.html
[mongo]: https://www.mongodb.com/
[trello]: https://trello.com/
[twilio]: https://www.twilio.com/

[packageLogrus]: https://pkg.go.dev/github.com/sirupsen/logrus
[packageMongo]: https://pkg.go.dev/go.mongodb.org/mongo-driver
[packageMux]: https://pkg.go.dev/github.com/gorilla/mux
[packageTrello]: https://pkg.go.dev/github.com/adlio/trello
[packageTwilio]: https://pkg.go.dev/github.com/twilio/twilio-go

[checksAhead]: https://github.com/tyler-cromwell/Forage/blob/master/api/mocks.go#L136
[checkExpirationsScheduled]: https://github.com/tyler-cromwell/Forage/blob/master/api/listen.go#L36
