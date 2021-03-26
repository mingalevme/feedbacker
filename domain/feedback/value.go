package feedback

import (
	"regexp"
	"strconv"
	"strings"
)

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

type FeedbackId interface {
	valueObject
	Get() int
}

func NewFeedbackId(id int) (FeedbackId, error) {
	if id < 1 {
		return nil, IdIsInvalid
	}
	return &feedbackId{
		value: id,
	}, nil
}

type feedbackId struct {
	value int
}

func (s *feedbackId) Equals(value interface{}) bool {
	target, ok := value.(FeedbackId)
	return ok && s.String() == target.String()
}

func (s *feedbackId) Get() int {
	return s.value
}

func (s *feedbackId) String() string {
	return strconv.Itoa(s.Get())
}

// ---------------------------------------------------------------------------------------------------------------------

type ServiceValue interface {
	stringValueObject
}

type serviceValue struct {
	abstractStringValueObject
}

func NewServiceValue(value string) (ServiceValue, error) {
	origin := value
	value = strings.TrimSpace(value)
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

func NewEditionValue(value string) (EditionValue, error) {
	origin := value
	value = strings.TrimSpace(value)
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

func NewTextValue(value string) (TextValue, error) {
	origin := value
	value = strings.TrimSpace(value)
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