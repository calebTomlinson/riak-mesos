// +build !rel

package cepm

//go:generate go-bindata -ignore=Makefile -o bindata_generated.go -pkg=cepm -prefix=data -debug data
import _ "github.com/jteeuwen/go-bindata"
