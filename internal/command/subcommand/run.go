/*
Copyright 2020 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package subcommand

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	runTmp     = "/tmp/quanxiang/lowcode/faas/run"
	pluginPath = ".quanxiang/faas/go116"
)

type Run struct {
	// function name.
	name string
	// group name.
	group string
	// handler dir path.
	src string
}

func newRun() *Run {
	return &Run{}
}

func NewRun() *cobra.Command {
	o := newRun()

	cmd := &cobra.Command{
		Use:                   "run",
		DisableFlagsInUseLine: true,
		Short:                 "",
		Long: `
`,
		Example: ``,

		Run: func(cmd *cobra.Command, args []string) {
			CheckErr(o.Complete(cmd, args))
			CheckErr(o.Validate(cmd))
			CheckErr(o.Run(cmd, args))
		},
	}

	cmd.Flags().StringVar(&o.src, "src", ".", "hander dir path")
	cmd.Flags().StringVar(&o.group, "group", "", "group name")

	return cmd
}

func (r *Run) Complete(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		r.name = args[0]
	}

	return nil
}

func (r *Run) Validate(cmd *cobra.Command) error {
	if r.name == "" {
		return UsageErrorf(cmd, "a name is required")
	}

	if r.group == "" {
		return UsageErrorf(cmd, "group name is required")
	}
	return nil
}

func (r *Run) Run(cmd *cobra.Command, args []string) error {
	// temporary space for caching
	cachePath := filepath.Join(runTmp, fmt.Sprintf("%s-%s", r.group, r.name))

	l := &GO{
		path:    cachePath,
		handler: r.src,
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		if err := r.initProject(cmd, l, cachePath); err != nil {
			return err
		}
	}

	// update dependencies
	if err := l.tidy(); err != nil {
		return UsageErrorf(cmd, "update dependencies fail: %v", err)
	}

	if err := l.run(); err != nil {
		return UsageErrorf(cmd, "run fail: %v", err)
	}

	return nil
}

func (r *Run) initProject(cmd *cobra.Command, l language, path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return UsageErrorf(cmd, "creating %s: %v", path, err)
	}

	// TODO multilingual support
	if err := l.initProject(); err != nil {
		return UsageErrorf(cmd, "init go project: %v", err)
	}

	return nil
}

type language interface {
	initProject() error
	tidy() error
	run() error
}

type GO struct {
	path    string
	handler string
}

func (g *GO) initProject() error {
	if err := g.templateExecute(template.Must(template.New("main").Parse(mainGOTemplate)), filepath.Join(g.path, "main.go"), nil); err != nil {
		return err
	}

	if err := g.templateExecute(template.Must(template.New("mod").Parse(goModTemplate)), filepath.Join(g.path, "go.mod"), map[string]string{
		"Handler": g.handler,
	}); err != nil {
		return err
	}

	// FIXME
	// copy plugins to workspace
	user, err := user.Current()
	if err != nil {
		return err
	}
	out := &bytes.Buffer{}
	cmd := exec.Command(
		"cp",
		"-r",
		filepath.Join(user.HomeDir, pluginPath)+"/",
		".",
	)

	cmd.Dir = g.path
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd.Run()
}

func (g *GO) tidy() error {
	out := &bytes.Buffer{}
	cmd := exec.Command(
		"go",
		"mod",
		"tidy",
		"-v",
	)

	cmd.Dir = g.path
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		fmt.Println(out.String())
		return err
	}

	fmt.Println(out.String())
	return nil
}

func (g *GO) run() error {
	out := &bytes.Buffer{}
	cmd := exec.Command(
		"go",
		"run",
		".",
	)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, []string{
		"FUNC_CLEAR_SOURCE=true",
		"FUNC_NAME=HelloWorld",
		"FUNC_CONTEXT={\"name\":\"demo13\",\"version\":\"v2.0.0\",\"runtime\":\"Knative\",\"port\":\"8080\",\"prePlugins\":[\"plugin-quanxiang-lowcode-client\"]}",
		"POD_NAME=xx",
		"POD_NAMESPACE=serving",
	}...)

	cmd.Dir = g.path
	cmd.Stdout = out
	cmd.Stderr = out

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func(cmd *exec.Cmd, c <-chan os.Signal) error {
		for {
			s := <-c
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			default:
			}

			return cmd.Process.Kill()
		}

	}(cmd, c)

	fmt.Println("Running...")
	if err := cmd.Run(); err != nil {
		fmt.Println(out.String())
		return err
	}

	fmt.Println(out.String())
	return nil
}

func (g *GO) templateExecute(t *template.Template, fileName string, data interface{}) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}
	return nil
}

const (
	mainGOTemplate = `package main

import (
	"context"

	h "quanxiang-cloud/faas/handler"

	pluginlcustom "quanxiang-cloud/faas/plugins/plugin-quanxiang-lowcode-client"
	"github.com/OpenFunction/functions-framework-go/framework"
	"github.com/OpenFunction/functions-framework-go/plugin"
	"k8s.io/klog"
)

func main() {
	ctx := context.Background()
	fwk, err := framework.NewFramework()
	if err != nil {
		klog.Exit(err)
	}
	fwk.RegisterPlugins(getLocalPlugins())
	if err := fwk.Register(ctx, h.Handler); err != nil {
		klog.Exit(err)
	}
	if err := fwk.Start(ctx); err != nil {
		klog.Exit(err)
	}
}

func getLocalPlugins() map[string]plugin.Plugin {
	localPlugins := map[string]plugin.Plugin{
		pluginlcustom.Name: pluginlcustom.New(),
	}

	if len(localPlugins) == 0 {
		return nil
	} else {
		return localPlugins
	}
}

`

	goModTemplate = `module quanxiang-cloud/faas

go 1.16

replace quanxiang-cloud/faas/handler => {{ .Handler }}
replace github.com/quanxiang-cloud/faas-lowcode => ./pkg/faas-lowcode
`
)
