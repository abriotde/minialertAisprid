package monitorer

import (
        "time"
        "fmt"
)

type Alert struct {
    Timestamp time.Time
    Name      string
    Value     int32
}


type Monitorer struct {
	alerts []Alert
}

func (monito Monitorer) GetAlerts () []Alert {
	return monito.alerts
}

func (monito Monitorer) Log (varname string, varvalue int32) {
	alert := monito.isAlert(varname, varvalue)
	if alert!="" {
		fmt.Println("New alert : ",alert," for ",varname, " = ",varvalue,".")
		alert := Alert{Timestamp:time.Now(), Name:alert, Value:varvalue}
	       	monito.alerts = append(monito.alerts, alert)
	}
	for _,a := range monito.alerts {
		fmt.Println("Alert: ", a.Name, " = ", a.Value)
	}
}

func (monito Monitorer) isAlert (varname string, varvalue int32) string {
	if varname=="cpu" {
		if varvalue>80 {
			return "cpu_hight";
		}
	} else if varname=="battery" {
		if varvalue<20 {
			return "battery_low";
		} else if varvalue>98 {
			return "battery_hight";
		}
	}
	return "";
}

