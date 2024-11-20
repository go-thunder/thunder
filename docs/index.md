# Introduction

Thunder is a collection of libraries and opinionated patterns to build
cloud-native services. The project provides different modules, which can be
used individually and replaced at any time.

Most of the modules use consolidated projects under the hood, the idea of this
project is to provide wrappers and connect all of these pieces seamlessly.

## How it works

This project provides modules and constructors for all of its components, so
one may use a dependency injection framework such as [Uber's
fx](https://uber-go.github.io/fx/), which is used in the docs, or manually
instantiate the components.

Take a look at the [modules page](./modules/index.md) to see the list of
available modules and the ones planned for the future.

## Why?

This project was created to solve some of the problems we faced when building
services at [Alternative Payments](https://www.alternativepayments.io/).

## Next

- [Getting started](./getting-started.md): how to use thunder right away.
- [Modules](./modules/index.md): list of available modules.
- [Cookbook](./cookbook.md): common use cases.
- [Conventions](./conventions.md): how to use thunder in your project.
- [Project structure](./project-structure.md): how to structure a thunder
  project.
