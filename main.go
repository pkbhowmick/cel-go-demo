package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resource struct {
	Self *runtime.RawExtension `json:"self"`
}

func main() {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("self", decls.Dyn)))
	if err != nil {
		panic(err)
	}

	ast, iss := env.Parse(`self.status.phase == 'Running'`)
	if iss.Err() != nil {
		panic(err)
	}

	checked, iss := env.Check(ast)
	if iss.Err() != nil {
		panic(iss.Err())
	}
	program, err := env.Program(checked)
	if err != nil {
		panic(err)
	}

	pod1 := &core.Pod{
		Status: core.PodStatus{
			Phase: core.PodRunning,
		},
	}

	resource := &Resource{Self: &runtime.RawExtension{Object: pod1}}
	bytes, err := json.Marshal(resource)
	if err != nil {
		panic(err)
	}

	obj := make(map[string]interface{})

	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		panic(err)
	}

	fmt.Println(obj)

	val, det, err := program.Eval(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(val) // true
	fmt.Println(det) // nil
}
