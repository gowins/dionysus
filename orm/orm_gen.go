package orm

import (
	"fmt"

	"github.com/gowins/dionysus/cmd"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	defaultOutPath = "./common/db/repository"
	defaultPkgPath = "./common/db/model"
)

// GenModel generate data necessary for QueryStructMeta
type GenModel struct {
	TableName string
	ModelName string
	ModelOpts []gen.ModelOpt
}

// OrmGenCmd generate orm model and repository command
type OrmGenCmd struct {
	DataTyMap map[string]func(string) string
	ModelOpts []gen.ModelOpt
	GenModels []GenModel
	Dsn       DsnInfo
	Cfg       gen.Config
}

// GetCmd implement GetCmd method of cmd.Commander
func (og *OrmGenCmd) GetCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "gorm_gen",
		Short: "generate gorm model and repository",
	}
	c.RunE = func(*cobra.Command, []string) error {
		return og.genGorm()
	}
	return c
}

// GetShutdownFunc implement GetShutdownFunc method of cmd.Commander
func (og *OrmGenCmd) GetShutdownFunc() cmd.StopFunc {
	return func() {}
}

// RegShutdownFunc implement RegShutdownFunc method of cmd.Commander
func (og *OrmGenCmd) RegShutdownFunc(stopSteps ...cmd.StopStep) {}

// genGorm call gen package to generate model and repository
func (og *OrmGenCmd) genGorm() error {
	if len(og.GenModels) == 0 {
		return fmt.Errorf("generate models is empty")
	}
	db, err := gorm.Open(og.Dsn.Dialector())
	if err != nil {
		return err
	}
	if og.Cfg.OutPath == "" {
		og.Cfg.OutPath = defaultOutPath
	}
	if og.Cfg.ModelPkgPath == "" {
		og.Cfg.ModelPkgPath = defaultPkgPath
	}
	g := gen.NewGenerator(og.Cfg)
	g.UseDB(db)
	if og.DataTyMap != nil {
		g.WithDataTypeMap(og.DataTyMap)
	}
	g.WithOpts(og.ModelOpts...)
	applies := make([]any, len(og.GenModels))
	for i, gm := range og.GenModels {
		if gm.ModelName == "" {
			applies[i] = g.GenerateModel(gm.TableName, gm.ModelOpts...)
		} else {
			applies[i] = g.GenerateModelAs(gm.TableName, gm.ModelName, gm.ModelOpts...)
		}
	}
	g.ApplyBasic(applies...)
	g.Execute()
	return nil
}
