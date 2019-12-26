package Site

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ipsc/Configuration"
	"ipsc/Page"
	"ipsc/Utils"

	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SiteModule struct {
	spp *SiteProject
	mpp *MarkdownPageModule
	hpp *HtmlPageModule
	lp  *LinkModule

	projectFolderPath string
}

func NewSiteModule() *SiteModule {
	var sm SiteModule
	var smp *SiteModule
	smp = &sm

	_spp := NewSiteProject()
	smp.spp = _spp

	var mpm = NewMarkdownPageModule(smp.spp, smp)
	smp.mpp = &mpm

	var hpm = NewHtmlPageModule(smp.spp, smp)
	smp.hpp = &hpm

	var lm = NewLinkModule(smp.spp)
	smp.lp = &lm

	return smp
}

func NewSiteModule_WithArgs(_projectFolderPath, _projectFileName string) *SiteModule {
	var sm SiteModule
	var smp *SiteModule
	smp = &sm
	smp.projectFolderPath = _projectFolderPath
	//fmt.Println("NewSMPointA")
	_, errOpen := smp.OpenSiteProject(_projectFolderPath, _projectFileName)
	//fmt.Println("NewSMPointB")
	if errOpen != nil {
		Utils.Logger.Println("SiteModule.NewSiteModule: Cannot create Site Module")
		return nil
	}

	var mpm = NewMarkdownPageModule(smp.spp, smp)
	smp.mpp = &mpm

	var hpm = NewHtmlPageModule(smp.spp, smp)
	smp.hpp = &hpm

	var lm = NewLinkModule(smp.spp)
	smp.lp = &lm

	return smp
}

func (smp *SiteModule) GetProjectFolderPath() string {
	return smp.projectFolderPath
}

func (smp *SiteModule) GetSrcFolderPath(projectFolderPath string) string {
	return filepath.Join(projectFolderPath, "Src")
}

func (smp *SiteModule) GetSrcMarkdownFolderPath(projectFolderPath string) string {
	return filepath.Join(smp.GetSrcFolderPath(projectFolderPath), "Markdown")
}

func (smp *SiteModule) GetSrcHtmlFolderPath(projectFolderPath string) string {
	return filepath.Join(smp.GetSrcFolderPath(projectFolderPath), "Html")
}

func (smp *SiteModule) GetOutputFolderPath(projectFolderPath string) string {
	return filepath.Join(projectFolderPath, "Output")
}

func (smp *SiteModule) GetOutputPagesFolderPath(projectFolderPath string) string {
	return filepath.Join(smp.GetOutputFolderPath(projectFolderPath), "Pages")
}

func (smp *SiteModule) GetOutputFilesFolderPath(projectFolderPath string) string {
	return filepath.Join(smp.GetOutputPagesFolderPath(projectFolderPath), "Files")
}

func (smp *SiteModule) GetTemplateFolderPath(projectFolderPath string) string {
	return filepath.Join(projectFolderPath, "Templates")
}

func (smp *SiteModule) GetSrcFilesFolderPath(projectFolderPath string) string {
	return filepath.Join(smp.GetSrcFolderPath(projectFolderPath), "Files")
}

