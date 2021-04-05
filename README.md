# ct - consistency tracker

ct is a command line tool I created to track some basic metrics daily, such as my body weight and walking distance.

The intent of ct is to help consistently track key metrics so that I can track my progress.

## Usage

Initialise a ct config file and database:

```bash
$ ct init
```

Create a metric called _weight_:
```bash
$ ct metric create weight
```

Create a log for the _weight_ metric:

```bash
$ ct log create weight 100
```

Create a log for the _weight_ metric with a timestamp that is in the past:

```bash
$ ct log create weight 90 --timestamp 2020-01-22 
```

Update an existing log:

```bash
$ ct log create weight 95 --timestamp 2020-01-22 --update
```

Produce a monthly report on weight:

```bash
$ ct report monthly weight
+---------+---------+-------+
|  MONTH  | AVERAGE | COUNT |
+---------+---------+-------+
| 2020-01 |     95  |     2 |
+---------+---------+-------+
```

Configure metric to be an integer:

```bash
$ ct configure weight --data-type int
```

Configure text to be shown when logging a metric without passing the value param:

```bash
$ ct configure weight --value-text "whats your weight?"
```

Create a log for the _weight_ metric without passing the value param and be prompted with the _value-text_ config option:

```bash
$ ct log weight
whats your weight? 
```

Dump your data to json:

```bash
$ ct dump
```

Prompt for all metrics that dont have logs for today:

```bash
$ ct log prompt
```
