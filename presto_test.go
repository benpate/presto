package presto

import (
	"context"
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/data/expression"
	"github.com/benpate/data/mockdb"
	"github.com/benpate/remote"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPresto(t *testing.T) {

	db := mockdb.New()

	go startTestServer(db)

	session := db.Session(context.TODO())

	// Verify that the server is running.
	if err := remote.Get("http://localhost:8080/").Send(); err != nil {
		err.Report()
		assert.Fail(t, "Error getting default route", err)
	}

	///////////////////
	// POST a record
	john := testPerson{
		PersonID: "jc123",
		Name:     "John Connor",
		Email:    "john@sky.net",
	}

	sarah := testPerson{
		PersonID: "sc456",
		Name:     "Sarah Connor",
		Email:    "sarah@sky.net",
	}

	person := testPerson{}

	// Post a record to the "remote" server
	t1 := remote.Post("http://localhost:8080/persons").
		JSON(john)

	if err := t1.Send(); err != nil {
		err.Report()
		assert.Fail(t, "Error posting to localhost", err)
	}

	// Confirm that the record was sent/saved correctly.
	criteria := expression.New("personId", "=", john.PersonID)

	if err := session.Load("Persons", criteria, &person); err != nil {
		err.Report()
		assert.Fail(t, "Error loading new record from db", err)
	}

	assert.Equal(t, john.PersonID, person.PersonID)
	assert.Equal(t, john.Name, person.Name)
	assert.Equal(t, john.Email, person.Email)

	//////////////////////////
	// PUT a record

	t2 := remote.Put("http://localhost:8080/persons/" + sarah.ID()).
		JSON(sarah)

	if err := t2.Send(); err != nil {
		err.Report()
		assert.Fail(t, "Error PUT-ing a record", sarah)
	}

	// Confirm that the record was sent/saved correctly.
	criteria = expression.New("personId", "=", sarah.PersonID)
	if err := session.Load("Persons", criteria, &person); err != nil {
		err.Report()
		assert.Fail(t, "Error loading new record", err)
	}

	assert.Equal(t, sarah.PersonID, person.PersonID)
	assert.Equal(t, sarah.Name, person.Name)
	assert.Equal(t, sarah.Email, person.Email)

	////////////////////////
	// GET records

	// Load John
	t3 := remote.Get("http://localhost:8080/persons/"+john.PersonID).
		Response(&person, nil)

	if err := t3.Send(); err != nil {
		err.Report()
		assert.Fail(t, "Error retrieving person from REST service")
	}

	assert.Equal(t, john.PersonID, person.PersonID)
	assert.Equal(t, john.Name, person.Name)
	assert.Equal(t, john.Email, person.Email)

	// Load Sarah

	t4 := remote.Get("http://localhost:8080/persons/"+sarah.PersonID).
		Response(&person, nil)

	if err := t4.Send(); err != nil {
		err.Report()
		assert.Fail(t, "Error retrieving person from REST service")
	}

	assert.Equal(t, sarah.PersonID, person.PersonID)
	assert.Equal(t, sarah.Name, person.Name)
	assert.Equal(t, sarah.Email, person.Email)

	{
		// UPDATE RECORDS

		sarah := testPerson{
			PersonID: "sc456",
			Name:     "Sarah Connor",
			Email:    "sarahs-new-email@sky.net",
		}

		txn := remote.Put("http://localhost:8080/persons/" + sarah.ID()).
			JSON(sarah)

		if err := txn.Send(); err != nil {
			err.Report()
			assert.Fail(t, "Error PUT-ing a record", sarah)
		}

		// Confirm that the record was sent/saved correctly.
		criteria = expression.New("personId", "=", sarah.PersonID)
		if err := session.Load("Persons", criteria, &person); err != nil {
			err.Report()
			assert.Fail(t, "Error loading new record", err)
		}

		assert.Equal(t, sarah.PersonID, person.PersonID)
		assert.Equal(t, sarah.Name, person.Name)
		assert.Equal(t, sarah.Email, person.Email)
	}

	{
		// DELETE RECORDS

		txn1 := remote.Delete("http://localhost:8080/persons/" + john.ID())

		if err1 := txn1.Send(); err1 != nil {
			err1.Report()
			assert.Fail(t, "Error DELETE-ing a record", sarah)
		}

		txn2 := remote.Get("http://localhost:8080/persons/" + john.ID())
		err2 := txn2.Send()

		assert.NotNil(t, err2)
		assert.Equal(t, 404, err2.Code)

		txn3 := remote.Delete("http://localhost:8080/persons/" + john.ID())
		err3 := txn3.Send()

		assert.NotNil(t, err3)
		assert.Equal(t, 404, err3.Code)

	}
}

// FACTORY OBJECT

func testFactory(db data.Datastore) ServiceFunc {

	return func(ctx context.Context) Service {

		return &testPersonService{
			session: db.Session(ctx),
		}
	}
}

func startTestServer(db data.Datastore) {

	UseScopes()

	factory := testFactory(db)

	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.NoContent(200)
	})

	UseRouter(e)

	NewCollection(factory, "/persons").
		UseToken("personId").
		UseScopes().
		Post().
		Get().
		Put().
		Patch().
		Delete()

	e.Start(":8080")
}
