package Site

import (
	"errors"
	"ipsc/Page"
	"ipsc/Utils"
	"os"
)

const MAXTITLEIMAGESIZE int64 = 30720 //Title Image must smaller than 30KB

type LinkModule struct {
	spp *SiteProject
}

func NewLinkModule(_spp *SiteProject) LinkModule {
	var lm LinkModule
	lm.spp = _spp

	return lm
}

func (lmp *LinkModule) GetSiteProjectP() *SiteProject {
	return lmp.spp
}

func (lmp *LinkModule) AddLink(title, description, author, url, titleImagePath string, isTop bool) (bool, string, error) {

	if lmp.LinkAlreadyExist(url, title) {
		var errMsg = "LinkModule.AddLink: Target Link Already Exist"
		Utils.Logger.Println(errMsg)
		return false, "", errors.New(errMsg)
	}

	var psf Page.PageSourceFile
	psf = Page.NewPageSourceFile()

	psf.Title = title
	psf.Author = author
	psf.Description = description
	psf.SourceFilePath = url
	psf.LastModified = Utils.CurrentTime()
	psf.Status = Page.ACTIVE
	psf.Type = Page.LINK
	psf.OutputFile = ""
	psf.IsTop = isTop
	if Utils.PathIsExist(titleImagePath) && Utils.PathIsImage(titleImagePath) {

		fileInfo, errFileInfo := os.Stat(titleImagePath)

		if errFileInfo != nil {
			var errMsg = "Cannot get file size of titleImage"
			Utils.Logger.Println(errMsg)
			Utils.Logger.Println(errFileInfo.Error())
			return false, "", errors.New(errMsg)
		}

		imageSize := fileInfo.Size()

		if imageSize > MAXTITLEIMAGESIZE {
			var errMsg = "Title Image bigger than 30KB"
			Utils.Logger.Println(errMsg)
			return false, "", errors.New(errMsg)
		}

		psf.TitleImage, _ = Utils.ReadImageAsBase64(titleImagePath)
	} else {
		psf.TitleImage = ""
	}

	bAdd, errAdd := lmp.spp.AddPageSourceFile(psf)
	return bAdd, psf.ID, errAdd
}

func (lmp *LinkModule) RemoveLink(psf Page.PageSourceFile, restore bool) (bool, error) {
	var outputFileID = psf.OutputFile
	if outputFileID != "" {
		var pofIndex = lmp.spp.GetPageOutputFile(outputFileID)
		var pof Page.PageOutputFile
		if pofIndex != -1 {
			pof = lmp.spp.OutputFiles[pofIndex]
			bDelete, errDelete := lmp.spp.RemovePageOutputFile(pof)
			if errDelete != nil {
				var errMsg = "LinkModule.RemoveLink: Cannot remove Page output File"
				Utils.Logger.Println(errMsg)
				return bDelete, errors.New(errMsg)
			}
		}
	}

	return lmp.spp.RemovePageSourceFile(psf, restore)
}

func (lmp *LinkModule) RestoreLink(ID string) (bool, error) {
	return lmp.spp.ResotrePageSourceFile(ID)
}

func (lmp *LinkModule) UpdateLink(psf Page.PageSourceFile) (bool, error) {
	return lmp.spp.UpdatePageSourceFile(psf)
}

func (lmp *LinkModule) GetLink(ID string) int {
	return lmp.spp.GetPageSourceFile(ID)
}

func (lmp *LinkModule) LinkAlreadyExist(linkUrl, linkTitle string) bool {
	for _, link := range lmp.spp.SourceFiles {
		if link.Type == Page.LINK && link.Title == linkTitle && link.SourceFilePath == linkUrl {
			return true
		}
	}
	return false
}

func (lmp *LinkModule) Compile(ID string) (int, error) {
	iFind := lmp.spp.GetPageSourceFile(ID)
	if iFind == -1 {
		var errMsg string
		errMsg = "Cannot find the source File with ID " + ID
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	psf := lmp.spp.SourceFiles[iFind]

	if psf.SourceFilePath == "" {
		var errMsg string
		errMsg = "Page Source File Url is emtpy"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	if psf.Status == Page.RECYCLED {
		var errMsg string
		errMsg = "Page Source File is in Recycled status, cannot Compile"
		Utils.Logger.Println(errMsg)
		return -1, errors.New(errMsg)
	}

	var _pofIndex int
	pofIndex := lmp.spp.GetPageOutputFile(psf.OutputFile)
	if psf.OutputFile == "" || pofIndex == -1 {
		pof := Page.NewPageOutputFile()
		pof.Author = psf.Author
		pof.Description = psf.Description
		pof.FilePath = psf.SourceFilePath
		pof.IsTop = psf.IsTop
		pof.Title = psf.Title
		pof.TitleImage = psf.TitleImage
		pof.Type = psf.Type
		pof.CreateTime = Utils.CurrentTime()

		_, errAdd := lmp.spp.AddPageOutputFile(pof)

		if errAdd != nil {
			return -1, errAdd
		}

		_pofIndex = lmp.spp.GetPageOutputFile(pof.ID)

		if _pofIndex == -1 {
			var errMsg = "LinkModule: Page Out File add Fail"
			Utils.Logger.Println(errMsg)
			return _pofIndex, errors.New(errMsg)
		}

		psf.OutputFile = pof.ID

	} else {
		pof := lmp.spp.OutputFiles[pofIndex]
		pof.Author = psf.Author
		pof.Description = psf.Description
		pof.FilePath = psf.SourceFilePath
		pof.IsTop = psf.IsTop
		pof.Title = psf.Title
		pof.TitleImage = psf.TitleImage
		pof.Type = psf.Type
		pof.CreateTime = Utils.CurrentTime()

		_, errUpdatePof := lmp.spp.UpdatePageOutputFile(pof)

		if errUpdatePof != nil {
			var errMsg = "LinkModule: Page Out File update Fail"
			Utils.Logger.Println(errMsg)
			return -1, errUpdatePof
		}

	}
	psf.LastCompiled = Utils.CurrentTime()

	lmp.spp.UpdatePageSourceFile(psf)

	return _pofIndex, nil
}
