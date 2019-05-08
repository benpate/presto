# Presto 🎩

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/benpate/presto)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/presto?style=flat-square)](https://goreportcard.com/report/github.com/benpate/presto)
[![Build Status](http://img.shields.io/travis/benpate/presto.svg?style=flat-square)](https://travis-ci.org/benpate/presto)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/presto.svg?style=flat-square)](https://codecov.io/gh/benpate/presto)

## Magical REST interfaces for Go

Presto is a thin wrapper library that helps structure and simplify the REST interfaces you create in [Go](https://golang.org). Its purpose is to encapsulate all of the boilerplate code that is commonly required to publish a server-side service via a REST interface.  Using Presto, your route configuration code looks like this:

```go

// Presto requires the echo router by Labstack.  Let's create a new instance of echo first.
e := echo.New()

// Define a new service to expose online as a REST collection. (Services, Factories, Scopes and Roles defined below)
presto.NewCollection(e, NoteFactory, "/notes").
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

Presto lays the groundword to implement a REST API according to the [CLEAN architecture, first published by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).  This means decoupling business logic and databases, by injecting dependencies down through your application.  To do this in a type-safe manner, Presto requires that your services and objects fit into its interfaces, which  describe minimal behavior that each must support in order to be used by Presto.

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

```go

// The following code will registera POST handler on the route `/cart/checkout`, using the function `CheckoutHandler`
presto.NewCollection(echo.Echo, factory.Cart, "/cart").
    Method("/checkout", CheckoutHandler, roles)
```

## Scopes and Database Criteria

Your REST server should be able to limit the records accessed though the website, for instance, hiding records that have been virtually deleted, or limiting users in a multi-tennant database to only see the records for their virtual account.  Presto accomplishes this using `scopes`, and `ScopeFuncs` which are functions that inspect the `echo.Context` and return a `data.Expression` that limits users access.  The [data](https://github.com/benpate/data) package is used to create an intermediate representation of the query criteria that can then be interpreted into the specific formats used by your database system.  Here's an example of some ScopeFunc functions.

```go

/// IN YOUR `SCOPES` PACKAGE

// NotDeleted filters out all records that have not been "virtually deleted" from the database.
func NotDeleted(ctx echo.Context) (data.Expression, *derp.Error) {
    return data.Expression{{"journal.deleteDate", data.OperatorEqual, 0}}, nil
}


// ByPersonID uses the route Param "personId" to limit requests to records that include that personId only.
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


/// IN YOUR PRESTO CONFIGURATION

// This overrides the default scopeing function, and uses the
// NotDeleted function for all routes in your API instead.
presto.WithScope(scope.NotDeleted)

// This configures this specific collection to limit all
// database queries using the `ByUsername` scope, in addition
// to the globally defined `NotDeleted` scope declared above.
presto.NewCollection(e, PersonFactory, "/person").
    WithScope(scope.ByUsername)
```

## User Roles

It's very likely that your API requires custom authentication and authorization for each endpoint.  Since this is very custom to your application logic and environment, Presto can't automate this for you.  But, Presto does make it very easy to organize the permissions for each endpoint into a single, readable location.  Authorization requirements for each endpoint are baked into common functions called roles, and then passed in to Presto during system configuration.


```go

/// IN YOUR ROUTE CONFIGURATION

presto.NewCollection(echo.Echo, NoteFactory, "/notes").
    Post(role.InRoom) // The user must have permissions to post into the Room into which they're posting.

// IN YOUR ROLES PACKAGE

// InRoom determines if the requester has access to the Room in which this object resides.
// If so, then access to it is valid, so return a TRUE..  If not, then return a FALSE.
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

### Cache

### ETags

### Using ETags for Optimistic Locking