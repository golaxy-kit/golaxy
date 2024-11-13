/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func genEventTab() {
	declFile := viper.GetString("decl_file")
	packageEventAlias := viper.GetString("package_event_alias")
	pkg := viper.GetString("package")
	dir := viper.GetString("dir")
	tabName := viper.GetString("name")

	// 解析事件定义
	eventDeclTab := EventDeclTab{}
	eventDeclTab.Parse()

	code := &bytes.Buffer{}

	// 生成注释
	{
		program := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
		args := strings.Join(os.Args[1:], " ")
		copyright := copyrightNotice

		if !viper.GetBool("copyright") {
			copyright = ""
		}

		fmt.Fprintf(code, `%s// Code generated by %s %s; DO NOT EDIT.

package %s
`, copyright, program, args, pkg)
	}

	// 生成import
	{
		importCode := &bytes.Buffer{}

		fmt.Fprintf(importCode, "\nimport (")

		fmt.Fprintf(importCode, `
	%s "%s"`, packageEventAlias, packageEventPath)

		fmt.Fprintf(importCode, "\n)\n")

		fmt.Fprintf(code, importCode.String())
	}

	// event包前缀
	eventPrefix := ""
	if packageEventAlias != "." {
		eventPrefix = packageEventAlias + "."
	}

	// 生成事件表接口
	{
		var eventsCode string

		for _, event := range eventDeclTab.Events {
			eventsCode += fmt.Sprintf("\t%s() %sIEvent\n", event.Name, eventPrefix)
		}

		fmt.Fprintf(code, `
type I%[1]s interface {
%[2]s}
`, strings.Title(tabName), eventsCode)
	}

	// 生成事件表
	{
		var eventsRecursionCode string

		for i, event := range eventDeclTab.Events {
			eventRecursion := "recursion"

			// 解析atti
			atti := parseGenAtti(event.Comment, "+event-tab-gen:")

			if atti.Has("recursion") {
				switch atti.Get("recursion") {
				case "allow":
					eventRecursion = eventPrefix + "EventRecursion_Allow"
				case "disallow":
					eventRecursion = eventPrefix + "EventRecursion_Disallow"
				case "discard":
					eventRecursion = eventPrefix + "EventRecursion_Discard"
				case "truncate":
					eventRecursion = eventPrefix + "EventRecursion_Truncate"
				case "deepest":
					eventRecursion = eventPrefix + "EventRecursion_Deepest"
				}
			}

			eventsRecursionCode += fmt.Sprintf("\t(*eventTab)[%d].Init(autoRecover, reportError, %s)\n", i, eventRecursion)
		}

		// 生成事件Id
		{
			fmt.Fprintln(code, `
var (`)
			fmt.Fprintf(code, `	_%[1]sId = %[2]sDeclareEventTabIdT[%[1]s]()
`, tabName, eventPrefix)

			for i, event := range eventDeclTab.Events {
				fmt.Fprintf(code, `	%[2]sId = _%[1]sId + %[3]d
`, tabName, event.Name, i)
			}

			fmt.Fprintln(code, ")")
		}

		fmt.Fprintf(code, `
type %[1]s [%[2]d]%[4]sEvent

func (eventTab *%[1]s) Init(autoRecover bool, reportError chan error, recursion %[4]sEventRecursion) {
%[3]s}

func (eventTab *%[1]s) Event(id uint64) %[4]sIEvent {
	if _%[1]sId != id & 0xFFFFFFFF00000000 {
		return nil
	}
	pos := id & 0xFFFFFFFF
	if pos >= uint64(len(*eventTab)) {
		return nil
	}
	return &(*eventTab)[pos]
}

func (eventTab *%[1]s) Open() {
	for i := range *eventTab {
		(*eventTab)[i].Open()
	}
}

func (eventTab *%[1]s) Close() {
	for i := range *eventTab {
		(*eventTab)[i].Close()
	}
}

func (eventTab *%[1]s) Clean() {
	for i := range *eventTab {
		(*eventTab)[i].Clean()
	}
}
`, tabName, len(eventDeclTab.Events), eventsRecursionCode, eventPrefix)
	}

	for i, event := range eventDeclTab.Events {
		fmt.Fprintf(code, `
func (eventTab *%[1]s) %[2]s() %[4]sIEvent {
	return &(*eventTab)[%[3]d]
}
`, tabName, event.Name, i, eventPrefix)
	}

	log.Printf("EventTab: %s", tabName)

	// 输出文件
	outFile := filepath.Join(dir, filepath.Base(strings.TrimSuffix(declFile, ".go"))+".tab.gen.go")

	os.MkdirAll(filepath.Dir(outFile), os.ModePerm)

	if err := ioutil.WriteFile(outFile, code.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}
