package service_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/highlanderticketing/service"
)

func TestCreateUserIntegration(t *testing.T) {
	user := &model.User{
		Email:      "test@example.com",
		Name:       "John",
		FamilyName: "Doe",
		IsAdmin:    false,
	}

	// Testfall 1: der benutzer wird angelegt
	err2 := service.CreateUser(user)

	assert.Nil(t, err2)
	assert.False(t, user.IsAdmin)

	if !reflect.DeepEqual(t, user) {
		t.Errorf("Expected %+v, but got %+v", t, err2)
	}

	// Testfall 2: es gibt bereits den Benutzer
	err3 := service.CreateUser(user)

	assert.Error(t, err3)
	assert.Equal(t, "Der Benutzer existiert bereits", err3.Error())

	if !reflect.DeepEqual(err2, user) {
		t.Errorf("Expected %+v, but got %+v", t, err2)
	}
}
