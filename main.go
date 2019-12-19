// IPSC project main.go
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ipsc/Page"
	"ipsc/Site"
	"ipsc/Utils"
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
	//fmt.Println("A")
	if cp.CurrentCommand == COMMAND_NEWSITE {
		//NewSiteProject no site project exist, cannot open and do operations
		if Utils.PathIsExist(cp.SiteFolderPath) == true {
			files, errorReadDir := ioutil.ReadDir(cp.SiteFolderPath)
			if errorReadDir != nil {
				Utils.Logger.Println("Dispatch: " + errorReadDir.Error())
			}
			for _, f := range files {
				var ext = filepath.Ext(f.Name())
				if ext == ".sp" {
					var errMsg = "Main.Dispatch: Cannot Create Site Project, there is a site project already exist at " + cp.SiteFolderPath
					Utils.Logger.Println(errMsg)
					return false, errors.New(errMsg)
				}
			}
		}
		var smp *Site.SiteModule
		smp = Site.NewSiteModule()

		bCreate, errCreate := smp.InitializeSiteProjectFolder(cp.SiteTitle, cp.SiteAuthor, cp.SiteDescription, cp.SiteFolderPath, cp.SiteOutputFolderPath)

		if errCreate != nil {
			Utils.Logger.Println("Main.Dispatch: " + errCreate.Error())
			return bCreate, errCreate
		} else {
			fmt.Println("Create Site Success")
		}

	} else if cp.CurrentCommand == COMMAND_HELP {
		DipslayHelp(cp.HelpType)
	} else {
		//Open site project
		if Utils.PathIsExist(cp.SiteFolderPath) == false {
			var errMsg string
			errMsg = "Main.Dispatch: Cannot find folder " + cp.SiteFolderPath
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		var siteProjectFileName string
		if cp.SiteTitle == "" {
			var spCount int
			files, errReadDir := ioutil.ReadDir(cp.SiteFolderPath)
			if errReadDir != nil {
				Utils.Logger.Println("Dispatch: " + errReadDir.Error())
			}
			for _, f := range files {
				var ext = filepath.Ext(f.Name())
				if ext == ".sp" {
					siteProjectFileName = f.Name()
					spCount = spCount + 1
				}
			}

			if spCount > 1 {
				Utils.Logger.Println("Main.More than 1 sp file")
				return false, errors.New("Main.More than 1 sp file")
			}
		} else {
			siteProjectFileName = cp.SiteTitle + ".sp"
		}

		if siteProjectFileName == "" {
			var errMsg = "Main.Dispatch: SiteTitle is empty and cannot find .sp file in root folder of " + cp.SiteFolderPath
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		var smp *Site.SiteModule

		var siteProjectFilePath = filepath.Join(cp.SiteFolderPath, siteProjectFileName)

		if Utils.PathIsExist(siteProjectFilePath) == false {
			var errSPFPath error
			siteProjectFileName, errSPFPath = Utils.Try2FindSpFile(cp.SiteFolderPath)
			if errSPFPath != nil || siteProjectFileName == "" {
				var errMsg = "Main.Dispatch: Cannot find site proejct file path at " + siteProjectFilePath
				Utils.Logger.Println(errMsg)
				return false, errors.New(errMsg)
			}
		}

		smp = Site.NewSiteModule_WithArgs(cp.SiteFolderPath, siteProjectFileName)

		if smp == nil {
			var errMsg = "Main.Dispatch: Cannot initialize Site Module"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
		//fmt.Println("B")
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

		case COMMAND_EXPORTSOURCEPAGES:
			bExport, errExport := ExportSourcePages(smp, cp.ExportFolderPath)
			if errExport != nil {
				Utils.Logger.Println(errExport.Error())
				return bExport, errExport
			}
		case COMMAND_COMPILE:
			var sitePageSize = IndexPageSizeConvert(cp.IndexPageSize)
			bCompile, errCompile := smp.Compile(sitePageSize)
			if errCompile == nil && bCompile {
				fmt.Println("Compile Summary:")
				DisplayCompileSummary("    ", smp.GetSiteProject().LastCompileSummary)
			} else {
				Utils.Logger.Println("Main.Dispatch: Compile " + errCompile.Error())
			}
			return bCompile, errCompile

		case COMMAND_CREATEMARKDOWN:
			bCreate, errCreate := smp.CreateMarkdown(cp.SiteFolderPath, cp.SourcePagePath, cp.MarkdownType)
			if errCreate == nil && bCreate {
				Utils.Logger.Println("CreateMarkdown Success " + cp.SourcePagePath)
			} else {
				Utils.Logger.Println("Main.Dispatch: Create Markdown Page " + errCreate.Error())
			}
			return bCreate, errCreate

		case COMMAND_ADDPAGE:
			var bAdd bool
			var pageID string
			var errAdd error
			isTop, errParseBool := strconv.ParseBool(cp.PageIsTop)
			if errParseBool != nil {
				Utils.Logger.Println("Dispatch.COMMAND_ADDPAGE: " + errParseBool.Error())
			}
			if cp.PageType == Page.MARKDOWN || cp.PageType == Page.HTML {
				bAdd, pageID, errAdd = smp.AddPage(cp.PageTitle, "", cp.PageAuthor, cp.SourcePagePath, cp.PageTitleImagePath, cp.PageType, isTop)
			} else if cp.PageType == Page.LINK {
				bAdd, pageID, errAdd = smp.AddPage(cp.PageTitle, "", cp.PageAuthor, cp.LinkUrl, cp.PageTitleImagePath, cp.PageType, isTop)
			}
			if errAdd == nil && bAdd {
				fmt.Println("Add Success, ID generated for added page is " + pageID)
			} else {
				Utils.Logger.Println("Main.Dispatch: Add Page " + errAdd.Error())
			}
			return bAdd, errAdd

		case COMMAND_UPDATEPAGE:
			var bUpdate bool
			var errUpdate error
			isTop, errParseBool := strconv.ParseBool(cp.PageIsTop)
			if errParseBool != nil {
				Utils.Logger.Println("Dispatch.COMMAND_ADDPAGE: " + errParseBool.Error())
			}
			if cp.SourcePagePath != "" {
				bUpdate, errUpdate = smp.UpdatePage(cp.PageID, cp.PageTitle, "", cp.PageAuthor, cp.SourcePagePath, cp.PageTitleImagePath, isTop)
			} else if cp.LinkUrl != "" {
				bUpdate, errUpdate = smp.UpdatePage(cp.PageID, cp.PageTitle, "", cp.PageAuthor, cp.LinkUrl, cp.PageTitleImagePath, isTop)
			} else if cp.SourcePagePath == "" && cp.LinkUrl == "" {
				bUpdate, errUpdate = smp.UpdatePage(cp.PageID, cp.PageTitle, "", cp.PageAuthor, "", cp.PageTitleImagePath, isTop)
			}
			if errUpdate == nil && bUpdate {
				fmt.Println("Update Success")
			} else {
				Utils.Logger.Println("Main.Dispatch: Update Source Page " + errUpdate.Error())
			}
			return bUpdate, errUpdate

		case COMMAND_DELETEPAGE:
			restorePage, errParseBool := strconv.ParseBool(cp.RestorePage)
			if errParseBool != nil {
				Utils.Logger.Println("Dispatch.COMMAND_ADDPAGE: " + errParseBool.Error())
			}
			bDelete, errDelete := smp.DeletePage(cp.PageID, restorePage)
			if errDelete == nil && bDelete {
				fmt.Println("Delete Success")
			} else {
				Utils.Logger.Println("Main.Dispatch: Delete Source Page " + errDelete.Error())
			}
			return bDelete, errDelete

		case COMMAND_LISTRECYCLEDPAGES:
			ListRecycledPages(smp)
			return true, nil

		case COMMAND_RESTORERECYCLEDPAGE:
			return smp.RestoreRecycledPageSourceFile(cp.PageID)

		case COMMAND_CLEARRECYCLEDPAGES:
			return smp.CleanRecycledPageSourceFiles()
		default:
			Utils.Logger.Println("Command not found " + cp.CurrentCommand)
			return false, errors.New("Main.Command not found " + cp.CurrentCommand)
		}
	}
	return true, nil
}

func DipslayHelp(helpType string) {
	helpType = strings.ToUpper(helpType)

	if helpType == FULLHELP {
		helpContent, errHelp := GetFullHelpInformation()
		if errHelp != nil {
			Utils.Logger.Println("Main.DisplayHelp: Cannot get full help information")
		} else {
			fmt.Println(helpContent)
		}
	} else {
		helpContent, errHelp := GetQuickHelpInformation()
		if errHelp != nil {
			Utils.Logger.Println("Main.DisplayHelp: Cannot get quick help information")
		} else {
			fmt.Println(helpContent)
		}
	}
}

func ExportSourcePages(smp *Site.SiteModule, exportFolderPath string) (bool, error) {
	return smp.ExportSourcePages(exportFolderPath)
}

/*
func SearchPage(smp *Site.SiteModule, propertyName, propertyValue string) {
	pageID, errSearch := smp.SearchPage(propertyName, propertyValue)
	if errSearch != nil {
		var errMsg string
		errMsg = "Cannot find page with " + propertyName + " : " + propertyValue
		Utils.Logger.Println(errMsg)
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
	fmt.Println("    Last Compiled: " + sp.LastCompiled)
	fmt.Println("    Last Compile Summary: ")
	DisplayCompileSummary("        ", sp.LastCompileSummary)
	fmt.Println("    Output Folder: " + sp.OutputFolderPath)
	fmt.Println("-----------------")
}

func DisplayCompileSummary(prefix, summary string) {
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
	Utils.InitLogger()
	var cp CommandParser
	bParse := cp.ParseCommand()
	if bParse == true {
		_, errRet := Dispatch(cp)
		if errRet != nil {
			Utils.Logger.Println(errRet.Error())
		}
	}
	//fmt.Println("RunE")
	fmt.Println("Done")
}

func main() {
	Run()
}
