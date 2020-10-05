# gaen

`gaen` is a simple cli to interact with the Google-Apple Exposure Notification system.

If you want to know more about the system please refer to the official Google and Apple documentations:
- https://www.google.com/covid19/exposurenotifications/
- https://covid19.apple.com/contacttracing


## Install

To install from source clone the project and run

```bash
go install ./...
```

## Usage

To download a TEK export run `gaen download` with a supported app (`immuni` or `swisscovid`, at the moment):

```
gaen download immuni
```

This will download the export in a `out/immuni/xxx` folder.

Then you can decode the export running

```
gaen decode out/immuni/xxx/export.bin
```

and this will output the first TEK with the first RPI:

```
TEK: [+KYwhAsB9QEq0PvaJfo4+Q==] - [F8 A6 30 84 B 1 F5 1 2A D0 FB DA 25 FA 38 F9]

RPI:
{
        "id": "q7WZAsXPgkuE+BEiadwLPQ==",
        "id_bytes": "AB B5 99 2 C5 CF 82 4B 84 F8 11 22 69 DC B 3D",
        "interval": "2020-10-03T02:00:00+02:00"
}
```
