---
title: "Getting Started With Templates"
description: "Learn how to get starting using templates."
category: narrative
tags:
  - templates
  - usage
weight: 10
---

The `$INSTALLATION_ROOT/templates` directory will be subsequently referred to as the **Default Template Location**.

Templates currently come in 2 flavours: `project` and `item`.

* A `project` is a template containing many files in a particular structure. They create a ready-to-run structure to start using a Puppet product. _These are great starting points._ You can create a boilerplate empty starter Puppet Module or a fully customized Puppet Module with specialized CI files and RSAPI providers.
* An `item` is a template that will supplement a project or existing content. These could be language features like a Puppet class or single files like a Git ignore file.

### pct new

PCT is available through the `pct new` command.

The `--list` or `-l` flag displays a list of locally available templates located in the **Default Template Location**. The list of templates is also available by calling `pct new` without flags.

``` bash
pct new
pct new --list
```

Example output:

<!-- This breaks glamour for some reason -->

    DISPLAYNAME                   | AUTHOR     | NAME                    | TYPE
    ──────────────────────────────┼────────────┼─────────────────────────┼─────────
    Bolt Plan                     | puppetlabs | bolt-plan               | item
    Bolt Project                  | puppetlabs | bolt-project            | project
    Bolt PowerShell Task          | puppetlabs | bolt-pwsh-task          | item
    Bolt YAML Plan                | puppetlabs | bolt-yaml-plan          | item
    Puppet Module Managed Gemfile | puppetlabs | git-attributes          | item
    Puppet Class                  | puppetlabs | puppet-class            | item
    Puppet Content Template       | puppetlabs | puppet-content-template | project
    Puppet Defined Type           | puppetlabs | puppet-defined-type     | item
    Puppet Fact                   | puppetlabs | puppet-fact             | item
    Puppet Module                 | puppetlabs | puppet-module           | project
    Puppet Resource API Provider  | puppetlabs | rsapi-provider          | item
    Puppet Resource API Transport | puppetlabs | puppet-transport        | item

Using the available templates above, its time to generate some content.

``` bash
pct new <author>/<template>
```

Replace `<author>` and `<template>` with the `author` and `name` of the template containing the content you want.

By default the `new <author>/<template>` function will use the directory name of your current working directory to "name" your new content.
To override this behaviour use the `--name` or `-n` flag.

``` bash
pct new <author>/<template> --name MyProject
```

By default the `new <author>/<template>` function will output the template content to the current working directory.
To override this behavour use the `--output` or `-o` flag.

``` bash
pct new <author>/<template> --output /path/to/your/project
```

> **Note:**
> Not all templates require a `name`.
> If a template doesn't require one, providing a value to the `--name` parameter will have no effect on the generated content.

### Example workflows

``` bash
> cd /home/me/projects/MyBoltProject
> pct new puppetlabs/bolt-project
```

``` bash
> pct new puppet-module -n MyNewProject -o /home/me/projects/
> cd /home/me/projects/MyNewProject
> pct new puppetlabs/puppet-fact -n ApplicationVersion
> pct new puppetlabs/rsapi-provider -n Awesomething
> pct new puppetlabs/puppet-transport -n AwesomethingApi
```

### Template Updates

At this time `pct new` will **NOT** update existing code to a newer version of a template.

If you run a `pct new` command using a `project` template, the project will replace the content within the output directory with the template code.

If you run a `pct new` command using an `item` template, the item will suppliment the content within the output directory with the template code. If files / folders that are named the same as the template content already exist, it will overwite this content.
