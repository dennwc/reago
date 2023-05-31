# Go HTML component engine

This is just a simple/naive experiment for making HTML components based on the Go templates.

Components are defined in `./component/<tag>.gohtml` and will replace `<tag attr1="val1">content</tag>` in the main HTML of the page.

## Simple component

For example, component can be defined like this:

```html
{{define "mybutton"}}
<button type="button" class="btn btn-{{ .Attrs.type }}">
    {{ .Content }}
</button>
{{end}}
```

And used like this:

```html
<body>
    <div>
        <MyButton type="primary">Hello!</MyButton>
    </div>
</body>
```

The resulting HTML will be rendered as:

```html
<body>
    <div>
        <button type="button" class="btn btn-primary">
            Hello!
        </button>
    </div>
</body>
```

## Accessing data

Components can also access external data, which is exposed as `.DB` in the template context by the library user.

For example, when using the library, data interface can be defined as:

```go
type ExampleRecord struct {
	ID   int
	Name string
}

type ExampleDB struct {
}

func (d *ExampleDB) Table() []ExampleRecord {
	return []ExampleRecord{
		{ID: 1, Name: "Foo"},
		{ID: 2, Name: "Bar"},
	}
}
```

Then we can define a table component that shows this data:

```html
{{define "mytable"}}
    <table class="table">
        <tbody>
        {{range .DB.Table}}
            <tr><td>{{.ID}}</td><td>{{.Name}}</td></tr>
        {{end}}
        </tbody>
    </table>
{{end}}
```

And component can be used on the main HTML as:

```html
<body>
    <div>
        <MyTable></MyTable>
    </div>
</body>
```

We can potentially let the user pass parameters to this component. For this we first need to add parameter to the method:

```go
func (d *ExampleDB) Table(limit int) []ExampleRecord
```

And pass this parameter from component tag attributes:

```html
{{range .DB.Table .Attrs.limit}}
```

Usage on the main html:

```html
<body>
    <div>
        <MyTable limit="10"></MyTable>
    </div>
</body>
```

## License

MIT