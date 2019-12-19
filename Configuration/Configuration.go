package Configuration

import (
	"errors"
	"io/ioutil"
	"ipsc/Page"
	"ipsc/Utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aWildProgrammer/fconf"
)

func GetCssFilePath() (string, error) {
	resourceFolderPath, errPath := GetResourcesFolderPath()

	if errPath != nil {
		var errMsg = "Cannot get resources folder path"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var cssFilePath = filepath.Join(resourceFolderPath, "news.css")
	return cssFilePath, nil
}

func GetResourcesFolderPath() (string, error) {
	executionPath, errPath := GetCurrentPath()

	if errPath != nil {
		var errMsg = "Cannot get path of current executable"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var resourceFolderPath = filepath.Join(executionPath, "Resources")
	return resourceFolderPath, nil
}

func GetIndexTemplateFilePath(indexPageSize string) (string, error) {
	resourceFolderPath, errPath := GetResourcesFolderPath()

	if errPath != nil {
		var errMsg = "Cannot get resources folder path"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var indexPageTemplateFilePath string
	if indexPageSize == Page.INDEX_PAGE_SIZE_5 {
		indexPageTemplateFilePath = filepath.Join(resourceFolderPath, "IndexPage5.md")
	} else if indexPageSize == Page.INDEX_PAGE_SIZE_10 {
		indexPageTemplateFilePath = filepath.Join(resourceFolderPath, "IndexPage10.md")
	} else if indexPageSize == Page.INDEX_PAGE_SIZE_20 {
		indexPageTemplateFilePath = filepath.Join(resourceFolderPath, "IndexPage20.md")
	} else if indexPageSize == Page.INDEX_PAGE_SIZE_30 {
		indexPageTemplateFilePath = filepath.Join(resourceFolderPath, "IndexPage30.md")
	}
	return indexPageTemplateFilePath, nil
}

func GetMoreTemplateFilePath(morePageSize string) (string, error) {
	resourceFolderPath, errPath := GetResourcesFolderPath()

	if errPath != nil {
		var errMsg = "Cannot get resources folder path"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var morePageTemplateFilePath string
	if morePageSize == Page.INDEX_PAGE_SIZE_5 {
		morePageTemplateFilePath = filepath.Join(resourceFolderPath, "MorePage5.md")
	} else if morePageSize == Page.INDEX_PAGE_SIZE_10 {
		morePageTemplateFilePath = filepath.Join(resourceFolderPath, "MorePage10.md")
	} else if morePageSize == Page.INDEX_PAGE_SIZE_20 {
		morePageTemplateFilePath = filepath.Join(resourceFolderPath, "MorePage20.md")
	} else if morePageSize == Page.INDEX_PAGE_SIZE_30 {
		morePageTemplateFilePath = filepath.Join(resourceFolderPath, "MorePage30.md")
	}
	return morePageTemplateFilePath, nil
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		Utils.Logger.Println(err.Error())
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		Utils.Logger.Println(err.Error())
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		Utils.Logger.Println(`error: Can't find "/" or "\".`)
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func getIniObject() (*fconf.Config, error) {
	configFilePath, errConfig := GetConfigurationFilePath()

	if errConfig != nil {
		var errMsg = "Cannot get file path of configuration path"
		Utils.Logger.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if Utils.PathIsExist(configFilePath) == false {
		var errMsg string
		errMsg = "Configuration file " + configFilePath + " not exist"
		Utils.Logger.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	return fconf.NewFileConf(configFilePath)
}

func GetConfigurationFilePath() (string, error) {
	executionPath, errPath := GetCurrentPath()

	if errPath != nil {
		var errMsg = "Cannot get path of current executable"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var configFilePath = filepath.Join(executionPath, "config.ini")
	return configFilePath, nil
}

func GetEmptyIndexItemTemplate() (string, error) {
	resourcesFolderPath, errResoruce := GetResourcesFolderPath()
	if errResoruce != nil {
		var errMsg = "Cannot get path of resource folder path"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var emptyIndexTemplateFilePath = filepath.Join(resourcesFolderPath, "EmptyItemTemplate.txt")

	if Utils.PathIsExist(emptyIndexTemplateFilePath) == false {
		var errMsg = "Cannot find empty Index Item Template setting file, its name is eit.txt, it should be along with ipsc"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	eit, errEit := ioutil.ReadFile(emptyIndexTemplateFilePath)

	if errEit != nil {
		var errMsg = "Read file content from empty Index Item Template file failed, please check its content, its name is EmptyItemTemplate.txt, it should be in the resources folder"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errEit.Error())
		return "", errors.New(errMsg)
	}

	return string(eit), nil
}

func GetEmptyImageItemTemplate() (string, error) {
	resourcesFolderPath, errResoruce := GetResourcesFolderPath()
	if errResoruce != nil {
		var errMsg = "Cannot get path of resource folder path"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var emptyImageTemplateFilePath = filepath.Join(resourcesFolderPath, "EmptyImageTemplate.txt")

	if Utils.PathIsExist(emptyImageTemplateFilePath) == false {
		var errMsg = "Cannot find empty Image Template setting file, its name is EmptyImageTemplate.txt, it should be in the resources folder"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	eit, errEit := ioutil.ReadFile(emptyImageTemplateFilePath)

	if errEit != nil {
		var errMsg = "Read file content from empty Index Item Template file failed, please check its content, its name is eit.txt, it should be along with ipsc"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errEit.Error())
		return "", errors.New(errMsg)
	}

	return string(eit), nil
}

func GetFullHelpPath() (string, error) {
	executionPath, errPath := GetCurrentPath()

	if errPath != nil {
		var errMsg = "Cannot get path of current executable"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var helpFilePath = filepath.Join(executionPath, "FullHelp.txt")
	if Utils.PathIsExist(helpFilePath) == false {
		var errMsg = "Cannot find FullHelp.txt at path " + helpFilePath
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}
	return helpFilePath, nil
}

func GetQuickHelpPath() (string, error) {
	executionPath, errPath := GetCurrentPath()

	if errPath != nil {
		var errMsg = "Cannot get path of current executable"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var helpFilePath = filepath.Join(executionPath, "QuickHelp.txt")
	if Utils.PathIsExist(helpFilePath) == false {
		var errMsg = "Cannot find QuickHelp.txt at path " + helpFilePath
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}
	return helpFilePath, nil
}

func GetTemplatesFolderPath() (string, error) {
	resourceFolderPath, errResource := GetResourcesFolderPath()
	if errResource != nil {
		Utils.Logger.Println(errResource.Error())
		return "", errResource
	}
	return filepath.Join(resourceFolderPath, "Templates"), nil
}
