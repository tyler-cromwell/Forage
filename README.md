# Forage

## What is Forage?
An Internet-of-Things (IoT) application written in [Go][go] using [MongoDB][mongo] that adds smart functionality to my kitchen. Goal is to utilize a barcode scanner and scale on my Raspberry Pi to automatically track the groceries I have stored to quickly tell me what recipes I can cook and when items will expire.

Its goals are to:
1. Quickly identify what recipes I can make from the food I have stocked.
2. Reduce/prevent food waste by alerting (via [Twilio][twilio] SMS) and creating a shopping list (as a [Trello][trello] card) when items are nearing their expiration dates. Food items are defined as [JSON][json] documents in a [MongoDB][mongo] database.

## Repository Hierarchy
- `api`: Code that defines the REST API and background job(s).
  - `api.go`: Primary API handler and background job functions.
  - `listen.go`: Code for the listening/serving of HTTP requests.
  - `respond.go`: Code for constructing standard error responses.
- `clients`: Code for interfacing with 3rd-party systems.
  - `mongo.go`: Code for interfacing with MongoDB.
  - `trello.go`: Code for interfacing with Trello.
  - `twilio.go`: Code for interfacing with Twilio.
- `config`: Definition of global configuration (including client interface definitions).
- `operations`: Various scripts intended to be run manually.
  - `populate.js`: Script to populate MongoDB with documents describing each item of food being tracked.
- `tests`: Unit tests and mocks.
- `utils`: Miscellaneous helper code.
- `vendor`: Vendored dependencies.

## Depenencies
See `go.mod`


[go]: https://golang.org/
[json]: https://www.json.org/json-en.html
[mongo]: https://www.mongodb.com/
[trello]: https://trello.com/
[twilio]: https://www.twilio.com/
