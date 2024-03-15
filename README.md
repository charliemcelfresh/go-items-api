### go-items-api

This project is part of a performance test among Sinatra (Ruby), Ruby on Rails, and Go. Here are the Rails and Sinatra projects:

- https://github.com/charliemcelfresh/rails-items-api
- https://github.com/charliemcelfresh/sinatra-items-api

Set up the Rails project first, because it creates, migrates, and seeds the database that all three projects use. See [here](https://github.com/charliemcelfresh/rails-items-api)

Once all three are set up, you can compare their performance by running the Go projects performance testing client. See below.

#### Run the app

* Set up and run `sinatra-items-api` and `rails-items-api`. Note that db setup for all three projects is in `rails-items-api`
* Copy .env.SAMPLE to .env
* Paste the URI of the MySQL db you created and seeded in `rails-items-api` as the value for `MYSQL_DATABASE_URL` in .env
* Run the Go server: `go run main.go`
* Run the performance tester: `go run internal/performance_tester/main.go`