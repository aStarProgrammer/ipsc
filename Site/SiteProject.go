// siteProject

package Site

import (
	"IPSC/Page"
	"IPSC/Utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type SiteProject struct {
	ID                 string
	Title              string
	Description        string
	Author             string
	CreateTime         string
	LastModified       string
	LastComplied       string
	LastComplieSummary string
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
		return false, errors.New("Argument jsonString is null")
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
		return "", errors.New("Pointer spp is nil")
	}

	if IsSiteProjectEmpty(*spp) {
		return "", errors.New("Site Project is empty")
	}

	_jsonbyte, err := json.Marshal(*spp)

	return string(_jsonbyte), err
}

func (spp *SiteProject) LoadFromFile(filePath string) (bool, error) {
	if "" == filePath {
		return false, errors.New("FilePath is empty")
	}

	bFileExist := Utils.PathIsExist(filePath)

	if false == bFileExist {
		return false, errors.New("File not exist")
	}

	_json, errRead := ioutil.ReadFile(filePath)

	if errRead != nil {
		return false, errors.New("Read File Fail")
	}

	_jsonString := string(_json)

	if "" == _jsonString {
		return false, errors.New("File is empty")
	}

	bUnMarshal, errUnMarshal := spp.FromJson(_jsonString)

	return bUnMarshal, errUnMarshal
}

func (spp *SiteProject) SaveToFile(filePath string) (bool, error) {
	if "" == filePath {
		return false, errors.New("FilePath is empty")
	}

	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	json, errMarshal := spp.ToJson()

	if errMarshal != nil {
		return false, errMarshal
	}

	var errFilePath error
	if !Utils.PathIsExist(filePath) {
		filePath, errFilePath = Utils.MakePath(filePath)
		if errFilePath != nil {
			return false, errors.New("Path nor exist and create parent folder failed")
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
		return false, errors.New("Write json to file failed")
	}

	return true, nil
}

func (spp *SiteProject) AddPageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
	}

	spp.SourceFiles = append(spp.SourceFiles, psf)
	return true, nil
}

func (spp *SiteProject) RemovePageSourceFile(psf Page.PageSourceFile, restore bool) (bool, error) {

	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
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
	return false, errors.New("Source Page not found")
}

func (spp *SiteProject) ResotrePageSourceFile(ID string) (bool, error) {
	var index = spp.GetPageSourceFile(ID)
	if index == -1 {
		return false, errors.New("Not find page source file with ID " + ID)
	}
	spp.SourceFiles[index].Status = Page.ACTIVE
	return true, nil
}

func (spp *SiteProject) UpdatePageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
	}

	for i, sf := range spp.SourceFiles {
		if sf.ID == psf.ID {
			spp.SourceFiles[i].Author = psf.Author
			spp.SourceFiles[i].CreateTime = psf.CreateTime
			spp.SourceFiles[i].Description = psf.Description
			spp.SourceFiles[i].LastComplied = psf.LastComplied
			spp.SourceFiles[i].LastModified = Utils.CurrentTime()
			spp.SourceFiles[i].OutputFile = psf.OutputFile
			spp.SourceFiles[i].SourceFilePath = psf.SourceFilePath
			spp.SourceFiles[i].Title = psf.Title
			spp.SourceFiles[i].Type = psf.Type
			spp.SourceFiles[i].IsTop = psf.IsTop
			return true, nil
		}
	}
	return false, errors.New("Source Page not found")
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
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageOutputFileEmpty(pof) {
		return false, errors.New("Output Page is empty")
	}

	spp.OutputFiles = append(spp.OutputFiles, pof)

	return true, nil
}

func (spp *SiteProject) RemovePageOutputFile(pof Page.PageOutputFile) (bool, error) {

	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageOutputFileEmpty(pof) {
		return false, errors.New("Output Page is empty")
	}

	for i, sf := range spp.OutputFiles {
		if sf.ID == pof.ID {
			spp.OutputFiles = append(spp.OutputFiles[:i], spp.OutputFiles[i+1:]...)
			return true, nil
		}
	}
	return false, errors.New("Output Page not found")
}

func (spp *SiteProject) UpdatePageOutputFile(pof Page.PageOutputFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageOutputFileEmpty(pof) {
		return false, errors.New("Output Page is empty")
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

	return false, errors.New("Output Page not found")
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
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
	}

	spp.MorePageSourceFiles = append(spp.MorePageSourceFiles, psf)
	return true, nil
}

func (spp *SiteProject) RemoveMorePageSourceFile(psf Page.PageSourceFile) (bool, error) {

	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
	}

	for i, sf := range spp.MorePageSourceFiles {
		if sf.ID == psf.ID {
			spp.MorePageSourceFiles = append(spp.MorePageSourceFiles[:i], spp.MorePageSourceFiles[i+1:]...)
			return true, nil
		}
	}
	return false, errors.New("Source Page not found")
}

func (spp *SiteProject) UpdateMorePageSourceFile(psf Page.PageSourceFile) (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
	}

	for i, sf := range spp.MorePageSourceFiles {
		if sf.ID == psf.ID {
			spp.MorePageSourceFiles[i].Author = psf.Author
			spp.MorePageSourceFiles[i].CreateTime = psf.CreateTime
			spp.MorePageSourceFiles[i].Description = psf.Description
			spp.MorePageSourceFiles[i].LastComplied = psf.LastComplied
			spp.MorePageSourceFiles[i].LastModified = psf.LastModified
			spp.MorePageSourceFiles[i].OutputFile = psf.OutputFile
			spp.MorePageSourceFiles[i].SourceFilePath = psf.SourceFilePath
			spp.MorePageSourceFiles[i].Title = psf.Title
			spp.MorePageSourceFiles[i].Type = psf.Type
			return true, nil
		}
	}
	return false, errors.New("Source Page not found")
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
		return false, errors.New("Site Project is empty")
	}

	if Page.IsPageSourceFileEmpty(psf) {
		return false, errors.New("Source Page is empty")
	}

	spp.IndexPageSourceFile = psf
	return true, nil
}

func (spp *SiteProject) CleanIndexPageSourceFile() (bool, error) {
	if IsSiteProjectEmpty(*spp) {
		return false, errors.New("Site Project is empty")
	}

	Page.ResetPageSourceFile(spp.IndexPageSourceFile)
	return true, nil
}

func (spp *SiteProject) PageStatistics() (string, error) {

	if IsSiteProjectEmpty(*spp) {
		return "", errors.New("Site Project is empty")
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
		return false, errors.New("Backup File is not exist")
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
