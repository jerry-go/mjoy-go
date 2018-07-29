package actioncontext

import (
	"mjoy.io/common/types"
)

type api interface {
	setCtx(ctx *Context)
}

type baseApi struct {
	ctx *Context
}

func (self *baseApi) setCtx(ctx *Context) {
	self.ctx = ctx
}

type authorizationApi struct {
	baseApi
}

func (self *authorizationApi) RequireAuth(address types.Address) {

}

func (self *authorizationApi) IsAddress(address types.Address) bool {
	return false
}

func (self *authorizationApi) IsAccount(address types.Address) bool {
	return false
}

func (self *authorizationApi) IsContract(address types.Address) bool {
	return false
}

/*func (self *authorizationApi) pushRecipient(address types.Address) {

}*/

type systemApi struct {
	baseApi
}

type assertApi struct {
	baseApi
}

func (self *assertApi) Assert(test bool, msg string) {
	assertEx(test, msg)
}

type consoleApi struct {
	baseApi
}

type crpytoApi struct {
	baseApi
}

type producerApi struct {
	baseApi
}

func (self *producerApi) Producer() types.Address {
	return types.Address{}
}

type actionApi struct {
	baseApi
}

func (self *actionApi) Contract() types.Address {
	return self.ctx.con.self
}

type databaseApi struct {
	baseApi
}

func (self *databaseApi) Emplace(key []byte, val []byte) {

}

func (self *databaseApi) Modify(key []byte, val []byte) {

}

func (self *databaseApi) Erase(key []byte) {

}

func (self *databaseApi) Find(key []byte) []byte {
	return nil
}

func (self *databaseApi) Get(key []byte) []byte {
	ret := self.Find(key)
	assertEx(ret != nil, "unable to find key")
	return ret
}
