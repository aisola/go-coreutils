GO-Coreutils
------------
This is a Go1 implementation of the GNU Coreutils.

All files are maintained by: Abram C. Isola

[Go Coreutils Build](https://magnum-ci.com/status/3ebdc7894c138b7d57e3993c0b302747.png)

### Installation

via go get...

    $ go get github.com/aisola/go-coreutils/...

via git...

    $ git clone https://github.com/aisola/go-coreutils.git
    $ cd go-coreutils
    $ go build ./...

### Known Issues

+ Incomplete flags : Not all commands have the flags you may expect.
+ Large Binary Size: If file size of the binary is a concern, use the
gcc-go compiler instead: go build -compiler gccgo. The downside,
however, is that gccgo compiles binaries that are slower than gc.
Binaries compiled with the default gc binary have approximately 2MiB
of overhead. New releases of Go are working on improving the binary
sizes, as can be seen with the recently released Go 1.3 which reduced
binary sizes by 18.73% for this project (63.2 MiB to 51.4 MiB).



### Legal
go-coreutils 0.1 is licensed under the GNU General Public License v3.

    go-coreutils v0.1
    Copyright (C) 2014, The GO-Coreutils Developers.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
