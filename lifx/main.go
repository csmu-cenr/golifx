package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

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

var bulb *golifx.Bulb
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
			var power_state bool
			var set_power_state bool
			var mac string
			var find_by_mac bool
			var sleep_float64 float64
			var sleep int64
			var sleep_awhile bool
			var bulbs_by_mac_address map[string]*golifx.Bulb

			bulbs_by_mac_address = make(map[string]*golifx.Bulb)

			children, _ := arguments.Children()
			for _, child := range children {
				path, ok = child.Path("export").Data().(string)
				if ok {
					export(path)
				}
				path, ok = child.Path("load").Data().(string)
				if ok {
					bulbs, err = load(path)
					for _, bulb := range bulbs {
						bulbs_by_mac_address[bulb.MacAddress()] = bulb
					}
				}
			}
			for _, child := range children {
				fmt.Println(child)
				mac, find_by_mac = child.Path("mac").Data().(string)
				power_state, set_power_state = child.Path("power_state").Data().(bool)
				sleep_float64, sleep_awhile = child.Path("sleep").Data().(float64)
				if sleep_awhile {
					fmt.Println("Sleeping", sleep, time.Now())
					sleep = int64(sleep_float64)
					time.Sleep(time.Duration(sleep) * time.Millisecond)
					fmt.Println("Awake after", sleep, time.Now())
				}
				if find_by_mac {
					bulb, ok := bulbs_by_mac_address[mac]
					if ok {
						if set_power_state {
							bulb.SetPowerState(power_state)
						}
					} else {
						fmt.Println("bulb.MacAddress: %s was not found.", mac)
					}
				}
			}

		} else {

			var error_json string = fmt.Sprintf("[{\"error\":\"%s\"}]", err)
			arguments, _ = gabs.ParseJSON([]byte(error_json))
			fmt.Println(arguments.StringIndent("", "  "))

		}

	}
}
