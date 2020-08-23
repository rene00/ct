# ct - consistency tracker

ct is a command line tool I created to track some basic metrics daily, such as my body weight and walking distance.

The intent of ct is to help consistently track key metrics so that I can track my progress.

## Usage

Initialise a ct config file and database:

```bash
$ ct init
```

Create and log new metric called weight:

```bash
$ ct log --metric weight --value 100
```

Log weight with a timestamp:

```bash
$ ct log --metric weight --timestamp 2020-01-22 --value 90
```

Edit existing log:

```bash
$ ct log --metric weight --timestamp 2020-01-22 --value 80 --edit
```

Report on weight:

```bash
$ ct report all --metrics weight
+----------------------+--------+-------+
|      TIMESTAMP       |  NAME  | VALUE |
+----------------------+--------+-------+
| 2020-01-22T00:00:00Z | weight |    90 |
| 2020-01-23T00:00:00Z | weight |   100 |
+----------------------+--------+-------+
```

Produce a monthly average report on weight:

```bash
$ ct report monthly-average --metrics weight
+---------+--------+-------+-------+
|  MONTH  |  NAME  | VALUE | COUNT |
+---------+--------+-------+-------+
| 2020-01 | weight |    95 |     2 |
+---------+--------+-------+-------+
```

Configure metric to be an integer:

```bash
$ ct configure --metric weight --data-type int
```

Configure text to be shown when logging a metric without passing the value param:

```bash
$ ct configure --metric weight --value-text "whats your weight?"
```

Log a metric without passing the value param and be prompted with the value-text config option:

```bash
$ ct log --metric weight
whats your weight? 
```

Dump your data to json:

```bash
$ ct dump
```

Have ct ask for a log when you open your shell but stop asking once you add a value:
```bash
$ ct configure --metric weight --value-text "whats your weight?"
```
and then add to your `~/.bashrc`:
```
ct log --metric weight 2>/dev/null
```
ct will continue to ask `whats your weight?` everytime you open a new shell and will stop once you submit a metric. It will then ask again the next day.
