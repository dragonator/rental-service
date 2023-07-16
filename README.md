# rental-service

The service exposes two endpoints for requesting data:

`GET /rentals/{id}` - get rental by id

`GET /rentals` - list filtered rentals

Filters for listing:

* `ids` - list of integers representing rental ids
* `price_min` - integer value to filter for minimum price
* `price_max` - integer value to filter for maximum price
* `near` - 2 float values representing a location
* `sort` - string value representing a field to order results by
* `limit` - integer value to specify a pagination limit
* `offset` - integer value to specify a pagination offset

    

## Run the tests

1. Install `github.com/moq/moq` to generate test mocks:

    > go install github.com/moq/moq@latest

2.  Generate the mocks:

    > make generate

3. Run the tests:

    > make test

## Run the service

1. Start the db:

    > make db-up

2. Run the service:

    > make server-start

## Stop the database

    > make db-down

## Stop the database and remove volumes

    > make volumes-down