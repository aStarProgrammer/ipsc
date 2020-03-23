package Site

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"ipsc/Configuration"
	"ipsc/Page"
	"ipsc/Utils"
	"os"
	"os/exec"
	"path/filepath"
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
		var errMsg = "MarkdownPageModule.FileIsMarkdown: Markdown file not exist"
		Utils.Logger.Println(errMsg)
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
		var errMsg = "MarkdownPageModule.AddMarkdown: Markdown file not exist"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	bMarkdown, errMarkdown := FileIsMarkdown(filePath)

	if errMarkdown != nil {
		var errMsg = "MarkdownPageModule.AddMarkdown: Cannot confirm file type"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	} else if bMarkdown == false {
		var errMsg = "MarkdownPageModule.AddMarkdown: File is not Markdown"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	_, fileName := filepath.Split(filePath)
	markdownSrc = filePath
	markdownDst = filepath.Join(mpmp.smp.GetSrcMarkdownFolderPath(mpmp.smp.GetProjectFolderPath()), fileName)

	if Utils.PathIsExist(markdownDst) {
		var errMsg = "MarkdownPageModule.AddMarkdown: Target Markdown File Already Exist"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	_, errCopy := Utils.CopyFile(markdownSrc, markdownDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.AddMarkdown: Copy File from " + markdownSrc + " to " + markdownDst + " Failed"
		Utils.Logger.Println(errMsg)
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
			var errMsg = "MarkdownPageModule.AddMarkdown: Cannot get file size of titleImage"
			Utils.Logger.Println(errMsg)
			return false, "", errors.New(errMsg)
		}

		titleImageSize := fileInfoTitleImage.Size()

		if titleImageSize > MAXTITLEIMAGESIZE {
			var errMsg = "MarkdownPageModule.AddMarkdown: Title Image bigger than 30KB"
			Utils.Logger.Println(errMsg)
			return false, "", errors.New(errMsg)
		}

		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}
	psf.IsTop = isTop
	psf.OutputFile = ""

	bAdd, errorAdd := mpmp.spp.AddPageSourceFile(psf) //Add to Source Pages

	if bAdd == false && errorAdd != nil {
		Utils.Logger.Println(errorAdd.Error())
		return false, "", errorAdd
	}

	return true, psf.ID, nil
}

