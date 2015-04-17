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
make
```

You can then find the cross compiled binary in the source directory.

## Installation

Copy lutron-spy to /usr/local/bin/ and S59lutron-spy to /etc/init.d/ on your
Wink hub.

## Configuration

```bash
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

Remote configuration is done via a json file which specifies which HTTP method,
content-type,
url and data to make a request with. Simply copy example-config.json to
/etc/remote-config.json and customize it. The serial number is not the one
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
/etc/init.d/S59lutron-spy restart
```

Due to interdependency with lutron-core (lutron-core must be stopped before
starting lutron-spy, then started after lutron-spy), start and stop are
assumed to only be run by rcS and rcK during startup and shutdown.
/etc/init.d/S59lutron-spy restart will start lutron-spy even if
it was not originally running.

Remotes can be paired by using the built in aprontest binary.

```bash
aprontest -a -r lutron
```

### OpenHAB Integration

Here's an example of how to set it up to control a light in [openHAB](http://www.openhab.org/), with all 5 buttons working:

- On turns the light on
- Off turns the light off
- Up turns up the brightness
- Down turns down the brightness
- Select (middle button) sets the brightness to max

#### OpenHAB files

##### Item definitions
First, set up light with two items, one for the level and one for the state.

```java
Switch itm_light_office      "Office Light"  (Office, Lights)
Dimmer itm_light_lvl_office  "Level [%.1f]"  (Office, Lights)
```

##### Wink script

I have a script called `wink.sh` that just passes commands to `aprontest` through ssh
  - you could also use the set_dev_val.php script that comes with it,
  - you could also use [@nashira](https://github.com/nashira)'s [blink API](https://github.com/nashira/blink).

- My light here is `m1`, replace with a different m-value for the light you want to control.
- `-t1` is the on/off status, `-t2` is the brightness (at least for GE link bulbs, other brands may vary.)

##### Rule definitions
```java
rule "Office Light ON"
  when
    Item itm_light_office received command ON
  then
    executeCommandLine("sudo -u pi /opt/openhab/scripts/wink.sh -u -m1 -t1 -v ON");
end

rule "Office Light OFF"
  when
    Item itm_light_office received command OFF
  then
    executeCommandLine("sudo -u pi /opt/openhab/scripts/wink.sh -u -m1 -t1 -v OFF");
end

rule "Office Light Level"
  when
    Item itm_light_lvl_office received command
  then
    executeCommandLine("echo receivedCommand=" + receivedCommand)
    if (receivedCommand == INCREASE) {
      if (itm_light_lvl_office.state > 90) {
        itm_light_lvl_office.state = new PercentType(100)
      } else {
        itm_light_lvl_office.state = new PercentType((itm_light_lvl_office.state as DecimalType) + 10)
      }
    } else if (receivedCommand == DECREASE) {
      if (itm_light_lvl_office.state < 10) {
        itm_light_lvl_office.state = new PercentType(0)
      } else {
        itm_light_lvl_office.state = new PercentType((itm_light_lvl_office.state as DecimalType) - 10)
      }
    }

    // Convert percent ('100') to max ('255')
    var     new_per = (itm_light_lvl_office.state as DecimalType).floatValue
    var int new_val = new DecimalType(255 * (new_per * 0.01)).intValue

    // Set the dim level:
    executeCommandLine("sudo -u pi /opt/openhab/scripts/wink.sh -u -m1 -t2 -v " + new_val);
end
```

##### Remote config:

```json
"remotes": {
  "YOURSERIAL": {
    "nickname": "office",
    "buttons": {
      "on": {
        "url": "http://ip_to_your_openhab/rest/items/itm_light_office",
        "data": "ON",
        "method": "POST",
        "type": "text/plain"
      },
      "off": {
        "url": "http://ip_to_your_openhab/rest/items/itm_light_office",
        "data": "OFF",
        "method": "POST",
        "type": "text/plain"
      },
      "up": {
        "url": "http://ip_to_your_openhab/rest/items/itm_light_lvl_office",
        "data": "INCREASE",
        "method": "POST",
        "type": "text/plain"
      },
      "down": {
        "url": "http://ip_to_your_openhab/rest/items/itm_light_lvl_office",
        "data": "DECREASE",
        "method": "POST",
        "type": "text/plain"
      },
      "select": {
        "url": "http://ip_to_your_openhab/rest/items/itm_light_lvl_office",
        "data": "100",
        "method": "POST",
        "type": "text/plain"
      }
    }
  }
}}
```
