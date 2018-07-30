////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: plugin.go
// @Date: 2018/07/30 17:34:30
////////////////////////////////////////////////////////////////////////////////

package interpreter

import "reflect"

type State int
const (
	Registered  = iota ///< the plugin is constructed but doesn't do anything
	Initialized        ///< the plugin has initialized any state required but is idle
)

type PluginImpl interface {
	Initialize()
	Getter
}

type Getter interface {
	Get() Interpreter
}

type Plugin interface {
	GetState() State
	Name() string

	PluginImpl
	//Register( /*p *Plugin*/ )
}

type PluginObj struct {
	pImpl PluginImpl
	state State
	flag bool
	name  pluginName
}

func newPluginObj(pImpl PluginImpl) *PluginObj {
	plugObj := &PluginObj{
		pImpl,
		Registered,
		false,
		pluginName{},
	}
	plugObj.name.Set(pImpl)
	return plugObj
}

func (obj *PluginObj) Initialize() {
	assert(obj.pImpl != nil)

	if obj.state == Registered {
		obj.state = Initialized
		obj.pImpl.Initialize()
	}
	assert(obj.state == Initialized)
}

func (obj PluginObj) Get() Interpreter {
	return obj.pImpl.Get()
}

func (obj PluginObj) GetState() State {
	return obj.state
}

func (obj PluginObj) Name() string {
	return obj.name.Name()
}

type pluginName struct {
	pImpl PluginImpl
	name  string
}

func (pn *pluginName) Set(pImpl interface{}) {
	ppImpl := reflect.TypeOf(pImpl).Elem()
	required := reflect.TypeOf((*PluginImpl)(nil)).Elem()
	assert(ppImpl.Implements(required))
	pn.pImpl, _ = ppImpl.(PluginImpl)
	pn.name = ppImpl.String()
	// TODO: The name only needs the string after the last ".". i.e., main.pkgname => pkgname
}

func (pn pluginName) Name() string {
	return pn.name
}
