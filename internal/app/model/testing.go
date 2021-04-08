// +build testing

// go test -v -cover -tags testing some/package/*.go

package model

import (
	"crypto/md5"
	"fmt"
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"math/rand"
	"strconv"
	"time"
)

func MakeFeedback() Feedback {
	context := MakeContext()
	customer := MakeCustomer()
	t := time.Unix(time.Now().Unix() - int64(rand.Intn(60*60*24*30)), 0)
	return Feedback{
		ID:        int(rand.Uint64()),
		Service:   "feedbacker",
		Edition:   strutils.StrToPointerStr("ru-ru"),
		Text:      "Hello, World!",
		Context:   &context,
		Customer:  &customer,
		CreatedAt: t,
		UpdatedAt: t,
	}
}

func MakeContext() Context {
	return Context{
		AppVersion:  strutils.StrToPointerStr(fmt.Sprintf("%d.%d.%d", rand.Intn(10)+1, rand.Intn(10)+1, rand.Intn(10)+1)),
		AppBuild:    strutils.StrToPointerStr(fmt.Sprintf("%d", rand.Intn(100)+1)),
		OsName:      strutils.StrToPointerStr("iOS"),
		OsVersion:   strutils.StrToPointerStr(fmt.Sprintf("%d.%d.%d", rand.Intn(10)+1, rand.Intn(10)+1, rand.Intn(10)+1)),
		DeviceBrand: strutils.StrToPointerStr("Apple"),
		DeviceModel: strutils.StrToPointerStr(fmt.Sprintf("iPhone %d", rand.Intn(20)+1)),
	}
}

func MakeCustomer() Customer {
	t := time.Now().Nanosecond()
	h := md5.Sum([]byte(strconv.Itoa(t)))
	return Customer{
		Email:          strutils.StrToPointerStr(fmt.Sprintf("%x@example.com", h)),
		InstallationID: strutils.StrToPointerStr(fmt.Sprintf("%x@example.com", h)),
	}
}
