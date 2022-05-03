package chromedp

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
	"io/ioutil"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
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

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 1*time.Minute)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		return err
	}

	var ret string
	err = chromedp.Run(taskCtx,
		fakeLocation(),
		network.Enable(),
		submit(url, `#username_input`, username, `#password-input`, password, `#login-button`, targetElementByCommand(command), &ret),
	)
	if err != nil {
		return err
	}
	return nil
}

func fakeLocation() chromedp.Tasks {
	permissions := &browser.GrantPermissionsParams{
		Permissions: []browser.PermissionType{browser.PermissionTypeGeolocation},
		//Origin:      "https://cloud.nueip.com/home",
	}
	geolocations := &emulation.SetGeolocationOverrideParams{
		Latitude:  25.02283095064086,
		Longitude: 121.54949954857622,
		Accuracy:  99.999,
	}
	return chromedp.Tasks{permissions, geolocations}
}

func submit(url, accountElement, accountVal, passwordElement, passwordVal, subElement, targetElement string, res *string) chromedp.Tasks {
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
