// siteProject

package Site

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"ipsc/Page"
	"ipsc/Utils"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type SiteProject struct {
	ID                 string
	Title              string
	Description        string
	Author             string
	CreateTime         string
	LastModified       string
	LastCompiled       string
	LastCompileSummary string
	OutputFolderPath   string

	SourceFiles         []Page.PageSourceFile
	OutputFiles         []Page.PageOutputFile
	IndexPageSourceFile Page.PageSourceFile
	MorePageSourceFiles []Page.PageSourceFile
}

func NewSiteProject() *SiteProject {
	var sp SiteProject
	var spp *SiteProject
	spp = &sp

	spp.ID = Utils.GUID()
	spp.CreateTime = Utils.CurrentTime()

	return spp
}

func NewSiteProject_WithArgs(title, description, author, outputFolderPath string) *SiteProject {
	var sp SiteProject
	var spp *SiteProject
	spp = &sp

	spp.ID = Utils.GUID()
	spp.Title = title
	spp.Description = description
	spp.Author = author
	spp.CreateTime = Utils.CurrentTime()
	spp.OutputFolderPath = outputFolderPath

	return spp
}

func ResetSiteProject(spp *SiteProject) {
	spp.ID = ""
}

func IsSiteProjectEmpty(sp SiteProject) bool {
	if sp.ID == "" {
		return true
	}
	return false
}

