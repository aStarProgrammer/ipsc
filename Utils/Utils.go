package Utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shamsher31/goimgtype"
)

//
func PathIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func CurrentTime() string {
	t := time.Now()
	str := t.Format("2006-01-02 15:04:05")

	return str
}

func PathIsMarkdown(filePath string) bool {
	if filePath == "" {
		return false
	}

	ext := filepath.Ext(filePath)

	if ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mmd" {
		return true
	}
	return false
}

func MakeFolder(sPath string) (bool, error) {
	if sPath == "" {
		return false, errors.New("MakeFolder: sPath is empty")
	}

	sFolderPath, errFolderPath := MakePath(sPath)

	if errFolderPath != nil {
		return false, errFolderPath
	}

	errFolderPath = os.Mkdir(sFolderPath, os.ModePerm)

	if errFolderPath != nil {
		return false, errFolderPath
	}

	return true, nil
}

func ReadImageAsBase64(imagePath string) (string, error) {

	var retImage string
	retImage = ""

	image, errRead := ioutil.ReadFile(imagePath)

	if errRead != nil {
		return "", errors.New("ReadImageAsBase64: Read Fail")
	}

	imageBase64 := base64.StdEncoding.EncodeToString(image)

	datatype, err2 := imgtype.Get(imagePath)
	if err2 != nil {
		return "", errors.New("ReadImageAsBase64: Cannot get image type")
	} else {
		retImage = "data:" + datatype + ";base64," + imageBase64
	}

	return retImage, nil
}

func PathIsImage(filePath string) bool {

	if filePath == "" {
		return false
	}

	_, err2 := imgtype.Get(filePath)
	if err2 != nil {
		return false
	}
	return true
}

func MakePath(sPath string) (string, error) {
	sfolder, sfile := filepath.Split(sPath)

	if sfolder == "" || sfile == "" {
		return "", errors.New("MakePath: folder or file name is empty folder " + sfolder + " file " + sfile)
	}

	sfolder = filepath.Clean(sfolder)

	if !PathIsExist(sfolder) {
		os.MkdirAll(sfolder, os.ModePerm)
	}

	return filepath.Join(sfolder, sfile), nil

}

func MakeSoftLink4Folder(srcFolder, linkFolder string) (bool, error) {
	srcExist := PathIsExist(srcFolder)

	if !srcExist {
		return false, errors.New("Make Soft Link 4 Folder: SrcFolder Not Exist " + srcFolder)
	}

	targetExist := PathIsExist(linkFolder)

	if targetExist {
		return false, errors.New("Make Soft Link 4 Folder:linkFolder Already Exist " + linkFolder)
	}

	sysType := runtime.GOOS

	var mkLinkCmd *exec.Cmd

	if "windows" == sysType {
		mkLinkCmd = exec.Command("cmd", "/c", "mklink /j "+linkFolder+"  "+srcFolder)
	} else if "linux" == sysType || "darwin" == sysType {
		mkLinkCmd = exec.Command("bash", "-c", "ln -s "+srcFolder+" "+linkFolder)
	} else { //Not support other platforms now
		return false, nil
	}

	_, err := mkLinkCmd.Output()

	if err != nil {
		return false, err
	}
	return true, nil
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		var errMsg = "CopyFile " + src + "is not a regular file"
		return 0, errors.New(errMsg)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func MoveFile(src, dst string) (int64, error) {
	iCopy, errCopy := CopyFile(src, dst)

	if errCopy != nil {
		return 0, errCopy
	}

	errRemove := os.Remove(src)

	if errRemove != nil {
		return 0, errRemove
	}

	return iCopy, nil
}

func DeleteFile(filePath string) bool {
	errRemove := os.Remove(filePath)

	if errRemove != nil {
		return false
	}
	return true
}

func Try2FindSpFile(siteFolderPath string) (string, error) {
	if PathIsExist(siteFolderPath) == false {
		return "", errors.New("Try2FindSpFile: Site Folder not exist")
	}

	var spCount int
	spCount = 0
	var spFileName string
	spFileName = ""

	files, _ := ioutil.ReadDir(siteFolderPath)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sp") {
			spFileName = f.Name()
			spCount++
			if spCount > 1 {
				return "", errors.New("Try2FindSpFile: More than 1 .sp file")
			}
		}
	}
	return spFileName, nil
}
