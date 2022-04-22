---
title: "About PCT"
description: "An overview of the PCT program."
category: concept
tags:
  - meta
---

## Overview

Puppet Content Templates (PCT) codify a structure to produce content for any Puppet Product that can be authored by Puppet Product Teams or external users without direct help of the PDK team.

PCT can create any type of a Puppet Product project: Puppet control repo, Puppet Module, Bolt project, etc. It can create one or more independent files, such as CI files or gitignores. This can be as simple as a name for a Puppet Class, a set of CI files to add to a Puppet Module, or as complex as a complete Puppet Control repo with roles and profiles.

These are meant to be ready-to-run, which means they put everything needed for a user to run the project from the moment after creation. This solves the 'blank page' problem, where a few files are in place but the user does not know what the next steps are.

> **Note:**
> PCT is currently in an EXPERIMENTAL phase and feedback is encouraged via [pct/discussions](https://github.com/puppetlabs/pct/discussions) and starting a `feedback` post.