func (smp *SiteModule) GetSiteProjectFilePath(projectFolderPath string) (string, error) {
	if nil != smp.spp && smp.spp.Title == "" {
		var errMsg = "SiteModule.GetSiteProjectFilePath: SiteProject Title is empty"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	return filepath.Join(projectFolderPath, smp.spp.Title) + ".sp", nil
}

func (smp *SiteModule) PathIsSiteProject(projectPath, projectName string) (bool, error) {
	if Utils.PathIsExist(projectPath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: " + projectPath + " is not exist"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var projectFilePath = filepath.Join(projectPath, projectName)

	if strings.HasSuffix(projectFilePath, ".sp") == false {

		projectFilePath += ".sp"
	}

	if Utils.PathIsExist(projectFilePath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find sp file in project " + projectPath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var sp SiteProject
	_, loadError := sp.LoadFromFile(projectFilePath)
	if loadError != nil {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot load sp file in project " + projectPath
		Utils.Logger.Println(errMsg)
		return false, errors.New("SiteModule.PathIsSiteProject: Cannot load sp file in project " + projectPath)
	}

	var srcFolderPath = smp.GetSrcFolderPath(projectPath)

	if Utils.PathIsExist(srcFolderPath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Src folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var markdownFolderPath = smp.GetSrcMarkdownFolderPath(projectPath)

	if Utils.PathIsExist(markdownFolderPath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Markdown folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var htmlFolderPath = smp.GetSrcHtmlFolderPath(projectPath)

	if Utils.PathIsExist(htmlFolderPath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Html folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var filesFolderPath = smp.GetSrcFilesFolderPath(projectPath)
	if Utils.PathIsExist(filesFolderPath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Files folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var outputFolderPath = smp.GetOutputFolderPath(projectPath)

	if Utils.PathIsExist(outputFolderPath) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Output folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var outputPagesSubFolder = smp.GetOutputPagesFolderPath(projectPath)

	if Utils.PathIsExist(outputPagesSubFolder) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Output/Pages folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)

	}

	var outputFilesSubFolder = smp.GetOutputFilesFolderPath(projectPath)

	if Utils.PathIsExist(outputFilesSubFolder) == false {
		var errMsg = "SiteModule.PathIsSiteProject: Cannot find Output/Pages/Files folder"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	return true, nil
}

func (smp *SiteModule) InitializeSiteProjectFolder(siteTitle, siteAuthor, siteDescription, _projectFolderPath, _outputFolderPath string) (bool, error) {
	if _projectFolderPath == "" {
		var errMsg = "SiteModule.InitializeSiteProjectFolder: Project Folder Path is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	//Create each foldrs

	//ProjectFolder
	var errProjectFolder error
	if !Utils.PathIsExist(_projectFolderPath) {
		_, errProjectFolder = Utils.MakeFolder(_projectFolderPath)

		if errProjectFolder != nil {
			Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errProjectFolder.Error())
			return false, errProjectFolder
		}

	}

	//ProjectFolder->Src
	var srcFolderPath = smp.GetSrcFolderPath(_projectFolderPath)
	var errSrcFolderPath error
	if !Utils.PathIsExist(srcFolderPath) {
		_, errSrcFolderPath = Utils.MakeFolder(srcFolderPath)

		if errSrcFolderPath != nil {
			Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errSrcFolderPath.Error())
			return false, errSrcFolderPath
		}
	}

	//ProjectFolder->Src->Markdown
	var srcMarkdownFolderPath = smp.GetSrcMarkdownFolderPath(_projectFolderPath)
	var errSrcMarkdownFoldrPath error
	if !Utils.PathIsExist(srcMarkdownFolderPath) {
		_, errSrcMarkdownFoldrPath = Utils.MakeFolder(srcMarkdownFolderPath)

		if errSrcMarkdownFoldrPath != nil {
			Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errSrcMarkdownFoldrPath.Error())
			return false, errSrcMarkdownFoldrPath
		}
	}

	//ProjectFolder->Src->Html
	var srcHtmlFolderPath = smp.GetSrcHtmlFolderPath(_projectFolderPath)
	var errSrcHtmlFolderPath error

	if !Utils.PathIsExist(srcHtmlFolderPath) {
		_, errSrcHtmlFolderPath = Utils.MakeFolder(srcHtmlFolderPath)

		if errSrcHtmlFolderPath != nil {
			Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errSrcHtmlFolderPath.Error())
			return false, errSrcHtmlFolderPath
		}
	}

	//ProjectFolder->Src->Files
	var srcFilesFolderPath = smp.GetSrcFilesFolderPath(_projectFolderPath)
	var errSrcFilesFolderPath error

	if !Utils.PathIsExist(srcFilesFolderPath) {
		_, errSrcFilesFolderPath = Utils.MakeFolder(srcFilesFolderPath)

		if errSrcFilesFolderPath != nil {
			Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errSrcFilesFolderPath.Error())
			return false, errSrcFilesFolderPath
		}
	}

	//ProjectFolder->Output
	var outputFolderPath = smp.GetOutputFolderPath(_projectFolderPath)
	var errOutputFolderPath error

	if outputFolderPath == _outputFolderPath || _outputFolderPath == "" {

		if !Utils.PathIsExist(outputFolderPath) {
			_, errOutputFolderPath = Utils.MakeFolder(outputFolderPath)

			if errOutputFolderPath != nil {
				Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errOutputFolderPath.Error())
				return false, errOutputFolderPath
			}

			//ProjectFolder->Output->Pages
			var outputPagesSubFolder = smp.GetOutputPagesFolderPath(_projectFolderPath)
			_, errOutputPagesFolder := Utils.MakeFolder(outputPagesSubFolder)

			if errOutputPagesFolder != nil {
				Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errOutputPagesFolder.Error())
				return false, errOutputPagesFolder
			}

			//ProjectFolder->Output->Pages->Files
			var outputFilesSubFolder = smp.GetOutputFilesFolderPath(_projectFolderPath)
			if Utils.PathIsExist(outputFilesSubFolder) == false {
				_, errOutputFilesFolder := Utils.MakeFolder(outputFilesSubFolder)

				if errOutputFilesFolder != nil {
					Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errOutputFilesFolder.Error())
					return false, errOutputPagesFolder
				}
			}

		}
	} else {
		if !Utils.PathIsExist(outputFolderPath) {
			_, errOutputFolderPath = Utils.MakeSoftLink4Folder(_outputFolderPath, outputFolderPath)

			if errOutputFolderPath != nil {
				Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errOutputFolderPath.Error())
				return false, errOutputFolderPath
			}

			//ProjectFolder->Output->Pages
			var outputPagesSubFolder = smp.GetOutputPagesFolderPath(_projectFolderPath)

			if Utils.PathIsExist(outputPagesSubFolder) == false {
				_, errOutputPagesFolder := Utils.MakeFolder(outputPagesSubFolder)

				if errOutputPagesFolder != nil {
					Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errOutputPagesFolder.Error())
					return false, errOutputPagesFolder
				}
			}

			//ProjectFolder->Output->Pages->Files
			var outputFilesSubFolder = smp.GetOutputFilesFolderPath(_projectFolderPath)
			if Utils.PathIsExist(outputFilesSubFolder) == false {
				_, errOutputFilesFolder := Utils.MakeFolder(outputFilesSubFolder)

				if errOutputFilesFolder != nil {
					Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errOutputFilesFolder.Error())
					return false, errOutputFilesFolder
				}
			}
		}
	}

	//Create Templates Path and copy templates from IPSC Resources folder
	//Project Folder->Templates
	var templateFolderPath = smp.GetTemplateFolderPath(_projectFolderPath)
	var errTemplateFolder error
	if !Utils.PathIsExist(templateFolderPath) {
		_, errTemplateFolder = Utils.MakeFolder(templateFolderPath)

		if errTemplateFolder != nil {
			Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errTemplateFolder.Error())
			return false, errTemplateFolder
		}
	}
	//Copy temlates from Resources
	srcTemplateFolder, errSrcTemplate := Configuration.GetTemplatesFolderPath()
	if errSrcTemplate != nil {
		Utils.Logger.Println("SiteModule.InitializeSiteProjectFolder: " + errSrcTemplate.Error())
		return false, errSrcTemplate
	}

	if Utils.PathIsExist(srcTemplateFolder) == false {
		var errMsg = "SiteModule.InitializeSiteProjectFolder: Try to copy tempaltes, src tempalte folder not exist " + srcTemplateFolder
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	files, _ := ioutil.ReadDir(srcTemplateFolder)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			srcTemplateFilePath := filepath.Join(srcTemplateFolder, f.Name())
			dstTemplateFilePath := filepath.Join(templateFolderPath, f.Name())

			_, errCopy := Utils.CopyFile(srcTemplateFilePath, dstTemplateFilePath)
			if errCopy != nil {
				var errMsg = "SiteModule.InitializeSiteProjectFolder: Cannot copy template file " + srcTemplateFilePath + " to " + dstTemplateFilePath
				Utils.Logger.Println(errMsg)
				return false, errors.New(errMsg)
			}
		}
	}
	//create empty project file

	var spp = smp.GetSiteProject()
	spp.Title = siteTitle
	spp.Author = siteAuthor
	spp.Description = siteDescription
	spp.OutputFolderPath = outputFolderPath
	spp.LastModified = Utils.CurrentTime()

	projectFilePath, errProjectFilePath := smp.GetSiteProjectFilePath(_projectFolderPath)

	if errProjectFilePath != nil {
		Utils.Logger.Println(errProjectFilePath.Error())
		return false, errProjectFilePath
	}

	bSaveToFile, errSaveToFile := smp.spp.SaveToFile(projectFilePath)

	if bSaveToFile == false || errSaveToFile != nil {
		Utils.Logger.Println(errSaveToFile.Error())
		return false, errSaveToFile
	}

	return true, nil
}

