[![Build Status](https://travis-ci.org/gregoryv/workdir.svg?branch=master)](https://travis-ci.org/gregoryv/workdir)
[![codecov](https://codecov.io/gh/gregoryv/workdir/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/workdir)


[workdir](https://godoc.org/github.com/gregoryv/workdir) - Go package defines WorkDir type for easy file operations

Note! This package is the precursor
of [working](https://godoc.org/github.com/gregoryv/working). I was not
quite happy with name workdir, but keeping it open for now.

This package is most useful if you need to do multiple file operations
within one directory. Eg. something like

    cd /tmp/dir
	mkdir child1 child2
	touch child1/A child2/B
	# do some stuff
	cd
	rm -rf /tmp/dir

Simply cast a path to a WorkDir type

    wd := workdir.New("/tmp/dir")
	wd.MkdirAll("child1", "child2")
	wd.TouchAll("child1/A", "child2/B")
	// do stuff
	wd.RemoveAll()


There is also some code here for merging git output with plain ls result.
