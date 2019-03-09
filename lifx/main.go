package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	//	"time"

	"github.com/Jeffail/gabs"
	"github.com/csmu-cenr/golifx"
)

func load(path string) ([]*golifx.Bulb, error) {

	var result []*golifx.Bulb = nil
	var err error
	var bulbs *gabs.Container
	var bulb *golifx.Bulb
	var children []*gabs.Container
	var label string
	var power_state bool
	var mac string
	var ok bool

	file, _ := ioutil.ReadFile(path)
	bulbs, err = gabs.ParseJSON(file)
	if err == nil {
		children, _ = bulbs.Children()
		for _, child := range children {
			fmt.Println(child.String())
			bulb = &golifx.Bulb{}
			label, ok = child.Path("label").Data().(string)
			if ok {
				bulb.Label = label
			}
			power_state, ok = child.Path("power_state").Data().(bool)
			if ok {
				bulb.PowerState = power_state
			}
			mac, ok = child.Path("mac").Data().(string)
			if ok {
				bulb.SetHardwareAddressFromMacAddress(mac)
			}
			result = append(result, bulb)
			fmt.Println(bulb)
		}
	} else {
		fmt.Println(err)
	}

	return result, err
}

func export(path string) (bool, error) {

	var result bool = false
	var err error
	var bulb *golifx.Bulb
	var bulbs []*golifx.Bulb
	var bulbs_byte_array []byte

	bulbs, err = golifx.LookupBulbs()

	if err == nil {
		for _, bulb = range bulbs {
			_, _ = bulb.GetLabel()
			_, _ = bulb.GetPowerState()
		}

		bulbs_byte_array, err = json.MarshalIndent(bulbs, "", "    ")
		if err == nil {
			err = ioutil.WriteFile(path, bulbs_byte_array, 0644)
			if err == nil {
				result = true
			}
		}
	}

	return result, err
}

var bulbs []*golifx.Bulb

func main() {

	var arguments *gabs.Container
	var err error

	if len(os.Args) == 1 {

		arguments, _ = gabs.ParseJSON([]byte(`[{"error":"JSON parameter missing. Please supply valid JSON as the first parameter to lifx"}]`))
		fmt.Println(arguments.StringIndent("", "  "))

	} else {

		arguments, err = gabs.ParseJSON([]byte(os.Args[1]))

		if err == nil {

			var ok bool
			var path string

			children, _ := arguments.Children()
			for _, child := range children {
				path, ok = child.Path("export").Data().(string)
				if ok {
					export(path)
				}
				path, ok = child.Path("load").Data().(string)
				if ok {
					bulbs, err = load(path)
					fmt.Println(bulbs)
				}
			}

		} else {

			var error_json string = fmt.Sprintf("[{\"error\":\"%s\"}]", err)
			arguments, _ = gabs.ParseJSON([]byte(error_json))
			fmt.Println(arguments.StringIndent("", "  "))

		}

	}
}
