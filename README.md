grepj
-----

grepj searches for a class within jar files.

Specification
-------------

    grepj <class-file-name> <file 1> ... <file n>

class-name is searched in the files provided. At least one file is required for a search.

The format of the class name is like foo.bar.Baz where foo.bar is the package name and Baz is the class name.

The files provided need to be valid jar, war, ear or zip files that can be read. Class files will be searched in nested locations as well. For example jar files in WEB-INF/lib of a war file will be searched. WEB-INF/classes will also be searched.

If the class file exists in the file provided, the file is output.

Download
--------

| Operating System | Architecture |                                                                           |
|------------------|--------------|---------------------------------------------------------------------------|
| Windows          | 64 bit       | [Download](http://grepj.s3.amazonaws.com/release/windows_amd64/grepj.exe) |
| Windows          | 32 bit       | [Download](http://grepj.s3.amazonaws.com/release/windows_386/grepj.exe)   |
| Linux            | 64 bit       | [Download](http://grepj.s3.amazonaws.com/release/linux_amd64/grepj)       |
| Linux            | 32 bit       | [Download](http://grepj.s3.amazonaws.com/release/linux_386/grepj)         |
| Mac              | 64 bit       | [Download](http://grepj.s3.amazonaws.com/release/darwin_amd64/grepj)      |
| Mac              | 32 bit       | [Download](http://grepj.s3.amazonaws.com/release/darwin_386/grepj)        |

Exit codes
----------

* 0 for class-file-name found in any of the files
* 1 for no match found
* 2 for any other errors

