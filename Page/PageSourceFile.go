package Page

import (
	"encoding/json"
	"errors"
	"fmt"
	"ipsc/Utils"
	"reflect"
	"strconv"
	"strings"
)

type PageSourceFile struct {
	ID             string
	Title          string
	Description    string
	TitleImage     string
	Type           string //Markdown Html or Link
	Author         string
	CreateTime     string
	LastModified   string
	LastCompiled   string
	SourceFilePath string
	Status         string
	IsTop          bool
	OutputFile     int
}

func NewPageSourceFileP() *PageSourceFile {
	var ps PageSourceFile
	var psp *PageSourceFile
	psp = &ps

	psp.ID = Utils.GUID()
	psp.CreateTime = Utils.CurrentTime()
	psp.LastModified = Utils.CurrentTime()
	psp.IsTop = false

	return psp
}

func NewPageSourceFile() PageSourceFile {
	var ps PageSourceFile

	ps.ID = Utils.GUID()
	ps.CreateTime = Utils.CurrentTime()
	ps.LastModified = Utils.CurrentTime()
	ps.IsTop = false

	return ps
}

func ResetPageSourceFile(ps PageSourceFile) {
	ps.ID = ""
}

func IsPageSourceFileEmpty(ps PageSourceFile) bool {
	if ps.ID == "" {
		return true
	}
	return false
}

func (ps *PageSourceFile) ToJson() (string, error) {
	var _jsonbyte []byte

	if ps == nil {
		var errMsg = "PageSourceFile->ToJson: Pointer ps is nil"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	if IsPageSourceFileEmpty(*ps) {
		var errMsg = "PageSourceFile->ToJson: Page Source File is empty"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	_jsonbyte, err := json.Marshal(*ps)

	return string(_jsonbyte), err
}

func (ps *PageSourceFile) ToString() string {
	var properties string
	properties = "ID: " + ps.ID
	properties += "|Title: " + ps.Title
	properties += "|Author: " + ps.Author
	properties += "|Type: " + ps.Type
	properties += "|CreateTime: " + ps.CreateTime
	properties += "|LastModified: " + ps.LastModified
	properties += "|LastCompiled: " + ps.LastCompiled
	properties += "|Status: " + ps.Status
	properties += "|IsTop: " + strconv.FormatBool(ps.IsTop)
	properties += "|SourceFilePath: " + ps.SourceFilePath
	return properties
}

func (ps *PageSourceFile) GetProperty(propertyName string) (string, error) {
	typeOfSiteProject := reflect.TypeOf(*ps)
	_, bFind := typeOfSiteProject.FieldByName(propertyName)

	if bFind == false {
		var errMsg = "PageSourceFile->GetProperty: Cannot find field " + propertyName
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}
	immutable := reflect.ValueOf(*ps)
	val := immutable.FieldByName(propertyName)
	return val.String(), nil
}

func (ps *PageSourceFile) SetProperty(propertyName, propertyValue string) (bool, error) {
	typeOfSiteProject := reflect.TypeOf(*ps)
	field, bFind := typeOfSiteProject.FieldByName(propertyName)

	if bFind == false {
		var errMsg = "PageSourceFile->GetProperty: Cannot find field " + propertyName
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	mutable := reflect.ValueOf(*ps).Elem()

	propertyType := field.Type.Name()
	propertyType = strings.ToUpper(propertyType)
	if propertyType == "STRING" {
		mutable.FieldByName(propertyName).SetString(propertyValue)
	} else if propertyType == "BOOL" {
		val, errVal := strconv.ParseBool(propertyValue)
		if errVal != nil {
			var errMsg = "PageSourceFile->GetProperty: Cannot parse property value " + propertyValue + " property " + propertyName + " to Bool"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
		mutable.FieldByName(propertyName).SetBool(val)
	} else if propertyType == "INT" {
		val, errVal := strconv.ParseInt(propertyValue, 10, 64)
		if errVal != nil {
			var errMsg = "PageSourceFile->GetProperty: Cannot parse property value " + propertyValue + " property " + propertyName + " to Int"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
		mutable.FieldByName(propertyName).SetInt(val)
	} else {
		var errMsg = "PropertyType set error"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	return true, nil
}

type PageSourceFileSlice []PageSourceFile

func (s PageSourceFileSlice) Len() int {
	return len(s)
}

func (s PageSourceFileSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s PageSourceFileSlice) Less(i, j int) bool {
	return s[i].CreateTime < s[j].CreateTime
}
