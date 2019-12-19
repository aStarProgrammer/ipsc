package Site

import (
	"errors"
	"ipsc/Page"
	"ipsc/Utils"
	"os"
	"path/filepath"
	"strings"
)

type HtmlPageModule struct {
	spp *SiteProject
	smp *SiteModule
}

func (hpmp *HtmlPageModule) GetSiteProjectP() *SiteProject {
	return hpmp.spp
}

func (hpmp *HtmlPageModule) GetSiteModuleP() *SiteModule {
	return hpmp.smp
}

func NewHtmlPageModule(_spp *SiteProject, _smp *SiteModule) HtmlPageModule {
	var hpm HtmlPageModule
	hpm.spp = _spp
	hpm.smp = _smp
	return hpm
}

func FileIsHtml(filePath string) (bool, error) {
	if Utils.PathIsExist(filePath) == false {
		var errMsg = "Html file not exist"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var extension = filepath.Ext(filePath)

	if extension == ".html" || extension == ".htm" {
		return true, nil
	}
	return false, nil
}

func (hpmp *HtmlPageModule) AddHtml(title, description, author, filePath, titleImagePath string, isTop bool) (bool, string, error) {

	var htmlSrc, htmlDst string

	if Utils.PathIsExist(filePath) == false {
		var errMsg = "Html file not exist"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	bHtml, errHtml := FileIsHtml(filePath)

	if errHtml != nil {
		var errMsg = "Cannot confirm file type"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	} else if bHtml == false {
		var errMsg = "File is not html"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	_, fileName := filepath.Split(filePath)
	htmlSrc = filePath
	htmlDst = filepath.Join(hpmp.smp.GetSrcHtmlFolderPath(hpmp.smp.GetProjectFolderPath()), fileName)

	if Utils.PathIsExist(htmlDst) {
		var errMsg = "HtmlPageModule.AddHtml: Target Html File Already Exist"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	_, errCopy := Utils.CopyFile(htmlSrc, htmlDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "Copy File from " + htmlSrc + " to " + htmlDst + " Failed"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	var psf Page.PageSourceFile
	psf = Page.NewPageSourceFile()

	psf.Title = title
	psf.Author = author
	psf.Description = description
	psf.SourceFilePath = htmlDst
	psf.LastModified = Utils.CurrentTime()
	psf.Status = Page.ACTIVE
	psf.Type = Page.HTML
	if Utils.PathIsExist(titleImagePath) && Utils.PathIsImage(titleImagePath) {
		fileInfoTitleImage, errFileInfoTitleImage := os.Stat(titleImagePath)

		if errFileInfoTitleImage != nil {
			var errMsg = "Cannot get file size of titleImage"
			Utils.Logger.Println(errMsg)
			return false, "", errors.New(errMsg)
		}

		titleImageSize := fileInfoTitleImage.Size()

		if titleImageSize > MAXTITLEIMAGESIZE {
			var errMsg = "Title Image bigger than 30KB"
			Utils.Logger.Println(errMsg)
			return false, "", errors.New(errMsg)
		}
		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}
	psf.IsTop = isTop
	psf.OutputFile = ""

	bAdd, errorAdd := hpmp.spp.AddPageSourceFile(psf) //Add to Source Pages

	if bAdd == false && errorAdd != nil {
		Utils.Logger.Println(errorAdd.Error())
		return false, "", errorAdd
	}

	return true, psf.ID, nil
}

func (hpmp *HtmlPageModule) RemoveHtml(psf Page.PageSourceFile, restore bool) (bool, error) {
	var outputFileID = psf.OutputFile
	if outputFileID != "" {
		var pofIndex = hpmp.spp.GetPageOutputFile(outputFileID)
		var pof Page.PageOutputFile
		if pofIndex != -1 {
			pof = hpmp.spp.OutputFiles[pofIndex]

			if restore == false {
				bDelOutput, errDeleteOutput := hpmp.spp.RemovePageOutputFile(pof)
				if errDeleteOutput != nil {
					return bDelOutput, errDeleteOutput
				}
				if pof.FilePath != "" {
					bDeleteOutputFile := Utils.DeleteFile(pof.FilePath)
					if bDeleteOutputFile == false {
						var errMsg = "Cannot delete output file " + pof.FilePath
						Utils.Logger.Println(errMsg)
						return false, errors.New(errMsg)
					}
				}
			}
		}
	}

	bRemove, errRemove := hpmp.spp.RemovePageSourceFile(psf, restore)
	if errRemove != nil {
		iFind := hpmp.spp.GetPageSourceFile(psf.ID)
		if iFind == -1 {
			hpmp.spp.AddPageSourceFile(psf)
		}
		return bRemove, errRemove
	}

	var filePath = psf.SourceFilePath

	if restore == false {
		if Utils.DeleteFile(filePath) == false {
			hpmp.spp.AddPageSourceFile(psf)
			var errMsg = "Delete File from Disk Fail"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
	}
	return true, nil
}

func (hpmp *HtmlPageModule) RestoreHtml(ID string) (bool, error) {
	return hpmp.spp.ResotrePageSourceFile(ID)
}

func (hpmp *HtmlPageModule) UpdateHtml(psf Page.PageSourceFile, filePath string) (bool, error) {
	_psfID := hpmp.spp.GetPageSourceFile(psf.ID)
	if _psfID == -1 {
		var errMsg = "File not found"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	psf_Backup := hpmp.spp.SourceFiles[_psfID]

	if filePath != psf.SourceFilePath {

		if Utils.PathIsExist(filePath) == false {
			var errMsg = "Html file not exist"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		bHtml, errHtml := FileIsHtml(filePath)

		if errHtml != nil {
			var errMsg = "Cannot confirm file type"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		} else if bHtml == false {
			var errMsg = "File is not html"
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}

		_, fileName := filepath.Split(filePath)

		var htmlSrc, htmlDst string
		htmlSrc = filePath
		htmlDst = filepath.Join(hpmp.smp.GetSrcHtmlFolderPath(hpmp.smp.GetProjectFolderPath()), fileName)
		psf.SourceFilePath = htmlDst

		bUpdate, errUpdate := hpmp.spp.UpdatePageSourceFile(psf)

		if errUpdate != nil {
			return bUpdate, errUpdate
		}

		if psf.SourceFilePath != psf_Backup.SourceFilePath {
			Utils.DeleteFile(psf_Backup.SourceFilePath)
		}

		_, errCopy := Utils.CopyFile(htmlSrc, htmlDst)

		if errCopy != nil {
			var errMsg string
			errMsg = "Copy File from " + htmlSrc + " to " + htmlDst + " Failed"
			//恢复被更新的内容
			hpmp.spp.UpdatePageSourceFile(psf_Backup)
			Utils.Logger.Println(errMsg)
			return false, errors.New(errMsg)
		}
	} else {
		bUpdate, errUpdate := hpmp.spp.UpdatePageSourceFile(psf)
		if errUpdate != nil {
			return bUpdate, errUpdate
		}
	}
	return true, nil
}

func (hpmp *HtmlPageModule) GetHtmlFile(ID string) string {
	iFind := hpmp.spp.GetPageSourceFile(ID)
	if iFind != -1 {
		psf := hpmp.spp.SourceFiles[iFind]

		if psf.SourceFilePath != "" {
			return psf.SourceFilePath
		}
	}

	return ""
}

func (hpmp *HtmlPageModule) GetHtmlInformation(ID string) int {
	return hpmp.spp.GetPageSourceFile(ID)
}

func (hpmp *HtmlPageModule) UpdateHtmlInformation(title, description, author, filePath, titleImagePath string) (bool, error) {
	var psf Page.PageSourceFile
	psf = Page.NewPageSourceFile()

	psf.Title = title
	psf.Author = author
	psf.Description = description
	psf.SourceFilePath = filePath
	psf.LastModified = Utils.CurrentTime()
	psf.Status = Page.ACTIVE
	psf.Type = Page.HTML
	if Utils.PathIsExist(titleImagePath) && Utils.PathIsImage(titleImagePath) {
		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}

	bUpdate, errorUpdate := hpmp.spp.UpdatePageSourceFile(psf) //Update Source Pages

	if bUpdate == false && errorUpdate != nil {
		Utils.Logger.Println(errorUpdate.Error())
		return false, errorUpdate
	}

	return true, nil
}

//Compile Html, just copy html from Src to Output folder, change source information and add PageOutputFile
func (hpmp *HtmlPageModule) Compile(ID string) (int, error) {
	iFind := hpmp.spp.GetPageSourceFile(ID)
	if iFind == -1 {
		var errMsg string
		errMsg = "Cannot find the source File with ID " + ID
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	psf := hpmp.spp.SourceFiles[iFind]

	if psf.SourceFilePath == "" {
		var errMsg string
		errMsg = "Page Source File FilePath is emtpy"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	if psf.Status == Page.RECYCLED {
		var errMsg string
		errMsg = "Page Source File is in Recycled status, cannot Compile"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	var htmlSrc, htmlDst string
	htmlSrc = psf.SourceFilePath

	if Utils.PathIsExist(htmlSrc) == false {
		var errMsg = "Source Html File not found on the disk"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	bHtml, errHtml := FileIsHtml(htmlSrc)

	if errHtml != nil {
		var errMsg = "Cannot confirm file type"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	} else if bHtml == false {
		var errMsg = "File is not html"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	_, fileName := filepath.Split(htmlSrc)
	ext := filepath.Ext(htmlSrc)
	fileNameOnly := strings.TrimSuffix(fileName, ext)
	newFileName := fileNameOnly + ".html"
	htmlDst = filepath.Join(hpmp.smp.GetOutputFolderPath(hpmp.smp.GetProjectFolderPath()), "Pages", newFileName)

	_, errCopy := Utils.CopyFile(htmlSrc, htmlDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "Copy File from " + htmlSrc + " to " + htmlDst + " Failed"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	var _pofIndex int

	if psf.OutputFile != "" {
		pofIndex := hpmp.spp.GetPageOutputFile(psf.OutputFile)

		pof := hpmp.spp.OutputFiles[pofIndex]

		pof.Author = psf.Author
		pof.Description = psf.Description
		pof.FilePath = htmlDst
		pof.IsTop = psf.IsTop
		pof.Title = psf.Title
		pof.TitleImage = psf.TitleImage
		pof.Type = psf.Type
		pof.CreateTime = Utils.CurrentTime()

		_, errUpdatePof := hpmp.spp.UpdatePageOutputFile(pof)

		if errUpdatePof != nil {
			Utils.DeleteFile(htmlDst) //Add fail,delete the file already copied
			var errMsg = "HtmlPageModule: Page Out File Update Fail"
			Utils.Logger.Println(errMsg)
			return -1, errUpdatePof
		}
	} else {
		pof := Page.NewPageOutputFile()
		pof.Author = psf.Author
		pof.Description = psf.Description
		pof.FilePath = htmlDst
		pof.IsTop = psf.IsTop
		pof.Title = psf.Title
		pof.TitleImage = psf.TitleImage
		pof.Type = psf.Type
		pof.CreateTime = Utils.CurrentTime()

		_, errAdd := hpmp.spp.AddPageOutputFile(pof)

		if errAdd != nil {
			Utils.DeleteFile(htmlDst) //Add fail,delete the file already copied
			return -1, errAdd
		}

		_pofIndex = hpmp.spp.GetPageOutputFile(pof.ID)

		if _pofIndex == -1 {
			Utils.DeleteFile(htmlDst) //Add fail,delete the file already copied
			var errMsg = "HtmlPageModule: Page Out File add Fail"
			Utils.Logger.Println(errMsg)
			return _pofIndex, errors.New(errMsg)
		}

		psf.OutputFile = pof.ID
	}
	psf.LastCompiled = Utils.CurrentTime()

	hpmp.spp.UpdatePageSourceFile(psf)

	return _pofIndex, nil
}
