package Site

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ipsc/Configuration"
	"ipsc/Page"
	"ipsc/Utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type MarkdownPageModule struct {
	spp             *SiteProject
	smp             *SiteModule
	outputPageFiles Page.PageOutputFileSlice
}

func (mpmp *MarkdownPageModule) GetSiteProjectP() *SiteProject {
	return mpmp.spp
}

func (mpmp *MarkdownPageModule) GetSiteModuleP() *SiteModule {
	return mpmp.smp
}

func NewMarkdownPageModule(_spp *SiteProject, _smp *SiteModule) MarkdownPageModule {
	var mpm MarkdownPageModule
	mpm.spp = _spp
	mpm.smp = _smp
	return mpm
}

//Markdown extension .md .markdown .mmd .mdown
func FileIsMarkdown(filePath string) (bool, error) {
	if Utils.PathIsExist(filePath) == false {
		var errMsg = "Markdown file not exist"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var extension = filepath.Ext(filePath)

	if extension == ".md" || extension == ".markdown" || extension == ".mmd" || extension == ".mdown" {
		return true, nil
	}
	return false, nil
}

func (mpmp *MarkdownPageModule) AddMarkdown(title, description, author, filePath, titleImagePath string, isTop bool) (bool, string, error) {

	var markdownSrc, markdownDst string

	if Utils.PathIsExist(filePath) == false {
		var errMsg = "Markdown file not exist"
		fmt.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	bMarkdown, errMarkdown := FileIsMarkdown(filePath)

	if errMarkdown != nil {
		var errMsg = "Cannot confirm file type"
		fmt.Println(errMsg)
		return false, "", errors.New(errMsg)
	} else if bMarkdown == false {
		var errMsg = "File is not Markdown"
		fmt.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	_, fileName := filepath.Split(filePath)
	markdownSrc = filePath
	markdownDst = filepath.Join(mpmp.smp.GetSrcMarkdownFolderPath(mpmp.smp.GetProjectFolderPath()), fileName)

	_, errCopy := Utils.CopyFile(markdownSrc, markdownDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "Copy File from " + markdownSrc + " to " + markdownDst + " Failed"
		fmt.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	var psf Page.PageSourceFile
	psf = Page.NewPageSourceFile()

	psf.Title = title
	psf.Author = author
	psf.Description = description
	psf.SourceFilePath = markdownDst
	psf.LastModified = Utils.CurrentTime()
	psf.Status = Page.ACTIVE
	psf.Type = Page.MARKDOWN
	if Utils.PathIsExist(titleImagePath) && Utils.PathIsImage(titleImagePath) {
		//test titleImagePath size before alll the operations, so will not waste time if titleImage bigger than 30KB
		fileInfoTitleImage, errFileInfoTitleImage := os.Stat(titleImagePath)

		if errFileInfoTitleImage != nil {
			var errMsg = "Cannot get file size of titleImage"
			fmt.Println(errMsg)
			return false, "", errors.New(errMsg)
		}

		titleImageSize := fileInfoTitleImage.Size()

		if titleImageSize > MAXTITLEIMAGESIZE {
			var errMsg = "Title Image bigger than 30KB"
			fmt.Println(errMsg)
			return false, "", errors.New(errMsg)
		}

		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}
	psf.IsTop = isTop
	psf.OutputFile = -1

	bAdd, errorAdd := mpmp.spp.AddPageSourceFile(psf) //Add to Source Pages

	if bAdd == false && errorAdd != nil {
		fmt.Println(errorAdd.Error())
		return false, "", errorAdd
	}

	return true, psf.ID, nil
}

func (mpmp *MarkdownPageModule) GetPageSourceTemplateFile(pageType, templateFolderPath string) (string, error) {
	if Utils.PathIsExist(templateFolderPath) == false {
		var errMsg = "Template Folder Path not exist"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	switch pageType {
	case Page.MARKDOWN_NEWS:
		return filepath.Join(templateFolderPath, "News.md"), nil
	}
	return filepath.Join(templateFolderPath, "Blank.md"), nil
}

func (mpmp *MarkdownPageModule) CreateMarkdown(pageFilePath, markdownType, templateFolderPath string) (bool, error) {
	templateFilePath, errTemplate := mpmp.GetPageSourceTemplateFile(markdownType, templateFolderPath)
	if errTemplate != nil {
		return false, errTemplate
	}
	if Utils.PathIsExist(pageFilePath) {
		var errMsg = pageFilePath + " already exist, cannot create again"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	nCopy, errCopy := Utils.CopyFile(templateFilePath, pageFilePath)
	return nCopy > 0, errCopy
}

func (mpmp *MarkdownPageModule) RemoveMarkdown(psf Page.PageSourceFile, restore bool) (bool, error) {
	var outputIndex = psf.OutputFile
	if outputIndex != -1 {
		var pof = mpmp.spp.OutputFiles[outputIndex]
		if restore == false {
			bDelOutput, errDeleteOutput := mpmp.spp.RemovePageOutputFile(pof)
			if errDeleteOutput != nil {
				return bDelOutput, errDeleteOutput
			}
			if pof.FilePath != "" {
				bDeleteOutputFile := Utils.DeleteFile(pof.FilePath)
				if bDeleteOutputFile == false {
					var errMsg = "Cannot delete output file " + pof.FilePath
					fmt.Println(errMsg)
					return false, errors.New(errMsg)
				}
			}
		}

	}

	bRemove, errRemove := mpmp.spp.RemovePageSourceFile(psf, restore)
	if errRemove != nil {
		iFind := mpmp.spp.GetPageOutputFile(psf.ID)
		if iFind == -1 {
			mpmp.spp.AddPageSourceFile(psf)
		}
		return bRemove, errRemove
	}

	var filePath = psf.SourceFilePath
	if restore == false {
		if Utils.DeleteFile(filePath) == false {
			mpmp.spp.AddPageSourceFile(psf)
			var errMsg = "Delete File from Disk Fail"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) RestoreMarkdown(ID string) (bool, error) {
	return mpmp.spp.ResotrePageSourceFile(ID)
}

func (mpmp *MarkdownPageModule) UpdateMarkdown(psf Page.PageSourceFile, filePath string) (bool, error) {
	_psfID := mpmp.spp.GetPageSourceFile(psf.ID)
	if _psfID == -1 {
		var errMsg = "File not found"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	psf_Backup := mpmp.spp.SourceFiles[_psfID]

	if filePath == "" {
		return true, nil
	}

	if Utils.PathIsExist(filePath) == false {
		var errMsg = "Markdown file not exist"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	bMarkdown, errMarkdown := FileIsMarkdown(filePath)

	if errMarkdown != nil {
		var errMsg = "Cannot confirm file type"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	} else if bMarkdown == false {
		var errMsg = "File is not Markdown"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	_, fileName := filepath.Split(filePath)

	var markdownSrc, markdownDst string
	markdownSrc = filePath
	markdownDst = filepath.Join(mpmp.smp.GetSrcMarkdownFolderPath(mpmp.smp.GetProjectFolderPath()), fileName)
	psf.SourceFilePath = markdownDst

	bUpdate, errUpdate := mpmp.spp.UpdatePageSourceFile(psf)

	if errUpdate != nil {
		return bUpdate, errUpdate
	}

	if psf.SourceFilePath != psf_Backup.SourceFilePath {
		Utils.DeleteFile(psf_Backup.SourceFilePath)
	}

	_, errCopy := Utils.CopyFile(markdownSrc, markdownDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "Copy File from " + markdownSrc + " to " + markdownDst + " Failed"
		//恢复被更新的内容
		mpmp.spp.UpdatePageSourceFile(psf_Backup)
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	return true, nil
}

func (mpmp *MarkdownPageModule) GetMarkdownFile(ID string) string {
	iFind := mpmp.spp.GetPageSourceFile(ID)
	if iFind != -1 {
		psf := mpmp.spp.SourceFiles[iFind]

		if psf.SourceFilePath != "" {
			return psf.SourceFilePath
		}
	}

	return ""
}

func (mpmp *MarkdownPageModule) GetMarkdownInformation(ID string) int {
	return mpmp.spp.GetPageSourceFile(ID)
}

func (mpmp *MarkdownPageModule) UpdateMarkdownInformation(title, description, author, filePath, titleImagePath string) (bool, error) {
	var psf Page.PageSourceFile
	psf = Page.NewPageSourceFile()

	psf.Title = title
	psf.Author = author
	psf.Description = description
	psf.SourceFilePath = filePath
	psf.LastModified = Utils.CurrentTime()
	psf.Status = Page.ACTIVE
	psf.Type = Page.MARKDOWN
	if Utils.PathIsExist(titleImagePath) && Utils.PathIsImage(titleImagePath) {
		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}

	bUpdate, errorUpdate := mpmp.spp.UpdatePageSourceFile(psf) //Update Source Pages

	if bUpdate == false && errorUpdate != nil {
		fmt.Println("UpdateMarkdown:" + errorUpdate.Error())
		return false, errorUpdate
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) Compile_Psf(psf Page.PageSourceFile) (int, error) {
	if psf.SourceFilePath == "" {
		var errMsg string
		errMsg = "Page Source File FilePath is emtpy"
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	var markdownSrc, markdownDst, cssFilePath string
	markdownSrc = psf.SourceFilePath

	if Utils.PathIsExist(markdownSrc) == false {
		var errMsg = "Source Markdown File not found on the disk"
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	bMarkdown, errMarkdown := FileIsMarkdown(markdownSrc)

	if errMarkdown != nil {
		var errMsg = "Cannot confirm file type"
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	} else if bMarkdown == false {
		var errMsg = "File is not html"
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	fileName := filepath.Base(markdownSrc)
	ext := filepath.Ext(markdownSrc)
	fileNameOnly := strings.TrimSuffix(fileName, ext)

	if psf.Type == Page.INDEX {
		newFileName := fileNameOnly + ".html"
		markdownDst = filepath.Join(mpmp.smp.GetOutputFolderPath(mpmp.smp.GetProjectFolderPath()), newFileName)
	} else {
		newFileName := fileNameOnly + "_" + Utils.GUID() + ".html"
		markdownDst = filepath.Join(mpmp.smp.GetOutputFolderPath(mpmp.smp.GetProjectFolderPath()), "Pages", newFileName)
	}
	var errCssFilePath error
	cssFilePath, errCssFilePath = Configuration.GetCssFilePath()

	if cssFilePath == "" || errCssFilePath != nil {
		var errMsg = "Css File Path is empty"
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	if Utils.PathIsExist(cssFilePath) == false {
		var errMsg string
		errMsg = "Css File Path " + cssFilePath + " not exist"
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	//Call pandoc to convert md to html

	var strPandocCmd = "pandoc -s --self-contained -c \"" + cssFilePath + "\" \"" + markdownSrc + "\" -o \"" + markdownDst + "\" --metadata pagetitle=\"" + psf.Title + "\""

	//fmt.Println(strPandocCmd)
	sysType := runtime.GOOS

	var pandocCmd *exec.Cmd

	if "windows" == sysType {
		pandocCmd = exec.Command("Powershell", strPandocCmd)
	} else if "linux" == sysType || "darwin" == sysType {
		pandocCmd = exec.Command("bash", "-c", strPandocCmd)
	} else { //Not support other platforms now
		var errMsg string
		errMsg = "Compile Markdown, not supported platform " + sysType
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	_, errPandoc := pandocCmd.Output()
	if errPandoc != nil {
		fmt.Println(errPandoc.Error())
		return -1, errPandoc
	}

	//Add pof
	pof := Page.NewPageOutputFile()
	pof.Author = psf.Author
	pof.Description = psf.Description
	pof.FilePath = markdownDst
	pof.IsTop = psf.IsTop
	pof.Title = psf.Title
	pof.TitleImage = psf.TitleImage
	pof.Type = psf.Type

	_, errAdd := mpmp.spp.AddPageOutputFile(pof)

	if errAdd != nil {
		Utils.DeleteFile(markdownDst) //Add fail,delete the file already copied
		return -1, errAdd
	}

	_pofID := mpmp.spp.GetPageOutputFile(pof.ID)

	if _pofID == -1 {
		Utils.DeleteFile(markdownDst) //Add fail,delete the file already copied
		var errMsg = "Page Output File add Fail"
		fmt.Println(errMsg)
		return _pofID, errors.New(errMsg)
	}

	psf.OutputFile = _pofID
	psf.LastCompiled = Utils.CurrentTime()

	mpmp.spp.UpdatePageSourceFile(psf)

	return _pofID, nil
}

//Compile Markdown, call pandoc to convert md to html to Output folder
//change sourceinformation and add PageOutputFile
func (mpmp *MarkdownPageModule) Compile(ID string) (int, error) {
	iFind := mpmp.spp.GetPageSourceFile(ID)
	if iFind == -1 {
		var errMsg string
		errMsg = "Cannot find the source File with ID " + ID
		fmt.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	psf := mpmp.spp.SourceFiles[iFind]
	return mpmp.Compile_Psf(psf)
}

func (mpmp *MarkdownPageModule) CreateIndexPage(indexPageSize string) (bool, error) {
	//Get template file path for index page
	indexTemplateFilePath, errIndexTemplateFilePath := Configuration.GetIndexTemplateFilePath(indexPageSize)

	if errIndexTemplateFilePath != nil {
		var errMsg string
		errMsg = "Cannot find index page template file for page size " + indexPageSize
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	//Set IndexPageSourceFile Properties
	mpmp.spp.IndexPageSourceFile = Page.NewPageSourceFile()
	mpmp.spp.IndexPageSourceFile.Author = mpmp.spp.Author
	mpmp.spp.IndexPageSourceFile.Description = mpmp.spp.Description
	mpmp.spp.IndexPageSourceFile.IsTop = false
	mpmp.spp.IndexPageSourceFile.LastModified = Utils.CurrentTime()
	mpmp.spp.IndexPageSourceFile.Status = Page.ACTIVE
	mpmp.spp.IndexPageSourceFile.Type = Page.INDEX
	mpmp.spp.IndexPageSourceFile.Title = mpmp.spp.Title
	mpmp.spp.IndexPageSourceFile.TitleImage = ""
	mpmp.spp.IndexPageSourceFile.SourceFilePath = filepath.Join(mpmp.smp.GetSrcMarkdownFolderPath(mpmp.smp.GetProjectFolderPath()), "index.md")

	//Copy index template file to markdown folder
	var srcIndexPageSourceFilePath = mpmp.spp.IndexPageSourceFile.SourceFilePath
	nByte, errCopy := Utils.CopyFile(indexTemplateFilePath, srcIndexPageSourceFilePath)

	if nByte == 0 && errCopy != nil {
		var errMsg string
		errMsg = "Copy Index Template File from " + indexTemplateFilePath + " to " + srcIndexPageSourceFilePath + " failed, will reset index page properties in site project file"
		fmt.Println(errMsg)
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = "Clean Index Page properties failed, please check site project file"
			fmt.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}

		return false, errors.New(errMsg)
	}

	//Start to modify the template md file
	//Sort the output file
	if len(mpmp.outputPageFiles) == 0 {
		topOutputPageFiles, errTop := mpmp.spp.GetSortedTopOutputFiles()
		if errTop != nil {
			var errMsg string
			errMsg = "Cannot get top Page Source File"
			fmt.Println(errMsg)

			bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

			if errClean != nil {
				var errCleanMsg string
				errCleanMsg = "Clean Index Page properties failed, please check site project file"
				fmt.Println(errCleanMsg)
				return bClean, errors.New(errCleanMsg)
			}

			bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
				fmt.Println(deleteMsg)
				fmt.Println(deleteMsg)
			}

			return false, errors.New(errMsg)
		}

		normalOutputPageFiles, errNormal := mpmp.spp.GetSortedNormalOutputFiles()

		if errNormal != nil {
			var errMsg string
			errMsg = "Cannot get normal Page Source File"
			fmt.Println(errMsg)

			bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

			if errClean != nil {
				var errCleanMsg string
				errCleanMsg = "Clean Index Page properties failed, please check site project file"
				fmt.Println(errCleanMsg)
				return bClean, errors.New(errCleanMsg)
			}

			bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
				fmt.Println(deleteMsg)
				fmt.Println(deleteMsg)
			}

			return false, errors.New(errMsg)
		}

		mpmp.outputPageFiles = append(topOutputPageFiles, normalOutputPageFiles...)
	}
	// get the first pagesize items
	nIndexPageSize, errNIndexPageSize := Page.ConvertPageSize2Int(indexPageSize)

	if errNIndexPageSize != nil {
		var errMsg string
		errMsg = "Cannot get page size"
		fmt.Println(errMsg)
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = "Clean Index Page properties failed, please check site project file"
			fmt.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}
		bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
			fmt.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}
	var indexOutputPageFiles Page.PageOutputFileSlice
	if nIndexPageSize <= len(mpmp.outputPageFiles) {
		indexOutputPageFiles = mpmp.outputPageFiles[:nIndexPageSize]
	} else {
		indexOutputPageFiles = mpmp.outputPageFiles
	}

	// modify copied index page  template md
	// Read file
	bIndexFileContent, errReadFile := ioutil.ReadFile(srcIndexPageSourceFilePath)

	if errReadFile != nil {
		var errMsg string
		errMsg = "Cannot read src Index md file"
		fmt.Println(errMsg)
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = " Read file fail,then Clean Index Page properties failed, please check site project file"
			fmt.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}
		bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
			fmt.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}

	//Update md file info

	indexFileContent := string(bIndexFileContent)

	indexFileContent = strings.Replace(indexFileContent, Page.INDEX_PAGE_TITLE_MARK, mpmp.spp.IndexPageSourceFile.Title, -1)

	for index, indexOutputPageFile := range indexOutputPageFiles {
		var indexNewsTitleMark, indexNewsUrlMark, indexNewsImageMark, indexNewsTimeMark string

		indexNewsTitleMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_TITLE_MARK
		indexNewsUrlMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_URL_MARK
		indexNewsImageMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_IMAGE_MARK
		indexNewsTimeMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_TIME_MARK

		if indexOutputPageFile.Title != "" {
			indexFileContent = strings.Replace(indexFileContent, indexNewsTitleMark, indexOutputPageFile.Title, 1)
		} else {
			var errMsg = "CreateIndexPage: Title of Item is empty: " + indexOutputPageFile.ID
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}

		if indexOutputPageFile.FilePath != "" {
			if indexOutputPageFile.Type == Page.MARKDOWN || indexOutputPageFile.Type == Page.HTML {
				_, indexOutputHtmlName := filepath.Split(indexOutputPageFile.FilePath)
				if indexOutputHtmlName != "" {
					indexOutputHtmlName = "./Pages/" + indexOutputHtmlName
					indexFileContent = strings.Replace(indexFileContent, indexNewsUrlMark, indexOutputHtmlName, 1)
				}
			} else if indexOutputPageFile.Type == Page.LINK {
				indexFileContent = strings.Replace(indexFileContent, indexNewsUrlMark, indexOutputPageFile.FilePath, 1)
			}
		} else {
			var errMsg = "CreateIndexPage: FilePath of Item is empty: " + indexOutputPageFile.ID
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}

		if indexOutputPageFile.TitleImage != "" {
			indexFileContent = strings.Replace(indexFileContent, indexNewsImageMark, indexOutputPageFile.TitleImage, 1)
		} else {
			fmt.Println("CreateIndexPage: TitleImage of Item is empty: " + indexOutputPageFile.FilePath + " This item will not have title image in index.html")

			emptyImageTemplate, errEmptyImage := Configuration.GetEmptyImageItemTemplate()
			if errEmptyImage != nil {
				var errMsg = "CreateIndexPage: Cannot get empty image template"
				fmt.Println(errMsg)
				return false, errors.New(errMsg)
			}
			emptyImageTemplate = strings.Replace(emptyImageTemplate, Page.INDEX_NEWS_IMAGE_MARK, indexNewsImageMark, 1)
			indexFileContent = strings.Replace(indexFileContent, emptyImageTemplate, "", -1)
		}

		if indexOutputPageFile.CreateTime != "" {
			indexFileContent = strings.Replace(indexFileContent, indexNewsTimeMark, indexOutputPageFile.CreateTime, 1)
		} else {
			var errMsg = "CreateIndexPage: CreateTime of Item is empty " + indexOutputPageFile.ID
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	//Delete item template in md that not used when indexPageSize > outputFiles.Count
	// Not enough output files to update the md file, so need to remove the item with orignial
	// item mask,looks like
	//<font size=4>[NEWSTITLE_E6F6DF62-5BC6-4172-86F1-1250F8618E0F](NEWSURL_1C387CE9-FFE9-469F-96E5-E4FAA83DF668) </font><img align="right" src="NEWSIMAGE_870BB9B8-20CB-45B0-86F7-BEC643321376" />

	//<br> NEWSTIME_EC093DDF-B972-4775-9F3E-44CB493E5D07

	if nIndexPageSize > len(mpmp.outputPageFiles) {
		//read empty item template
		emptyItemTemplate, errEmptyItemTemplate := Configuration.GetEmptyIndexItemTemplate()

		if errEmptyItemTemplate != nil {
			var errMsg string
			errMsg = "Cannot read empty item template from item"
			fmt.Println(errMsg)

			bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

			if errClean != nil {
				var errCleanMsg string
				errCleanMsg = " Read file fail,then Clean Index Page properties failed, please check site project file"
				fmt.Println(errCleanMsg)
				return bClean, errors.New(errCleanMsg)
			}
			bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
				fmt.Println(deleteMsg)
			}
			return false, errors.New(errMsg)
		}

		var emptyStartIndex = len(mpmp.outputPageFiles)
		for emptyIndex := emptyStartIndex; emptyIndex < nIndexPageSize; emptyIndex++ {
			var indexNewsTitleMark, indexNewsUrlMark, indexNewsImageMark, indexNewsTimeMark string

			indexNewsTitleMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_TITLE_MARK
			indexNewsUrlMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_URL_MARK
			indexNewsImageMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_IMAGE_MARK
			indexNewsTimeMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_TIME_MARK

			//build emptyItem for each emptyItem
			var emptyItem string
			emptyItem = emptyItemTemplate

			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_TITLE_MARK, indexNewsTitleMark, -1)
			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_URL_MARK, indexNewsUrlMark, -1)
			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_IMAGE_MARK, indexNewsImageMark, -1)
			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_TIME_MARK, indexNewsTimeMark, -1)

			//Replace emptyItem with ""
			indexFileContent = strings.Replace(indexFileContent, emptyItem, "", 1)
			//fmt.Println(indexFileContent)
		}

	}
	//fmt.Println(indexFileContent)
	// save file
	errWriteFile := ioutil.WriteFile(srcIndexPageSourceFilePath, []byte(indexFileContent), 0x666)

	if errWriteFile != nil {
		var errMsg string
		errMsg = "Cannot Save content to index md file"
		fmt.Println(errMsg)
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = " Read file fail,then Clean Index Page properties failed, please check site project file"
			fmt.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}
		bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
			fmt.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) CreateMorePage(indexPageSize string, startIndex, pageNo int) (bool, error) {
	//Get template file path for index page
	moreTemplateFilePath, errMoreTemplateFilePath := Configuration.GetMoreTemplateFilePath(indexPageSize)

	if errMoreTemplateFilePath != nil {
		var errMsg string
		errMsg = "Cannot find more page template file for page size " + indexPageSize
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	var morePageSourceFile = Page.NewPageSourceFile()

	//Set morePageSourceFile Properties
	morePageSourceFile = Page.NewPageSourceFile()
	morePageSourceFile.Author = mpmp.spp.Author
	morePageSourceFile.Description = mpmp.spp.Description
	morePageSourceFile.IsTop = false
	morePageSourceFile.LastModified = Utils.CurrentTime()
	morePageSourceFile.Status = Page.ACTIVE
	morePageSourceFile.Type = Page.INDEX
	morePageSourceFile.Title = mpmp.spp.Title
	morePageSourceFile.TitleImage = ""
	var morePageName = "more" + strconv.Itoa(pageNo) + ".md"
	morePageSourceFile.SourceFilePath = filepath.Join(mpmp.smp.GetSrcMarkdownFolderPath(mpmp.smp.GetProjectFolderPath()), morePageName)

	_, errAddMorePage := mpmp.spp.AddMorePageSourceFile(morePageSourceFile)

	if errAddMorePage != nil {
		var errMsg string
		errMsg = "Cannot add More Page Source File"
		fmt.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "CreateMorePage: Cannot delete it"
			fmt.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}

		return false, errors.New(errMsg)
	}

	//Copy more template file to markdown folder
	var srcMorePageSourceFilePath = morePageSourceFile.SourceFilePath
	nByte, errCopy := Utils.CopyFile(moreTemplateFilePath, srcMorePageSourceFilePath)

	if nByte == 0 && errCopy != nil {
		var errMsg string
		errMsg = "Copy More Template File from " + moreTemplateFilePath + " to " + srcMorePageSourceFilePath + " failed, will remove this more page in site project file"
		fmt.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "Create More Page: Cannot delete more page properties from site project"
			fmt.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}

		return false, errors.New(errMsg)
	}

	//Start to modify the template md file
	//Sort the output file
	if len(mpmp.outputPageFiles) == 0 {
		topOutputPageFiles, errTop := mpmp.spp.GetSortedTopOutputFiles()
		if errTop != nil {
			var errMsg string
			errMsg = "Cannot get top Output Page Files"
			fmt.Println(errMsg)

			bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

			if errRemove != nil {
				var errRemoveMsg string
				errRemoveMsg = "Cannot delete more page properties from site project"
				fmt.Println(errRemoveMsg)
				return bRemove, errors.New(errRemoveMsg)
			}
			bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
				fmt.Println(deleteMsg)
			}
			return false, errors.New(errMsg)
		}

		normalOutputPageFiles, errNormal := mpmp.spp.GetSortedNormalOutputFiles()

		if errNormal != nil {
			var errMsg string
			errMsg = "Cannot get sorted normal output page files"
			fmt.Println(errMsg)

			bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

			if errRemove != nil {
				var errRemoveMsg string
				errRemoveMsg = "Cannot delete more page properties from site project"
				fmt.Println(errRemoveMsg)
				return bRemove, errors.New(errRemoveMsg)
			}
			bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
				fmt.Println(deleteMsg)
			}
			return false, errors.New(errMsg)
		}

		mpmp.outputPageFiles = append(topOutputPageFiles, normalOutputPageFiles...)
	}
	// get the first pagesize items
	nMorePageSize, errNMorePageSize := Page.ConvertPageSize2Int(indexPageSize)

	if errNMorePageSize != nil {
		var errMsg string
		errMsg = "Cannot get page size"
		fmt.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "Cannot delete more page properties from site project"
			fmt.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}
		bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
			fmt.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}

	var moreOutputPageFiles Page.PageOutputFileSlice
	if startIndex+nMorePageSize <= len(mpmp.outputPageFiles) {
		moreOutputPageFiles = mpmp.outputPageFiles[startIndex : startIndex+nMorePageSize]
	} else {
		moreOutputPageFiles = mpmp.outputPageFiles[startIndex:]
	}

	// modify copied index page  template md
	// Read file
	bMoreFileContent, errReadFile := ioutil.ReadFile(srcMorePageSourceFilePath)

	if errReadFile != nil {
		var errMsg string
		errMsg = "Cannot read More Page md file"
		fmt.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "Cannot delete more page properties from site project"
			fmt.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}
		bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
			fmt.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}

	//Update md file
	moreFileContent := string(bMoreFileContent)

	moreFileContent = strings.Replace(moreFileContent, Page.INDEX_PAGE_TITLE_MARK, morePageSourceFile.Title, -1)

	for index, moreOutputPageFile := range moreOutputPageFiles {
		var moreNewsTitleMark, moreNewsUrlMark, moreNewsImageMark, moreNewsTimeMark string

		moreNewsTitleMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_TITLE_MARK
		moreNewsUrlMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_URL_MARK
		moreNewsImageMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_IMAGE_MARK
		moreNewsTimeMark = strconv.Itoa(index) + "_" + Page.INDEX_NEWS_TIME_MARK

		if moreOutputPageFile.Title != "" {
			moreFileContent = strings.Replace(moreFileContent, moreNewsTitleMark, moreOutputPageFile.Title, 1)
		} else {
			return false, errors.New("CreateMorePage: Title of Item is empty" + moreOutputPageFile.ID)
			moreFileContent = strings.Replace(moreFileContent, moreNewsTitleMark, "", 1)
		}

		if moreOutputPageFile.FilePath != "" {
			if moreOutputPageFile.Type == Page.MARKDOWN || moreOutputPageFile.Type == Page.HTML {
				_, moreOutputHtmlName := filepath.Split(moreOutputPageFile.FilePath)
				if moreOutputHtmlName != "" {
					moreOutputHtmlName = "./Pages/" + moreOutputHtmlName
					moreFileContent = strings.Replace(moreFileContent, moreNewsUrlMark, moreOutputHtmlName, 1)
				}
			} else if moreOutputPageFile.Type == Page.LINK {
				moreFileContent = strings.Replace(moreFileContent, moreNewsUrlMark, moreOutputPageFile.FilePath, 1)
			}
		} else {
			var errMsg = "CreateMorePage: FilePath of Item is empty: " + moreOutputPageFile.ID
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}

		if moreOutputPageFile.TitleImage != "" {
			moreFileContent = strings.Replace(moreFileContent, moreNewsImageMark, moreOutputPageFile.TitleImage, 1)
		} else {
			fmt.Println("CreateMorePage: TitleImage of Item is empty: " + moreOutputPageFile.FilePath + " This item will not have title image in moreXX.html")

			emptyImageTemplate, errEmptyImage := Configuration.GetEmptyImageItemTemplate()
			if errEmptyImage != nil {
				var errMsg = "CreateMorePage: Cannot get empty image template"
				fmt.Println(errMsg)
				return false, errors.New(errMsg)
			}
			emptyImageTemplate = strings.Replace(emptyImageTemplate, Page.INDEX_NEWS_IMAGE_MARK, moreNewsImageMark, 1)
			moreFileContent = strings.Replace(moreFileContent, emptyImageTemplate, "", -1)
		}

		if moreOutputPageFile.CreateTime != "" {
			moreFileContent = strings.Replace(moreFileContent, moreNewsTimeMark, moreOutputPageFile.CreateTime, 1)
		} else {
			var errMsg = "CreateMorePage: CreateTime of Item is empty: " + moreOutputPageFile.ID
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	if startIndex+nMorePageSize > len(mpmp.outputPageFiles) {
		//read empty item template
		emptyItemTemplate, errEmptyItemTemplate := Configuration.GetEmptyIndexItemTemplate()

		if errEmptyItemTemplate != nil {
			var errMsg string
			errMsg = "Cannot read empty Item template file"
			fmt.Println(errMsg)

			bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

			if errRemove != nil {
				var errRemoveMsg string
				errRemoveMsg = "Cannot delete more page properties from site project"
				fmt.Println(errRemoveMsg)
				return bRemove, errors.New(errRemoveMsg)
			}
			bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
				fmt.Println(deleteMsg)
			}
			return false, errors.New(errMsg)
		}

		var emptyStartIndex = len(mpmp.outputPageFiles) - startIndex
		for emptyIndex := emptyStartIndex; emptyIndex < nMorePageSize; emptyIndex++ {
			var moreNewsTitleMark, moreNewsUrlMark, moreNewsImageMark, moreNewsTimeMark string

			moreNewsTitleMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_TITLE_MARK
			moreNewsUrlMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_URL_MARK
			moreNewsImageMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_IMAGE_MARK
			moreNewsTimeMark = strconv.Itoa(emptyIndex) + "_" + Page.INDEX_NEWS_TIME_MARK

			//build emptyItem for each emptyItem
			var emptyItem string
			emptyItem = emptyItemTemplate

			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_TITLE_MARK, moreNewsTitleMark, -1)
			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_URL_MARK, moreNewsUrlMark, -1)
			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_IMAGE_MARK, moreNewsImageMark, -1)
			emptyItem = strings.Replace(emptyItem, Page.INDEX_NEWS_TIME_MARK, moreNewsTimeMark, -1)

			//Replace emptyItem with ""
			moreFileContent = strings.Replace(moreFileContent, emptyItem, "", 1)
		}

	}

	// save file
	errWriteFile := ioutil.WriteFile(srcMorePageSourceFilePath, []byte(moreFileContent), 0x666)

	if errWriteFile != nil {
		var errMsg string
		errMsg = "Cannot Save modified md file"
		fmt.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "Cannot delete more page properties from site project"
			fmt.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}
		bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
			fmt.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) CreateNavigationForIndexPage() (bool, error) {
	//Get template file path for index page

	if mpmp.spp.IndexPageSourceFile.SourceFilePath == "" {
		var errMsg = "Index Page file not created,please run CreateIndexPage firstly"
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	srcIndexPageSourceFilePath := mpmp.spp.IndexPageSourceFile.SourceFilePath
	// modify copied index page  template md
	// Read file
	bIndexFileContent, errReadFile := ioutil.ReadFile(srcIndexPageSourceFilePath)

	if errReadFile != nil {
		var errMsg string
		errMsg = "Cannot read src Index md file"
		fmt.Println(errMsg)

		return false, errors.New(errMsg)
	}

	//Update md file info

	indexFileContent := string(bIndexFileContent)

	if len(mpmp.spp.MorePageSourceFiles) == 0 {
		var linkmoreStr = "[More...](" + Page.INDEX_LINK_MORE_MARK + ")"
		indexFileContent = strings.Replace(indexFileContent, linkmoreStr, "", -1)
	} else {
		indexFileContent = strings.Replace(indexFileContent, Page.INDEX_LINK_MORE_MARK, "more1.html", -1)
	}
	// save file
	errWriteFile := ioutil.WriteFile(srcIndexPageSourceFilePath, []byte(indexFileContent), 0x666)

	if errWriteFile != nil {
		var errMsg string
		errMsg = "Cannot Save content to index md file"
		fmt.Println(errMsg)

		return false, errors.New(errMsg)
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) CreateNavigationForMorePages() (bool, error) {
	var morePageCount = len(mpmp.spp.MorePageSourceFiles)

	if morePageCount == 0 {
		return false, nil
	}

	var navigationString string
	navigationString = ""
	//Create navigation mark txt
	for c := 1; c <= morePageCount; c++ {
		navigationString = navigationString + "[[" + strconv.Itoa(c) + "](more" + strconv.Itoa(c) + ".html)]   "
	}

	navigationString = strings.TrimRight(navigationString, " ")

	// modify more page md
	// Read file
	for _, mpsf := range mpmp.spp.MorePageSourceFiles {
		srcMorePageSourceFilePath := mpsf.SourceFilePath
		bMoreFileContent, errReadFile := ioutil.ReadFile(srcMorePageSourceFilePath)

		if errReadFile != nil {
			var errMsg string
			errMsg = "Cannot read More Page md file " + srcMorePageSourceFilePath
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}

		//Update md file
		moreFileContent := string(bMoreFileContent)

		moreFileContent = strings.Replace(moreFileContent, Page.MORE_LINK_INDEX_MARK, "index.html", -1)
		moreFileContent = strings.Replace(moreFileContent, Page.MORE_PAGE_LINK_MARK, navigationString, -1)

		// save file
		errWriteFile := ioutil.WriteFile(srcMorePageSourceFilePath, []byte(moreFileContent), 0x666)

		if errWriteFile != nil {
			var errMsg string
			errMsg = "Cannot Save modified md file"
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	return true, nil
}
