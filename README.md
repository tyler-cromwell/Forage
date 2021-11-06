# Forage

## What is Forage?
Forage is a [Go][go] app that monitors the various expiration dates of food that keep in my kitchen pantry, refrigerator, and freezer. Its goal is to reduce/prevent food waste by alerting (via [Twilio][twilio] SMS) and creating a shopping list (as a [Trello][trello] card) when items are nearing their expiration dates. Food items are defined as [JSON][json] documents in a [MongoDB][mongo] database.

## Repository Hierarchy
- `api`: Code that defines the REST API and background job(s).
  - `api.go`: Primary API handler and background job functions.
  - `respond.go`: Code for constructing standard error responses.
- `clients`: Code for interfacing with 3rd-party systems.
  - `mongo.go`: Code for interfacing with MongoDB.
  - `trello.go`: Code for interfacing with Trello.
  - `twilio.go`: Code for interfacing with Twilio.
- `operations`: Various scripts intended to be run manually.
  - `populate.js`: Script to populate MongoDB with documents describing each item of food being tracked.
- `utils`: Miscellaneous helper code.
- `vendor`: Vendored dependencies.

## Depenencies
See `go.mod`


[go]: https://golang.org/
[json]: https://www.json.org/json-en.html
[mongo]: https://www.mongodb.com/
[trello]: https://trello.com/
[twilio]: https://www.twilio.com/