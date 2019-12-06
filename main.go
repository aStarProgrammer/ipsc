// IPSC project main.go
package main

import (
	"ipsc/Page"
	"ipsc/Site"
	"ipsc/Utils"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func IndexPageSizeConvert(strPageSize string) string {

	if strPageSize == "" {
		return Page.INDEX_PAGE_SIZE_20
	}
	strPageSize = strings.ToUpper(strPageSize)
	switch strPageSize {
	case Site.PAGESIZE_NORMAL:
		return Page.INDEX_PAGE_SIZE_20
	case Site.PAGESIZE_SMALL:
		return Page.INDEX_PAGE_SIZE_10
	case Site.PAGESIZE_VERYSMALL:
		return Page.INDEX_PAGE_SIZE_5
	case Site.PAGESIZE_BIG:
		return Page.INDEX_PAGE_SIZE_30
	}
	return Page.INDEX_PAGE_SIZE_20
}

func Dispatch(cp CommandParser) (bool, error) {

	if cp.CurrentCommand == COMMAND_NEWSITE {
		//NewSiteProject no site project exist, cannot open and do operations
		if Utils.PathIsExist(cp.SiteFolderPath) == true {
			files, _ := ioutil.ReadDir(cp.SiteFolderPath)
			for _, f := range files {
				var ext = filepath.Ext(f.Name())
				if ext == ".sp" {
					var errMsg = "Cannot Create Site Project, there is a site project already exist at " + cp.SiteFolderPath
					return false, errors.New(errMsg)
				}
			}
		}
		var smp *Site.SiteModule
		smp = Site.NewSiteModule()

		bCreate, errCreate := smp.InitializeSiteProjectFolder(cp.SiteTitle, cp.SiteAuthor, cp.SiteDescription, cp.SiteFolderPath, cp.SiteOutputFolderPath)

		if errCreate != nil {
			return bCreate, errCreate
		}

	} else if cp.CurrentCommand == COMMAND_HELP {
		DipslayHelp(cp.HelpType)
	} else {
		//Open site project
		if Utils.PathIsExist(cp.SiteFolderPath) == false {
			var errMsg string
			errMsg = "Cannot find folder " + cp.SiteFolderPath
			fmt.Fprintln(os.Stderr, errMsg)
			return false, errors.New(errMsg)
		}

		var siteProjectFileName string
		if cp.SiteTitle == "" {
			var spCount int
			files, _ := ioutil.ReadDir(cp.SiteFolderPath)
			for _, f := range files {
				var ext = filepath.Ext(f.Name())
				if ext == ".sp" {
					siteProjectFileName = f.Name()
					spCount = spCount + 1
				}
			}

			if spCount > 1 {
				fmt.Fprintln(os.Stderr, "More than 1 sp file")
				return false, errors.New("More than 1 sp file")
			}
		} else {
			siteProjectFileName = cp.SiteTitle + ".sp"
		}

		if siteProjectFileName == "" {
			return false, errors.New("SiteTitle is empty and cannot find .sp file in root folder of " + cp.SiteFolderPath)
		}

		var smp *Site.SiteModule

		var siteProjectFilePath = filepath.Join(cp.SiteFolderPath, siteProjectFileName)

		if Utils.PathIsExist(siteProjectFilePath) == false {
			var errSPFPath error
			siteProjectFileName, errSPFPath = Utils.Try2FindSpFile(cp.SiteFolderPath)
			if errSPFPath != nil || siteProjectFileName == "" {
				return false, errors.New("Cannot find site proejct file path at " + siteProjectFilePath)
			}
		}

		smp = Site.NewSiteModule_WithArgs(cp.SiteFolderPath, siteProjectFileName)

		if smp == nil {
			return false, errors.New("Cannot initialize Site Module")
		}

		//Start dispatch message
		switch cp.CurrentCommand {
		case COMMAND_UPDATESITE:
			return smp.UpdateSiteProject(cp.SiteFolderPath, cp.SiteTitle, cp.SiteAuthor, cp.SiteDescription)

		case COMMAND_GETSITEPROPERTY:
			DisplaySiteProperties(smp)

		case COMMAND_LISTSOURCEPAGES:
			DisplaySourcePages(smp)

		case COMMAND_LISTOUTPUTPAGES:
			DisplayOutputPages(smp)

		case COMMAND_LISTPAGE:
			DisplayPage(smp, cp.PageID)

		case COMMAND_COMPILE:
			var sitePageSize = IndexPageSizeConvert(cp.IndexPageSize)
			bComplie, errComplie := smp.Compile(sitePageSize)
			if errComplie == nil {
				fmt.Println("COMPILE Summary:")
				DisplayComplieSummary("    ", smp.GetSiteProject().LastComplieSummary)
			}
			return bComplie, errComplie

		case COMMAND_CREATEMARKDOWN:
			return smp.CreateMarkdown(cp.SiteFolderPath, cp.SourcePagePath, cp.MarkdownType)

		case COMMAND_ADDPAGE:
			var bAdd bool
			var pageID string
			var errAdd error
			if cp.PageType == Page.MARKDOWN || cp.PageType == Page.HTML {
				bAdd, pageID, errAdd = smp.AddPage(cp.PageTitle, "", cp.PageAuthor, cp.SourcePagePath, cp.PageTitleImagePath, cp.PageType, cp.PageIsTop)
			} else if cp.PageType == Page.LINK {
				bAdd, pageID, errAdd = smp.AddPage(cp.PageTitle, "", cp.PageAuthor, cp.LinkUrl, cp.PageTitleImagePath, cp.PageType, cp.PageIsTop)
			}
			if errAdd == nil {
				fmt.Println("Add Success, ID generated for added page is " + pageID)
			}
			return bAdd, errAdd

		case COMMAND_UPDATEPAGE:
			var bUpdate bool
			var errUpdate error
			if cp.SourcePagePath != "" {
				bUpdate, errUpdate = smp.UpdatePage(cp.PageID, cp.PageTitle, "", cp.PageAuthor, cp.SourcePagePath, cp.PageTitleImagePath, cp.PageIsTop)
			} else if cp.LinkUrl != "" {
				bUpdate, errUpdate = smp.UpdatePage(cp.PageID, cp.PageTitle, "", cp.PageAuthor, cp.LinkUrl, cp.PageTitleImagePath, cp.PageIsTop)
			}
			if errUpdate == nil {
				fmt.Println("Update Success")
			}
			return bUpdate, errUpdate

		case COMMAND_DELETEPAGE:
			return smp.DeletePage(cp.PageID, cp.RestorePage)

		case COMMAND_LISTRECYCLEDPAGES:
			ListRecycledPages(smp)
			return true, nil

		case COMMAND_RESTORERECYCLEDPAGE:
			return smp.RestoreRecycledPageSourceFile(cp.PageID)

		case COMMAND_CLEARRECYCLEDPAGES:
			return smp.CleanRecycledPageSourceFiles()
		}
	}
	return true, nil
}

func DipslayHelp(helpType string) {
	helpType = strings.ToUpper(helpType)

	if helpType == FULLHELP {
		helpContent, errHelp := GetFullHelpInformation()
		if errHelp != nil {
			fmt.Println("Cannot get full help information")
		} else {
			fmt.Println(helpContent)
		}
	} else {
		helpContent, errHelp := GetQuickHelpInformation()
		if errHelp != nil {
			fmt.Println("Cannot get quick help information")
		} else {
			fmt.Println(helpContent)
		}
	}
}

/*
func SearchPage(smp *Site.SiteModule, propertyName, propertyValue string) {
	pageID, errSearch := smp.SearchPage(propertyName, propertyValue)
	if errSearch != nil {
		var errMsg string
		errMsg = "Cannot find page with " + propertyName + " : " + propertyValue
		fmt.Println(errMsg)
		return
	}
	var resultMsg string
	resultMsg = "Page with " + propertyName + " : " + propertyValue + " found, PageID is " + pageID
	fmt.Println(resultMsg)
}
*/

func DisplaySourcePages(smp *Site.SiteModule) {
	var allpages = smp.GetAllPages()

	var sActive = allpages[0]

	active, _ := strconv.Atoi(sActive)
	if active == 1 {
		fmt.Println("There is 1 source page ")
	} else {
		fmt.Println("There are " + strconv.Itoa(active) + " source pages ")
	}
	fmt.Println("=============")

	var index int
	var count int
	count = 1
	for index = 3; index < 3+active; index++ {
		fmt.Println("    Page " + strconv.Itoa(count) + " :")
		count++
		DisplayPageProperties(allpages[index])
		fmt.Println("    --------------")
	}
}

func DisplayPageProperties(strPageProperteis string) {
	if strPageProperteis == "" {
		return
	}

	var properties = strings.Split(strPageProperteis, "|")
	for _, property := range properties {
		fmt.Println("    " + property)
	}
}

func DisplayOutputPages(smp *Site.SiteModule) {
	var allpages = smp.GetAllPages()

	var sActive = allpages[0]
	var sRecycled = allpages[1]
	var sOutput = allpages[2]

	active, _ := strconv.Atoi(sActive)
	recycled, _ := strconv.Atoi(sRecycled)
	output, _ := strconv.Atoi(sOutput)

	source := active + recycled
	if output == 1 {
		fmt.Println("There is 1 output page ")
	} else {
		fmt.Println("There are " + strconv.Itoa(output) + " output pages")
	}

	fmt.Println("==============")
	var count int
	count = 1
	for index := 3 + source; index < len(allpages); index++ {
		fmt.Println("    Page " + strconv.Itoa(count) + " :")
		count++
		DisplayPageProperties(allpages[index])
		fmt.Println("    --------------")
	}
}

func DisplayPage(smp *Site.SiteModule, pageID string) {
	var allpages = smp.GetAllPages()

	for _, page := range allpages {
		if strings.Contains(page, pageID) {
			fmt.Println("Page Found:")
			fmt.Println("=============")
			DisplayPageProperties(page)
		}
	}

}

func DisplaySiteProperties(smp *Site.SiteModule) {
	var sp = smp.GetSiteProject()
	fmt.Println("Site Properties:")
	fmt.Println("-----------------")
	fmt.Println("    Title: " + sp.Title)
	fmt.Println("    Description: " + sp.Description)
	fmt.Println("    Author: " + sp.Author)
	fmt.Println("    Create Time: " + sp.CreateTime)
	fmt.Println("    Last Modified: " + sp.LastModified)
	fmt.Println("    Last Compiled: " + sp.LastComplied)
	fmt.Println("    Last Compile Summary: ")
	DisplayComplieSummary("        ", sp.LastComplieSummary)
	fmt.Println("    Output Folder: " + sp.OutputFolderPath)
	fmt.Println("-----------------")
}

func DisplayComplieSummary(prefix, summary string) {
	var items = strings.Split(summary, "_")

	for _, item := range items {
		fmt.Println(prefix + item)
	}
}

func ListRecycledPages(smp *Site.SiteModule) {
	var allpages = smp.GetAllPages()

	var sActive = allpages[0]
	var sRecycled = allpages[1]

	active, _ := strconv.Atoi(sActive)
	recycled, _ := strconv.Atoi(sRecycled)

	if recycled == 1 {
		fmt.Println("There is 1 recycled page ")
	} else {
		fmt.Println("There are " + strconv.Itoa(recycled) + " recycled pages")
	}
	fmt.Println("==============")
	var count int
	count = 1
	for index := 3 + active; index < 3+active+recycled; index++ {
		fmt.Println("    Page " + strconv.Itoa(count) + " :")
		count++
		DisplayPageProperties(allpages[index])
		fmt.Println("    --------------")
	}
}

func Run() {
	var cp CommandParser
	bParse := cp.ParseCommand()
	if bParse == true {
		_, errRet := Dispatch(cp)
		if errRet != nil {
			fmt.Println(errRet.Error())
		}
		fmt.Println("Done")
	}
}

func main() {
	var softwareinfo = "IPSC 0.1.0.1"
	fmt.Println(softwareinfo)
	Run()
}
