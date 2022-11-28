# reconciler

This repository contains a simple library to make reconciling Kubernetes
resources a bit less tedious. It wraps the typical get-update-else-create
loop in a convenient interface and relies on code generation. A set of
default reconciler functions are shipped with this library.
