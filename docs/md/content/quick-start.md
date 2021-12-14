---
title: "Quick Start Guide"
description: "Quick start guide to using PCT."
category: narrative
tags:
  - fundamentals
  - quickstart
weight: 20
---

This quick start guide will show you how to:

* Create a "bare bones" Puppet module from the `puppet-module-base` project template
* Add a Puppet Class to the module
* Add a Github Actions Workflow to test the module against the currently supported Puppet versions

### **STEP 1:** Create a Puppet Module

Let's name our module `test_module` using the `-n` flag:
![new_project_template](https://github.com/puppetlabs/pdkgo/blob/main/docs/_resources/new_module.gif?raw=true)

### **STEP 2:** Add a New Class

If we `cd` in to the module root dir, everything will get deployed with the correct layout:
![new_class](https://github.com/puppetlabs/pdkgo/blob/main/docs/_resources/new_class.gif?raw=true)

### **STEP 3:** Add a Github Actions Workflow

Want to know what configurable parameters are availble for a template and their defaults?
Run `pct new --info <TEMPLATE_AUTHOR>/<TEMPLATE_ID>`:

![new_info](https://github.com/puppetlabs/pdkgo/blob/main/docs/_resources/new_info.gif?raw=true)

We're happy with those defaults, so let's deploy this item.

Since we're outside the module root dir, we'll use the `-o` option to point at the root dir:

![new_info](https://github.com/puppetlabs/pdkgo/blob/main/docs/_resources/new_ghactions.gif?raw=true)
