# Why?
Allows for generating code implementing a simple LRUCache. This makes it easy to not to have to use a package which
stores keys and values as interfaces. 

## Usage
```
go run main.go -package="test" -key="string" -value="*Example" > test.go
```
Usage is pretty self explanatory and just involves running main.go from the command line passing in the target
package, key type and value type as in the example above.

## Limitations
Keys can be any type, however values should be a pointer to any type of the users choosing. As this allows us to simply 
return nil when a value is not found, rather than an interface.