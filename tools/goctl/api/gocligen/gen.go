package gocligen

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/core/logx"
	apiformat "github.com/zeromicro/go-zero/tools/goctl/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	tmpDir = path.Join(os.TempDir(), "goctl")
	// VarStringDir describes the directory.
	VarStringDir string
	// VarStringAPI describes the API.
	VarStringAPI string
	// VarStringHome describes the go home.
	VarStringHome string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
)

// GoCliCommand gen go cli project files from command line
func GoCliCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	dir := VarStringDir
	home := VarStringHome
	namingStyle := VarStringStyle

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	return DoGenProject(apiFile, dir, namingStyle)
}

// DoGenProject gen go project files with api file
func DoGenProject(apiFile, dir, style string) error {
	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	logx.Must(pathx.MkdirIfNotExist(dir))
	rootPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	logx.Must(genTypes(dir, cfg, api))
	logx.Must(genHandleResponse(dir, cfg))
	logx.Must(genLogic(dir, rootPkg, cfg, api))

	//if err := backupAndSweep(apiFile); err != nil {
	//	return err
	//}

	if err := apiformat.ApiFormatByPath(apiFile, false); err != nil {
		return err
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}
