GO-Coreutils
------------
This is a Go1 implimentation of the GNU Coreutils. In general, as the commands were made for linux, 
even though the sources are written in Go1, they may not all be cross-platform.

### Installation

via goget...

    $ go get github.com/aisola/go-coreutils
    $ cd $GOPATH/src/github.com/aisola/go-coreutils
    $ python make build
    
via git...

    $ git clone https://github.com/aisola/go-coreutils.git
    $ cd go-coreutils
    $ python make build

### Legal
go-coreutils 0.1 is licensed under the GNU General Public License v3.
    
    go-coreutils v0.1
    Copyright (C) 2014, Abram C. Isola

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