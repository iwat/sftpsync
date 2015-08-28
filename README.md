sftpsync
========

Sync local file system to SFTP.

Usage
-----

    $ sftpsync -h
    Usage of sftpsync:
      -dryrun
            Dry-run
      -v    Verbose
      -vv
            Very verbose

    $ sftpsync sftp://username:password@example.com:22/var/www/html
    NFO 2015/08/28 13:25:43 sftpsync.go:41: Dialing sftp://username:password@example.com:22/var/www/html
    NFO 2015/08/28 13:25:44 sftpsync.go:47: Listing files
    .................................................................................
    .................................................................................
    .................................................................................
    ...........................................
    NFO 2015/08/28 13:26:25 action.go:86: PUT -rw-r--r--  59781 2014-May-23 13:24:50 media/js/TableTools.js
    NFO 2015/08/28 13:26:25 action.go:86: PUT -rw-r--r--  59781 2014-May-23 13:24:50 media/js/TableTools.js
    NFO 2015/08/28 13:26:25 action.go:86: PUT -rw-r--r--  28892 2014-May-23 13:24:50 media/js/TableTools.min.js
    NFO 2015/08/28 13:26:25 action.go:86: PUT -rw-r--r--   2047 2014-May-23 13:24:50 media/swf/copy_cvs_xls.swf
    NFO 2015/08/28 13:26:25 action.go:86: PUT -rw-r--r--  11093 2014-May-23 13:24:50 media/js/ZeroClipboard.js
    NFO 2015/08/28 13:26:25 action.go:86: PUT -rw-r--r--  58660 2014-May-23 13:24:50 media/swf/copy_cvs_xls_pdf.swf

Legal
-----

This application is developed under MIT license, and can be used for open and
proprietary projects.

Copyright (c) 2015 Chaiwat Shuetrakoonpaiboon

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