func (spp *SiteProject) FromJson(_jsonString string) (bool, error) {
	if "" == _jsonString {
		var errMsg = "Argument jsonString is null"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	errUnmarshal := json.Unmarshal([]byte(_jsonString), spp)
	if errUnmarshal != nil {
		return false, errUnmarshal
	}
	return true, nil
}

func (spp *SiteProject) ToJson() (string, error) {
	var _jsonbyte []byte

	if spp == nil {
		var errMsg = "Pointer spp is nil"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	_jsonbyte, err := json.Marshal(*spp)

	return string(_jsonbyte), err
}

func (spp *SiteProject) LoadFromFile(filePath string) (bool, error) {
	if "" == filePath {
		var errMsg = "FilePath is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	bFileExist := Utils.PathIsExist(filePath)

	if false == bFileExist {
		fmt.Println("File not exist")
		return false, errors.New("File not exist")
	}

	_json, errRead := ioutil.ReadFile(filePath)

	if errRead != nil {
		fmt.Println(filePath)
		fmt.Println("Read File Fail")
		return false, errors.New("Read File Fail")
	}

	_jsonString := string(_json)

	if "" == _jsonString {
		var errMsg = "File is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	bUnMarshal, errUnMarshal := spp.FromJson(_jsonString)

	return bUnMarshal, errUnMarshal
}

func (spp *SiteProject) SaveToFile(filePath string) (bool, error) {
	if "" == filePath {
		var errMsg = "FilePath is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	json, errMarshal := spp.ToJson()

	if errMarshal != nil {
		return false, errMarshal
	}

	var errFilePath error
	if !Utils.PathIsExist(filePath) {
		filePath, errFilePath = Utils.MakePath(filePath)
		if errFilePath != nil {
			var errMsg = "Path nor exist and create parent folder failed"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}
	//路径分为绝对路径和相对路径
	//create，文件存在则会覆盖原始内容（其实就相当于清空），不存在则创建
	fp, error := os.Create(filePath)
	if error != nil {
		return false, error
	}
	//延迟调用，关闭文件
	defer fp.Close()

	_, errWriteFile := fp.WriteString(json)

	if errWriteFile != nil {
		var errMsg = "Write json to file failed"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	return true, nil
}

func (spp *SiteProject) AddPageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	spp.SourceFiles = append(spp.SourceFiles, psf)
	return true, nil
}

func (spp *SiteProject) RemovePageSourceFile(psf Page.PageSourceFile, restore bool) (bool, error) {

	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, sf := range spp.SourceFiles {
		if sf.ID == psf.ID {
			if restore {
				spp.SourceFiles[i].Status = Page.RECYCLED
				return true, nil
			}

			spp.SourceFiles = append(spp.SourceFiles[:i], spp.SourceFiles[i+1:]...)
			return true, nil
		}
	}
	var errMsg = "Source Page not found"
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) ResotrePageSourceFile(ID string) (bool, error) {
	var index = spp.GetPageSourceFile(ID)
	if index == -1 {
		var errMsg = "Not find page source file with ID " + ID
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	spp.SourceFiles[index].Status = Page.ACTIVE
	return true, nil
}

func (spp *SiteProject) UpdateMoreSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "SiteProject.UpdateMoreSourceFile: Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "SiteProject.UpdateMoreSourceFile: Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, sf := range spp.MorePageSourceFiles {
		if sf.ID == psf.ID {
			spp.SourceFiles[i].Author = psf.Author
			spp.SourceFiles[i].CreateTime = psf.CreateTime
			spp.SourceFiles[i].Description = psf.Description
			spp.SourceFiles[i].LastCompiled = psf.LastCompiled
			spp.SourceFiles[i].LastModified = Utils.CurrentTime()
			spp.SourceFiles[i].OutputFile = psf.OutputFile
			spp.SourceFiles[i].SourceFilePath = psf.SourceFilePath
			spp.SourceFiles[i].Title = psf.Title
			spp.SourceFiles[i].Type = psf.Type
			spp.SourceFiles[i].IsTop = psf.IsTop
			return true, nil
		}
	}
	var errMsg = "SiteProject.UpdateMoreSourceFile: Source Page not found"
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) UpdateIndexSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "SiteProject.UpdateIndexSourceFile: Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "SiteProject.UpdateIndexSourceFile: Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	spp.IndexPageSourceFile.Author = psf.Author
	spp.IndexPageSourceFile.CreateTime = psf.CreateTime
	spp.IndexPageSourceFile.Description = psf.Description
	spp.IndexPageSourceFile.LastCompiled = psf.LastCompiled
	spp.IndexPageSourceFile.LastModified = Utils.CurrentTime()
	spp.IndexPageSourceFile.OutputFile = psf.OutputFile
	spp.IndexPageSourceFile.SourceFilePath = psf.SourceFilePath
	spp.IndexPageSourceFile.Title = psf.Title
	spp.IndexPageSourceFile.IsTop = psf.IsTop
	return true, nil

}

func (spp *SiteProject) UpdatePageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
<<<<<<< HEAD
		var errMsg = "SiteProject.UpdatePageSourceFile: Site Project is empty"
=======
<<<<<<< HEAD
		var errMsg = "SiteProject.UpdatePageSourceFile: Site Project is empty"
=======
		var errMsg = "Site Project is empty"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
<<<<<<< HEAD
		var errMsg = "SiteProject.UpdatePageSourceFile: Source Page is empty"
=======
<<<<<<< HEAD
		var errMsg = "SiteProject.UpdatePageSourceFile: Source Page is empty"
=======
		var errMsg = "Source Page is empty"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, sf := range spp.SourceFiles {
		if sf.ID == psf.ID {
			spp.SourceFiles[i].Author = psf.Author
			spp.SourceFiles[i].CreateTime = psf.CreateTime
			spp.SourceFiles[i].Description = psf.Description
			spp.SourceFiles[i].LastCompiled = psf.LastCompiled
			spp.SourceFiles[i].LastModified = Utils.CurrentTime()
			spp.SourceFiles[i].OutputFile = psf.OutputFile
			spp.SourceFiles[i].SourceFilePath = psf.SourceFilePath
			spp.SourceFiles[i].Title = psf.Title
			spp.SourceFiles[i].Type = psf.Type
			spp.SourceFiles[i].IsTop = psf.IsTop
			return true, nil
		}
	}
<<<<<<< HEAD
	var errMsg = "SiteProject.UpdatePageSourceFile: Source Page not found"
=======
<<<<<<< HEAD
	var errMsg = "SiteProject.UpdatePageSourceFile: Source Page not found"
=======
	var errMsg = "Source Page not found"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) GetPageSourceFile(ID string) int {

	if IsSiteProjectEmpty(*spp) {
		return -1
	}

	if ID == "" {
		return -1
	}

	for i, sf := range spp.SourceFiles {
		if sf.ID == ID {
			return i
		}
	}
	return -1
}

func (spp *SiteProject) AddPageOutputFile(pof Page.PageOutputFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageOutputFileEmpty(pof) {
		var errMsg = "Output Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	spp.OutputFiles = append(spp.OutputFiles, pof)

	return true, nil
}

func (spp *SiteProject) RemovePageOutputFile(pof Page.PageOutputFile) (bool, error) {

	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageOutputFileEmpty(pof) {
		var errMsg = "Output Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, sf := range spp.OutputFiles {
		if sf.ID == pof.ID {
			spp.OutputFiles = append(spp.OutputFiles[:i], spp.OutputFiles[i+1:]...)
			return true, nil
		}
	}
	var errMsg = "Output Page not found"
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) UpdatePageOutputFile(pof Page.PageOutputFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageOutputFileEmpty(pof) {
		var errMsg = "Output Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, of := range spp.OutputFiles {
		if of.ID == pof.ID {
			spp.OutputFiles[i].Author = pof.Author
			spp.OutputFiles[i].CreateTime = pof.CreateTime
			spp.OutputFiles[i].Description = pof.Description
			spp.OutputFiles[i].FilePath = pof.FilePath
			spp.OutputFiles[i].IsTop = pof.IsTop
			spp.OutputFiles[i].Title = pof.Title
			spp.OutputFiles[i].Type = pof.Type

			return true, nil
		}
	}

	var errMsg = "Output Page not found"
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) GetPageOutputFile(ID string) int {
	if IsSiteProjectEmpty(*spp) {
		return -1
	}

	if ID == "" {
		return -1
	}

	for i, of := range spp.OutputFiles {
		if of.ID == ID {
			return i
		}
	}
	return -1
}

func (spp *SiteProject) AddMorePageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	spp.MorePageSourceFiles = append(spp.MorePageSourceFiles, psf)
	return true, nil
}

func (spp *SiteProject) RemoveMorePageSourceFile(psf Page.PageSourceFile) (bool, error) {

	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, sf := range spp.MorePageSourceFiles {
		if sf.ID == psf.ID {
			spp.MorePageSourceFiles = append(spp.MorePageSourceFiles[:i], spp.MorePageSourceFiles[i+1:]...)
			return true, nil
		}
	}
	var errMsg = "Source Page not found"
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) UpdateMorePageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	for i, sf := range spp.MorePageSourceFiles {
		if sf.ID == psf.ID {
			spp.MorePageSourceFiles[i].Author = psf.Author
			spp.MorePageSourceFiles[i].CreateTime = psf.CreateTime
			spp.MorePageSourceFiles[i].Description = psf.Description
			spp.MorePageSourceFiles[i].LastCompiled = psf.LastCompiled
			spp.MorePageSourceFiles[i].LastModified = psf.LastModified
			spp.MorePageSourceFiles[i].OutputFile = psf.OutputFile
			spp.MorePageSourceFiles[i].SourceFilePath = psf.SourceFilePath
			spp.MorePageSourceFiles[i].Title = psf.Title
			spp.MorePageSourceFiles[i].Type = psf.Type
			return true, nil
		}
	}
	var errMsg = "Source Page not found"
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
}

func (spp *SiteProject) GetMorePageSourceFile(ID string) int {

	if IsSiteProjectEmpty(*spp) {
		return -1
	}

	if ID == "" {
		return -1
	}

	for i, sf := range spp.MorePageSourceFiles {
		if sf.ID == ID {
			return i
		}
	}
	return -1
}

func (spp *SiteProject) SetIndexPageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Page.IsPageSourceFileEmpty(psf) {
		var errMsg = "Source Page is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	spp.IndexPageSourceFile = psf
	return true, nil
}

func (spp *SiteProject) CleanIndexPageSourceFile() (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	Page.ResetPageSourceFile(spp.IndexPageSourceFile)
	return true, nil
}

func (spp *SiteProject) PageStatistics() (string, error) {

	if IsSiteProjectEmpty(*spp) {
		var errMsg = "Site Project is empty"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var msg string
	msg = "Pages in " + spp.Title + ": \r\n"
	msg += "\tSource Pages:\r\n"

	var srcMdS, srcHtmlS, srcLinkS int
	srcMdS = 0
	srcHtmlS = 0
	srcLinkS = 0

	for _, source := range spp.SourceFiles {
		if source.Type == Page.MARKDOWN {
			srcMdS = srcMdS + 1
		} else if source.Type == Page.HTML {
			srcHtmlS = srcHtmlS + 1
		} else if source.Type == Page.LINK {
			srcLinkS = srcLinkS + 1
		}
	}

	msg += "\t\tMarkdown: " + strconv.Itoa(srcMdS) + " Html: " + strconv.Itoa(srcHtmlS) + " Link: " + strconv.Itoa(srcLinkS) + "\r\n"

	var outMdS, outHtmlS, outLinkS int
	outMdS = 0
	outHtmlS = 0
	outLinkS = 0

	for _, output := range spp.OutputFiles {
		if output.Type == Page.MARKDOWN {
			outMdS += 1
		} else if output.Type == Page.HTML {
			outHtmlS += 1
		} else if output.Type == Page.LINK {
			outLinkS += 1
		}
	}

	msg += "\tOutput:\r\n"
	msg += "\t\tMarkdown: " + strconv.Itoa(outMdS) + " Html: " + strconv.Itoa(outHtmlS) + " Link: " + strconv.Itoa(outLinkS) + "\r\n"

	msg += "\tMore: " + strconv.Itoa(len(spp.MorePageSourceFiles)) + "\r\n"

	if !Page.IsPageSourceFileEmpty(spp.IndexPageSourceFile) {
		msg += "\tIndex: 1\r\n"
	} else {
		msg += "\tIndex: 0\r\n"
	}

	return msg, nil
}

func (spp *SiteProject) GetActivePageSources() []string {
	var pages []string

	for _, psf := range spp.SourceFiles {
		if psf.Status == Page.ACTIVE {
			psfStr := psf.ToString()
			pages = append(pages, psfStr)
		}
	}

	return pages
}

func (spp *SiteProject) GetRecycledPageSources() []string {
	var pages []string

	for _, psf := range spp.SourceFiles {
		if psf.Status == Page.RECYCLED {
			psfStr := psf.ToString()
			pages = append(pages, psfStr)
		}
	}

	return pages
}

func (spp *SiteProject) GetAllPageOutputs() []string {
	var pages []string

	for _, pof := range spp.OutputFiles {
		pofStr := pof.ToString()
		pages = append(pages, pofStr)
	}

	return pages
}

func (spp *SiteProject) BackupSiteProjectFile(siteProjectFilePath string) (bool, error) {
	var siteProjectFileBackupPath string
	siteProjectFileBackupPath = siteProjectFilePath + ".backup"

	_, errCopy := Utils.CopyFile(siteProjectFilePath, siteProjectFileBackupPath)

	if errCopy != nil {
		return false, errCopy
	}
	return true, nil
}

func (spp *SiteProject) RestoreSiteProjectFile(siteProjectFilePath string) (bool, error) {
	var siteProjectFileBackupPath string
	siteProjectFileBackupPath = siteProjectFilePath + ".backup"

	if Utils.PathIsExist(siteProjectFileBackupPath) == false {
		var errMsg = "Backup File is not exist"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	_, errCopy := Utils.CopyFile(siteProjectFileBackupPath, siteProjectFilePath)

	if errCopy != nil {
		return false, errCopy
	}
	return true, nil
}

func (spp *SiteProject) GetSortedTopOutputFiles() (Page.PageOutputFileSlice, error) {
	var outputFileSlice Page.PageOutputFileSlice

	for _, outputFile := range spp.OutputFiles {
		if outputFile.IsTop == true && outputFile.Type != Page.INDEX {
			outputFileSlice = append(outputFileSlice, outputFile)
		}
	}

	sort.Sort(sort.Reverse(outputFileSlice))
	return outputFileSlice, nil
}

func (spp *SiteProject) GetSortedNormalOutputFiles() (Page.PageOutputFileSlice, error) {
	var outputFileSlice Page.PageOutputFileSlice

	for _, outputFile := range spp.OutputFiles {
		if outputFile.IsTop == false && outputFile.Type != Page.INDEX {
			outputFileSlice = append(outputFileSlice, outputFile)
		}
	}

	sort.Sort(sort.Reverse(outputFileSlice))
	return outputFileSlice, nil

}

func (spp *SiteProject) ExportSourcePages(exportFolderPath string) (bool, error) {
	if exportFolderPath == "" {
		var errMsg = "Export Source File: Export Folder Path is empty"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Utils.PathIsExist(exportFolderPath) == false {
		var errMsg = "Export Source File: Export Folder Path doesn't exist"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	//Make target folder
	var targetMarkdownFolder = filepath.Join(exportFolderPath, "Markdown")
	var targetHtmlFolder = filepath.Join(exportFolderPath, "Html")
	var targetLinkFolder = filepath.Join(exportFolderPath, "Link")
	var targetLinkFile = filepath.Join(targetLinkFolder, "Link.txt")

	if Utils.PathIsExist(targetMarkdownFolder) == false {
		_, errorMakeFolder := Utils.MakeFolder(targetMarkdownFolder)
		if errorMakeFolder != nil {
			var errMsg = "Export Source File : Cannot Make Markdown Folder"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	if Utils.PathIsExist(targetHtmlFolder) == false {
		_, errorMakeFolder := Utils.MakeFolder(targetHtmlFolder)
		if errorMakeFolder != nil {
			var errMsg = "Export Source File : Cannot Make Html Folder"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	if Utils.PathIsExist(targetLinkFolder) == false {
		_, errorMakeFolder := Utils.MakeFolder(targetLinkFolder)
		if errorMakeFolder != nil {
			var errMsg = "Export Source File : Cannot Make Link Folder"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	var links []string

	//Copy Source File (Markdown,Html) to export folder
	for _, psf := range spp.SourceFiles {
		if psf.Type == Page.MARKDOWN {
			if psf.SourceFilePath != "" && Utils.PathIsExist(psf.SourceFilePath) {
				var fileName = filepath.Base(psf.SourceFilePath)
				if strings.Contains(fileName, ".") == false {
					var errMsg = "Export Source File: File Name error " + fileName
					fmt.Println(errMsg)
					return false, errors.New(errMsg)
				}
				var targetFilePath = filepath.Join(targetMarkdownFolder, fileName)

				_, errCopy := Utils.CopyFile(psf.SourceFilePath, targetFilePath)
				if errCopy != nil {
					return false, errCopy
				}
			}
		} else if psf.Type == Page.HTML {
			if psf.SourceFilePath != "" && Utils.PathIsExist(psf.SourceFilePath) {
				var fileName = filepath.Base(psf.SourceFilePath)
				if strings.Contains(fileName, ".") == false {
					var errMsg = "Export Source File: File Name error " + fileName
					fmt.Println(errMsg)
					return false, errors.New(errMsg)
				}
				var targetFilePath = filepath.Join(targetHtmlFolder, fileName)

				_, errCopy := Utils.CopyFile(psf.SourceFilePath, targetFilePath)
				if errCopy != nil {
					return false, errCopy
				}
			}

		} else if psf.Type == Page.LINK {
<<<<<<< HEAD

			links = append(links, psf.SourceFilePath+"|"+psf.ID+"|"+psf.Title+"|"+strconv.FormatBool(psf.IsTop))
=======
			links = append(links, psf.SourceFilePath+"|"+psf.ID+"|"+psf.Title)
>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6
		}
	}

	//Write Links to Link.txt

	json, errMarshal := json.Marshal(links)

	if errMarshal != nil {
		return false, errMarshal
	}

	//create，文件存在则会覆盖原始内容（其实就相当于清空），不存在则创建
	fp, error := os.Create(targetLinkFile)
	if error != nil {
		return false, error
	}
	//延迟调用，关闭文件
	defer fp.Close()

	_, errWriteFile := fp.WriteString(string(json))

	if errWriteFile != nil {
		var errMsg = "Export Source File: Write Links to Link.txt failed"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	//Export Title Image of source files to source folder
	for _, psf := range spp.SourceFiles {
		if psf.TitleImage != "" {
<<<<<<< HEAD
=======
			var fileName = psf.ID

>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6
			fileType, errFileType := Utils.GetImageType(psf.TitleImage)
			if errFileType != nil {
				var errMsg = "Export Source File: Title Image Content is not correct"
				fmt.Println(errMsg)
				return false, errors.New(errMsg)
			}
<<<<<<< HEAD
			var fileName string

			if psf.Type == Page.MARKDOWN || psf.Type == Page.HTML {
				fileName = filepath.Base(psf.SourceFilePath)
				var fileExt = filepath.Ext(psf.SourceFilePath)

				fileName = strings.Replace(fileName, fileExt, "."+fileType, -1)
			} else if psf.Type == Page.LINK {
				fileName = psf.Title + "." + fileType
			}
=======

			fileName = fileName + "." + fileType
>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6

			var targetFilePath string
			targetFilePath = ""
			if psf.Type == Page.MARKDOWN {
				targetFilePath = filepath.Join(targetMarkdownFolder, fileName)
			} else if psf.Type == Page.HTML {
				targetFilePath = filepath.Join(targetHtmlFolder, fileName)
			} else if psf.Type == Page.LINK {
				targetFilePath = filepath.Join(targetLinkFolder, fileName)
			}

			if targetFilePath == "" {
				var errMsg = "Export Source File: Source File should be MARKDOWN,HTML or LINK"
				fmt.Println(errMsg)
				return false, errors.New(errMsg)
			}

			bSave, errSave := Utils.SaveBase64AsImage(psf.TitleImage, targetFilePath)

			if errSave != nil {
				return bSave, errSave
			}

		}
	}

	//Output the file information  FilePath|ID|LastModified

	var outputMdFiles []Page.PageSourceFile
	var outputHtmlFiles []Page.PageSourceFile
	var outputLinkFiles []Page.PageSourceFile

	for _, psf := range spp.SourceFiles {
		if psf.Type == Page.MARKDOWN {
			outputMdFiles = append(outputMdFiles, psf)
		} else if psf.Type == Page.HTML {
			outputHtmlFiles = append(outputHtmlFiles, psf)
		} else if psf.Type == Page.LINK {
			outputLinkFiles = append(outputLinkFiles, psf)
		}
	}

	var mdLength = len(outputMdFiles)
	var htmlLength = len(outputHtmlFiles)
	var linkLength = len(outputLinkFiles)
	var allCount = mdLength + htmlLength + linkLength

	fmt.Println("Exported:" + strconv.Itoa(allCount) + "`")
	fmt.Println("Markdown:" + strconv.Itoa(mdLength) + "`")
	for _, psf := range outputMdFiles {
		var str = psf.SourceFilePath + "|" + psf.ID + "|" + psf.LastModified + "`"
		fmt.Println(str)
	}
	fmt.Println("Html:" + strconv.Itoa(htmlLength) + "`")
	for _, psf := range outputHtmlFiles {
		var str = psf.SourceFilePath + "|" + psf.ID + "|" + psf.LastModified + "`"
		fmt.Println(str)
	}
	fmt.Println("Link:" + strconv.Itoa(linkLength) + "`")
	for _, psf := range outputLinkFiles {
<<<<<<< HEAD
		var str = psf.SourceFilePath + "|" + psf.ID + "|" + psf.Title + "|" + strconv.FormatBool(psf.IsTop) + "`"
=======
		var str = psf.SourceFilePath + "|" + psf.ID + "|" + psf.Title + "`"
>>>>>>> a006a98be915ab73501f3eaea84732339c2d4ca6
		fmt.Println(str)
	}
	return true, nil
}

/*
func (spp *SiteProject) GetProjectProperty(propertyName string) (string, error) {
	typeOfSiteProject := reflect.TypeOf(*spp)
	_, bFind := typeOfSiteProject.FieldByName(propertyName)

	if bFind == false {
		return "", errors.New("Cannot find field " + propertyName)
	}
	immutable := reflect.ValueOf(*spp)
	val := immutable.FieldByName(propertyName).String()
	return val, nil
}
*/
