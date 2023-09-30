package prettify

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/HexmosTech/gabs/v2"
	"github.com/rs/zerolog/log"
)

func Prettify(parsedAPI *gabs.Container, context map[string]bool, markRange map[string]int, content string, fPath string) {
	defer func() {
		if err := recover(); err != nil { //catch
			log.Debug().Msg("Potential issue with prettify")
			// os.Exit(0)
		}
	}()
	parsedAPIblocks := parsedAPI.S("value").Data().(*gabs.Container).Children()
	// Prettification procedure:
	// 1. Scan from bottom of contents
	// 2. For every JSON fragment:
	// 		2.1. reformat
	//      2.2. replace
	//      2.3. write output
	l := len(parsedAPIblocks)
	markMax := (len(markRange) / 2) - 1
	for i := l - 1; i >= 0; i-- {
		block := parsedAPIblocks[i]
		// vm := cmdexec.GetJSVm()
		blockType := block.S("type").Data().(string)
		if blockType == "processor" {
			fmt.Println("Skipping processor block")
		} else if blockType == "Lama2File" {
			// preprocess.ProcessVarsInBlock(block, vm)
			isJSON, ok := context["JSON"]
			if !ok {
				return
			}

			if !isJSON {
				return
			}

			fmt.Println("Prettifying JSON in the l2 file")
			jsonObj := block.S("details", "ip_data")

			idxStr := strconv.Itoa(markMax)

			content = content[:markRange["DataStart"+idxStr]] + jsonObj.StringIndent("", "  ") + "\n" + content[markRange["DataEnd"+idxStr]:]
			markMax -= 1
		}

		var re = regexp.MustCompile(`(?m)"626f4c60-([^"]+)"`)
		var substitution = "$1"
		res2 := re.ReplaceAllString(content, substitution)

		os.WriteFile(fPath, []byte(res2), 0644)
	}
}
