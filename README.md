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

and this will output the decoded JSON

```json
{
    "ID": "+wK7aDl5cTC2wbZ4Ux6bvw==",
    "Date": "2020-10-02",
    "RPIs": [
        [
            {
                "ID": "GkUh9M/fYslxaxucp0ayWg==",
                "Interval": "2020-09-29T02:00:00+02:00"
            }
        ],
        ...
        [
            {
                "ID": "N99nmdxpfAvk0ByUGWI6EQ==",
                "Interval": "2020-09-29T02:40:00+02:00"
            }
        ]
    ],
    ...
}
```

### query

`gaen` implements a `--query` flag that follows the [JMESPath specification](https://jmespath.org/) that you can use to filter the output.

For example if you want to get the first TEK with its first RPI you can run:

```bash
gaen decode out/immuni/167/export.bin --query '[0].{ ID:ID, Date:Date, RPIS:RPIs[0] }'
```
```json
{
    "Date": "2020-09-29",
    "ID": "+wK7aDl5cTC2wbZ4Ux6bvw==",
    "RPIS": {
        "ID": "GkUh9M/fYslxaxucp0ayWg==",
        "Interval": "2020-09-29T02:00:00+02:00"
    }
}
```

or get the first 5 RPIs like this:

```
gaen decode out/immuni/167/export.bin --query '[0].{ ID:ID, Date:Date, RPIs:RPIs[:5].[{ ID:ID, Interval:Interval }] }'
```

```json
{
    "Date": "2020-09-29",
    "ID": "+wK7aDl5cTC2wbZ4Ux6bvw==",
    "RPIs": [
        [
            {
                "ID": "GkUh9M/fYslxaxucp0ayWg==",
                "Interval": "2020-09-29T02:00:00+02:00"
            }
        ],
        [
            {
                "ID": "72eL1vRaRXuRfTFTnR6gDA==",
                "Interval": "2020-09-29T02:10:00+02:00"
            }
        ],
        [
            {
                "ID": "itB7DTxs6aCl3FWz5QxQVw==",
                "Interval": "2020-09-29T02:20:00+02:00"
            }
        ],
        [
            {
                "ID": "h77WS2cu5x7SHoUw6Tgdfg==",
                "Interval": "2020-09-29T02:30:00+02:00"
            }
        ],
        [
            {
                "ID": "N99nmdxpfAvk0ByUGWI6EQ==",
                "Interval": "2020-09-29T02:40:00+02:00"
            }
        ]
    ]
}
```