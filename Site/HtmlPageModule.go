package Site

import (
	"IPSC/Page"
	"IPSC/Utils"
	"errors"
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
		return false, errors.New("Html file not exist")
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
		return false, "", errors.New("Html file not exist")
	}

	bHtml, errHtml := FileIsHtml(filePath)

	if errHtml != nil {
		return false, "", errors.New("Cannot confirm file type")
	} else if bHtml == false {
		return false, "", errors.New("File is not html")
	}

	_, fileName := filepath.Split(filePath)
	htmlSrc = filePath
	htmlDst = filepath.Join(hpmp.smp.GetSrcHtmlFolderPath(hpmp.smp.GetProjectFolderPath()), fileName)

	_, errCopy := Utils.CopyFile(htmlSrc, htmlDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "Copy File from " + htmlSrc + " to " + htmlDst + " Failed"
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
			return false, "", errors.New("Cannot get file size of titleImage")
		}

		titleImageSize := fileInfoTitleImage.Size()

		if titleImageSize > MAXTITLEIMAGESIZE {
			return false, "", errors.New("Title Image bigger than 30KB")
		}
		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}
	psf.IsTop = isTop
	psf.OutputFile = -1

	bAdd, errorAdd := hpmp.spp.AddPageSourceFile(psf) //Add to Source Pages

	if bAdd == false && errorAdd != nil {
		return false, "", errorAdd
	}

	return true, psf.ID, nil
}

func (hpmp *HtmlPageModule) RemoveHtml(psf Page.PageSourceFile, restore bool) (bool, error) {
	var outputIndex = psf.OutputFile
	if outputIndex != -1 {
		var pof = hpmp.spp.OutputFiles[outputIndex]
		if restore == false {
			bDelOutput, errDeleteOutput := hpmp.spp.RemovePageOutputFile(pof)
			if errDeleteOutput != nil {
				return bDelOutput, errDeleteOutput
			}
			if pof.FilePath != "" {
				bDeleteOutputFile := Utils.DeleteFile(pof.FilePath)
				if bDeleteOutputFile == false {
					return false, errors.New("Cannot delete output file " + pof.FilePath)
				}
			}
		}
	}

	bRemove, errRemove := hpmp.spp.RemovePageSourceFile(psf, restore)
	if errRemove != nil {
		iFind := hpmp.spp.GetPageOutputFile(psf.ID)
		if iFind == -1 {
			hpmp.spp.AddPageSourceFile(psf)
		}
		return bRemove, errRemove
	}

	var filePath = psf.SourceFilePath

	if restore == false {
		if Utils.DeleteFile(filePath) == false {
			hpmp.spp.AddPageSourceFile(psf)
			return false, errors.New("Delete File from Disk Fail")
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
		return false, errors.New("File not found")
	}

	psf_Backup := hpmp.spp.SourceFiles[_psfID]

	if filePath == "" {
		return true, nil
	}

	if Utils.PathIsExist(filePath) == false {
		return false, errors.New("Html file not exist")
	}

	bHtml, errHtml := FileIsHtml(filePath)

	if errHtml != nil {
		return false, errors.New("Cannot confirm file type")
	} else if bHtml == false {
		return false, errors.New("File is not html")
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
		return false, errors.New(errMsg)
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
		return false, errorUpdate
	}

	return true, nil
}

//Complie Html, just copy html from Src to Output folder, change sourceinformation and add PageOutputFile
func (hpmp *HtmlPageModule) Complie(ID string) (int, error) {
	iFind := hpmp.spp.GetPageSourceFile(ID)
	if iFind == -1 {
		var errMsg string
		errMsg = "Cannot find the source File with ID " + ID
		return -1, errors.New(errMsg)
	}

	psf := hpmp.spp.SourceFiles[iFind]

	if psf.SourceFilePath == "" {
		var errMsg string
		errMsg = "Page Source File FilePath is emtpy"
		return -1, errors.New(errMsg)
	}

	if psf.Status == Page.RECYCLED {
		var errMsg string
		errMsg = "Page Source File is in Recycled status, cannot complie"
		return -1, errors.New(errMsg)
	}

	var htmlSrc, htmlDst string
	htmlSrc = psf.SourceFilePath

	if Utils.PathIsExist(htmlSrc) == false {
		return -1, errors.New("Source Html File not found on the disk")
	}

	bHtml, errHtml := FileIsHtml(htmlSrc)

	if errHtml != nil {
		return -1, errors.New("Cannot confirm file type")
	} else if bHtml == false {
		return -1, errors.New("File is not html")
	}

	_, fileName := filepath.Split(htmlSrc)
	ext := filepath.Ext(htmlSrc)
	fileNameOnly := strings.TrimSuffix(fileName, ext)
	newFileName := fileNameOnly + "_" + Utils.GUID() + ".html"
	htmlDst = filepath.Join(hpmp.smp.GetOutputFolderPath(hpmp.smp.GetProjectFolderPath()), "Pages", newFileName)

	_, errCopy := Utils.CopyFile(htmlSrc, htmlDst)

	if errCopy != nil {
		var errMsg string
		errMsg = "Copy File from " + htmlSrc + " to " + htmlDst + " Failed"
		return -1, errors.New(errMsg)
	}

	pof := Page.NewPageOutputFile()
	pof.Author = psf.Author
	pof.Description = psf.Description
	pof.FilePath = htmlDst
	pof.IsTop = psf.IsTop
	pof.Title = psf.Title
	pof.TitleImage = psf.TitleImage
	pof.Type = psf.Type

	_, errAdd := hpmp.spp.AddPageOutputFile(pof)

	if errAdd != nil {
		Utils.DeleteFile(htmlDst) //Add fail,delete the file already copied
		return -1, errAdd
	}

	_pofID := hpmp.spp.GetPageOutputFile(pof.ID)

	if _pofID == -1 {
		Utils.DeleteFile(htmlDst) //Add fail,delete the file already copied
		return _pofID, errors.New("Page Out File add Fail")
	}

	psf.OutputFile = _pofID
	psf.LastComplied = Utils.CurrentTime()

	hpmp.spp.UpdatePageSourceFile(psf)

	return _pofID, nil
}
