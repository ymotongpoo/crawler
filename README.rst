.. -*- coding: utf-8 -*-
   Date: Tue Apr 30 16:57:57 2013
   Author: ymotongpoo (Yoshifumi YAMAGUCHI, ymotongpoo AT gmail.com)

.. _README:

=========
 crawler
=========

1. About
========

2ch crawler written in Go.

Original version is in Python written by @mopemope.
https://gist.github.com/mopemope/5464814


2. Future work
==============

Datastore implementation
------------------------

As of now, this only crawls all threads and do nothing further than that.
Currently, mongoDB is most possible candidate as a backend of datastore package.

http://labix.org/mgo

Better concurrent implementation
--------------------------------

https://gist.github.com/moriyoshi/5487253


3. License
==========

New BSD License:

Copyright (c) 2013, Yutaka Matsubara, Yoshi Yamaguchi
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice,
  this list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the "Yoshi Yamaguchi" nor the names of its contributors
  may be used to endorse or promote products derived from this software
  without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE
GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
THE POSSIBILITY OF SUCH DAMAGE.
