# Puppet Content Templates

* [Overview](#overview)
* [Getting Started](#getting-started)
  * [pct new](#pct-new)
  * [Template Updates](#template-updates)
  * [Tab Completion](#tab-completion)
* [Writing Your Own Templates](#writing-templates)
  * [Dos and Don'ts](#dos-and-donts)
  * [Structure](#structure)
    * [pct-config.yml](#pct-configyml)
  * [Templating Language](#templating-language)
* [Overriding Template Defaults](#overriding-template-defaults)

## Overview

Puppet Content Templates (PCT) codify a structure to produce content for any Puppet Product that can be authored by Puppet Product Teams or external users without direct help of the PDK team.

PCT can create any type of a Puppet Product project: Puppet control repo, Puppet Module, Bolt project, etc. It can create one or more independent files, such as CI files or gitignores. This can be as simple as a name for a Puppet Class, a set of CI files to add to a Puppet Module, or as complex as a complete Puppet Control repo with roles and profiles.

These are meant to be ready-to-run, which means they put everything needed for a user to run the project from the moment after creation. This solves the 'blank page' problem, where a few files are in place but the user does not know what the next steps are.

> :warning: PCT is currently in an EXPERIMENTAL phase and feedback is encouraged via [pdkgo/discussions](https://github.com/puppetlabs/pdkgo/discussions) and starting a `feedback` post.

## Getting Started

Grab the `experimental` PCT release available from [github.com/puppetlabs/pdkgo/releases/](https://github.com/puppetlabs/pdkgo/releases/).

Uncompress the archive to a location of your choosing - this will be refered to as `$INSTALLATION_ROOT` subsequently.

This should contain the following file and directory:

```bash
pct[.exe]
templates/
```

The `$INSTALLATION_ROOT/templates` directory will be subsequently referred to as the **Default Template Location**.

Templates currently come in 2 flavours: `project` and `item`.

* A `project` is a template containing many files in a particular structure. They create a ready-to-run structures to start using a Puppet product. _These are great starting points._ You can create a boilerplate empty starter Puppet Module or a fully customized Puppet Module with specialized CI files and RSAPI providers.
* An `item` is a template that will supplement a project or existing content. These could be language features like a Puppet class or single files like a Git ignore file.

### pct new

PCT is available through the `pct new` command.

The `--list` or `-l` flag displays a list of locally available templates located in the **Default Template Location**. The list of templates is also available by calling `pct new` without flags.

``` bash
pct new
pct new --list
```

Example output:

```bash
           DISPLAYNAME          |          NAME           |  TYPE
--------------------------------+-------------------------+----------
                                |                         |
  Bolt Plan                     | bolt-plan               | item
  Bolt Project                  | bolt-project            | project
  Bolt PowerShell Task          | bolt-pwsh-task          | item
  Bolt YAML Plan                | bolt-yaml-plan          | item
  Puppet Module Managed Gemfile | git-attributes          | item
  Puppet Class                  | puppet-class            | item
  Puppet Content Template       | puppet-content-template | project
  Puppet Defined Type           | puppet-defined-type     | item
  Puppet Fact                   | puppet-fact             | item
  Puppet Module                 | puppet-module           | project
  Puppet Resource API Provider  | rsapi-provider          | item
  Puppet Resource API Transport | puppet-transport        | item
```

Using the available templates above, its time to generate some content.

``` bash
pct new <template>
```

Replace `<template>` with the `name` of the template containing the content you want.

By default the `new <template>` function will use the directory name of your current working directory to "name" your new content.
To override this behaviour use the `--name` or `-n` flag.

``` bash
pct new <template> --name MyProject
```

By default the `new <template>` function will output the template content to the current working directory.
To override this behavour use the `--output` or `-o` flag.

``` bash
pct new <template> --output /path/to/your/project
```

> :memo: Not all templates require a `name`. If a template doesn't require one, providing a value to the `--name` parameter will have no effect on the generated content.

### Example workflows

``` bash
> cd /home/me/projects/MyBoltProject
> pct new bolt-project
```

``` bash
> pct new puppet-module -n MyNewProject -o /home/me/projects/
> cd /home/me/projects/MyNewProject
> pct new puppet-fact -n ApplicationVersion
> pct new rsapi-provider -n Awesomething
> pct new puppet-transport -n AwesomethingApi
```

### Template Updates

At this time `pct new` will **NOT** update existing code to a newer version of a template.

If you run a `pct new` command using a `project` template, the project will replace the content within the output directory with the template code.

If you run a `pct new` command using an `item` template, the item will suppliment the content within the output directory with the template code. If files / folders that are named the same as the template content already exist, it will overwite this content.

### Tab Completion

PCT has built in tab completion support. You can enable it in the following shells: `bash`, `zsh`, `fish` and `powershell`

To view the install instructions, access the `--help` menu in `pct completion` and follow the instructions.

```bash
pct completion --help
```

## Writing Templates

### Structure

A PCT is an archive containing a templated set of files and folders that represent a completed set of content. Files and folders stored in the template aren't limited to formal Puppet project types. Source files and folders may consist of any content that you wish to create when the template is used, even if the template engine produces just one file as its output.

### Location

You can specify the location of your templates using the `--templatepath` option:

```bash
pct new my-custom-project --templatepath /home/me/templates
```

### Composition

A PCT must contain a `pct-config.yml` in the root directory, alongside a `content` directory. The root directory must be named the same as the `id` value defined in  your `pct-config.yml`

The `content` directory contains the files and folders required to produce the `project` or `item`.

To mark a file as a template, use the `.tmpl` extension. Templated files can also use the global variable of `{{pct_name}}` to access the input from the `--name` cli argument.

> :memo: Folders within the `content` directory can also use the `{{pct_name}}` variable

Example template file names:

``` bash
myConfig.json.tmpl
{{pct_name}}_spec.rb
```

> :memo: One, all or none of the files can be templated.

#### pct-config.yml

Format of pct-config.yml

``` yaml
---
template:
  id: <a unique name>
  type: <'item' or 'project'>
  display: <a human readable name>
  version: <semver>
  url: <url to project repo>

<template parameters>
```

> :memo: Template `id` must not contain spaces or special characters. We recommend using a hyphen to break up the identifier.

Example pct-config.yml:

``` yaml
---
template:
  id: example-template
  type: project
  display: Example
  version: 0.1.0
  url: https://github.com/puppetlabs/pct-example
```

Example structure for `example-template`:

``` bash
> tree ~/templates/example-template-params
/Users/me/templates/example-template-params
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
  id: example-template-params
  type: project
  display: Example with Parameters
  version: 0.1.0
  url: https://github.com/puppetlabs/pct-example-params


example_template:
  foo: "bar"
  isPuppet: true
  colours:
  - "Red"
  - "Blue"
  - "Green"

```

In the above template `example-template-params` the parameters can be accessed in a `.tmpl` file like so:

``` go
{{.example_template.foo}}
{{.example_template.isPuppet}}
{{.example_template.colours}}
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

## Overriding Template Defaults

Perhaps you use a template often and find that you set the same values over and over?
As a template user, you can choose to override the default values specified by a template author.

> :memo: To view the default parameters for a template, look at it's `pct-config.yaml`.

To override these defaults you need to create a `pct.yml` within `$HOME/.pdk/`.  This file can contain overrides for multiple templates.

Example:

``` yaml
example_template:
  foo: "wibble"
  isPuppet: false
  colours:
  - "Red"
  - "Blue"

another_template:
  key: "value"
```

> :memo: You don't need to override everything
>
> ``` yaml
> example_template:
>  isPuppet: false
> ```
>
