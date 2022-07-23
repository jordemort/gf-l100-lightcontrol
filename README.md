# gf-l100-lightcontrol

I own a [GF-L100 Floodlight Camera](https://fccid.io/2AL56GF-L100PRO/User-Manual/Users-Manual-3990248).
This is sometimes also marketed as the [Escam QF608](http://www.escam.cn/product/77-en.html).

In this thing's favor:

- It is cheap
- It has RTSP streams
- It is more-or-less ONVIF

Against:

- Controlling the floodlight is only possible through some proprietary app that doesn't integrate with Home Assistant

I had a fun time reverse-engineering the software on the camera and figured out that it controls the light by sending some weird ASCII-encoded hexadecimal numbers to a serial port.
So, I wrote this; it emulates just enough of [kankun-json](https://github.com/homedash/kankun-json) that it can be used with Home Assistant's [Kankun integration](https://www.home-assistant.io/integrations/kankun/) (I already had some Kankun plugs around).

## Example

Add this to your Home Assistant `configuration.yaml`

```yaml
switch:
  platform: kankun
  switches:
    floodlight_light:
      host: your.camera.host.or.ip
      port: 8090
      path: /light

    floodlight_motion:
      host: your.camera.host.or.ip
      port: 8090
      path: /motion
```
