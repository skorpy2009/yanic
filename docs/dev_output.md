# Add new output type

Write a new package to implement the interface [output.Output:](https://chaos.expert/FreifunkBremen/yanic/blob/master/output/output.go)

```go
type Output interface {
	Save(nodes *runtime.Nodes)
}
```

**Save** a pre-filtered state of the Nodes



For startup, you need to bind your output type by calling
 `output.RegisterAdapter("typeofoutput",Register)`

it should be in the `func init() {}` of your package.



The _typeofoutput_ is used as mapping in the configuration `[[nodes.output.typeofoutput]]` the `map[string]interface{}` of the content are parsed to the _Register_ and on of your implemented `Output` or a `error` is needed as result.



Short: the function signature of _Register_ should be `func Register(configuration map[string]interface{}) (Output, error)`



At last add you import string to compile the your database as well in this [all](https://chaos.expert/FreifunkBremen/yanic/blob/master/output/all/main.go) package.
