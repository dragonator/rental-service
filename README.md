# rental-service

The service exposes two endpoints for requesting data:

`GET /rentals/{id}` - get rental by id

`GET /rentals` - list filtered rentals

#### Filters for listing:

* `ids` - list of integers representing rental ids
* `price_min` - integer value to filter for minimum price
* `price_max` - integer value to filter for maximum price
* `near` - 2 float values representing a location
* `sort` - string value representing a field to order results by
* `limit` - integer value to specify a pagination limit
* `offset` - integer value to specify a pagination offset

#### Example queries:
    rentals?ids=3,4,5
    rentals?price_min=9000&price_max=75000
    rentals?limit=3&offset=6
    rentals?near=33.64,-117.93
    rentals?sort=price
    rentals?near=33.64,-117.93&price_min=9000&price_max=75000&limit=3&offset=6&sort=price

## Run the tests

#### Install `github.com/moq/moq` to generate test mocks:

    go install github.com/moq/moq@latest

#### Init environment file:

    make init

#### Generate the mocks:

    make generate

#### Run the tests:

    make test

## Run the service

#### Init environment file:

    make init

#### Start the db:

    make db-up

#### Run the service:

    make server-start

## Stop the database

    make db-down

## Stop the database and remove volumes

    make volumes-down