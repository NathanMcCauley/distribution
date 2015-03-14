// Package health provides a generic health checking framework.
// The health package works expvar style. By importing the package
// the debug server is getting a "/debug/health" endpoint that
// returns the status of all registered checks.
//
// A Check can either be run synchonously, or asynchronously.
// We recommend that most checks are registered as an asynchronous
// check, so a call to the "/debug/health" endpoint always returns
// immediately.
// This pattern is particularly useful for checks that verify
// upstream connectivity or database status, since they might take
// a long time to return/timeout.
//
// To install health, just import it in your application:
//
// import "github.com/docker/distribution/health"
//
// You can also (optionally) import health/api that will add two
// convenience endpoints: /debug/health/down and /debug/health/up.
// These endpoints add "manual" checks that allow the service to
// quickly be brought in/out of rotation.
//
// import _ "github.com/docker/distribution/registry/health/api"
//
// After importing these packages to your main application, you can
// start registering checks.
//
// The lowest-level way to register a check is by calling Register.
// Register allows you to pass in an arbitrary string and a Checker
// method that runs your check. If your method returns nil, it is
// considered a healthy check, otherwise it will make the health
// check endpoint "/debug/health" start returning a 500 and list
// the specific check that failed.
//
// Assuming you wish to register method currentMinuteEvenCheck you
// could do that by doing:
//
//  health.Register("even_minute", health.CheckFunc(currentMinuteEvenCheck))
//
// CheckFunc is a convenience type that implements Checker.
//
// Another way of registering a check could be by using an anonymous
// function and the convenience method RegisterFunc. An example that
// makes the status endpoint always return an error:
//
//  health.RegisterFunc("my_check", func() error {
//   return Errors.new("This is an error!")
// }))
//
// The reccomended way of registering checks, however, is by using a
// periodic Check. PeriodicChecks run on a certain schedulle and
// asynchronously update the status of the check. This allows CheckStatus()
// to return without blocking on an expensive check.
//
// A trivial example of a check that runs every 5 seconds and
// shuts down our server if the current minute is even, could be
// added as follows:
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
//
// Finally, you could also use the healthchecker mechanism to ensure
// your application only comes up if certain conditions are met,
// or to allow the developer to take the service out of rotation
// immediatelly. An example that checks database connectivity:
//
//  updater = health.NewStatusUpdater()
//   health.RegisterFunc("database_check", func() error {
//    return updater.Check()
//  }))
//
//  conn, err := Connect(...) // database call here
//  if err != nil {
//    updater.Update(errors.New("There was an error connecting to the database: " + err.Error()))
//  }
package health
