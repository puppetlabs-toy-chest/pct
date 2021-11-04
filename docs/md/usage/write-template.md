---
title: "Writing Templates"
description: "Learn how to write PCT templates."
weight: 20
---

## Writing Templates

### Structure

A PCT is an archive containing a templated set of files and folders that represent a completed set of content. Files and folders stored in the template aren't limited to formal Puppet project types. Source files and folders may consist of any content that you wish to create when the template is used, even if the template engine produces just one file as its output.

### Location

You can specify the location of your templates using the `--templatepath` option:

```bash
pct new my-name/my-custom-project --templatepath /home/me/templates
```

### Composition

A PCT must contain a `pct-config.yml` in the root directory, alongside a `content` directory.

The `content` directory contains the files and folders required to produce the `project` or `item`.

To mark a file as a template, use the `.tmpl` extension. Templated files can also use the global variable of `{{pct_name}}` to access the input from the `--name` cli argument.

{{% alert color="light" %}}
Folders within the `content` directory can also use the `{{pct_name}}` variable
{{% /alert %}}

Example template file names:

``` bash
myConfig.json.tmpl
{{pct_name}}_spec.rb
```

{{% alert color="light" %}}
One, all or none of the files can be templated.
{{% /alert %}}

#### pct-config.yml

Format of pct-config.yml

``` yaml
---
template:
  id: <a unique name>
  author: <name|username|orgname|handle|etc>
  type: <'item' or 'project'>
  display: <a human readable name>
  version: <semver>
  url: <url to project repo>

<template parameters>
```

{{% alert color="light" %}}
Template `id` and `author` must not contain spaces or special characters. We recommend using a hyphen to break up the identifier.
{{% /alert %}}

Example pct-config.yml:

``` yaml
---
template:
  id: example-template
  author: myorgname
  type: project
  display: Example
  version: 0.1.0
  url: https://github.com/puppetlabs/example-template
```

Example structure for `example-template`:

``` bash
> tree ~/templates/example-template
/Users/me/templates/example-template
├── content
│   └── example.txt.tmpl
└── pct-config.yml
```

### Templating Language

PCT uses [Go's templating language](https://golang.org/pkg/text/template/#hdr-Actions).

Example pct-config.yml with parameters:

``` yaml
---
template:
  id: example-template-with-params
  author: myorgname
  type: project
  display: Example with Parameters
  version: 0.1.0
  url: https://github.com/puppetlabs/pct-example-with-params


example_params:
  foo: "bar"
  isPuppet: true
  colours:
  - "Red"
  - "Blue"
  - "Green"

```

In the above template `example-template-with-params` the parameters can be accessed in a `.tmpl` file like so:

``` go
{{.example_params.foo}}
{{.example_params.isPuppet}}
{{.example_params.colours}}
```

outputs:

``` text
bar
<no value>
[Red Blue Green]
```

As a template author you can chose your own parameters and parameter structure so long as it is [valid YAML](https://yaml.org/spec/1.2/spec.html). Then utilise the GO templating language to display or iterate over these.

For most templates, we believe that you can do most of the things you would want with these common template controls:

``` go
// Outputs the value of `foo` defined within pct.yml
{{.example_template.foo}}

// A conditional
{{if .example_template.isPuppet}}
 "boo :("
{{else}}
 "yay!"
{{end}}

// Loops over all "colours" and renders each using {{.}}
{{range .example_template.colours}} {{.}} {{end}}
```

For more examples look at the existing templates provided in the **Default Template Location**.

### Dos and Don'ts

* `project` templates should provide all the code necessary to create a project from scratch and no more.
* Do not include configuration files that can be added via an `item` template later by an end user, for example, CI job configuration.
* Templates should be self documenting to help guide new users on how to use the file that has been created.
