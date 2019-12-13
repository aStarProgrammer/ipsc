package Site

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ipsc/Configuration"
	"ipsc/Page"
	"ipsc/Utils"

	//"os"
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
<<<<<<< HEAD
		fmt.Println("SiteModule.NewSiteModule: Cannot create Site Module")
=======
		fmt.Println("NewSiteModule: Cannot create Site Module")
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
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

func (smp *SiteModule) GetTemplateFolderPath(projectFolderPath string) string {
	return filepath.Join(projectFolderPath, "Templates")
}

func (smp *SiteModule) GetSiteProjectFilePath(projectFolderPath string) (string, error) {
	if nil != smp.spp && smp.spp.Title == "" {
<<<<<<< HEAD
		var errMsg = "SiteModuole.GetSiteProjectFilePath: SiteProject Title is empty"
=======
		var errMsg = "SiteProject Title is empty"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	return filepath.Join(projectFolderPath, smp.spp.Title) + ".sp", nil
}

func (smp *SiteModule) PathIsSiteProject(projectPath, projectName string) (bool, error) {
	if Utils.PathIsExist(projectPath) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: " + projectPath + " is not exist"
=======
		var errMsg = projectPath + " is not exist"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var projectFilePath = filepath.Join(projectPath, projectName)

	if strings.HasSuffix(projectFilePath, ".sp") == false {

		projectFilePath += ".sp"
	}

	if Utils.PathIsExist(projectFilePath) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: Cannot find sp file in project " + projectPath
=======
		var errMsg = "Cannot find sp file in project " + projectPath
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var sp SiteProject
	_, loadError := sp.LoadFromFile(projectFilePath)
	if loadError != nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: Cannot load sp file in project " + projectPath
		fmt.Println(errMsg)
		return false, errors.New("SiteModuole.PathIsSiteProject: Cannot load sp file in project " + projectPath)
=======
		var errMsg = "Cannot load sp file in project " + projectPath
		fmt.Println(errMsg)
		return false, errors.New("Cannot load sp file in project " + projectPath)
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
	}

	var srcFolderPath = smp.GetSrcFolderPath(projectPath)

	if Utils.PathIsExist(srcFolderPath) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: Cannot find src folder"
=======
		var errMsg = "Cannot find src folder"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var markdownFolderPath = smp.GetSrcMarkdownFolderPath(projectPath)

	if Utils.PathIsExist(markdownFolderPath) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: Cannot find markdown folder"
=======
		var errMsg = "Cannot find markdown folder"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var htmlFolderPath = smp.GetSrcHtmlFolderPath(projectPath)

	if Utils.PathIsExist(htmlFolderPath) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: Cannot find html folder"
=======
		var errMsg = "Cannot find html folder"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var outputFolderPath = smp.GetOutputFolderPath(projectPath)

	if Utils.PathIsExist(outputFolderPath) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.PathIsSiteProject: Cannot find output folder"
=======
		var errMsg = "Cannot find output folder"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	return true, nil
}

