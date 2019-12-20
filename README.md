# Performance Test with Electron for Xdelta for Go

This repository performs a performance test on the [Xdelta library for Go](http://github.com/konsorten/go-xdelta). 

It is using binary releases of [Electron](https://github.com/electron/electron).

To run all the tests, call the following command:

```
go test -v
```

The library assumes that the Xdelta for Go source code is located in the sibling-directory `../go-xdelta`.

## License

(The Apache 2 License)

Copyright 2019 marvin + konsorten GmbH (open-source@konsorten.de)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
