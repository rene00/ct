# ct - a command line metric tracker

ct is a personal project I created to track some basic metrics such as my body weight and daily walking distance.

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

Report on weight:

```bash
$ ct report --metric weight
2020-01-22T00:00:00Z weight 90
2020-01-23T00:00:00Z weight 100
```

Configure metric to be set daily:

```bash
$ ct configure --metric weight --frequency daily
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
