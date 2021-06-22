package template

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/txzdream/csp-translator/parser"
)

type procInfo struct {
	Exist   bool
	Name    string
	IsEvent bool
}

type templateData struct {
	IsEvent      bool
	IsCatchEvent bool
	FunctionName string
	Next         procInfo
	End          procInfo
	Init         []string
	Enable       []string
}

type Template struct {
	Init  []string
	Procs []templateData
}

func getNonEventName(id int) string {
	return "proc_" + strconv.Itoa(id)
}

func Build(data *parser.CSPData, modelPath string) (string, error) {
	tData := Template{
		Procs: make([]templateData, 0),
	}
	for _, initEvent := range data.InitProcs {
		_, isEvent := data.Event[initEvent]
		name := getNonEventName(initEvent)
		if isEvent {
			name = data.Procs[initEvent]
		}
		tData.Init = append(tData.Init, name)
	}
	for proc, name := range data.Procs {
		tmp := templateData{
			FunctionName: name,
		}

		_, isEvent := data.Event[proc]
		tmp.IsEvent = isEvent
		if !isEvent {
			tmp.FunctionName = getNonEventName(proc)
		}

		if tmp.IsEvent {
			_, isCatchEvent := data.CatchEvents[proc]
			tmp.IsCatchEvent = isCatchEvent
		}

		if enables, ok := data.Enable[proc]; ok {
			enablesStr := make([]string, 0)
			for _, enable := range enables {
				enableStr := data.Procs[enable]
				enablesStr = append(enablesStr, enableStr)
			}
			tmp.Enable = enablesStr
		}

		if next, ok := data.Next[proc]; ok {
			tmp.Next.Exist = true
			tmp.Next.Name = data.Procs[proc]
			_, isEvent := data.Event[next]
			tmp.Next.IsEvent = isEvent
			if !isEvent {
				tmp.Next.Name = getNonEventName(next)
			}
		}

		if end, ok := data.End[proc]; ok {
			tmp.End.Exist = true
			tmp.End.Name = data.Procs[proc]
			_, isEvent := data.Event[end]
			tmp.End.IsEvent = isEvent
			if !isEvent {
				tmp.End.Name = getNonEventName(end)
			}
		}

		if init, ok := data.Init[proc]; ok {
			for k := range init {
				name := data.Procs[k]
				if _, isEvent := data.Event[k]; !isEvent {
					name = getNonEventName(k)
				}
				tmp.Init = append(tmp.Init, name)
			}
		}

		tData.Procs = append(tData.Procs, tmp)
	}

	tmpl, err := template.ParseFiles(modelPath)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, tData); err != nil {
		return "", err
	}
	return buff.String(), nil
}
