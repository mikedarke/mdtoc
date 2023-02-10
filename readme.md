# Markdown Table of Contents Generator 

Given a directory containing one or more markdown files this tool will generate a new file containing a "Table of contents" i.e. links to each of the headings within the markdown files.

## Controlling order 
When there are multiple markdown files within the folder you will often want to control the order they are listed within the table of contents. To do this you can add a HTML comment within the file that sets the order like this...

```html
<!-- index[1] -->
```

The files are sorted by the numerical value within `[]` in ascending order.

## Deep linking

To provide deep links to sections within each of the files you must specify `<div>` tags with `id` attributes next to the header e.g.

```html
<div id="main-heading">

# Main heading
```

## Options

The following options allow you to customize the output

| Option | Default Value | Descriptiion                                        |
| ------ | ------------- | --------------------------------------------------- |
| o     | index.md      | The output file name to write the table of contents |
| t     | Table of contents | The title text included at the top of the table of contents |

## Tasks

These tasks follow [eXeCute](https://github.com/Joe-Davidson1802/xc) syntax, therefore can be ran with `xc [taskname]`.

### test

Run all tests for the project.

```shell
go test ./...
```

### build

Tests and builds mdtoc

Requires: test

```shell
go build
```