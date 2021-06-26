package webrobby

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

type Element struct {
	el             selenium.WebElement
	currentTimeout time.Duration
}

func (e *Element) Click() {
	err := e.el.Click()
	if err != nil {
		panic(errors.Wrap(err, "while clicking on element"))
	}
}

func (e *Element) Type(text string) {
	err := e.el.SendKeys(text)
	if err != nil {
		panic(errors.Wrap(err, "while sending keys"))
	}
}
