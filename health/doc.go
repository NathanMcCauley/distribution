// Package health provides a generic health checking framework. The health
// package works expvar style. By importing the package the debug server is
// getting a /debug/health endpoint that returns the current status of the
// application. If there are no errors, /debug/health will return a HTTP 200
// status, together with an empty JSON reply {}. If there are any checks with
// errors, the JSON reply will include all the failed checks, and the response
// will be have a HTTP 500 status.
//
// A Check can either be run synchronously, or asynchronously. We recommend
// that most checks are registered as an asynchronous check, so a call to the
// /debug/health endpoint always returns immediately. This pattern is
// particularly useful for checks that verify upstream connectivity or database
// status, since they might take a long time to return/timeout.
//
// To install health, just import it in your application:
//
// import "github.com/docker/distribution/health"
//
// You can also (optionally) import health/api that will add two convenience
// endpoints: /debug/health/down and /debug/health/up. These endpoints add
// "manual" checks that allow the service to quickly be brought in/out of
// rotation.
//
// import _ "github.com/docker/distribution/registry/health/api"
//
// # curl localhost:5001/debug/health
// {}
// # curl -X POST localhost:5001/debug/health/down
// # curl localhost:5001/debug/health
// {"manual_http_status":"Manual Check"}
// After importing these packages to your main application, you can start
// registering checks.
//
// The lowest-level way to register a check is by calling Register. Register
// allows you to pass in an arbitrary string and a Checker method that runs
// your check. If your method returns nil, it is considered a healthy check,
// otherwise it will make the health check endpoint /debug/health start
// returning a 500 and list the specific check that failed.
//
// Assuming you wish to register a method called currentMinuteEvenCheck()
// error you could do that by doing:
//
// health.Register("even_minute", health.CheckFunc(currentMinuteEvenCheck))
// CheckFunc is a convenience type that implements Checker.
//
// Another way of registering a check could be by using an anonymous function
// and the convenience method RegisterFunc. An example that makes the status
// endpoint always return an error:
//
//  health.RegisterFunc("my_check", func() error {
//   return Errors.new("This is an error!")
// }))
// The recommended way of registering is, however, using a periodic Check.
// PeriodicChecks run on a certain schedule and asynchronously update the
// status of the check. This allows CheckStatus() to return without blocking
// on an expensive check.
//
// A trivial example of a check that runs every 5 seconds and shuts down our
// server if the current minute is even, could be added as follows:
//
//  func currentMinuteEvenCheck() error {
//    m := time.Now().Minute()
//    if m%2 == 0 {
//      return errors.New("Current minute is even!")
//    }
//    return nil
//  }
//
//  health.RegisterPeriodicFunc("minute_even", currentMinuteEvenCheck, time.Second*5)
// You could also use the health checker mechanism to ensure your application
// only comes up if certain conditions are met, or to allow the developer to
// take the service out of rotation immediately. An example that checks database
// connectivity and immediately takes the server out of rotation on err:
//
//  updater = health.NewStatusUpdater()
//   health.RegisterFunc("database_check", func() error {
//    return updater.Check()
//  }))
//
//  conn, err := Connect(...) // database call here
//  if err != nil {
//    updater.Update(errors.New("Error connecting to the database: " + err.Error()))
//  }
package health
