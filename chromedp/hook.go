package chromedp

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type Hooker struct {
	ReqHooks  map[string]map[string]interface{}
	RespHooks map[string]map[string]interface{}
}

func NewHooker(url, username, password, command string) (error, *Hooker) {
	hooker := &Hooker{
		ReqHooks:  make(map[string]map[string]interface{}),
		RespHooks: make(map[string]map[string]interface{}),
	}
	err := hooker.Hooking(url, username, password, command)
	if err != nil {
		return err, nil
	}
	return nil, hooker
}

func (h *Hooker) Hooking(url, username, password, command string) error {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		//chromedp.Flag("headless", false),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("window-size", "600,600"),
		chromedp.UserDataDir(dir),
	)

	times := 3
	for i := 0; i <= times; i++ {
		if i == times {
			return errors.New("punch failed too many times, please contact helpers")
		}

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		// also set up a custom logger
		taskCtx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		// create a timeout
		taskCtx, cancel = context.WithTimeout(taskCtx, 10*time.Second)
		defer cancel()

		// ensure that the browser process is started
		if err := chromedp.Run(taskCtx); err != nil {
			return err
		}
		err = chromedp.Run(taskCtx,
			fakeLocation(),
			network.Enable(),
			submit(url, `#username_input`, username, `#password-input`, password, `#login-button`, targetElementByCommand(command)),
		)
		if err != nil {
			continue
		}
		err = chromedp.Run(taskCtx, check(command))
		if err != nil {
			continue
		}
		if err == nil {
			break
		}
	}
	return nil
}

func fakeLocation() chromedp.Tasks {
	min := int64(-5000000000)
	max := int64(5000000000)
	float := math.Pow(10, -14)
	rand.Seed(time.Now().UnixNano())

	latRandomNumber := rand.Int63n(max-min) + min
	latRandomNumberFloat64 := float64(latRandomNumber) * float // -0.00005~0.00005, 緯度

	lngRandomNumber := rand.Int63n(max-min) + min
	lngRandomNumberFloat64 := float64(lngRandomNumber) * float // -0.00005~0.00005, 經度

	permissions := &browser.GrantPermissionsParams{
		Permissions: []browser.PermissionType{browser.PermissionTypeGeolocation},
	}

	geolocations := &emulation.SetGeolocationOverrideParams{
		Latitude:  (25.02283095064086) + latRandomNumberFloat64,
		Longitude: 121.54949954857622 + lngRandomNumberFloat64,
		Accuracy:  100,
	}
	return chromedp.Tasks{permissions, geolocations}
}

func submit(url, accountElement, accountVal, passwordElement, passwordVal, subElement, targetElement string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(accountElement, chromedp.ByID),
		chromedp.SendKeys(accountElement, accountVal, chromedp.ByID),
		chromedp.WaitVisible(passwordElement, chromedp.ByID),
		chromedp.SendKeys(passwordElement, passwordVal, chromedp.ByID),
		chromedp.WaitVisible(passwordElement, chromedp.ByID),
		chromedp.Submit(subElement),
		chromedp.WaitVisible(targetElement, chromedp.ByID),
		chromedp.Click(targetElement, chromedp.ByID),
		chromedp.WaitVisible(passwordElement, chromedp.ByID),
	}
}

func check(command string) chromedp.Tasks {
	element := ""
	if command == "/punch_in" {
		element = "/html/body/div[3]/div/div[1]/div[5]/div[1]/div/div[2]/div[1]/div[1]/div/div"
	}
	if command == "/punch_out" {
		element = "/html/body/div[3]/div/div[1]/div[5]/div[1]/div/div[2]/div[1]/div[2]/div/div"
	}
	return chromedp.Tasks{
		chromedp.WaitVisible(element, chromedp.BySearch),
	}
}

func targetElementByCommand(command string) string {
	if command == "/punch_in" {
		return "#clockin"
	}
	if command == "/punch_out" {
		return "#clockout"
	}
	fmt.Errorf("not found target, command=%s\n", command)
	return ""
}
