## Clean up Orphaned Docker Volumes

### Purpose:
This tool is created to periodically remove orphaned volumes where the contents are not going to need future retrieval. In some environments where containers are treated ephemerally and  have high turnover there is the potential to generate a lot of space consuming orphaned volumes. It is VERY difficult to track down the origin of a particular volume once the container it was attached to is gone.

### Quick Usage

```
docker run -d -v /var/lib/docker:/var/lib/docker -v /var/run/docker.sock:/var/run/docker.sock cloudnautique/vol-cleanup [-n] -i 5
```

If you run with `-n` it will log what it would delete.

### Building

To generate a statically linked (Linux) go binary

```$ ./scripts/build```

It will be placed in the `./dist` directory.

### Packaging into a container

To generate a small Busybox container with the binary installed run:

`$ IMAGE=<username>/vol-cleanup ./scripts/package`

The default command of the image is to display the help page. Figured it was safer then starting to clean up files. 

### Command usage

```
NAME:
   vol-cleanup - Clean up orphaned Docker volumes

USAGE:
   vol-cleanup [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --noop, -n		Run in a noop mode and log to screen
   --interval, -i "5"	Number of minutes between runs
   --help, -h		show help
   --version, -v	print the version
```


## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.