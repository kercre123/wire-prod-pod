package wirepod

import (
	"os"
	"plugin"
	"strings"
)

var PluginList []*plugin.Plugin
var PluginUtterances []*[]string
var PluginFunctions []func(string, string) string
var PluginNames []string

func LoadPlugins() {
	logger("Loading plugins")
	entries, err := os.ReadDir("./plugins")
	if err != nil {
		logger("Not loading plugins because the plugins directory doesn't exist")
	}
	for _, file := range entries {
		if strings.Contains(file.Name(), ".so") {
			plugin, err := plugin.Open("./plugins/" + file.Name())
			if err != nil {
				logger("Error loading plugin: " + file.Name())
				logger(err)
				continue
			} else {
				logger("Loading plugin: " + file.Name())
			}
			u, err := plugin.Lookup("Utterances")
			if err != nil {
				logger("Error loading Utterances []string from plugin file " + file.Name())
				logger(err)
				continue
			} else {
				if _, ok := u.(*[]string); ok {
					logger("Utterances []string in plugin " + file.Name() + " are OK")
				} else {
					logger("Error: Utterances in plugin " + file.Name() + " are not of type []string")
					continue
				}
			}
			a, err := plugin.Lookup("Action")
			if err != nil {
				logger("Error loading Action func from plugin file " + file.Name())
				continue
			} else {
				if _, ok := a.(func(string, string) string); ok {
					logger("Action func in plugin " + file.Name() + " is OK")
				} else {
					logger("Error: Action func in plugin " + file.Name() + " is not of type func(string, string) string")
					continue
				}
			}
			n, err := plugin.Lookup("Name")
			if err != nil {
				logger("Error loading Name string from plugin file " + file.Name())
				continue
			} else {
				if _, ok := n.(*string); ok {
					logger("Name string in plugin " + *n.(*string) + " is OK")
				} else {
					logger("Error: Name string in plugin " + file.Name() + " is not of type string")
					continue
				}
			}
			PluginUtterances = append(PluginUtterances, u.(*[]string))
			PluginFunctions = append(PluginFunctions, a.(func(string, string) string))
			PluginNames = append(PluginNames, *n.(*string))
			PluginList = append(PluginList, plugin)
			logger(file.Name() + " loaded successfully")
		} else {
			logger("Not loading " + file.Name() + ". Plugins must be built with 'go build -buildmode=plugin' and must end in '.so'.")
		}
	}
}
