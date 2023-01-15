// Package controller coordinates all the other
// components in the `Lama2` project. The high
// level overview of command execution is easily
// understood from this package
package contoller

import (
	"fmt"
	"os"
	"path"

	"github.com/HexmosTech/gabs/v2"
	"github.com/HexmosTech/lama2/cmdexec"
	"github.com/HexmosTech/lama2/cmdgen"
	"github.com/HexmosTech/lama2/importer"
	"github.com/HexmosTech/lama2/lama2cmd"
	outputmanager "github.com/HexmosTech/lama2/outputManager"
	"github.com/HexmosTech/lama2/parser"
	"github.com/HexmosTech/lama2/preprocess"
	"github.com/HexmosTech/lama2/utils"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/rs/zerolog/log"
)

// Process initiates the following tasks in the given order:
// 1. Parse command line arguments
// 2. Read API file contents
// 3. Expand environment variables in API file
// 4. Parse the API contents
// 5. Generate API request command
// 6. Execute command & retrieve results
// 7. Optionally, post-process and write results to a JSON file
func Process(version string) {
	o := lama2cmd.GetAndValidateCmd(os.Args)
	if o.Version {
		fmt.Println(version)
		return
	}
	if o.Update {
		utils.UpdateSelf()
		return
	}
	if len(o.PostmanFile) > 0 {
		if len(o.LamaDir) > 0 {
			importer.PostmanImporter(o.PostmanFile, o.LamaDir)
			return
		}
		log.Fatal().Msg("To convert Postman export to Lama2, try: l2 -p PostmanFile -l Lama2Dir")
	}
	if len(o.LamaDir) > 0 {
		if len(o.PostmanFile) > 0 {
			importer.PostmanImporter(o.PostmanFile, o.LamaDir)
			return
		}
		log.Fatal().Msg("To convert Postman export to Lama2, try: l2 -p PostmanFile -l Lama2Dir")
	}
	apiContent := preprocess.GetLamaFileAsString(o.Positional.LamaAPIFile)
	_, dir, _ := utils.GetFilePathComponents(o.Positional.LamaAPIFile)
	preprocess.LoadElfEnv(path.Join(dir, "l2.env"))
	p := parser.NewLama2Parser()
	parsedAPI, e := p.Parse(apiContent)

	if e != nil {
		log.Fatal().
			Str("Type", "Controller").
			Str("LamaFile", o.Positional.LamaAPIFile).
			Str("Error", e.Error()).
			Msg(fmt.Sprint("Parse Error"))
	}
	log.Debug().Str("Parsed API", parsedAPI.String()).Msg("")
	fmt.Println("***")
	parsedAPIblocks := parsedAPI.S("value").Data().(*gabs.Container).Children()
	fmt.Println(parsedAPIblocks)
	fmt.Println("***")
	vm := goja.New()
	new(require.Registry).Enable(vm)
	console.Enable(vm)
	for i, block := range parsedAPIblocks {
		fmt.Println(">> ", i)
		fmt.Println(block)
		blockType := block.S("type").Data().(string)
		if blockType == "processor" {
			b := block.S("value").Data().(*gabs.Container)
			fmt.Println(b)
			script := b.Data().(string)
			vm.RunString(script)
		} else if blockType == "Lama2File" {
			b := block.S("url", "value").Data().(string)
			fmt.Println(b)
			url := preprocess.ExpandEnv(b, vm)
			e1 := block.Delete("url", "value")
			fmt.Println(block)
			fmt.Println(e1)
			_, e = block.Set(url, "url", "value")
			fmt.Println(e)
			fmt.Println(block)
			fmt.Println(url)
			// TODO - replace stuff in headers, and varjson and json as well
			cmd := cmdgen.ConstructCommand(block, o)
			retStr := cmdexec.ExecCommand(cmd, dir)
			fmt.Println("----------xxxxxxxxx-----------")
			fmt.Println(retStr)
			parsedOutput := outputmanager.RequestLogParser(retStr)
			poStr := parsedOutput.S("body").Data().(string)
			fmt.Println("----------zzzzzzzzz-----------")
			fmt.Println(poStr)

		}
	}
	return

	apiContent, apiDir := preprocess.LamaFile(o.Positional.LamaAPIFile)
	cmdStr := cmdgen.ConstructCommand(parsedAPI, o)
	log.Info().Msg("COMMAND:\n" + cmdStr)
	op := cmdexec.ExecCommand(cmdStr, apiDir)
	if o.Output != "" {
		outputmanager.WriteJSONOutput(op, o.Output)
	}
}
