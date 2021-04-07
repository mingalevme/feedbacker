package feedback

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

//var ErrUnprocessableEntity = errors.New("unprocessable entity")

//var IdIsInvalid = errors.Wrap(ErrUnprocessableEntity, "id is invalid")
//var ServiceIsInvalid = errors.Wrap(ErrUnprocessableEntity, "service is invalid")
//var EditionIsInvalid = errors.Wrap(ErrUnprocessableEntity, "edition is invalid")
//var TextIsInvalid = errors.Wrap(ErrUnprocessableEntity, "text is invalid")
//var TextIsEmpty = errors.Wrap(ErrUnprocessableEntity, "text is empty")

var IdIsInvalid = errors.New("id is invalid")
var ServiceIsInvalid = errors.New("service is invalid")
var EditionIsInvalid = errors.New("edition is invalid")
var TextIsInvalid = errors.New("text is invalid")
var TextIsEmpty = errors.New("text is empty")

// ---------------------------------------------------------------------------------------------------------------------

type valueObject interface {
	Equals(value interface{}) bool
	String() string
}

// ---------------------------------------------------------------------------------------------------------------------

type stringValueObject interface {
	valueObject
	Get() string
}

type abstractStringValueObject struct {
	value string
}

func (s abstractStringValueObject) Equals(value interface{}) bool {
	target, ok := value.(abstractStringValueObject)
	return ok && s.value == target.value
}

func (s abstractStringValueObject) Get() string {
	return s.value
}

func (s abstractStringValueObject) String() string {
	return s.Get()
}

// ---------------------------------------------------------------------------------------------------------------------

type FeedbackID interface {
	valueObject
	Get() int
}

func NewFeedbackID(id interface{}) (FeedbackID, error) {
	var value int
	var err error
	switch v := id.(type) {
	case *feedbackID:
		if v == nil {
			return nil, IdIsInvalid
		}
		return v, nil
	case int:
		value = v
	case string:
		value, err = strconv.Atoi(v)
		if err != nil {
			return nil, IdIsInvalid
		}
	default:
		return nil, IdIsInvalid
	}
	if value < 1 {
		return nil, IdIsInvalid
	}
	return &feedbackID{
		value: value,
	}, nil
}

type feedbackID struct {
	value int
}

func (s *feedbackID) Equals(value interface{}) bool {
	target, ok := value.(FeedbackID)
	return ok && s.String() == target.String()
}

func (s *feedbackID) Get() int {
	return s.value
}

func (s *feedbackID) String() string {
	return strconv.Itoa(s.Get())
}

// ---------------------------------------------------------------------------------------------------------------------

type ServiceValue interface {
	stringValueObject
}

type serviceValue struct {
	abstractStringValueObject
}

func NewServiceValue(service interface{}) (ServiceValue, error) {
	var origin string
	switch v := service.(type) {
	case *serviceValue:
		if v == nil {
			return nil, ServiceIsInvalid
		}
		return v, nil
	case string:
		origin = v
	default:
		return nil, ServiceIsInvalid
	}
	value := strings.TrimSpace(origin)
	value = strings.Trim(value, "-_")
	value = strings.ToLower(value)
	if value != origin {
		return nil, ServiceIsInvalid
	}
	match, _ := regexp.MatchString("^[a-z-_]+$", value)
	if !match {
		return nil, ServiceIsInvalid
	}
	return &serviceValue{
		abstractStringValueObject{value: value},
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type EditionValue interface {
	stringValueObject
}

func NewEditionValue(edition interface{}) (EditionValue, error) {
	var origin string
	switch v := edition.(type) {
	case *editionValue:
		if v == nil {
			return nil, EditionIsInvalid
		}
		return v, nil
	case string:
		origin = v
	default:
		return nil, EditionIsInvalid
	}
	value := strings.TrimSpace(origin)
	value = strings.Trim(value, "-_")
	value = strings.ToLower(value)
	if value != origin {
		return nil, EditionIsInvalid
	}
	match, _ := regexp.MatchString(`^\w+-\w+$`, value)
	if !match {
		return nil, EditionIsInvalid
	}
	return &editionValue{
		abstractStringValueObject{value: value},
	}, nil
}

type editionValue struct {
	abstractStringValueObject
}

// ---------------------------------------------------------------------------------------------------------------------

type TextValue interface {
	stringValueObject
}

func NewTextValue(text interface{}) (TextValue, error) {
	var origin string
	switch v := text.(type) {
	case *textValue:
		if v == nil {
			return nil, TextIsInvalid
		}
		return v, nil
	case string:
		origin = v
	default:
		return nil, TextIsInvalid
	}
	value := strings.TrimSpace(origin)
	if value != origin {
		return nil, TextIsInvalid
	}
	if value == "" {
		return nil, TextIsEmpty
	}
	return &textValue{
		abstractStringValueObject{value: value},
	}, nil
}

type textValue struct {
	abstractStringValueObject
}
