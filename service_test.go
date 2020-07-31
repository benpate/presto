package presto

import (
	"context"
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/data/expression"
	"github.com/benpate/data/mockdb"
	"github.com/benpate/data/option"
	"github.com/benpate/derp"
	"github.com/stretchr/testify/assert"
)

// TEST METHODS

func TestServiceList(t *testing.T) {

	s := getTestPersonService()

	it, _ := s.ListObjects(nil, option.SortAsc("name"))

	person := s.New()

	it.Next(person)
	assert.Equal(t, "John Connor", person.Name)

	it.Next(person)
	assert.Equal(t, "Kyle Reese", person.Name)

	it.Next(person)
	assert.Equal(t, "Sara Connor", person.Name)

	assert.False(t, it.Next(person))
}

func TestServiceLoad(t *testing.T) {

	s := getTestPersonService()

	{
		person, err := s.LoadObject(expression.New("personId", "=", "john"))

		john := person.(*testPerson)
		assert.Nil(t, err)
		assert.Equal(t, "john", john.PersonID)
		assert.Equal(t, "John Connor", john.Name)
		assert.Equal(t, "john@sky.net", john.Email)
	}

	{
		person, err := s.LoadObject(expression.New("personId", "=", "sara"))

		sara := person.(*testPerson)
		assert.Nil(t, err)
		assert.Equal(t, "sara", sara.PersonID)
		assert.Equal(t, "Sara Connor", sara.Name)
		assert.Equal(t, "sara@sky.net", sara.Email)
	}

	{
		person, err := s.LoadObject(expression.New("personId", "=", "kyle"))

		kyle := person.(*testPerson)
		assert.Nil(t, err)
		assert.Equal(t, "kyle", kyle.PersonID)
		assert.Equal(t, "Kyle Reese", kyle.Name)
		assert.Equal(t, "kyle@resistance.mil", kyle.Email)
	}
}

// SERVICE OBJECT
type testPersonService struct {
	session data.Session
}

func (service *testPersonService) New() *testPerson {
	return &testPerson{}
}

func (service *testPersonService) NewObject() data.Object {
	return &testPerson{}
}

func (service *testPersonService) LoadObject(criteria expression.Expression) (data.Object, *derp.Error) {

	person := service.NewObject()

	if err := service.session.Collection("Perons").Load(criteria, person); err != nil {
		return nil, derp.Wrap(err, "testPersonService.Load", "Error Loading Person")
	}

	return person, nil
}

func (service *testPersonService) SaveObject(person data.Object, note string) *derp.Error {

	if err := service.session.Collection("Perons").Save(person, note); err != nil {
		return derp.Wrap(err, "testPersonService.Save", "Error Saving Person", person)
	}

	return nil
}

func (service *testPersonService) DeleteObject(person data.Object, note string) *derp.Error {

	if err := service.session.Collection("Perons").Delete(person, note); err != nil {
		return derp.Wrap(err, "testPersonService.Delete", "Error Deleting Person", person)
	}

	return nil
}

func (service *testPersonService) ListObjects(criteria expression.Expression, options ...option.Option) (data.Iterator, *derp.Error) {

	return service.session.Collection("Perons").List(criteria, options...)
}

func (service *testPersonService) Close() {}

// Prepopulate Database
func getTestPersonService() *testPersonService {

	session, _ := mockdb.New().Session(context.TODO())
	service := testPersonService{session: session}

	{
		person := service.New()
		person.PersonID = "john"
		person.Name = "John Connor"
		person.Email = "john@sky.net"
		service.SaveObject(person, "Created")
	}

	{
		person := service.New()
		person.PersonID = "sara"
		person.Name = "Sara Connor"
		person.Email = "sara@sky.net"
		service.SaveObject(person, "Created")
	}

	{
		person := service.New()
		person.PersonID = "kyle"
		person.Name = "Kyle Reese"
		person.Email = "kyle@resistance.mil"
		service.SaveObject(person, "Created")
	}

	return &testPersonService{session: session}
}