func (mpmp *MarkdownPageModule) GetPageSourceTemplateFile(pageType, templateFolderPath string) (string, error) {
	if Utils.PathIsExist(templateFolderPath) == false {
		var errMsg = "MarkdownPageModule.GetPageSourceTemplateFile: Template Folder Path not exist"
		Utils.Logger.Println(errMsg)
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
		var errMsg = "MarkdownPageModule.CreateMarkdown: " + pageFilePath + " already exist, cannot create again"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	nCopy, errCopy := Utils.CopyFile(templateFilePath, pageFilePath)
	return nCopy > 0, errCopy
}

func (mpmp *MarkdownPageModule) RemoveMarkdown(psf Page.PageSourceFile, restore bool) (bool, error) {
	var outputID = psf.OutputFile
	if outputID != "" {
		var pofIndex = mpmp.spp.GetPageOutputFile(outputID)
		var pof = mpmp.spp.OutputFiles[pofIndex]
		if restore == false {
			bDelOutput, errDeleteOutput := mpmp.spp.RemovePageOutputFile(pof)
			if errDeleteOutput != nil {
				return bDelOutput, errDeleteOutput
			}
			if pof.FilePath != "" {
				bDeleteOutputFile := Utils.DeleteFile(pof.FilePath)
				if bDeleteOutputFile == false {
					var errMsg = "MarkdownPageModule.RemoveMarkdown: Cannot delete output file " + pof.FilePath
					Utils.Logger.Println(errMsg)
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
			var errMsg = "MarkdownPageModule.RemoveMarkdown: Delete File from Disk Fail"
			Utils.Logger.Println(errMsg)
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
		var errMsg = "MarkdownPageModule.UpdateMarkdown: File not found"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	psf_Backup := mpmp.spp.SourceFiles[_psfID]

	if filePath != psf.SourceFilePath {

		if Utils.PathIsExist(filePath) == false {
			var errMsg = "MarkdownPageModule.UpdateMarkdown: Markdown file not exist"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		bMarkdown, errMarkdown := FileIsMarkdown(filePath)

		if errMarkdown != nil {
			var errMsg = "MarkdownPageModule.UpdateMarkdown: Cannot confirm file type"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		} else if bMarkdown == false {
			var errMsg = "MarkdownPageModule.UpdateMarkdown: File is not Markdown"
			Utils.Logger.Println(errMsg)
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
			errMsg = "MarkdownPageModule.UpdateMarkdown: Copy File from " + markdownSrc + " to " + markdownDst + " Failed"
			//恢复被更新的内容
			mpmp.spp.UpdatePageSourceFile(psf_Backup)
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
	} else {
		bUpdate, errUpdate := mpmp.spp.UpdatePageSourceFile(psf)
		if errUpdate != nil {
			return bUpdate, errUpdate
		}
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
		Utils.Logger.Println("MarkdownPageModule.UpdateMarkdown: " + errorUpdate.Error())
		return false, errorUpdate
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) Compile_Psf(psf Page.PageSourceFile) (int, error) {
	//fmt.Println("A")
	if psf.SourceFilePath == "" {
		var errMsg string
		errMsg = "MarkdownPageModule.Compile_Psf: Page Source File FilePath is emtpy"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	var markdownSrc, markdownDst, cssFilePath string
	markdownSrc = psf.SourceFilePath

	if Utils.PathIsExist(markdownSrc) == false {
		var errMsg = "MarkdownPageModule.Compile_Psf: Source Markdown File not found on the disk"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	bMarkdown, errMarkdown := FileIsMarkdown(markdownSrc)
	//fmt.Println("B")
	if errMarkdown != nil {
		var errMsg = "MarkdownPageModule.Compile_Psf: Cannot confirm file type"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	} else if bMarkdown == false {
		var errMsg = "MarkdownPageModule.Compile_Psf: File is not html"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	fileName := filepath.Base(markdownSrc)
	ext := filepath.Ext(markdownSrc)
	fileNameOnly := strings.TrimSuffix(fileName, ext)
	newFileName := fileNameOnly + ".html"

	if psf.Type == Page.INDEX {
		markdownDst = filepath.Join(mpmp.smp.GetOutputFolderPath(mpmp.smp.GetProjectFolderPath()), newFileName)
	} else {
		markdownDst = filepath.Join(mpmp.smp.GetOutputFolderPath(mpmp.smp.GetProjectFolderPath()), "Pages", newFileName)
	}
	var errCssFilePath error
	cssFilePath, errCssFilePath = Configuration.GetCssFilePath()

	if cssFilePath == "" || errCssFilePath != nil {
		var errMsg = "MarkdownPageModule.Compile_Psf: Css File Path is empty"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	if Utils.PathIsExist(cssFilePath) == false {
		var errMsg string
		errMsg = "MarkdownPageModule.Compile_Psf: Css File Path " + cssFilePath + " not exist"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}
	//fmt.Println("C")
	//Call pandoc to convert md to html

	var pandocCmd *exec.Cmd
	pandocCmd = exec.Command("pandoc")

	pandocCmd.Args = append(pandocCmd.Args, "-s")
	pandocCmd.Args = append(pandocCmd.Args, "--self-contained")
	pandocCmd.Args = append(pandocCmd.Args, "-c")
	pandocCmd.Args = append(pandocCmd.Args, cssFilePath)
	pandocCmd.Args = append(pandocCmd.Args, markdownSrc)
	pandocCmd.Args = append(pandocCmd.Args, "-o")
	pandocCmd.Args = append(pandocCmd.Args, markdownDst)
	pandocCmd.Args = append(pandocCmd.Args, "--metadata")
	pandocCmd.Args = append(pandocCmd.Args, "pagetitle="+psf.Title)

	var stdoutput bytes.Buffer
	var stderr bytes.Buffer

	pandocCmd.Stdout = &stdoutput
	pandocCmd.Stderr = &stderr

	errPandocCmd := pandocCmd.Run()
	if errPandocCmd != nil {
		Utils.Logger.Println(fmt.Sprint(errPandocCmd) + " : " + stderr.String())
		return -1, errPandocCmd
	}

	var _pofIndex int
	pofIndex := mpmp.spp.GetPageOutputFile(psf.OutputFile)
	if psf.OutputFile == "" || pofIndex == -1 {
		pof := Page.NewPageOutputFile()
		pof.Author = psf.Author
		pof.Description = psf.Description
		pof.FilePath = markdownDst
		pof.IsTop = psf.IsTop
		pof.Title = psf.Title
		pof.TitleImage = psf.TitleImage
		pof.Type = psf.Type
		pof.CreateTime = Utils.CurrentTime()

		_, errAdd := mpmp.spp.AddPageOutputFile(pof)

		if errAdd != nil {
			Utils.DeleteFile(markdownDst) //Add fail,delete the file already copied
			return -1, errAdd
		}

		_pofIndex := mpmp.spp.GetPageOutputFile(pof.ID)

		if _pofIndex == -1 {
			Utils.DeleteFile(markdownDst) //Add fail,delete the file already copied
			var errMsg = "MarkdownPageModule.Compile_Psf: Page Output File add Fail"
			Utils.Logger.Println(errMsg)
			return _pofIndex, errors.New(errMsg)
		}

		psf.OutputFile = pof.ID

	} else {

		pof := mpmp.spp.OutputFiles[pofIndex]
		pof.Author = psf.Author
		pof.Description = psf.Description
		pof.FilePath = markdownDst
		pof.IsTop = psf.IsTop
		pof.Title = psf.Title
		pof.TitleImage = psf.TitleImage
		pof.Type = psf.Type
		pof.CreateTime = Utils.CurrentTime()

		_, errUpdatePof := mpmp.spp.UpdatePageOutputFile(pof)

		if errUpdatePof != nil {
			Utils.DeleteFile(markdownDst) //Add fail,delete the file already copied
			return -1, errUpdatePof
		}

	}
	psf.LastCompiled = Utils.CurrentTime()

	if psf.Type == Page.INDEX {
		if psf.ID == mpmp.spp.IndexPageSourceFile.ID {
			_, errUpdateIndexPage := mpmp.spp.UpdateIndexSourceFile(psf)
			if errUpdateIndexPage != nil {
				var errMsg = "MarkdownPageModule.Compile_Psf: Cannot update Index Source File"
				Utils.Logger.Println(errMsg)
				return _pofIndex, errors.New(errMsg)
			}
		} else {
			_, errUpdateMorePage := mpmp.spp.UpdateMorePageSourceFile(psf)
			if errUpdateMorePage != nil {
				var errMsg = "MarkdownPageModule.Compile_Psf: Cannot update More Source File"
				Utils.Logger.Println(errMsg)
				return _pofIndex, errors.New(errMsg)
			}
		}
	} else {
		mpmp.spp.UpdatePageSourceFile(psf)
	}

	return _pofIndex, nil
}

//Compile Markdown, call pandoc to convert md to html to Output folder
//change sourceinformation and add PageOutputFile
func (mpmp *MarkdownPageModule) Compile(ID string) (int, error) {
	iFind := mpmp.spp.GetPageSourceFile(ID)
	if iFind == -1 {
		var errMsg string
		errMsg = "MarkdownPageModule.Compile: Cannot find the source File with ID " + ID
		Utils.Logger.Println(errMsg)
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
		errMsg = "MarkdownPageModule.CreateIndexPage: Cannot find index page template file for page size " + indexPageSize
		Utils.Logger.Println(errMsg)
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
	mpmp.spp.IndexPageSourceFile.OutputFile = ""

	//Copy index template file to markdown folder
	var srcIndexPageSourceFilePath = mpmp.spp.IndexPageSourceFile.SourceFilePath
	nByte, errCopy := Utils.CopyFile(indexTemplateFilePath, srcIndexPageSourceFilePath)

	if nByte == 0 && errCopy != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.CreateIndexPage: Copy Index Template File from " + indexTemplateFilePath + " to " + srcIndexPageSourceFilePath + " failed, will reset index page properties in site project file"
		Utils.Logger.Println(errMsg)
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = "MarkdownPageModule.CreateIndexPage: Clean Index Page properties failed, please check site project file"
			Utils.Logger.Println(errCleanMsg)
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
			errMsg = "MarkdownPageModule.CreateIndexPage: Cannot get top Page Source File"
			Utils.Logger.Println(errMsg)

			bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

			if errClean != nil {
				var errCleanMsg string
				errCleanMsg = "MarkdownPageModule.CreateIndexPage: Clean Index Page properties failed, please check site project file"
				Utils.Logger.Println(errCleanMsg)
				return bClean, errors.New(errCleanMsg)
			}

			bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "MarkdownPageModule.CreateIndexPgae: Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
				Utils.Logger.Println(deleteMsg)
			}

			return false, errors.New(errMsg)
		}

		normalOutputPageFiles, errNormal := mpmp.spp.GetSortedNormalOutputFiles()

		if errNormal != nil {
			var errMsg string
			errMsg = "MarkdownPageModule.CreateIndexPage: Cannot get normal Page Source File"
			Utils.Logger.Println(errMsg)

			bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

			if errClean != nil {
				var errCleanMsg string
				errCleanMsg = "MarkdownPageModule.CreateIndexPage: Clean Index Page properties failed, please check site project file"
				Utils.Logger.Println(errCleanMsg)
				return bClean, errors.New(errCleanMsg)
			}

			bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "MarkdownPageModule.CreateIndexPage: Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
				Utils.Logger.Println(deleteMsg)
			}

			return false, errors.New(errMsg)
		}

		mpmp.outputPageFiles = append(topOutputPageFiles, normalOutputPageFiles...)
	}
	// get the first pagesize items
	nIndexPageSize, errNIndexPageSize := Page.ConvertPageSize2Int(indexPageSize)

	if errNIndexPageSize != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.CreateIndexPage: Cannot get page size"
		Utils.Logger.Println(errMsg)
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = "MarkdownPageModule.CreateIndexPage: Clean Index Page properties failed, please check site project file"
			Utils.Logger.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}
		bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "MarkdownPageModule.CreateIndexPage: Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
			Utils.Logger.Println(deleteMsg)
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
		errMsg = "MarkdownPageModule.CreateIndexPage: Cannot read src Index md file"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errReadFile.Error())
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = "MarkdownPageModule.CreateIndexPage:  Read file fail,then Clean Index Page properties failed, please check site project file"
			Utils.Logger.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}
		bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "MarkdownPageModule.CreateIndexPage: Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
			Utils.Logger.Println(deleteMsg)
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
			var errMsg = "MarkdownPageModule.CreateIndexPage: Title of Item is empty: " + indexOutputPageFile.ID
			Utils.Logger.Println(errMsg)
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
			var errMsg = "MarkdownPageModule.CreateIndexPage: FilePath of Item is empty: " + indexOutputPageFile.ID
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		if indexOutputPageFile.TitleImage != "" {
			indexFileContent = strings.Replace(indexFileContent, indexNewsImageMark, indexOutputPageFile.TitleImage, 1)
		} else {
			Utils.Logger.Println("MarkdownPageModule.CreateIndexPage: TitleImage of Item is empty: " + indexOutputPageFile.FilePath + " This item will not have title image in index.html")

			emptyImageTemplate, errEmptyImage := Configuration.GetEmptyImageItemTemplate()
			if errEmptyImage != nil {
				var errMsg = "MarkdownPageModule.CreateIndexPage: Cannot get empty image template"
				Utils.Logger.Println(errMsg)
				return false, errors.New(errMsg)
			}
			emptyImageTemplate = strings.Replace(emptyImageTemplate, Page.INDEX_NEWS_IMAGE_MARK, indexNewsImageMark, 1)
			indexFileContent = strings.Replace(indexFileContent, emptyImageTemplate, "", -1)
		}

		if indexOutputPageFile.CreateTime != "" {
			indexFileContent = strings.Replace(indexFileContent, indexNewsTimeMark, indexOutputPageFile.CreateTime, 1)
		} else {
			var errMsg = "MarkdownPageModule.CreateIndexPage: CreateTime of Item is empty " + indexOutputPageFile.ID
			Utils.Logger.Println(errMsg)
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
			errMsg = "MarkdownPageModule.CreateIndexPage: Cannot read empty item template from item"
			Utils.Logger.Println(errMsg)

			bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

			if errClean != nil {
				var errCleanMsg string
				errCleanMsg = "MarkdownPageModule.CreateIndexPage: Read file fail,then Clean Index Page properties failed, please check site project file"
				Utils.Logger.Println(errCleanMsg)
				return bClean, errors.New(errCleanMsg)
			}
			bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "MarkdownPageModule.CreateIndexPage: Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
				Utils.Logger.Println(deleteMsg)
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
		errMsg = "MarkdownPageModule.CreateIndexPage: Cannot Save content to index md file"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errWriteFile.Error())
		bClean, errClean := mpmp.spp.CleanIndexPageSourceFile()

		if errClean != nil {
			var errCleanMsg string
			errCleanMsg = "MarkdownPageModule.CreateIndexPage: Read file fail,then Clean Index Page properties failed, please check site project file"
			Utils.Logger.Println(errCleanMsg)
			return bClean, errors.New(errCleanMsg)
		}
		bDelete := Utils.DeleteFile(srcIndexPageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "MarkdownPageModule.CreateIndexPage: Delete md file failed,please delete it manully, path " + srcIndexPageSourceFilePath
			Utils.Logger.Println(deleteMsg)
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
		errMsg = "MarkdownPageModule.CreateMorePage: Cannot find more page template file for page size " + indexPageSize
		Utils.Logger.Println(errMsg)
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
	morePageSourceFile.OutputFile = ""
	_, errAddMorePage := mpmp.spp.AddMorePageSourceFile(morePageSourceFile)

	if errAddMorePage != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.CreateMorePage: Cannot add More Page Source File"
		Utils.Logger.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete it"
			Utils.Logger.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}

		return false, errors.New(errMsg)
	}

	//Copy more template file to markdown folder
	var srcMorePageSourceFilePath = morePageSourceFile.SourceFilePath
	nByte, errCopy := Utils.CopyFile(moreTemplateFilePath, srcMorePageSourceFilePath)

	if nByte == 0 && errCopy != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.CreateMorePage: Copy More Template File from " + moreTemplateFilePath + " to " + srcMorePageSourceFilePath + " failed, will remove this more page in site project file"
		Utils.Logger.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "MarkdownPageModule.Create More Page: Cannot delete more page properties from site project"
			Utils.Logger.Println(errRemoveMsg)
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
			errMsg = "MarkdownPageModule.CreateMorePage: Cannot get top Output Page Files"
			Utils.Logger.Println(errMsg)

			bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

			if errRemove != nil {
				var errRemoveMsg string
				errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete more page properties from site project"
				Utils.Logger.Println(errRemoveMsg)
				return bRemove, errors.New(errRemoveMsg)
			}
			bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "MarkdownPageModule.CreateMorePage: Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
				Utils.Logger.Println(deleteMsg)
			}
			return false, errors.New(errMsg)
		}

		normalOutputPageFiles, errNormal := mpmp.spp.GetSortedNormalOutputFiles()

		if errNormal != nil {
			var errMsg string
			errMsg = "MarkdownPageModule.CreateMorePage: Cannot get sorted normal output page files"
			Utils.Logger.Println(errMsg)

			bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

			if errRemove != nil {
				var errRemoveMsg string
				errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete more page properties from site project"
				Utils.Logger.Println(errRemoveMsg)
				return bRemove, errors.New(errRemoveMsg)
			}
			bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "MarkdownPageModule.CreateMorePage: Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
				Utils.Logger.Println(deleteMsg)
			}
			return false, errors.New(errMsg)
		}

		mpmp.outputPageFiles = append(topOutputPageFiles, normalOutputPageFiles...)
	}
	// get the first pagesize items
	nMorePageSize, errNMorePageSize := Page.ConvertPageSize2Int(indexPageSize)

	if errNMorePageSize != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.CreateMorePage: Cannot get page size"
		Utils.Logger.Println(errMsg)

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete more page properties from site project"
			Utils.Logger.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}
		bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "MarkdownPageModule.CreateMorePage: Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
			Utils.Logger.Println(deleteMsg)
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
		errMsg = "MarkdownPageModule.CreateMorePage: Cannot read More Page md file"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errReadFile.Error())

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete more page properties from site project"
			Utils.Logger.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}
		bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "MarkdownPageModule.CreateMorePage: Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
			Utils.Logger.Println(deleteMsg)
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
			return false, errors.New("MarkdownPageModule.CreateMorePage: Title of Item is empty" + moreOutputPageFile.ID)
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
			var errMsg = "MarkdownPageModule.CreateMorePage: FilePath of Item is empty: " + moreOutputPageFile.ID
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		if moreOutputPageFile.TitleImage != "" {
			moreFileContent = strings.Replace(moreFileContent, moreNewsImageMark, moreOutputPageFile.TitleImage, 1)
		} else {
			Utils.Logger.Println("MarkdownPageModule.CreateMorePage: TitleImage of Item is empty: " + moreOutputPageFile.FilePath + " This item will not have title image in moreXX.html")

			emptyImageTemplate, errEmptyImage := Configuration.GetEmptyImageItemTemplate()
			if errEmptyImage != nil {
				var errMsg = "MarkdownPageModule.CreateMorePage: Cannot get empty image template"
				Utils.Logger.Println(errMsg)
				return false, errors.New(errMsg)
			}
			emptyImageTemplate = strings.Replace(emptyImageTemplate, Page.INDEX_NEWS_IMAGE_MARK, moreNewsImageMark, 1)
			moreFileContent = strings.Replace(moreFileContent, emptyImageTemplate, "", -1)
		}

		if moreOutputPageFile.CreateTime != "" {
			moreFileContent = strings.Replace(moreFileContent, moreNewsTimeMark, moreOutputPageFile.CreateTime, 1)
		} else {
			var errMsg = "MarkdownPageModule.CreateMorePage: CreateTime of Item is empty: " + moreOutputPageFile.ID
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}

	if startIndex+nMorePageSize > len(mpmp.outputPageFiles) {
		//read empty item template
		emptyItemTemplate, errEmptyItemTemplate := Configuration.GetEmptyIndexItemTemplate()

		if errEmptyItemTemplate != nil {
			var errMsg string
			errMsg = "MarkdownPageModule.CreateMorePage: Cannot read empty Item template file"
			Utils.Logger.Println(errMsg)

			bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

			if errRemove != nil {
				var errRemoveMsg string
				errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete more page properties from site project"
				Utils.Logger.Println(errRemoveMsg)
				return bRemove, errors.New(errRemoveMsg)
			}
			bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

			if bDelete == false {
				var deleteMsg = "MarkdownPageModule.CreateMorePage: Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
				Utils.Logger.Println(deleteMsg)
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
		errMsg = "MarkdownPageModule.CreateMorePage: Cannot Save modified md file"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errWriteFile.Error())

		bRemove, errRemove := mpmp.spp.RemoveMorePageSourceFile(morePageSourceFile)

		if errRemove != nil {
			var errRemoveMsg string
			errRemoveMsg = "MarkdownPageModule.CreateMorePage: Cannot delete more page properties from site project"
			Utils.Logger.Println(errRemoveMsg)
			return bRemove, errors.New(errRemoveMsg)
		}
		bDelete := Utils.DeleteFile(srcMorePageSourceFilePath)

		if bDelete == false {
			var deleteMsg = "MarkdownPageModule.CreateMorePage: Delete md file failed,please delete it manully, path " + srcMorePageSourceFilePath
			Utils.Logger.Println(deleteMsg)
		}
		return false, errors.New(errMsg)
	}

	return true, nil
}

