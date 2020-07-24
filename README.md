# Presto 🎩

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/benpate/presto)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/presto?style=flat-square)](https://goreportcard.com/report/github.com/benpate/presto)
[![Build Status](http://img.shields.io/travis/benpate/presto.svg?style=flat-square)](https://travis-ci.org/benpate/presto)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/presto.svg?style=flat-square)](https://codecov.io/gh/benpate/presto)
![Version](https://img.shields.io/github/v/release/benpate/presto?include_prereleases&style=flat-square&color=brightgreen)

## Magical REST interfaces for Go

Presto is a thin wrapper library that helps structure and simplify the REST interfaces you create in [Go](https://golang.org). Its purpose is to encapsulate all of the boilerplate code that is commonly required to publish a server-side service via a REST interface.  Using Presto, your route configuration code looks like this:

#### main.go
```go
// ROUTER CONFIGURATION

// Presto requires the echo router by LabStack.  So first, let's pass in a new instance of echo.
presto.UseRouter(echo.New())

// Define a new service to expose online as a REST collection. (Services, Factories, Scopes and Roles defined below)
presto.NewCollection(NoteFactory, "/notes").
    List().
    Post(role.InRoom).
    Get(role.InRoom).
    Put(role.InRoom, role.Owner).
    Patch(role.InRoom, role.Owner).
    Delete(role.InRoom, role.Owner).
    Method("action-name", customHandler, role.InRoom, role.CustomValue)
```

## Design Philosophy

### Clean Architecture

Presto lays the groundwork to implement a REST API according to the [CLEAN architecture, first published by "Uncle Bob" Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).  This means decoupling business logic and databases, by injecting dependencies down through your application.  To do this in a type-safe manner, Presto requires that your services and objects fit into its interfaces, which  describe minimal behavior that each must support in order to be used by Presto.

Presto also uses the [data package](https://github.com/benpate/data) as an abstract representation of some common database concepts, such as query criteria.  This allows you to swap in any database by building an adapter that implements the `data` interfaces.  Once Presto is able to work with your business logic in an abstract way, the rest of the common code is repeated for each API endpoint you need to create.

### REST API Design Rulebook

Presto works hard to implement REST APIs according to the patterns laid out in the ["REST API Design Rulebook", by Mark Massé](https://smile.amazon.com/REST-Design-Rulebook-Mark-Masse/dp/1449310508/). This means:

* Clear route names
* Using HTTP methods (GET, PUT, POST, PATCH, DELETE) to determine the action being taken
* Using POST and URL route parameters for other API endpoints that don't fit neatly into the standard HTTP method definitions.

### Minimal Dependencies

Presto's only dependency is on the [fast and fabulous Echo router](https://github.com/labstack/echo), which is an open-source package for creating HTTP servers in Go.  Our ultimate goal with this package is to remove this as a hard dependency eventually, and refactor this code to work with multiple routers in the Go ecosystem.

## Services

Presto does not replace your application business logic.  It only exposes your internal services via a REST API. Each endpoint must be linked to a corresponding service (that matches Presto's required interface) to handle the actual loading, saving, and deleting of objects.

## Factories

The specific work of creating services and objects is pushed out to a Factory object, which provides a map of your complete domain.  The factories also manage dependencies (such as a live database connection) for each service that requires it.  Here's an example factory:

## REST Endpoints: Defaults

Presto implements six standard REST endpoints that are defined in the REST API Design Rulebook, and should serve a majority of your needs.

### List

### Post

### Get

### Put

### Patch

### Delete

## REST Endpoints: Custom Methods

There are many cases where these six default endpoints are not enough, such as when you have to initiate a specific transaction.  A good example of this is a "checkout" function in a shopping cart.  The REST API Design Rulebook labels these actions as "Methods", and states that these transactions should always be registered as a POST handler.  Presto helps you to manage these functions as well, using the following calls:

#### main.go
```go
// The following code will register a POST handler on the
// route `/cart/checkout`, using the function `CheckoutHandler`
presto.NewCollection(echo.Echo, factory.Cart, "/cart").
    Method("/checkout", CheckoutHandler, roles)
```

## Scopes and Database Criteria

Your REST server should be able to limit the records accessed though the website, for instance, hiding records that have been virtually deleted, or limiting users in a multi-tenant database to only see the records for their virtual account.  Presto accomplishes this using `scopes`, and `ScopeFuncs` which are functions that inspect the `echo.Context` and return a `data.Expression` that limits users access.  The [data](https://github.com/benpate/data) package is used to create an intermediate representation of the query criteria that can then be interpreted into the specific formats used by your database system.  Here's an example of some ScopeFunc functions.

#### main.go

```go
// This overrides the default scoping function, and uses the
// NotDeleted function for all routes in your API instead.
presto.UseScope(scope.NotDeleted)

// This configures this specific collection to limit all
// database queries using the `ByUsername` scope, in addition
// to the globally defined `NotDeleted` scope declared above.
presto.NewCollection(e, PersonFactory, "/person").
    UseScope(scope.ByUsername)
```


#### scopes/scopes.go
```go
// NotDeleted filters out all records that have not been
// "virtually deleted" from the database.
func NotDeleted(ctx echo.Context) (data.Expression, *derp.Error) {
    return data.Expression{{"journal.deleteDate", data.OperatorEqual, 0}}, nil
}

// ByPersonID uses the route Param "personId" to limit
// requests to records that include that personId only.
func Route(ctx echo.Context) (data.Expression, *derp.Error) {

    personID := ctx.Param("personId")

    // If the personID is empty, then return an error to the caller..
    if personID == "" {
        return data.Expression{}, derp.New(derp.CodeBadRequestError, "example.Route", "Empty PersonID", personID)
    }

    // Convert the parameter value into a bson.ObjectID and return the expression
    if personID, err := primitive.ObjectIDFromHex(personID); err != nil {
        return data.Expression{{"personId", data.OperatorEqual, personId}}, nil
    }

    // Fall through to here means that we couldn't convert the personID into a valid ObjectID.  Return an error.
    return data.Expression{}, derp.New(derp.CodeBadRequestError, "example.Route", "Invalid PersonID", personID)
}
```

## User Roles

It's very likely that your API requires custom authentication and authorization for each endpoint.  Since this is very custom to your application logic and environment, Presto can't automate this for you.  But, Presto does make it very easy to organize the permissions for each endpoint into a single, readable location.  Authorization requirements for each endpoint are baked into common functions called roles, and then passed in to Presto during system configuration.

#### main.go
```go
// Sets up a new collection, where the user must have permissions
// to post into the Room.  This is handled by the `InRoom` function.
presto.NewCollection(echo.Echo, NoteFactory, "/notes").
    Post(role.InRoom)
```

#### roles/roles.go
```go
// InRoom determines if the requester has access to the Room in
// which this object resides. If so, then access to it is valid,
// so return a TRUE.
func InRoom(ctx echo.Context, object Object) bool {

    // Get the list of rooms that this user has access to..
    // For example, using JWT tokens in the context request headers.
    allowedRoomIDs := getRoomListFromContext(ctx)

    // Uses a type switch to retrieve the roomID from the Object interface.
    roomID, err := getRoomIDFromObject(object)

    if err != nil {
        return false
    }

    // Try to find the object.RoomID in the list of allowed rooms.
    for _, allowedRoomID := range allowedRoomIDs {
        if allowedRoomID == roomID {
            return true // If so, then you're in.
        }
    }

    // Otherwise, you are not permitted to access this object.
    return false;
}
```

## Performance: Caching, ETag Support

Presto uses ETags to dramatically improve performance and consistency of your REST API.  This requires client support as well, so if your client does not include ETag information with your REST requests, then this code is effectively skipped.

### 304 Not Modified

HTTP includes a great way to minimize bandwidth and latency, using [304 Not Modified](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/304) responses.  Presto can use ETags to determine if a resource has not been changed since it was last delivered to the client, and will send **304 Not Modified** responses when it can.

### Pluggable Cache Engines

Presto provides an interface for you to plug in your own caching system.  Caches only store resource URIs and the most recent ETag.  If a request's ETags match the value in the cache, then Presto can skip the database load entirely and deliver a simple 304 status code.

### Using ETags for Optimistic Locking

ETags are also useful to implement [optimistic locking](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) on records.  If the client sends ETag information along with a PUT, PATCH, or DELETE method, then this ETag is compared with the current value in the record.  If the ETags do not match, then the record has been modified since the client's last read, and the transaction is rejected.

Remember, this is an optional feature.  If your client does not include ETags with these transactions, then the logic for optimistic locking is simply skipped.

### Implementing ETags in your Domain Model

The [data library](https://github.com/benpate/data) includes an optional `Journal` object that implements *most* of the `Object` interface that Presto needs in order to operate.  The `data.Journal` object also includes a simple mechanism for reading and writing ETags into every record you create.  You're welcome to use this implementation, or to create one that suits your needs better.

## Pull Requests Welcome

Original versions of this library have been used in production on commercial applications for years, and greatly reduced the amount of work required to **create** and **maintain** a well-structured REST API.

This new, open-sourced version of PRESTO will greatly benefit from your experience reports, use cases, and contributions.  If you have an idea for making Rosetta better, send in a pull request.  We're all in this together! 🎩