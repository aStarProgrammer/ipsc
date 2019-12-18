package main

import (
	"flag"
	"fmt"
	"ipsc/Page"
	"os/user"
	"strings"
)

type CommandParser struct {
	CurrentCommand       string
	SiteTitle            string
	SiteDescription      string
	SiteFolderPath       string
	ExportFolderPath     string
	SiteAuthor           string
	SiteOutputFolderPath string
	PropertyName         string
	PropertyValue        string
	IndexPageSize        string
	StopCompileWithError bool
	VerboseLog           bool
	HelpType             string
	PageID               string
	PageTitle            string
	PageAuthor           string
	PageIsTop            string
	PageType             string
	SourcePagePath       string
	LinkUrl              string
	PageTitleImagePath   string
	RestorePage          string
	MarkdownType         string
}

func (cpp *CommandParser) ParseCommand() bool {
	//Set All Arguments
	flag.StringVar(&cpp.CurrentCommand, "Command", "", GetFieldHelpMsg("Command"))
	flag.StringVar(&cpp.SiteTitle, "SiteTitle", "", GetFieldHelpMsg("SiteTitle"))
	flag.StringVar(&cpp.SiteDescription, "SiteDescription", "", GetFieldHelpMsg("SiteDescription"))
	flag.StringVar(&cpp.SiteFolderPath, "SiteFolder", "", GetFieldHelpMsg("SiteFolder"))
	flag.StringVar(&cpp.ExportFolderPath, "ExportFolder", "", GetFieldHelpMsg("ExportFolderPath"))
	flag.StringVar(&cpp.SiteAuthor, "SiteAuthor", "", GetFieldHelpMsg("SiteAuthor"))
	flag.StringVar(&cpp.SiteOutputFolderPath, "OutputFolder", "", GetFieldHelpMsg("OutputFolder"))
	flag.StringVar(&cpp.PropertyName, "PropertyName", "", GetFieldHelpMsg("PropertyName"))
	flag.StringVar(&cpp.PropertyValue, "PropertyValue", "", GetFieldHelpMsg("PropertyValue"))
	flag.StringVar(&cpp.IndexPageSize, "IndexPageSize", "Normal", GetFieldHelpMsg("IndexPageSize"))
	flag.StringVar(&cpp.HelpType, "HelpType", "QuickHelp", GetFieldHelpMsg("HelpType"))
	flag.StringVar(&cpp.PageID, "PageID", "", GetFieldHelpMsg("PageID"))
	flag.StringVar(&cpp.PageTitle, "PageTitle", "", GetFieldHelpMsg("PageTitle"))
	flag.StringVar(&cpp.PageAuthor, "PageAuthor", "", GetFieldHelpMsg("PageAuthor"))
	flag.StringVar(&cpp.PageIsTop, "IsTop", "false", GetFieldHelpMsg("IsTop"))
	flag.StringVar(&cpp.PageType, "PageType", "MARKDOWN", GetFieldHelpMsg("PageType"))
	flag.StringVar(&cpp.SourcePagePath, "PagePath", "", GetFieldHelpMsg("PagePath"))
	flag.StringVar(&cpp.LinkUrl, "LinkUrl", "", GetFieldHelpMsg("LinkUrl"))
	flag.StringVar(&cpp.PageTitleImagePath, "TitleImage", "", GetFieldHelpMsg("TitleImage"))
	flag.StringVar(&cpp.RestorePage, "RestorePage", "true", GetFieldHelpMsg("RestorePage"))
	flag.StringVar(&cpp.MarkdownType, "MarkdownType", "News", GetFieldHelpMsg("MarkdownType"))
	//Parse
	flag.Parse()

	//Trim all String properties
	cpp.CurrentCommand = strings.TrimSpace(cpp.CurrentCommand)
	cpp.HelpType = strings.TrimSpace(cpp.HelpType)
	cpp.LinkUrl = strings.TrimSpace(cpp.LinkUrl)
	cpp.MarkdownType = strings.TrimSpace(cpp.MarkdownType)
	cpp.PageAuthor = strings.TrimSpace(cpp.PageAuthor)
	cpp.PageID = strings.TrimSpace(cpp.PageID)
	cpp.PageTitle = strings.TrimSpace(cpp.PageTitle)
	cpp.PageTitleImagePath = strings.TrimSpace(cpp.PageTitleImagePath)
	cpp.PageType = strings.TrimSpace(cpp.PageType)
	cpp.PropertyName = strings.TrimSpace(cpp.PropertyName)
	cpp.PropertyValue = strings.TrimSpace(cpp.PropertyValue)
	cpp.SiteAuthor = strings.TrimSpace(cpp.SiteAuthor)
	cpp.SiteFolderPath = strings.TrimSpace(cpp.SiteFolderPath)
	cpp.ExportFolderPath = strings.TrimSpace(cpp.ExportFolderPath)
	cpp.SiteDescription = strings.TrimSpace(cpp.SiteDescription)
	cpp.SiteOutputFolderPath = strings.TrimSpace(cpp.SiteOutputFolderPath)
	cpp.SiteTitle = strings.TrimSpace(cpp.SiteTitle)
	cpp.SourcePagePath = strings.TrimSpace(cpp.SourcePagePath)
	cpp.IndexPageSize = strings.TrimSpace(cpp.IndexPageSize)

	//To Upper
	cpp.CurrentCommand = strings.ToUpper(cpp.CurrentCommand)
	cpp.MarkdownType = strings.ToUpper(cpp.MarkdownType)
	cpp.PageType = strings.ToUpper(cpp.PageType)
	cpp.HelpType = strings.ToUpper(cpp.HelpType)

	cpp.RestorePage = strings.ToUpper(cpp.RestorePage)
	cpp.PageIsTop = strings.ToUpper(cpp.PageIsTop)

	//Check whether command is help, if it is help,jump other operations
	if cpp.CurrentCommand == "" {
		cpp.CurrentCommand = "HELP"
	}

	if cpp.CurrentCommand == "HELP" {
		return true
	}

	//Check Above 3 parameters
	if cpp.CheckMarkdownType(cpp.MarkdownType) == false {
		fmt.Println("CommandParse: MarkdownType parameter not current, must 'Blank' or 'News' or Empty")
		return false
	}

	if cpp.CheckHelpType(cpp.HelpType) == false {
		fmt.Println("CommandParse: HelpType parameter not current, must 'QuickHelp' or 'FullHelp' or Empty")
		return false
	}

	if cpp.CheckPageType(cpp.PageType) == false {
		fmt.Println("CommandParse: PageType parameter not current, must 'Markdown' 'Html' or 'Link' or Empty")
		return false
	}

	//Get Command
	//Don't input Command

	if cpp.SiteFolderPath == "" {
		fmt.Println("CommandParse: Site Folder is empty")
		return false
	}

	if cpp.SiteTitle == "" {
		fmt.Println("Site title is empty, if not create new site, will try to load .sp from the root of site project folder.if theree are more than 1 .sp file, will open fail with empty site title")
	}

	var ret bool
	ret = true
	cpp.CurrentCommand = strings.ToUpper(cpp.CurrentCommand)
	//Check Properties of New Site Project
	switch cpp.CurrentCommand {
	case COMMAND_NEWSITE:
		if cpp.SiteTitle == "" {
			fmt.Println("CommandParse: SiteTitle is empty, cannot create site ")
			ret = false
		}

		if cpp.SiteDescription == "" {
			fmt.Println("CommandParse: Site description is empty")
			ret = false
		}

		if cpp.SiteAuthor == "" {
			fmt.Println("CommandParse: Site author is empty,will use current login user")
			currentUser, errUser := user.Current()
			if errUser != nil {
				fmt.Println("CommandParse: User is empty, and cannot get user information from system")
				ret = false
			}
			cpp.SiteAuthor = currentUser.Username
		}

		if cpp.SiteOutputFolderPath == "" {
			fmt.Println("Output folder is empty,will create Output folder under site project folder")
		}

	case COMMAND_UPDATESITE:
		if cpp.SiteTitle == "" && cpp.SiteAuthor == "" && cpp.SiteDescription == "" {
			fmt.Println("CommandParse: Title Author and Description of site are all empty, will not udpate site property")
			ret = false
		}

	case COMMAND_GETSITEPROPERTY:
	case COMMAND_LISTSOURCEPAGES:
	case COMMAND_LISTOUTPUTPAGES:
	case COMMAND_LISTPAGE:
		if cpp.PageID == "" {
			fmt.Println("CommandParse: Page ID is empty,don't know which page to restore")
			ret = false
		}
	case COMMAND_EXPORTSOURCEPAGES:
		if cpp.ExportFolderPath == "" {
			fmt.Println("CommandParse: Export Folder Path cannot be empty")
			ret = false
		}
	case COMMAND_COMPILE:
	case COMMAND_LISTRECYCLEDPAGES:
	case COMMAND_RESTORERECYCLEDPAGE:
		if cpp.PageID == "" {
			fmt.Println("CommandParse: Page ID is empty,don't know which page to restore")
			ret = false
		}
	case COMMAND_CLEARRECYCLEDPAGES:
	case COMMAND_ADDPAGE:
		if cpp.PageTitle == "" {
			fmt.Println("CommandParse: Page title is empty")
			ret = false
		}

		if (cpp.PageType == Page.MARKDOWN || cpp.PageType == Page.HTML) && cpp.SourcePagePath == "" {
			fmt.Println("CommandParse: Path of source page file is empty")
			ret = false
		}
		if cpp.PageType == Page.LINK && cpp.LinkUrl == "" {
			fmt.Println("CommandParse: Url of link is empty")
			ret = false
		}
		if cpp.PageAuthor == "" {
			currentUser, errUser := user.Current()
			if errUser != nil {
				fmt.Println("CommandParse: User is empty, and cannot get user information from system")
				ret = false
			}
			cpp.PageAuthor = currentUser.Username
		}
		if cpp.PageTitleImagePath == "" {
			fmt.Println("Title image of page source file is empty,will not display image for this page in index page")
		}

		if cpp.PageIsTop != "TRUE" && cpp.PageIsTop != "FALSE" {
			fmt.Println("IsTop should be true or false in string")
			ret = false
		}

	case COMMAND_CREATEMARKDOWN:
		if cpp.SourcePagePath == "" {
			fmt.Println("CommandParse: Path of source page file is empty")
			ret = false
		}

	case COMMAND_UPDATEPAGE:
		if cpp.PageID == "" {
			fmt.Println("CommandParse: Page ID is empty,don't know which page to restore")
			ret = false
		}
		if cpp.PageIsTop != "TRUE" && cpp.PageIsTop != "FALSE" {
			fmt.Println("IsTop should be true or false in string")
			ret = false
		}
	case COMMAND_DELETEPAGE:
		if cpp.PageID == "" {
			fmt.Println("CommandParse: Page ID is empty,don't know which page to restore")
			ret = false
		}

		if cpp.RestorePage != "TRUE" && cpp.RestorePage != "FALSE" {
			fmt.Println("RestorePage should be true or false in string")
			ret = false
		}
	}
	return ret
}

func (cpp *CommandParser) CheckPageType(pageType string) bool {
	if pageType == Page.MARKDOWN || pageType == Page.HTML || pageType == Page.LINK || pageType == "" {
		return true
	}
	return false
}

func (cpp *CommandParser) CheckHelpType(help string) bool {
	if help == QUICKHELP || help == FULLHELP || help == "" {
		return true
	}
	return false
}

func (cpp *CommandParser) CheckMarkdownType(markdown string) bool {
	if markdown == Page.MARKDOWN_BLANK || markdown == Page.MARKDOWN_NEWS || markdown == "" {
		return true
	}
	return false
}
