package main

import (

	"github.com/YeSZ1520/biotools/internal/qpcr/service"

	"encoding/json"
)


func main() {
	data,err:=service.LoadExperimentalData("data2.xlsx")
	if err!=nil{
		panic(err)
	}

	formatData := service.FormatExperimentalData(data)

	formattedJsonData,err:=json.MarshalIndent(formatData,"","  ")
	if err!=nil{
		panic(err)
	}
	println(string(formattedJsonData))

	results,err:=service.Calculate(formatData,"b-actin-1","LFD",2)
	if err!=nil{
		panic(err)
	}

	service.SaveData("output.xlsx",results)
}