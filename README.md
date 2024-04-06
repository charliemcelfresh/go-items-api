### go-items-api

This project is part of a performance test among Postgrest (a PostgreSQL plugin), Sinatra (Ruby), Ruby on Rails, and Go. Here are the Postgrest, Rails and Sinatra projects:

- https://github.com/charliemcelfresh/postgrest-items-api
- https://github.com/charliemcelfresh/rails-items-api
- https://github.com/charliemcelfresh/sinatra-items-api

Set up the Postgrest project first, because it creates, migrates, and seeds the database that all three projects use. See [here](https://github.com/charliemcelfresh/postgrest-items-api).

Once all four projects are set up, you can compare their performance by running the Go projects performance testing client. See below.

#### Run the Server and Performance Tester

* Set up and run all three projects above. Note that db setup for all three projects is in [postgrest-items-api](https://github.com/charliemcelfresh/postgrest-items-api)
* Copy .env.SAMPLE to .env
* Paste the URI of the PostgreSQL db you created and seeded in [postgrest-items-api](https://github.com/charliemcelfresh/postgrest-items-api) as the value for `DATABASE_URL` in .env
* Run the Go server: `go run main.go`
* Run the performance tester: `go run internal/performance_tester/main.go`