func (smp *SiteModule) OpenSiteProject(projectFolderPath, projectName string) (bool, error) {

	if projectFolderPath == "" {
		Utils.Logger.Println("SiteModule.OpenSiteProject: Project Folder path is empty")
		return false, errors.New("SiteModule.OpenSiteProject: Project Folder path is empty")
	}
	//fmt.Println("OpenSPPointA")
	bIsSP, errIsSP := smp.PathIsSiteProject(projectFolderPath, projectName)

	if errIsSP != nil || false == bIsSP {
		var errMsg = "SiteModule.OpenSiteProject: Path " + projectFolderPath + " doesn't contain a IPSC Site"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	//fmt.Println("OpenSPPointB")
	var siteProjectFilePath = filepath.Join(projectFolderPath, projectName)
	if strings.HasSuffix(siteProjectFilePath, ".sp") == false {
		siteProjectFilePath += ".sp"
	}
	//fmt.Println("OpenSPPointC")

	var sp SiteProject
	_, loadError := sp.LoadFromFile(siteProjectFilePath)
	if loadError != nil {
		var errMsg = "SiteModule.OpenSiteProject: Cannot load sp file in project " + siteProjectFilePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.spp = &sp

	return true, nil
}

func (smp *SiteModule) GetSiteInformation() (string, error) {
	return smp.spp.ToJson()
}

func (smp *SiteModule) GetSiteProject() *SiteProject {
	return smp.spp
}

func (smp *SiteModule) UpdateSiteProject(siteFolder, siteTitle, siteAuthor, siteDescription string) (bool, error) {
	var oldSiteProjectFilePath = filepath.Join(siteFolder, smp.spp.Title+".sp")
	var siteProjectFilePath = oldSiteProjectFilePath
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "SiteModule.UpdateSiteProject: siteFolder " + siteFolder + " doesn't exist"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if smp.spp == nil {
		var errMsg = "SiteModule.UpdateSiteProject: Site Project is nil"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if siteTitle != "" && smp.spp.Title != siteTitle {
		if Utils.PathIsExist(oldSiteProjectFilePath) {
			var newSiteProjectFilePath = filepath.Join(siteFolder, siteTitle+".sp")

			_, errMove := Utils.MoveFile(oldSiteProjectFilePath, newSiteProjectFilePath)
			if errMove != nil {
				Utils.Logger.Println(errMove.Error())
				return false, errMove
			}
			siteProjectFilePath = newSiteProjectFilePath
		}
		smp.spp.Title = siteTitle
		smp.spp.LastModified = Utils.CurrentTime()
	}

	if siteAuthor != "" && smp.spp.Author != siteAuthor {
		smp.spp.Author = siteAuthor
		smp.spp.LastModified = Utils.CurrentTime()
	}

	if siteDescription != "" && smp.spp.Description != siteDescription {
		smp.spp.Description = siteDescription
		smp.spp.LastModified = Utils.CurrentTime()
	}

	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		Utils.Logger.Println(errSave.Error())
		return bSave, errSave
	}

	return true, nil
}

func (smp *SiteModule) GetAllPages() []string {
	var allpages, active, recycled, outputs []string

	active = smp.spp.GetActivePageSources()

	recycled = smp.spp.GetRecycledPageSources()

	outputs = smp.spp.GetAllPageOutputs()

	allpages = append(allpages, strconv.Itoa(len(active)))
	allpages = append(allpages, strconv.Itoa(len(recycled)))
	allpages = append(allpages, strconv.Itoa(len(outputs)))

	allpages = append(allpages, active...)
	allpages = append(allpages, recycled...)
	allpages = append(allpages, outputs...)

	return allpages
}

func (smp *SiteModule) GetAllRecycledPageSourceFiles() []string {
	return smp.spp.GetRecycledPageSources()
}

func (smp *SiteModule) RestoreRecycledPageSourceFile(ID string) (bool, error) {
	if ID == "" {
		var errMsg = "SiteModule.RestoreRecycledPageSourceFile: RestoreRecycledPageSourceFile: " + "ID is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	index := smp.spp.GetPageSourceFile(ID)

	if index == -1 {
		var errMsg = "SiteModule.RestoreRecycledPageSourceFile: Cannot find Page Source File with ID " + ID
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	bResotre, errRestore := smp.spp.ResotrePageSourceFile(ID)
	if errRestore != nil {
		return bResotre, errRestore
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
		var errMsg = "SiteModule.RestoreRecycledPageSourceFile: Cannot got site project file path "
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		return false, errSave
	}

	return true, nil
}

func (smp *SiteModule) CleanRecycledPageSourceFiles() (bool, error) {
	var deleteSlice []Page.PageSourceFile
	for _, psf := range smp.spp.SourceFiles {
		if psf.Status == Page.RECYCLED {
			deleteSlice = append(deleteSlice, psf)
		}
	}

	for _, delPsf := range deleteSlice {
		if delPsf.Type == Page.MARKDOWN {
			bM, errM := smp.mpp.RemoveMarkdown(delPsf, false)
			if errM != nil {
				return bM, errM
			}
		} else if delPsf.Type == Page.HTML {
			bH, errH := smp.hpp.RemoveHtml(delPsf, false)
			if errH != nil {
				return bH, errH
			}
		} else if delPsf.Type == Page.LINK {
			bL, errL := smp.lp.RemoveLink(delPsf, false)
			if errL != nil {
				return bL, errL
			}
		}
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
		var errMsg = "SiteModule.RestoreRecycledPageSourceFile: Cannot got site project file path "
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		return false, errSave
	}

	return true, nil
}

func (smp *SiteModule) Compile(indexPageSize string) (bool, error) {
	//Remove old index and more output file from spp.outputFiles
	var deletedOutputIndexs []Page.PageOutputFile
	for _, oldIndexOutput := range smp.spp.OutputFiles {
		if oldIndexOutput.Type == Page.INDEX {
			deletedOutputIndexs = append(deletedOutputIndexs, oldIndexOutput)
		}
	}

	for _, delPof := range deletedOutputIndexs {
		bDelOldIndex, errDelOldIndex := smp.spp.RemovePageOutputFile(delPof)

		if errDelOldIndex != nil {
			return bDelOldIndex, errDelOldIndex
		}
	}

	var mdCount, htmlCount, linkCount int
	mdCount = 0
	htmlCount = 0
	linkCount = 0
	//fmt.Println("A")
	for _, sp := range smp.spp.SourceFiles {
		if sp.Status == Page.ACTIVE {
			//Never Compiled
			// OR
			//Compiled, but source file changed
			if (sp.LastCompiled == "" && sp.OutputFile == "") || (sp.OutputFile != "" && sp.LastCompiled != "" && sp.LastModified != "" && sp.LastCompiled < sp.LastModified) {
				if sp.Type == Page.MARKDOWN {
					_, errCompileMd := smp.mpp.Compile(sp.ID)
					if errCompileMd != nil {
						Utils.Logger.Println(errCompileMd.Error())
						return false, errCompileMd
					}
					mdCount++
				} else if sp.Type == Page.HTML {
					_, errCompileHtml := smp.hpp.Compile(sp.ID)
					if errCompileHtml != nil {
						Utils.Logger.Println(errCompileHtml.Error())
						return false, errCompileHtml
					}
					htmlCount++
				} else if sp.Type == Page.LINK {
					_, errCompileLink := smp.lp.Compile(sp.ID)
					if errCompileLink != nil {
						Utils.Logger.Println(errCompileLink.Error())
						return false, errCompileLink
					}
					linkCount++
				}
			}
		}
	}
	//Create Index Page
	bIndex, errIndex := smp.mpp.CreateIndexPage(indexPageSize)
	//fmt.Println("B")
	if errIndex != nil {
		return bIndex, errIndex
	}

	var nIndexPageSize, _ = Page.ConvertPageSize2Int(indexPageSize)
	var nOutputFileLength = len(smp.spp.OutputFiles)

	var moreCount int
	moreCount = 0
	//Create more pages when the count of output files is bigger than index page size
	if nIndexPageSize < nOutputFileLength {
		//Delete More Pages created last Compile
		var deletedSourceIndexs []Page.PageSourceFile
		for _, oldIndexSource := range smp.spp.MorePageSourceFiles {
			if oldIndexSource.Type == Page.INDEX {
				deletedSourceIndexs = append(deletedSourceIndexs, oldIndexSource)
			}
		}

		for _, delPsf := range deletedSourceIndexs {
			bDelOldIndex, errDelOldIndex := smp.spp.RemoveMorePageSourceFile(delPsf)

			if errDelOldIndex != nil {
				return bDelOldIndex, errDelOldIndex
			}
		}

		//Create More Pages
		var moreOutputFileLength = nOutputFileLength - nIndexPageSize
		var moreOutputPageCount = moreOutputFileLength / nIndexPageSize
		var temp = moreOutputFileLength % nIndexPageSize
		if temp > 0 {
			moreOutputPageCount = moreOutputPageCount + 1
		}

		for index := 1; index <= moreOutputPageCount; index++ {
			var startIndex = index * nIndexPageSize
			bMore, errMore := smp.mpp.CreateMorePage(indexPageSize, startIndex, index)
			if errMore != nil {
				return bMore, errMore
			}
			moreCount++
		}
	}
	//Create Navigation of index page and more pages
	bNavigationIndex, errNavigationIndex := smp.mpp.CreateNavigationForIndexPage()
	//fmt.Println("C")
	if errNavigationIndex != nil {
		return bNavigationIndex, errNavigationIndex
	}

	bNavigationMore, errNavigationMore := smp.mpp.CreateNavigationForMorePages()

	if errNavigationMore != nil {
		return bNavigationMore, errNavigationMore
	}

	//fmt.Println("D")
	//Compile Index Page and More Pages
	_, errCompileIndex := smp.mpp.Compile_Psf(smp.spp.IndexPageSourceFile)

	if errCompileIndex != nil {
		return false, errCompileIndex
	}

	for _, morePsf := range smp.spp.MorePageSourceFiles {
		_, errCompileMore := smp.mpp.Compile_Psf(morePsf)
		if errCompileMore != nil {
			return false, errCompileMore
		}
	}
	//fmt.Println("E")

	fileCount := smp.CompileNormalFile()

	//Get Summary and write to spp
	var CompileSummary string
	CompileSummary = "Index: 1"
	CompileSummary += "_More: " + strconv.Itoa(moreCount)
	CompileSummary += "_Markdown: " + strconv.Itoa(mdCount)
	CompileSummary += "_Html: " + strconv.Itoa(htmlCount)
	CompileSummary += "_Link: " + strconv.Itoa(linkCount)
	CompileSummary += "_File: " + strconv.Itoa(fileCount)

	smp.spp.LastCompileSummary = CompileSummary
	smp.spp.LastCompiled = Utils.CurrentTime()

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
		var errMsg = "SiteModule.Compile: Cannot got site project file path "
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		var errMsg = "SiteModule.Compile: Cannot save site project file "
		Utils.Logger.Println(errMsg)
		return bSave, errors.New(errMsg)
	}
	return true, nil
}

func (smp *SiteModule) CompileNormalFile() int {
	var fileCount int
	fileCount = 0
	srcFileList := smp.GetSrcNormalFileList()
	outputFileList := smp.GetOutputNormalFileList()

	if len(srcFileList) == 0 {
		return 0
	}

	//Add or Update File
	for _, srcFile := range srcFileList {
		var iFind = Page.GetNormalFile(outputFileList, srcFile.FilePath)
		if iFind == -1 {
			//New File,add
			var srcFolderPath = smp.GetSrcFolderPath(smp.projectFolderPath)
			var srcFullPath = srcFolderPath + srcFile.FilePath

			var dstFolderPath = smp.GetOutputPagesFolderPath(smp.projectFolderPath)
			var dstFullPath = dstFolderPath + srcFile.FilePath

			if Utils.PathIsDir(srcFullPath) {
				if Utils.PathIsExist(dstFullPath) == false {
					Utils.MakeFolder(dstFullPath)
				}
			} else if Utils.PathIsFile(srcFullPath) {

				Utils.CopyFile(srcFullPath, dstFullPath)
			}

			fileCount = fileCount + 1
		} else {
			var dstFile = outputFileList[iFind]

			if srcFile.LastModified > dstFile.LastModified {
				var srcFolderPath = smp.GetSrcFolderPath(smp.projectFolderPath)
				var srcFullPath = srcFolderPath + srcFile.FilePath

				var dstFolderPath = smp.GetOutputPagesFolderPath(smp.projectFolderPath)
				var dstFullPath = dstFolderPath + srcFile.FilePath

				if Utils.PathIsDir(srcFullPath) {
					if Utils.PathIsExist(dstFullPath) == false {
						Utils.MakeFolder(dstFullPath)
					}
				} else if Utils.PathIsFile(srcFullPath) {

					Utils.CopyFile(srcFullPath, dstFullPath)
				}

				fileCount = fileCount + 1
			}
		}
	}

	//Delete

	for _, dstFile := range outputFileList {
		var iFind = Page.GetNormalFile(srcFileList, dstFile.FilePath)
		if iFind == -1 {
			var dstFolderPath = smp.GetOutputPagesFolderPath(smp.projectFolderPath)
			var dstFullPath = dstFolderPath + dstFile.FilePath

			if Utils.PathIsFile(dstFullPath) {
				Utils.DeleteFile(dstFullPath)
			}

			fileCount = fileCount + 1
		}
	}

	//Clearup empty sub folder of output/Pages/Files
	smp.clearUpEmptyFolderUnderOutputFiles()
	return fileCount
}

func (smp *SiteModule) clearUpEmptyFolderUnderOutputFiles() {
	var outputFilesFolder = smp.GetOutputFilesFolderPath(smp.projectFolderPath)

	files, _ := ioutil.ReadDir(outputFilesFolder)

	for _, f := range files {
		if f.IsDir() {
			var subFolderPath = filepath.Join(outputFilesFolder, f.Name())
			subFiles, _ := ioutil.ReadDir(subFolderPath)
			if len(subFiles) == 0 {
				Utils.DeleteFile(subFolderPath)
			} else {
				smp.clearUpEmptyFolder(subFolderPath)
			}
		}
	}

}

func (smp *SiteModule) clearUpEmptyFolder(folderPath string) {
	files, _ := ioutil.ReadDir(folderPath)

	for _, f := range files {
		if f.IsDir() {
			var subFolderPath = filepath.Join(folderPath, f.Name())
			subFiles, _ := ioutil.ReadDir(subFolderPath)
			if len(subFiles) == 0 {
				Utils.DeleteFile(subFolderPath)
			} else {
				smp.clearUpEmptyFolder(subFolderPath)
			}
		}
	}
}

func (smp *SiteModule) GetSrcNormalFileList() []Page.NormalFile {
	var srcFolderPath = smp.GetSrcFolderPath(smp.projectFolderPath)
	var srcFilesFolder = smp.GetSrcFilesFolderPath(smp.projectFolderPath)

	var filesList []Page.NormalFile

	filepath.Walk(srcFilesFolder, func(path string, info os.FileInfo, err error) error {
		var fileName = info.Name()
		var relativePath = path[len(srcFolderPath):]
		var lastModified = info.ModTime().Format("2006-01-02 15:04:05")

		if fileName != "Files" {
			var normalFile Page.NormalFile
			normalFile.FileName = fileName
			normalFile.FilePath = relativePath
			normalFile.LastModified = lastModified

			filesList = append(filesList, normalFile)
		}

		return nil
	})

	return filesList
}

func (smp *SiteModule) GetOutputNormalFileList() []Page.NormalFile {
	var outputFolderPath = smp.GetOutputPagesFolderPath(smp.projectFolderPath)
	var outputFilesFolder = smp.GetOutputFilesFolderPath(smp.projectFolderPath)

	var filesList []Page.NormalFile

	filepath.Walk(outputFilesFolder, func(path string, info os.FileInfo, err error) error {
		var fileName = info.Name()
		var relativePath = path[len(outputFolderPath):]
		var lastModified = info.ModTime().Format("2006-01-02 15:04:05")

		if fileName != "Files" {

			var normalFile Page.NormalFile
			normalFile.FileName = fileName
			normalFile.FilePath = relativePath
			normalFile.LastModified = lastModified

			filesList = append(filesList, normalFile)
		}
		return nil
	})

	return filesList
}

func (smp *SiteModule) AddPage(title, description, author, filePath, titleImagePath, pageType string, isTop bool) (bool, string, error) {
	var bAdd bool
	var ID string
	var errAdd error

	pageType = strings.ToUpper(pageType)
	if pageType == Page.MARKDOWN {
		bAdd, ID, errAdd = smp.mpp.AddMarkdown(title, description, author, filePath, titleImagePath, isTop)
	} else if pageType == Page.HTML {
		bAdd, ID, errAdd = smp.hpp.AddHtml(title, description, author, filePath, titleImagePath, isTop)
	} else if pageType == Page.LINK {
		bAdd, ID, errAdd = smp.lp.AddLink(title, description, author, filePath, titleImagePath, isTop)
	}

	if errAdd != nil {
		Utils.Logger.Println(errAdd.Error())
		return bAdd, "-1", errAdd
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
		var errMsg = "SiteModule.AddPage: Cannot got site project file path "
		Utils.Logger.Println(errMsg)
		return false, "-1", errors.New(errMsg)
	}

	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		var errMsg = "SiteModule.AddPage: Cannot save site project file "
		Utils.Logger.Println(errMsg)
		return bSave, "-1", errors.New(errMsg)
	}
	return true, ID, nil
}

func (smp *SiteModule) CreateMarkdown(projectFolderPath, pageFilePath, markdownType string) (bool, error) {
	var templateFolderPath = smp.GetTemplateFolderPath(projectFolderPath)

	return smp.mpp.CreateMarkdown(pageFilePath, markdownType, templateFolderPath)
}

func (smp *SiteModule) UpdatePage(pageID, title, description, author, filePath, titleImagePath string, isTop bool) (bool, error) {

	var index = smp.spp.GetPageSourceFile(pageID)

	if index == -1 {
		var errMsg = "SiteModule.UpdatePage: Cannot find page with ID " + pageID
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var psf = smp.spp.SourceFiles[index]
	var bFile bool
	var errFile error

	pageType := strings.ToUpper(psf.Type)

	if filePath != "" {
		switch pageType {
		case Page.MARKDOWN:
			bFile, errFile = FileIsMarkdown(filePath)
		case Page.HTML:
			bFile, errFile = FileIsHtml(filePath)

		}
		if errFile != nil {
			return bFile, errFile
		}
	} else {
		filePath = psf.SourceFilePath
	}

	if title != "" {
		psf.Title = title
	}

	if author != "" {
		psf.Author = author
	}

	if description != "" {
		psf.Description = description
	}

	if Utils.PathIsExist(titleImagePath) && Utils.PathIsImage(titleImagePath) {
		titleImage, errImage := Utils.ReadImageAsBase64(titleImagePath)
		if errImage == nil {
			psf.TitleImage = titleImage
		}
	}

	if psf.IsTop != isTop {
		psf.IsTop = isTop
	}

	psf.LastModified = Utils.CurrentTime()

	var bUpdate bool
	var errUpdate error

	switch pageType {
	case Page.MARKDOWN:
		bUpdate, errUpdate = smp.mpp.UpdateMarkdown(psf, filePath)
	case Page.HTML:
		bUpdate, errUpdate = smp.hpp.UpdateHtml(psf, filePath)
	case Page.LINK:
		psf.SourceFilePath = filePath
		bUpdate, errUpdate = smp.lp.UpdateLink(psf)
	}

	if errUpdate != nil {
		return bUpdate, errUpdate
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
		var errMsg = "Cannot got site project file path "
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		return false, errSave
	}

	return true, nil
}

func (smp *SiteModule) DeletePage(pageID string, restore bool) (bool, error) {

	var index = smp.spp.GetPageSourceFile(pageID)

	if index == -1 {
		var errMsg = "SiteModule.DeletePage: Cannot find page with ID " + pageID
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var psf = smp.spp.SourceFiles[index]

	var bDelete bool
	var errDelete error

	switch psf.Type {
	case Page.MARKDOWN:
		bDelete, errDelete = smp.mpp.RemoveMarkdown(psf, restore)
	case Page.HTML:
		bDelete, errDelete = smp.hpp.RemoveHtml(psf, restore)
	case Page.LINK:
		bDelete, errDelete = smp.lp.RemoveLink(psf, restore)
	}

	if errDelete != nil {
		return bDelete, errDelete
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
		var errMsg = "SiteModule.DeletePage: Cannot got site project file path "
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if bSave == false || errSave != nil {
		return false, errSave
	}

	return true, nil
}

func (smp *SiteModule) ExportSourcePages(exportFolderPath string) (bool, error) {
	if nil != smp.spp {
		var targetFilesFolder = filepath.Join(exportFolderPath, "Files")
		if Utils.PathIsExist(targetFilesFolder) == false {
			_, errorMakeFolder := Utils.MakeFolder(targetFilesFolder)
			if errorMakeFolder != nil {
				var errMsg = "Export Source File : Cannot Make Files Folder"
				Utils.Logger.Println(errMsg)
				return false, errors.New(errMsg)
			}
		}
		Utils.CopyFolder(smp.GetSrcFilesFolderPath(smp.projectFolderPath), targetFilesFolder, true)
		return smp.spp.ExportSourcePages(exportFolderPath)
	}

	var errMsg = "SiteModule.ExportSourcePages: Site Project is empty"
	Utils.Logger.Println(errMsg)
	return false, errors.New(errMsg)
}

func (smp *SiteModule) AddFile(filePath string, addForce bool) (bool, error) {
	if Utils.PathIsExist(filePath) == false {
		var errMsg = "SiteModule.AddFile: File not exist " + filePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	targetFolderPath := smp.GetSrcFilesFolderPath(smp.projectFolderPath)

	if Utils.PathIsFile(filePath) {
		var fileName = filepath.Base(filePath)
		var targetFilePath = filepath.Join(targetFolderPath, fileName)

		if Utils.PathIsExist(targetFilePath) {
			if addForce {
				Utils.CopyFile(filePath, targetFilePath)
			} else {
				Utils.CopyFileWithConfirm(filePath, targetFilePath)
			}
		} else {
			Utils.CopyFile(filePath, targetFilePath)
		}
		return true, nil
	} else if Utils.PathIsDir(filePath) {
		var srcFolder = filepath.Base(filePath)
		var targetFolderPath = filepath.Join(targetFolderPath, srcFolder)
		return Utils.CopyFolder(filePath, targetFolderPath, addForce)
	}

	var errMsg = "SiteModule.AddFile: " + filePath + " is no file or folder, add file fail"
	Utils.Logger.Println(errMsg)
	return false, errors.New(errMsg)
}

func (smp *SiteModule) DeleteFile(filePath string) (bool, error) {
	if filePath == ".\\Files" || filePath == "./Files" {
		filePath = smp.GetSrcFilesFolderPath(smp.projectFolderPath)
		Utils.ClearFolder(filePath)
	} else {
		if strings.HasPrefix(filePath, ".\\Files\\") {
			filePath = filePath[8:]
		}

		if strings.HasPrefix(filePath, "./Files/") {
			filePath = filePath[8:]
		}
		var srcFilesFolder = smp.GetSrcFilesFolderPath(smp.projectFolderPath)
		filePath = filepath.Join(srcFilesFolder, filePath)

		if Utils.PathIsExist(filePath) == false {
			var errMsg = "SiteModule.DeleteFile: File not exist " + filePath
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
		var bDelete bool
		if Utils.PathIsFile(filePath) {
			bDelete = Utils.DeleteFile(filePath)
		} else if Utils.PathIsDir(filePath) {
			bDelete = Utils.DeleteFolder(filePath)
		}
		if bDelete == false {

			var errMsg = "SiteModule.Delete: " + filePath + " Delete file fail"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}
	return true, nil
}

func (smp *SiteModule) ListFile() {
	var srcFolderPath = smp.GetSrcFolderPath(smp.projectFolderPath)
	var srcFilesFolder = smp.GetSrcFilesFolderPath(smp.projectFolderPath)

	var nameLength int
	nameLength = 0
	var relativePathLength int
	relativePathLength = 0

	filepath.Walk(srcFilesFolder, func(path string, info os.FileInfo, err error) error {
		var fileName = info.Name()
		if len(fileName) > nameLength {
			nameLength = len(fileName)
		}
		var relativePath = path[len(srcFolderPath):]
		if len(relativePath) > relativePathLength {
			relativePathLength = len(relativePath)
		}
		return nil
	})

	var formatFileName string
	formatFileName = "%-" + strconv.Itoa(nameLength) + "s"
	var formatFileRelativePath string
	formatFileRelativePath = "%-" + strconv.Itoa(relativePathLength+3) + "s"

	fmt.Println("Files in Src/Files folder, will list file name and relative path, you can use this relative path as src/href in you md file ")
	fmt.Printf(formatFileName, "Name")
	fmt.Printf(formatFileRelativePath, "|  Relative Path ")
	fmt.Println("| Last Modified")

	var seperatorLength = nameLength + relativePathLength + 24
	var seperator string

	for index := 0; index < seperatorLength; index = index + 1 {
		seperator = seperator + "-"
	}
	fmt.Println(seperator)
	filepath.Walk(srcFilesFolder, func(path string, info os.FileInfo, err error) error {
		var fileName = info.Name()
		var relativePath = path[len(srcFolderPath):]
		var lastModified = info.ModTime().Format("2006-01-02 15:04:05")

		fmt.Printf(formatFileName, fileName)
		fmt.Printf(formatFileRelativePath, "| ."+relativePath)
		fmt.Println("| " + lastModified)
		return nil
	})
	fmt.Println(seperator)
}

func (smp *SiteModule) FilesStatics() int {
	var srcFilesFolder = smp.GetSrcFilesFolderPath(smp.projectFolderPath)
	var count int
	count = 0
	filepath.Walk(srcFilesFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			count = count + 1
		}
		return nil
	})

	return count
}

/*
func (smp *SiteModule) SearchPage(propertyName, propertyValue string) (string, error) {
	if len(smp.spp.SourceFiles) == 0 {
		return "", errors.New("No Pages")
	}

	for _, page := range smp.spp.SourceFiles {
		pValue, errPValue := page.GetProperty(propertyName)
		if errPValue != nil {
			return "", errors.New("Page doesn't have property " + propertyName)
		}

		if strings.Contains(pValue, propertyValue) == true {
			return page.ID, nil
		}
	}
	return "", errors.New("Not found")
}
*/
