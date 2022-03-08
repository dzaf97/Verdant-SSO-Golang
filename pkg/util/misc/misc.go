package misc

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	uuid "github.com/nu7hatch/gouuid"
)

/*
Engine Oil Temperature
Engine Coolant Temperature
Battery Voltage
Moving Flag
Fuel Level
Speed
RPM
Fuel Consumption
Mileage
*/

type (
	CanbusData struct {
		EngineOilTemp   float64
		EngineCoolTemp  float64
		BatteryVoltage  float64
		FuelLevel       float64
		Speed           float64
		RPM             int64
		FualConsumption float64
		Millage         int64
		GearPosition    string
		MovingFlag      bool
		Latitude        float64
		Longitude       float64
	}
)

func GenerateRandUUID() string {

	u, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	log.Println(u.String())

	return u.String()
}

func GenerateRandInt() string {
	// seed := rand.NewSource(time.Now().UnixNano())
	// randid := rand.New(123456789)

	return strconv.Itoa(rand.Int())
}
func CheckPer(rule []interface{}, user []interface{}) bool {

	for i := 0; i < len(rule); i++ {
		for b := 0; b < len(user); b++ {
			fmt.Println(user[b])
			if user[b] == rule[i] {
				return true
			}
		}
	}

	return false

}

func parseFloat(val int64, err error) float64 {
	return float64(val)
}

func DecodeCanbusData(payload string) *CanbusData {
	fulldata := &CanbusData{}

	//Engine EngineOilTemp
	fulldata.EngineOilTemp = parseFloat(strconv.ParseInt(payload[:2], 16, 64))
	//Engine EngineCoolTemp
	fulldata.EngineCoolTemp = parseFloat(strconv.ParseInt(payload[2:4], 16, 64))
	//Engine BatteryVoltage
	fulldata.BatteryVoltage = parseFloat(strconv.ParseInt(payload[4:8], 16, 64)) / 10
	//Engine FuelLevel
	fulldata.FuelLevel = parseFloat(strconv.ParseInt(payload[8:10], 16, 64))
	//Engine Speed
	fulldata.Speed = parseFloat(strconv.ParseInt(payload[10:14], 16, 64))
	//Engine FualConsumption
	fulldata.FualConsumption = parseFloat(strconv.ParseInt(payload[14:18], 16, 64)) / 10
	//Gear position
	pos, _ := hex.DecodeString(payload[18:20])
	fulldata.GearPosition = string(pos)
	//Movie flag
	if moveFlag, _ := strconv.ParseInt(payload[20:22], 16, 64); moveFlag == 1 {
		fulldata.MovingFlag = true
	} else {
		fulldata.MovingFlag = false
	}
	//Engine RPM
	fulldata.RPM, _ = strconv.ParseInt(payload[22:26], 16, 32)
	//Engine Millage
	fulldata.Millage, _ = strconv.ParseInt(payload[26:32], 16, 64)
	// Latitude
	fulldata.Latitude = parseFloat(strconv.ParseInt(payload[33:39], 16, 64)) / 100000
	// Longitude
	fulldata.Longitude = parseFloat(strconv.ParseInt(payload[40:], 16, 64)) / 100000

	// log.Println(fulldata)
	return fulldata

}

func DecodeAlert() {
	/*
		Battery Low
		Battery Current Leakage
		Engine Fail to Start
		DTC appears
		Car on harsh impact/ accident/ towed
	*/

}