func (smp *SiteModule) InitializeSiteProjectFolder(siteTitle, siteAuthor, siteDescription, _projectFolderPath, _outputFolderPath string) (bool, error) {
	if _projectFolderPath == "" {
<<<<<<< HEAD
		var errMsg = "SiteModuole.InitializeSiteProjectFolder: Project Folder Path is empty"
=======
		var errMsg = "Project Folder Path is empty"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	//Create each foldrs

	//ProjectFolder
	var errProjectFolder error
	if !Utils.PathIsExist(_projectFolderPath) {
		_, errProjectFolder = Utils.MakeFolder(_projectFolderPath)

		if errProjectFolder != nil {
<<<<<<< HEAD
			fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errProjectFolder.Error())
=======
			fmt.Println(errProjectFolder.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
			return false, errProjectFolder
		}

	}

	//ProjectFolder->Src
	var srcFolderPath = smp.GetSrcFolderPath(_projectFolderPath)
	var errSrcFolderPath error
	if !Utils.PathIsExist(srcFolderPath) {
		_, errSrcFolderPath = Utils.MakeFolder(srcFolderPath)

		if errSrcFolderPath != nil {
<<<<<<< HEAD
			fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errSrcFolderPath.Error())
=======
			fmt.Println(errSrcFolderPath.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
			return false, errSrcFolderPath
		}
	}

	//ProjectFolder->Src->Markdown
	var srcMarkdownFolderPath = smp.GetSrcMarkdownFolderPath(_projectFolderPath)
	var errSrcMarkdownFoldrPath error
	if !Utils.PathIsExist(srcMarkdownFolderPath) {
		_, errSrcMarkdownFoldrPath = Utils.MakeFolder(srcMarkdownFolderPath)

		if errSrcMarkdownFoldrPath != nil {
<<<<<<< HEAD
			fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errSrcMarkdownFoldrPath.Error())
=======
			fmt.Println(errSrcMarkdownFoldrPath.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
			return false, errSrcMarkdownFoldrPath
		}
	}

	//ProjectFolder->Src->Html
	var srcHtmlFolderPath = smp.GetSrcHtmlFolderPath(_projectFolderPath)
	var errSrcHtmlFolderPath error

	if !Utils.PathIsExist(srcHtmlFolderPath) {
		_, errSrcHtmlFolderPath = Utils.MakeFolder(srcHtmlFolderPath)

		if errSrcHtmlFolderPath != nil {
<<<<<<< HEAD
			fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errSrcHtmlFolderPath.Error())
=======
			fmt.Println(errSrcHtmlFolderPath.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
			return false, errSrcHtmlFolderPath
		}
	}

	//ProjectFolder->Output
	var outputFolderPath = smp.GetOutputFolderPath(_projectFolderPath)
	var errOutputFolderPath error

	if outputFolderPath == _outputFolderPath || _outputFolderPath == "" {

		if !Utils.PathIsExist(outputFolderPath) {
			_, errOutputFolderPath = Utils.MakeFolder(outputFolderPath)

			if errOutputFolderPath != nil {
<<<<<<< HEAD
				fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errOutputFolderPath.Error())
=======
				fmt.Println(errOutputFolderPath.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
				return false, errOutputFolderPath
			}

			var outputPagesSubFolder = filepath.Join(outputFolderPath, "Pages")
			_, errOutputPagesFolder := Utils.MakeFolder(outputPagesSubFolder)

			if errOutputPagesFolder != nil {
<<<<<<< HEAD
				fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errOutputPagesFolder.Error())
=======
				fmt.Println(errOutputPagesFolder.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
				return false, errOutputPagesFolder
			}

		}
	} else {
		if !Utils.PathIsExist(outputFolderPath) {
			_, errOutputFolderPath = Utils.MakeSoftLink4Folder(_outputFolderPath, outputFolderPath)

			if errOutputFolderPath != nil {
<<<<<<< HEAD
				fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errOutputFolderPath.Error())
=======
				fmt.Println(errOutputFolderPath.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
				return false, errOutputFolderPath
			}

			var outputPagesSubFolder = filepath.Join(outputFolderPath, "Pages")

			if Utils.PathIsExist(outputPagesSubFolder) == false {
				_, errOutputPagesFolder := Utils.MakeFolder(outputPagesSubFolder)

				if errOutputPagesFolder != nil {
<<<<<<< HEAD
					fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errOutputPagesFolder.Error())
=======
					fmt.Println(errOutputPagesFolder.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
					return false, errOutputPagesFolder
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
<<<<<<< HEAD
			fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errTemplateFolder.Error())
=======
			fmt.Println(errTemplateFolder.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
			return false, errTemplateFolder
		}
	}
	//Copy temlates from Resources
	srcTemplateFolder, errSrcTemplate := Configuration.GetTemplatesFolderPath()
	if errSrcTemplate != nil {
<<<<<<< HEAD
		fmt.Println("SiteModuole.InitializeSiteProjectFolder: " + errSrcTemplate.Error())
=======
		fmt.Println(errSrcTemplate.Error())
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		return false, errSrcTemplate
	}

	if Utils.PathIsExist(srcTemplateFolder) == false {
<<<<<<< HEAD
		var errMsg = "SiteModuole.InitializeSiteProjectFolder: Try to copy tempaltes, src tempalte folder not exist " + srcTemplateFolder
=======
		var errMsg = "Try to copy tempaltes, src tempalte folder not exist " + srcTemplateFolder
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	files, _ := ioutil.ReadDir(srcTemplateFolder)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			srcTemplateFilePath := filepath.Join(srcTemplateFolder, f.Name())
			dstTemplateFilePath := filepath.Join(templateFolderPath, f.Name())

			_, errCopy := Utils.CopyFile(srcTemplateFilePath, dstTemplateFilePath)
			if errCopy != nil {
<<<<<<< HEAD
				var errMsg = "SiteModuole.InitializeSiteProjectFolder: Cannot copy template file " + srcTemplateFilePath + " to " + dstTemplateFilePath
=======
				var errMsg = "Cannot copy template file " + srcTemplateFilePath + " to " + dstTemplateFilePath
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
				fmt.Println(errMsg)
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
		fmt.Println(errProjectFilePath.Error())
		return false, errProjectFilePath
	}

	bSaveToFile, errSaveToFile := smp.spp.SaveToFile(projectFilePath)

	if bSaveToFile == false || errSaveToFile != nil {
		fmt.Println(errSaveToFile.Error())
		return false, errSaveToFile
	}

	return true, nil
}

func (smp *SiteModule) OpenSiteProject(projectFolderPath, projectName string) (bool, error) {

	if projectFolderPath == "" {
<<<<<<< HEAD
		fmt.Println("SiteModuole.OpenSiteProject: Project Folder path is empty")
		return false, errors.New("SiteModuole.OpenSiteProject: Project Folder path is empty")
=======
		fmt.Println("Project Folder path is empty")
		return false, errors.New("Project Folder path is empty")
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
	}
	//fmt.Println("OpenSPPointA")
	bIsSP, errIsSP := smp.PathIsSiteProject(projectFolderPath, projectName)

	if errIsSP != nil || false == bIsSP {
<<<<<<< HEAD
		var errMsg = "SiteModuole.OpenSiteProject: Path " + projectFolderPath + " doesn't contain a IPSC Site"
=======
		var errMsg = "Path " + projectFolderPath + " doesn't contain a IPSC Site"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
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
<<<<<<< HEAD
		var errMsg = "SiteModuole.OpenSiteProject: Cannot load sp file in project " + siteProjectFilePath
=======
		var errMsg = "Cannot load sp file in project " + siteProjectFilePath
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
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
<<<<<<< HEAD
		var errMsg = "SiteModuole.UpdateSiteProject: siteFolder " + siteFolder + " doesn't exist"
=======
		var errMsg = "siteFolder " + siteFolder + " doesn't exist"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if smp.spp == nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.UpdateSiteProject: Site Project is nil"
=======
		var errMsg = "Site Project is nil"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if siteTitle != "" && smp.spp.Title != siteTitle {
		if Utils.PathIsExist(oldSiteProjectFilePath) {
			var newSiteProjectFilePath = filepath.Join(siteFolder, siteTitle+".sp")

			_, errMove := Utils.MoveFile(oldSiteProjectFilePath, newSiteProjectFilePath)
			if errMove != nil {
				fmt.Println(errMove.Error())
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

	if errSave != nil {
		fmt.Println(errSave.Error())
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
<<<<<<< HEAD
		var errMsg = "SiteModuole.RestoreRecycledPageSourceFile: RestoreRecycledPageSourceFile: " + "ID is empty"
=======
		var errMsg = "RestoreRecycledPageSourceFile: " + "ID is empty"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	index := smp.spp.GetPageSourceFile(ID)

	if index == -1 {
<<<<<<< HEAD
		var errMsg = "SiteModuole.RestoreRecycledPageSourceFile: Cannot find Page Source File with ID " + ID
=======
		var errMsg = "Cannot find Page Source File with ID " + ID
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}

	bResotre, errRestore := smp.spp.ResotrePageSourceFile(ID)
	if errRestore != nil {
		return bResotre, errRestore
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.RestoreRecycledPageSourceFile: Cannot got site project file path "
=======
		var errMsg = "Cannot got site project file path "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	return smp.spp.SaveToFile(siteProjectFilePath)
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
<<<<<<< HEAD
		var errMsg = "SiteModuole.RestoreRecycledPageSourceFile: Cannot got site project file path "
=======
		var errMsg = "Cannot got site project file path "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	return smp.spp.SaveToFile(siteProjectFilePath)
}

func (smp *SiteModule) Compile(indexPageSize string) (bool, error) {
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
			if (sp.LastCompiled == "" && sp.OutputFile == -1) || (sp.OutputFile != -1 && sp.LastCompiled != "" && sp.LastModified != "" && sp.LastCompiled < sp.LastModified) {
				if sp.Type == Page.MARKDOWN {
					_, errCompileMd := smp.mpp.Compile(sp.ID)
					if errCompileMd != nil {
						return false, errCompileMd
					}
					mdCount++
				} else if sp.Type == Page.HTML {
					_, errCompileHtml := smp.hpp.Compile(sp.ID)
					if errCompileHtml != nil {
						return false, errCompileHtml
					}
					htmlCount++
				} else if sp.Type == Page.LINK {
					_, errCompileLink := smp.lp.Compile(sp.ID)
					if errCompileLink != nil {
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
	//fmt.Println("D")
	//Compile Index Page and More Pages
	_, errCompileIndex := smp.mpp.Compile_Psf(smp.spp.IndexPageSourceFile)

<<<<<<< HEAD
=======
	_, errCompileIndex := smp.mpp.Compile_Psf(smp.spp.IndexPageSourceFile)

>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
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
	//Get Summary and write to spp
	var CompileSummary string
	CompileSummary = "Index: 1"
	CompileSummary += "_More: " + strconv.Itoa(moreCount)
	CompileSummary += "_Markdown: " + strconv.Itoa(mdCount)
	CompileSummary += "_Html: " + strconv.Itoa(htmlCount)
	CompileSummary += "_Link: " + strconv.Itoa(linkCount)

	smp.spp.LastCompileSummary = CompileSummary
	smp.spp.LastCompiled = Utils.CurrentTime()

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.Compile: Cannot got site project file path "
=======
		var errMsg = "Cannot got site project file path "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if errSave != nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.Compile: Cannot save site project file "
=======
		var errMsg = "Cannot save site project file "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return bSave, errors.New(errMsg)
	}
	return true, nil
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
		fmt.Println(errAdd.Error())
		return bAdd, "-1", errAdd
	}

	siteProjectFilePath, errPath := smp.GetSiteProjectFilePath(smp.projectFolderPath)

	if errPath != nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.AddPage: Cannot got site project file path "
=======
		var errMsg = "Cannot got site project file path "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, "-1", errors.New(errMsg)
	}

	bSave, errSave := smp.spp.SaveToFile(siteProjectFilePath)

	if errSave != nil {
<<<<<<< HEAD
		var errMsg = "SiteModuole.AddPage: Cannot save site project file "
=======
		var errMsg = "Cannot save site project file "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
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
<<<<<<< HEAD
		var errMsg = "SiteModuole.UpdatePage: Cannot find page with ID " + pageID
=======
		var errMsg = "Cannot find page with ID " + pageID
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
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
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	return smp.spp.SaveToFile(siteProjectFilePath)

}

func (smp *SiteModule) DeletePage(pageID string, restore bool) (bool, error) {

	var index = smp.spp.GetPageSourceFile(pageID)

	if index == -1 {
<<<<<<< HEAD
		var errMsg = "SiteModuole.DeletePage: Cannot find page with ID " + pageID
=======
		var errMsg = "Cannot find page with ID " + pageID
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
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
<<<<<<< HEAD
		var errMsg = "SiteModuole.DeletePage: Cannot got site project file path "
=======
		var errMsg = "Cannot got site project file path "
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
		fmt.Println(errMsg)
		return false, errors.New(errMsg)
	}
	return smp.spp.SaveToFile(siteProjectFilePath)
}

func (smp *SiteModule) ExportSourcePages(exportFolderPath string) (bool, error) {
	if nil != smp.spp {
		return smp.spp.ExportSourcePages(exportFolderPath)
	}
<<<<<<< HEAD
	var errMsg = "SiteModuole.ExportSourcePages: Site Project is empty"
=======
	var errMsg = "Site Project is empty"
>>>>>>> 71276fde19654e48a3fd9f74fefda5cbdd634d5a
	fmt.Println(errMsg)
	return false, errors.New(errMsg)
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
