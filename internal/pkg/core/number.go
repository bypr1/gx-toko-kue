package core

import (
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
	"time"
)

type NumberGeneratorInterface interface {
	Generate() string
	Prefix() string
}

type NumberGenerator struct {
	Model schema.Tabler
}

func (nm NumberGenerator) Generate(generator NumberGeneratorInterface) string {
	var number string

	increment := GetIncrementMonthly(nm.Model)

	number += time.Now().Format("010602")
	number += StrPadLeft(strconv.Itoa(int(increment)), 4, '0')
	number += strconv.Itoa(RandInt(111, 999))

	return strings.ToUpper(generator.Prefix() + number)
}
