# json

This is a fork of [golang's encoding/json](https://github.com/golang/go/tree/master/src/encoding/json) with an extra method `MarshalSafeCollections` that marshals nil slices and maps into `[]` and `{}` respectfully instead of `null`. Additionally added `SetNilSafeCollection` to the `Encoder` in order to turn on the flag (default: false).

The semantic versioning for this repository follows that of the go src version that it forks from.

---

### Add support for nil slice force ignore omitempty
Details see: [example_proto_json_test.go](https://github.com/momopluto/json/blob/v1.15.x-ignore-omitempty/example_proto_json_test.go)