func (mpmp *MarkdownPageModule) CreateNavigationForIndexPage() (bool, error) {
	//Get template file path for index page

	if mpmp.spp.IndexPageSourceFile.SourceFilePath == "" {
		var errMsg = "MarkdownPageModule.CreateNavigationForIndexPage: Index Page file not created,please run CreateIndexPage firstly"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	srcIndexPageSourceFilePath := mpmp.spp.IndexPageSourceFile.SourceFilePath
	// modify copied index page  template md
	// Read file
	bIndexFileContent, errReadFile := ioutil.ReadFile(srcIndexPageSourceFilePath)

	if errReadFile != nil {
		var errMsg string
		errMsg = "MarkdownPageModule.CreateNavigationForIndexPage: Cannot read src Index md file"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errReadFile.Error())
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
		errMsg = "MarkdownPageModule.CreateNavigationForIndexPage: Cannot Save content to index md file"
		Utils.Logger.Println(errMsg)
		Utils.Logger.Println(errWriteFile.Error())
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
			errMsg = "MarkdownPageModule.CreateNavigationForIndexPage: Cannot read More Page md file " + srcMorePageSourceFilePath
			Utils.Logger.Println(errMsg)
			Utils.Logger.Println(errReadFile.Error())
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
			errMsg = "MarkdownPageModule.CreateNavigationForIndexPage: Cannot Save modified md file"
			Utils.Logger.Println(errMsg)
			Utils.Logger.Println(errWriteFile.Error())
			return false, errors.New(errMsg)
		}
	}

	return true, nil
}
