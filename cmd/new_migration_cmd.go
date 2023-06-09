package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go-rest-boilerplate/cmd/prompts"
	"go-rest-boilerplate/config"
	"go-rest-boilerplate/internal"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
)

const upMigrationTemplate = `-- vim: filetype=SQL
CREATE TABLE dummy (
	id TEXT PRIMARY KEY,
	ctime TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
const downMigrationTemplate = `-- vim: filetype=SQL
DROP TABLE dummy;
`

var newMigrationCmd = &cobra.Command{
	Use:   "new-migration (name)",
	Short: "Creates a new migration file with timestamps and the given name",
	RunE:  runNewMigrationCmd,
}

func runNewMigrationCmd(cmd *cobra.Command, args []string) (err error) {
	defer internal.WrapErr("new-migration", &err)

	cfg := config.MustConfigure()
	prompt := prompts.New(cfg, args)
	name := prompt.Str("name of migration")

	name = time.Now().Format("200601021504") + "_" + flect.Underscore(name)
	upname, downname := name+".up.sql", name+".down.sql"

	uppath, err := filepath.Abs(filepath.Join("./data/migrations", upname))
	if err != nil {
		return err
	}
	downpath, err := filepath.Abs(filepath.Join("./data/migrations", downname))
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, uppath)
	fmt.Fprintln(os.Stdout, downpath)
	if !prompt.YesNo("create these files") {
		log.Fatalln("aborted")
	}

	if err := os.WriteFile(uppath, []byte(upMigrationTemplate), 0644); err != nil {
		return err
	} else if err := os.WriteFile(downpath, []byte(downMigrationTemplate), 0644); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if strings.TrimSpace(editor) == "" {
		editor = "/usr/bin/vi"
	}

	proc := exec.Command(editor, uppath, downpath)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	return proc.Run()
}
