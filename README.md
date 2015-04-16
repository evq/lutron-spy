# Lutron Spy

Lutron Spy lets you use 
[Lutron Pico remotes](http://www.amazon.com/Lutron-PJ2-WALL-WH-L01-Remote-Control-Mounting/dp/B00JR202JQ/)
along with a rooted Wink hub to control any REST interface. It
listens to the output of a sniffed serial port over which communication
with the Lutron radio takes place.

## Prerequisites

You must have slsnif already compiled for the Wink hub. I have made
a buildroot available in the form of a [docker container.](https://github.com/evq/imx28-buildroot)

## How to compile

Install [goxc](https://github.com/laher/goxc)

```bash
goxc
```

You can then find the cross compiled binary in 

```
$GOPATH/bin/lutron-spy-xc/snapshot/linux_arm/lutron-spy
```

## Configuration

Configuration is done via a json file which specifies which HTTP method,
content-type,
url and data to make a request with. Simply copy example-config.json to
remote-config.json and customize it. The serial number is not the one
listed on the back of the remotes, it can either be found in the sqlite3 
database and converted to hex or from lutron-spy output. Pressing a button
on an unconfigured remote will result in the serial and button number being
printed, but no further action.

```
serial:  FFA265
button: 2
```

## Usage

```bash
/etc/init.d/S60lutron-core stop
vi /etc/lutron.d/lutron.conf
```

Change the contents to:

```json
{
   "DatabasePath":"/database/lutron-db.sqlite",
   "ServerSocketPath":"/tmp/lutron-core.sock",
   "LoggingPriority":0,
   "HardwareInterfaceType":0,
   "HardwareInterfaceParams":[[0,"/dev/ttyp0"],[1,"115200"]],
   "SqliteDebugEnabled":1
}
```

```bash
./slsnif /dev/ttySP2 -x | ./lutron-spy &
/etc/init.d/S60lutron-core start
```

Remotes can be paired by using the built in aprontest binary.

```bash
aprontest -a -r lutron
```
