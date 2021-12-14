---
title: "Sharing Templates"
description: "Learn how to share PCT templates."
category: narrative
tags:
  - templates
  - usage
---

After you've written your own template you may wish to share it with other members of your team or the wider Puppet community. Work is underway to improve this initial functionality.

### pct build

This command will attempt to package the current working directory. You can change the directory to pack by providing `--sourcedir`.

``` bash
pct build [--sourcedir <dir>][--targetdir <dir>]
```

The `build` command will ensure that the directory that you are attempting to package will produce a valid Puppet Content Template by looking for a `pct-config.yml` and a `content` directory.

The resulting `tar.gz` package will be created by default in `$cwd/pkg`. You can change the directory the package is created in by providing `--targetdir`.

### Installing template packages

Packages created using the `build` command can be installed by extracting the `tar.gz` into  the **Default Template Location**.

#### Local archive

Packages created using the `build` command can also be installed with the `pct install` command.

For example, this command:

```bash
pct install ~/my-template-1.2.3.tar.gz
```

Will install the template contained in `my-template-1.2.3.tar.gz` to the default template location.

#### Remote archive

Packages created using the `build` command can be automatically downloaded and extracted with `pct install` so long as you know the URL to where the archive is.

For example, this command:

```bash
pct install https://packages.mycompany.com/pct/my-template-1.2.3.tar.gz
```

Will attempt to download the PCT template from the specified url and then afterward install it like any other locally available PCT template archive.

#### Remote Git Repository

**Git** must be installed for this feature to work. The git repository must contain only one template and must be structured with the `pct-config.yml` file and the `content` directory in the root directory of the repository.

For more information on template structures see the [Writing Templates](https://github.com/puppetlabs/pdkgo#writing-templates) section in the `README`.

For example, this command:

```bash
pct install --git-uri https://github.com/myorg/myawesometemplate
```

This will attempt to clone the PCT template from the git repository at the specified URI and install to the default template location.
