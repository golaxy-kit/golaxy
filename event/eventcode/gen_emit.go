package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func genEmit(ctx *_CommandContext) {
	emitFile := ctx.EmitDir

	if emitFile == "" {
		emitFile = strings.TrimSuffix(ctx.DeclFile, ".go") + "_emit_code.go"
	} else {
		emitFile = strings.Join([]string{filepath.Dir(ctx.DeclFile), ctx.EmitDir, filepath.Base(strings.TrimSuffix(ctx.DeclFile, ".go")) + "_emit_code.go"}, string(filepath.Separator))
	}

	emitCode := &bytes.Buffer{}

	// 生成注释
	{
		program := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
		args := strings.Join(os.Args[1:], " ")

		fmt.Fprintf(emitCode, `// Code generated by %s %s; DO NOT EDIT.

package %s
`, program, args, ctx.EmitPackage)
	}

	// 生成import
	{
		importCode := &bytes.Buffer{}

		fmt.Fprintf(importCode, "\nimport (")

		fmt.Fprintf(importCode, `
	"fmt"
	%s "%s"
	%s "%s"`, ctx.PackageEventAlias, packageEventPath, ctx.PackageIfaceAlias, packageIfacePath)

		for _, imp := range ctx.FileAst.Imports {
			begin := ctx.FileSet.Position(imp.Pos())
			end := ctx.FileSet.Position(imp.End())

			impStr := string(ctx.FileData[begin.Offset:end.Offset])

			switch imp.Path.Value {
			case fmt.Sprintf(`"%s"`, packageEventPath):
				if imp.Name == nil {
					if ctx.PackageEventAlias == "event" {
						continue
					}
				} else {
					if imp.Name.Name == ctx.PackageEventAlias {
						continue
					}
				}
			case fmt.Sprintf(`"%s"`, packageIfacePath):
				if imp.Name == nil {
					if ctx.PackageIfaceAlias == "iface" {
						continue
					}
				} else {
					if imp.Name.Name == ctx.PackageIfaceAlias {
						continue
					}
				}
			}

			fmt.Fprintf(importCode, "\n\t%s", impStr)
		}

		fmt.Fprintf(importCode, "\n)\n")

		fmt.Fprintf(emitCode, importCode.String())
	}

	// 解析事件定义
	eventDeclTab := _EventDeclTab{}
	eventDeclTab.Parse(ctx)

	// event包前缀
	eventPrefix := ""
	if ctx.PackageEventAlias != "." {
		eventPrefix = ctx.PackageEventAlias + "."
	}

	// iface包前缀
	ifacePrefix := ""
	if ctx.PackageIfaceAlias != "." {
		ifacePrefix = ctx.PackageIfaceAlias + "."
	}

	// 生成事件发送代码
	for _, eventDecl := range eventDeclTab {
		// 是否导出事件发送代码
		exportEmitStr := "emit"
		if ctx.EmitDefExport {
			exportEmitStr = "Emit"
		}

		if strings.Contains(eventDecl.Comment, "[EmitExport]") {
			exportEmitStr = "Emit"
		} else if strings.Contains(eventDecl.Comment, "[EmitUnExport]") {
			exportEmitStr = "emit"
		}

		auto := ctx.EmitDefAuto

		if strings.Contains(eventDecl.Comment, "[EmitAuto]") {
			auto = true
		} else if strings.Contains(eventDecl.Comment, "[EmitManual]") {
			auto = false
		}

		// 生成代码
		if auto {
			if eventDecl.FuncHasRet {
				fmt.Fprintf(emitCode, `
type iAuto%[1]s interface {
	%[1]s() %[6]sIEvent
}

func Bind%[1]s(auto iAuto%[1]s, delegate %[2]s%[8]s) %[6]sHook {
	if auto == nil {
		panic(fmt.Errorf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	return %[6]sBindEvent[%[2]s%[8]s](auto.%[1]s(), delegate)
}

func Bind%[1]sWithPriority(auto iAuto%[1]s, delegate %[2]s%[8]s, priority int32) %[6]sHook {
	if auto == nil {
		panic(fmt.Errorf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	return %[6]sBindEventWithPriority[%[2]s%[8]s](auto.%[1]s(), delegate, priority)
}

func %[9]s%[1]s%[7]s(auto iAuto%[1]s%[4]s) {
	if auto == nil {
		panic(fmt.Errorf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(delegate %[10]sCache) bool {
		return %[10]sCache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, ifacePrefix)

			} else {
				fmt.Fprintf(emitCode, `
type iAuto%[1]s interface {
	%[1]s() %[6]sIEvent
}

func Bind%[1]s(auto iAuto%[1]s, delegate %[2]s%[8]s) %[6]sHook {
	if auto == nil {
		panic(fmt.Errorf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	return %[6]sBindEvent[%[2]s%[8]s](auto.%[1]s(), delegate)
}

func Bind%[1]sWithPriority(auto iAuto%[1]s, delegate %[2]s%[8]s, priority int32) %[6]sHook {
	if auto == nil {
		panic(fmt.Errorf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	return %[6]sBindEventWithPriority[%[2]s%[8]s](auto.%[1]s(), delegate, priority)
}

func %[9]s%[1]s%[7]s(auto iAuto%[1]s%[4]s) {
	if auto == nil {
		panic(fmt.Errorf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(delegate %[10]sCache) bool {
		%[10]sCache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, ifacePrefix)
			}
		} else {
			if eventDecl.FuncHasRet {
				fmt.Fprintf(emitCode, `
func %[9]s%[1]s%[7]s(evt %[6]sIEvent%[4]s) {
	if evt == nil {
		panic(fmt.Errorf("%%w: %%w: evt is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	%[6]sUnsafeEvent(evt).Emit(func(delegate %[10]sCache) bool {
		return %[10]sCache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, ifacePrefix)

			} else {
				fmt.Fprintf(emitCode, `
func %[9]s%[1]s%[7]s(evt %[6]sIEvent%[4]s) {
	if evt == nil {
		panic(fmt.Errorf("%%w: %%w: evt is nil", %[6]sErrEvent, %[6]sErrArgs))
	}
	%[6]sUnsafeEvent(evt).Emit(func(delegate %[10]sCache) bool {
		%[10]sCache2Iface[%[2]s%[8]s](delegate).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, ifacePrefix)
			}
		}

		fmt.Printf("Emit: %s\n", eventDecl.Name)
	}

	os.MkdirAll(filepath.Dir(emitFile), os.ModePerm)

	if err := ioutil.WriteFile(emitFile, emitCode.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}