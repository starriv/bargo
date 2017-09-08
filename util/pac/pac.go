package pac

import (
	"bufio"
	"strings"
	"regexp"
	"path/filepath"
	"encoding/base64"
	"os"
	"time"
	"io/ioutil"
	"net/http"
)

// 数据源
const GFWLIST_URL = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

//  补充黑名单数据
var supRules = []string{
	"google.",
	"youtube.",
	"facebook.",
	"twitter.",
}

// 黑名单列表
var rules []string

func init()  {
	rules = getBlackRule()
	// 补充规则
	rules = append(rules, supRules...)
}

// 新增规则 竖线为分隔符
func AddRules(userRules string)  {
	newRules := strings.Split(userRules, "|")
	if len(newRules) > 0 {
		rules = append(rules, newRules...)
	}
}

// 是否包含在黑名单
func InBlack(domain string) bool {
	domain = strings.ToLower(domain)
	for _, v := range rules {
		v = strings.ToLower(v)
		if strings.Contains(domain, v) && len(v) > 0 {
			return true
		}
	}
	return false
}

// 获得黑名单列表
func getBlackRule() []string {
	pacFile, _ := filepath.Abs(os.TempDir()+"/bargo_pac.txt")
	// 缓存是否过期
	fileinfo, err := os.Stat(pacFile)
	if err != nil || fileinfo.ModTime().Add(24*7*time.Hour).Before(time.Now()) {
		err := updatePacFile(pacFile)
		if err != nil {
			panic(err)
		}
	}
	data, err := ioutil.ReadFile(pacFile)
	if err != nil {
		panic(err)
	}
	rules := strings.Split(string(data), "\n")

	return rules
}

// 更新缓存文件
func updatePacFile(pacFile string) error  {
	var gfwBase64String string
	// 远程获取新数据
	resp, err := http.Get(GFWLIST_URL)
	if err != nil {
		// 获取失败则读取默认的
		gfwBase64String = DEFAULT_GFWLIST
	} else {
		gfwBase64, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		gfwBase64String = string(gfwBase64)
	}

	// 打开缓存文件
	cachefile, err := os.OpenFile(pacFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cachefile.Close()
	// 解析gfwlist
	gfwRule, err := base64.StdEncoding.DecodeString(gfwBase64String)
	if err != nil {
		return err
	}
	gfwRuleReader := bufio.NewReader(strings.NewReader(string(gfwRule)))
	// 保存结果用于去重
	result := make(map[string]int)
	for {
		// 开始匹配每一行
		line, _, err := gfwRuleReader.ReadLine()
		if err != nil {
			break
		}
		lineStr := string(line)
		if len(line) == 0 {
			continue
		}
		reg := regexp.MustCompile(`[!\[].*?`)
		isComment := reg.Match(line)
		if isComment {
			continue
		}
		isWrite := strings.HasPrefix(lineStr, "@@")
		// 匹配域名和ip
		reg = regexp.MustCompile(`(?:(?:[a-zA-Z0-9\-]{1,61}\.)+[a-zA-Z]{2,6}|(?:\d{1,3}\.){3}\d{1,3})`)
		domain := reg.FindAllStringSubmatch(lineStr, 1)
		if !isWrite && len(domain) > 0 {
			if _, ok := result[domain[0][0]]; ok {
				continue
			}
			cachefile.WriteString(domain[0][0]+"\n")
			result[domain[0][0]] = 1
		}
	}

	return nil
}
