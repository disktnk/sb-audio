# SensorBee audio plugin

## Setup

### Require

* PortAudio http://www.portaudio.com/
    * v19
* [gordonklaus/portaudio](https://github.com/gordonklaus/portaudio)
    * Go bidings for PortAudio
* SensorBee http://sensorbee.io/
    * v0.5.2 or later

### Build

```yaml
plugins:
- github.com/disktnk/sb-audio/plugin
```

```bash
$ build_sensorbee
```

## Usage

```sql
CREATE SOURCE audio TYPE audio_device WITH tick=3;
```

TODO

* discuss about map structure from source
* support wave and other file style
* support to select other audio device
* support rich delimiter, like volume param
