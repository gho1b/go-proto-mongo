package main

import (
	go_proto_mongo_plugin "github.com/gho1b/go-proto-mongo/plugin"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	gen := generator.New()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err)
	}

	if err := proto.Unmarshal(data, gen.Request); err != nil {
		gen.Error(err)
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("No file ready to generate")
	}

	useGogoImport := false
	// Match parsing algorithm from Generator.CommandLineParameters
	for _, parameter := range strings.Split(gen.Request.GetParameter(), ",") {
		kvp := strings.SplitN(parameter, "=", 2)
		// We only care about key-value pairs where the key is "gogoimport"
		if len(kvp) != 2 || kvp[0] != "gogoimport" {
			continue
		}
		useGogoImport, err = strconv.ParseBool(kvp[1])
		if err != nil {
			gen.Error(err, "parsing gogoimport option")
		}
	}

	gen.CommandLineParameters(gen.Request.GetParameter())

	gen.WrapTypes()
	gen.SetPackageNames()
	gen.BuildTypeNameMap()
	gen.GeneratePlugin(plugin.NewPlugin(useGogoImport))

	for i := 0; i < len(gen.Response.File); i++ {
		gen.Response.File[i].Name = proto.String(strings.Replace(*gen.Response.File[i].Name, ".pb.go", ".validator.pb.go", -1))
	}

	// Send back the results.
	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "failed to write output proto")
	}
}
