// 2026 Jan Provaznik (jan@provaznik.pro)
//

package main

import "os"
import "fmt"
import "time"
import "github.com/jan-provaznik/sus"
import "github.com/NVIDIA/go-nvml/pkg/nvml"

const upperDrawLimit = 105.00
const lowerDrawLimit =  10.00
const matchDrawLimit =   0.75

func main () {
	defer nvml.Shutdown()

	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		fmt.Println("nvmlInit failed")
		os.Exit(1)
	}

	list, err := sus.FindAstralDevices()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for index, device := range list {
		fmt.Printf("Detected device (%d) identified by (%s)\n",
			index, device.Identifier())
	}
	fmt.Println()

	for {
		for index, device := range list {
			err := deviceMonitor(index, device)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func deviceMonitor (index int, device sus.AstralDevice) error {
	// ... load, as reported via asus interface
	pins, err := sus.ReadAstralDevicePins(device)
	if err != nil {
		return err
	}

	// ... calculate statistics
	totalDraw := 0.0
	upperDraw := 0.0
	lowerDraw := 1e6

	for _, pin := range pins {
		value := pin.Drawing()
		if value > upperDraw {
			upperDraw = value
		}
		if value < lowerDraw {
			lowerDraw = value
		}
		totalDraw = totalDraw + value
	}

	// ... calculate draw match
	matchDraw := lowerDraw / upperDraw

	// ... emergency actions (pin overload)
	if upperDraw > upperDrawLimit {
		fmt.Printf("Device (%d) identified by (%s)\n",
			index, device.Identifier())
		fmt.Printf("... detected overload %.1f (limit %.1f)\n",
			upperDraw, upperDrawLimit)

		limit, err := sus.LimitAstralDeviceLoad(device)
		if err != nil {
			return err
		}

		fmt.Printf("... limiting power draw to %.1f W\n", 
			limit)
	}

	// ... emergency actions (pin mismatch min-max draw)
	if lowerDraw > lowerDrawLimit {
		if matchDraw < matchDrawLimit {
			fmt.Printf("Device (%d) identified by (%s)\n",
				index, device.Identifier())
			fmt.Printf("... detected mismatch %.2f (limit %.2f)\n",
				matchDraw, matchDrawLimit)

			limit, err := sus.LimitAstralDeviceFreq(device)
			if err != nil {
				return err
			}

			fmt.Printf("... attempting to limit device frequency to %d MHz\n", 
				limit)
		}
	}

	return nil
}

