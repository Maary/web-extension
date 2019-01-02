package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
)

type ConfigInfo struct {
	ConfigName           string   `json:"-"`
	ProName              string   `json:"name"`
	ProDescription       string   `json:"description"`
	ProPath              string   `json:"path"`
	ProType              string   `json:"type"`
	ProAllowedExtensions []string `json:"allowed_extensions"`
}

func (c *ConfigInfo) toJsonStr() (string, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (c *ConfigInfo) Name(name string) *ConfigInfo {
	c.ConfigName = name
	return c
}

const _OSX_FIREFOX_STORAGE_PATH_FORMAT_ = `%s/Library/Application Support/Mozilla/NativeMessagingHosts/%s.json`
const _LINUX_FIREFOX_STORAGE_PATH_FORMAT = `%s/.mozilla/native-messaging-hosts/%s.json`
const _WINDOWS_FIREFOX_STORAGE_PATH_FORMAT_ = `%s\Add-on\%s.json`
const _KEY_OF_ADD_ON_ = `Software\Mozilla\NativeMessagingHosts`
const _KEY_OF_ADD_ON_FORMAT = `Software\Mozilla\NativeMessagingHosts\%s`

const _OSX_APP_PATH_FORMAT_ = `%s/Applications/%s`
const _LINUX_APP_PATH_FORMAT_ = `%s%s`
const _WINDOWS_APP_PATH_FORMAT_ = `%s\Add-on\%s.exe`
const _WINDOWS_APP_DIR_ = `%s\Add-on`

const _LINUX_BINARY_FILE_PATH_ = `./%s`
const _OSX_BINARY_FILE_PATH_ = `./%s`
const _WIN_BINARY_FILE_PATH_ = `./%s.exe`

func (c *ConfigInfo) CreateConfig() (bool, error) {
	ur, err := user.Current()
	if err != nil {
		return false, err
	}
	if c.ConfigName == "" {
		log.Fatalln("config must have a name")
	}
	var (
		fileName, appPath, binaryFilePath string
	)
	switch runtime.GOOS {
	case "darwin":
		fileName = fmt.Sprintf(_OSX_FIREFOX_STORAGE_PATH_FORMAT_, ur.HomeDir, c.ConfigName)
		appPath = fmt.Sprintf(_OSX_APP_PATH_FORMAT_, ur.HomeDir, c.ConfigName)
		fmt.Println("当前安装路径为：", appPath)
		binaryFilePath = fmt.Sprintf(_OSX_BINARY_FILE_PATH_, c.ConfigName)
	case "linux":
		fileName = fmt.Sprintf(_LINUX_FIREFOX_STORAGE_PATH_FORMAT, ur.HomeDir, c.ConfigName)
		appPath = fmt.Sprintf(_LINUX_APP_PATH_FORMAT_, ur.HomeDir, c.ConfigName)
		fmt.Println("当前安装路径为：", appPath)
		binaryFilePath = fmt.Sprintf(_LINUX_BINARY_FILE_PATH_, c.ConfigName)
	case "windows":
		fileName = fmt.Sprintf(_WINDOWS_FIREFOX_STORAGE_PATH_FORMAT_, ur.HomeDir, c.ConfigName)
		appPath = fmt.Sprintf(_WINDOWS_APP_PATH_FORMAT_, ur.HomeDir, c.ConfigName)
		path := fmt.Sprintf(_WINDOWS_APP_DIR_, ur.HomeDir)
		fmt.Println("当前安装路径为： ", path)
		binaryFilePath = fmt.Sprintf(_WIN_BINARY_FILE_PATH_, c.ConfigName)
		fmt.Println("为安装程序创建文件目录...")
		err := WinMkAppDir(path)
		if err != nil {
			log.Fatalln(err.Error())
		}
		/*
		ok, err := windowsCreateRegistryForAddOn(c.ConfigName, path)
		if err != nil {
			log.Fatal("write registry failed, pleas check permission")
		}
		if ok {
			log.Println("注册成功")
		}
		*/
	}

	c.ProName = c.ConfigName
	c.ProPath = appPath
	c.ProType = "stdio"
	c.ProDescription = "native message"
	c.ProAllowedExtensions = []string{"starrymanjasper@gmail.com"}

	fp, err := os.Create(fileName)
	defer func() {
		err := fp.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		fmt.Println(err)
	}
	content, err := c.toJsonStr()
	if err != nil {
		log.Fatalln("config file error")
	}
	_, err = fp.Write([]byte(content))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("复制应用程序到用户目录")
	return copyBinaryFile(binaryFilePath, appPath)
}

func copyBinaryFile(source, target string) (bool, error) {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(target, input, 0777)
	if err != nil {
		return false, err
	}
	return true, nil
}

func WinMkAppDir(path string) error {
	ok, err := isExist(path)
	if ok {
		return err
	}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func isExist(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

/*
func windowsCreateRegistryForAddOn(appName, path string) (bool, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, _KEY_OF_ADD_ON_, registry.ALL_ACCESS)
	_, _, err = k.GetStringsValue(appName)
	if err != nil {
		n, x, err := registry.CreateKey(registry.LOCAL_MACHINE, fmt.Sprintf(_KEY_OF_ADD_ON_FORMAT, appName), registry.ALL_ACCESS)
		defer func() {
			err = n.Close()
			if err != nil {
				fmt.Println("写入注册表失败， 请重新打开安装程序。")
			}
		}()
		err = n.SetStringValue("", path+`\`+appName+".json")
		if err != nil {
			return false, err
		}
		if err != nil {
			log.Fatal("打开注册表失败")
		}
		if x == false {
			fmt.Println("正在为火狐浏览器新加注册表项")
		} else {
			fmt.Println("正在为火狐浏览器修改注册表项")
		}
		err = n.SetStringValue("", path+`\`+appName+".json")

		if err != nil {
			return false, err
		}
	} else {
		if err := k.SetStringValue(appName, path); err != nil {
			return false, err
		}
	}
	return true, nil
}
*/
