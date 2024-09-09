# DNP encoder

[![Go Reference](https://pkg.go.dev/badge/github.com/OI4/oi4-oec-service-go/dnp.svg)](https://pkg.go.dev/github.com/OI4/oi4-oec-service-go/dnp)

The package dnp provides encoding and decoding functions for DNP strings.
The DNP is basically a URL encoding, replacing the '%' with an ',',
so that the resulting string is  DIN SPEC 91406 compliant

The basic encoding and decoding is taken from the net/url package.

Copyright 2009 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
 license that can be found in the LICENSE file.
