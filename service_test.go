package presto

import (
	"context"
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/data/expression"
	"github.com/benpate/data/mock"
	"github.com/benpate/data/option"
	"github.com/benpate/derp"
)

// TEST METHODS

func TestServiceList(t *testing.T) {

	s := getTestPersonService()

	it, _ := s.ListObjects(nil, option.SortAsc("name"))

	person := s.New()

	for it.Next(person) {
		t.Log(person)
	}

	t.Fail()
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

	if err := service.session.Load("Persons", criteria, person); err != nil {
		return nil, derp.Wrap(err, "testPersonService.Load", "Error Loading Person")
	}

	return person, nil
}

func (service *testPersonService) SaveObject(person data.Object, note string) *derp.Error {

	if err := service.session.Save("Persons", person, note); err != nil {
		return derp.Wrap(err, "testPersonService.Save", "Error Saving Person", person)
	}

	return nil
}

func (service *testPersonService) DeleteObject(person data.Object, note string) *derp.Error {

	if err := service.session.Delete("Persons", person, note); err != nil {
		return derp.Wrap(err, "testPersonService.Delete", "Error Deleting Person", person)
	}

	return nil
}

func (service *testPersonService) ListObjects(criteria expression.Expression, options ...option.Option) (data.Iterator, *derp.Error) {

	return service.session.List("Persons", criteria, options...)
}

func (service *testPersonService) Close() {}

// Prepopulate Database
func getTestPersonService() *testPersonService {

	session := mock.New().Session(context.TODO())

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
