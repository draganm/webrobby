package webrobby

import (
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

type Webrobby struct {
	wd             selenium.WebDriver
	currentTimeout time.Duration
}

func New(wd selenium.WebDriver) *Webrobby {
	return &Webrobby{
		wd:             wd,
		currentTimeout: 4 * time.Second,
	}
}

func (w *Webrobby) Visit(url string) {
	err := w.wd.Get(url)
	if err != nil {
		panic(errors.Wrapf(err, "while visiting %s", url))
	}
}

func (w *Webrobby) currentTimeoutBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = w.currentTimeout
	b.InitialInterval = 50 * time.Millisecond
	return b
}

func (w *Webrobby) FindElement(cssSelector string) *Element {
	var el selenium.WebElement
	err := backoff.Retry(
		func() error {
			var err error
			el, err = w.wd.FindElement(selenium.ByCSSSelector, cssSelector)
			return err
		},
		w.currentTimeoutBackoff(),
	)

	if err != nil {
		panic(errors.Wrapf(err, "while waiting to find"))
	}

	return &Element{
		el:             el,
		currentTimeout: w.currentTimeout,
	}
}

func (w *Webrobby) FindElementWithText(cssSelector, text string) *Element {
	var el selenium.WebElement
	err := backoff.Retry(
		func() error {
			elements, err := w.wd.FindElements(selenium.ByCSSSelector, cssSelector)
			if err != nil {
				return err
			}

			for _, e := range elements {
				txt, err := e.Text()
				if err != nil {
					return errors.Wrap(err, "while finding text")
				}
				if strings.Contains(txt, text) {
					el = e
					return nil
				}
			}

			return errors.New("no elements found")
		},
		w.currentTimeoutBackoff(),
	)

	if err != nil {
		panic(errors.Wrapf(err, "while waiting to find"))
	}

	return &Element{
		el:             el,
		currentTimeout: w.currentTimeout,
	}

}

func (w *Webrobby) DeleteAllCookies() {
	err := w.wd.DeleteAllCookies()
	if err != nil {
		panic(errors.Wrap(err, "while deleting all cookies"))
	}
}

func (w *Webrobby) Title() string {
	t, err := w.wd.Title()
	if err != nil {

		panic(errors.Wrap(err, "while getting title"))
	}
	return t
}
