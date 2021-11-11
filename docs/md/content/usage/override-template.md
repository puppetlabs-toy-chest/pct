---
title: "Overriding Template Defaults"
description: "Learn how to override default template values."
---

## Overriding Template Defaults

Perhaps you use a template often and find that you set the same values over and over?
As a template user, you can choose to override the default values specified by a template author.

{{% alert color="light" %}}
To view the default parameters for a template run `pct new --info <TEMPLATE_ID>`.
{{% /alert %}}

To override these defaults you need to create a `pct.yml` containing the template id along with the values you wish to override.

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


{{% alert color="light" %}}
You don't need to override everything

``` yaml
example_template:
isPuppet: false
```
{{% /alert %}}
### User level configuration

Placing a `pct.yml` within `$HOME/.pdk/` allows you to create global overrides. Everytime you generate content from a template the configuration will be used.

### Workspace configuration

You may also place a `pct.yml` within a workspace.

Running `pct new` within a directory makes the current working directory your workspace.
If you specify an `--outputdir` that location is your workspace.

The configuration specified in a workspace `pct.yml` will override any configuration found within the user level configuration at `$HOME/.pdk/pct.yml`
