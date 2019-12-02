package Page

import (
	"IPSC/Utils"
	"encoding/json"
	"errors"
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
	LastComplied   string
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
		return "", errors.New("Pointer ps is nil")
	}

	if IsPageSourceFileEmpty(*ps) {
		return "", errors.New("Page Source File is empty")
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
	properties += "|LastComplied: " + ps.LastComplied
	properties += "|Status: " + ps.Status
	properties += "|IsTop: " + strconv.FormatBool(ps.IsTop)
	properties += "|SourceFilePath: " + ps.SourceFilePath
	return properties
}

func (ps *PageSourceFile) GetProperty(propertyName string) (string, error) {
	typeOfSiteProject := reflect.TypeOf(*ps)
	_, bFind := typeOfSiteProject.FieldByName(propertyName)

	if bFind == false {
		return "", errors.New("Cannot find field " + propertyName)
	}
	immutable := reflect.ValueOf(*ps)
	val := immutable.FieldByName(propertyName)
	return val.String(), nil
}

func (ps *PageSourceFile) SetProperty(propertyName, propertyValue string) (bool, error) {
	typeOfSiteProject := reflect.TypeOf(*ps)
	field, bFind := typeOfSiteProject.FieldByName(propertyName)

	if bFind == false {
		return false, errors.New("Cannot find field " + propertyName)
	}

	mutable := reflect.ValueOf(*ps).Elem()

	propertyType := field.Type.Name()
	propertyType = strings.ToUpper(propertyType)
	if propertyType == "STRING" {
		mutable.FieldByName(propertyName).SetString(propertyValue)
	} else if propertyType == "BOOL" {
		val, errVal := strconv.ParseBool(propertyValue)
		if errVal != nil {
			return false, errors.New("Cannot parse property value " + propertyValue + " property " + propertyName + " to Bool")
		}
		mutable.FieldByName(propertyName).SetBool(val)
	} else if propertyType == "INT" {
		val, errVal := strconv.ParseInt(propertyValue, 10, 64)
		if errVal != nil {
			return false, errors.New("Cannot parse property value " + propertyValue + " property " + propertyName + " to Int")
		}
		mutable.FieldByName(propertyName).SetInt(val)
	} else {
		return false, errors.New("PropertyType set error")
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
