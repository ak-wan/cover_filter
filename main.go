package main

import (
	"bufio"
	"flag"
	"fmt"
	"lpm/tool"
	"lpm/tool/logs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

var args Args

// Args 覆盖率文件参数
type Args struct {
	content []byte
	count   int
	markers []string
	orgPath string
	outPath string
}

// MarkPosition 标记信息，key：文件名(string)，value：标记所在文件的行号([]int)
type MarkPosition map[string][]int

func defaultPrint() {
	fmt.Println("" +
		`     ______                            ______ _  __ __					` + "\n" +
		`    / ____/____  _   __ ___   _____   / ____/(_)/ // /_ ___   _____	` + "\n" +
		`   / /    / __ \| | / // _ \ / ___/  / /_   / // // __// _ \ / ___/	` + "\n" +
		`  / /___ / /_/ /| |/ //  __// /     / __/  / // // /_ /  __// /		` + "\n" +
		`  \____/ \____/ |___/ \___//_/     /_/    /_//_/ \__/ \___//_/			` + "\n")
	fmt.Fprintf(os.Stderr, "coverage_filter ———— 覆盖率报告修改工具，根据标记过滤代码块（最小单位：function）\n")
	flag.PrintDefaults()
}

// 初始化命令参数
func init() {
	var marker string
	flag.StringVar(&marker, "marker", "no-cover", "设置过滤标记，可设置多个标签，以逗号(,)分隔")
	flag.StringVar(&args.orgPath, "file", "result.out", "通过go test生成覆盖率的源文件文件")
	flag.StringVar(&args.outPath, "out", args.orgPath+".coverfilter", "设置输出的文件名")
	flag.IntVar(&args.count, "count", 1, "指定注释标记块被统计次数，负数表示不统计改代码块")
	flag.Parse()

	args.markers = strings.Split(marker, ",")
	flag.Usage = defaultPrint
}

func main() {
	args.content = []byte("mode: set\n")
	profileList := args.AllProfiles()
	gopath := os.Getenv("GOPATH")
	for _, profile := range profileList {
		fileName := filepath.Join(gopath, "src", profile.FileName)
		args.filterLine(fileName, profile)
		args.content = append(args.content, formatProfile(profile)...)
	}
	tool.WriteFile(args.outPath, args.content)
}

// AllProfiles 根据覆盖率报告，获取每个源文件对应的覆盖率描述
func (a *Args) AllProfiles() []*cover.Profile {
	coverProfiles, err := cover.ParseProfiles(a.orgPath)
	if err != nil {
		logs.Errorln("Failed to parse cover profile file.", err)
		os.Exit(1)
	}
	return coverProfiles
}

// 过滤标记行
func (a *Args) filterLine(path string, profile *cover.Profile) {
	markerLines := a.getMarkerPosition(path)
	for i, block := range profile.Blocks {
		if isWithin(block, markerLines) && block.Count == 0 {
			profile.Blocks[i].Count = a.count
		}
	}
}

func isWithin(src cover.ProfileBlock, lines []int) bool {
	for _, line := range lines {
		if src.StartLine <= line && src.EndLine >= line {
			return true
		}
	}
	return false
}

func formatProfile(cov *cover.Profile) []byte {
	var out string
	for _, b := range cov.Blocks {
		if b.Count >= 0 {
			out += fmt.Sprintf("%s:%d.%d,%d.%d %d %d\n", cov.FileName, b.StartLine, b.StartCol, b.EndLine, b.EndCol, b.NumStmt, b.Count)
		}
	}
	return []byte(out)
}

// 获取文件中所有marker标记的行号
func (a *Args) getMarkerPosition(filePath string) (row []int) {
	file, err := os.Open(filePath)
	if err != nil {
		logs.Errorln("Failed to open file. FileName:", filePath, "Error:", err)
		return
	}
	defer file.Close()

	line := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, m := range a.markers {
			if strings.Contains(scanner.Text(), m) {
				row = append(row, line)
				break
			}
		}
		line++
	}
	return
}
