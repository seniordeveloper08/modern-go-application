#!/bin/bash

project=`basename $PWD`
package=`make var-PACKAGE`

boilerplatePackage="github.com/sagikazarmark/go-service-project-boilerplate"

rep() {
    sed -E -e $1 $2 > $2.new
    mv -f $2.new $2
}

mv .idea/project.iml .idea/${project}.iml
rep "s|.idea/project.iml|.idea/${project}.iml|g" .idea/modules.xml

# Run configurations
rep "s|name=\"project\"|name=\"${project}\"|" .idea/runConfigurations/Debug.xml
rep "s|name=\"project\"|name=\"${project}\"|" .idea/runConfigurations/All.xml
rep "s|name=\"project\"|name=\"${project}\"|" .idea/runConfigurations/Unit.xml
rep "s|name=\"project\"|name=\"${project}\"|" .idea/runConfigurations/Acceptance.xml
rep "s|name=\"project\"|name=\"${project}\"|" .idea/runConfigurations/Integration.xml

rep "s|${boilerplatePackage}|${package}|" cmd/app.go
rep "s|${boilerplatePackage}|${package}|" cmd/config.go
rep "s|${boilerplatePackage}|${package}|" cmd/main.go
