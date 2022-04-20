# termdict

termdict is a small dictionary tool for the command line, utilizing the [Free Dictionary API](https://github.com/meetDeveloper/freeDictionaryAPI).

You can use it to define words:

```bash
$ termdict define synthesis
synthesis
[noun] A deduction from the general to the particular.
[noun] The combination of thesis and antithesis.
[noun] (grammar) The uniting of ideas into a sentence.
[noun] The reunion of parts that have been divided.
```

or use it to manage your own vocab list for conventient access:

```bash
$ termdict add ameliorate entropy

$ termdict list
ameliorate
entropy

$ termdict random
ameliorate
[verb] To become better; improve.
```

## Install

Install via the below command:

```bash
go install github.com/caproven/termdict@latest
```

## Usage

Run `termdict` to see a list of available commands. Use the `--help` on any command to see all options.

## Storage

termdict will store your vocab list under `{USER_CONFIG_DIR}/termdict/`.
