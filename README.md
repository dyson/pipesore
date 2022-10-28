# Pipesore: A command-line text processor that nobody asked for

![version](https://img.shields.io/github/v/tag/dyson/pipesore?label=version)
[![release](https://github.com/dyson/pipesore/actions/workflows/release.yml/badge.svg)](https://github.com/dyson/pipesore/actions/workflows/release.yml)
[![build](https://github.com/dyson/pipesore/actions/workflows/build.yml/badge.svg)](https://github.com/dyson/pipesore/actions/workflows/build.yml)
[![test coverage](https://coveralls.io/repos/github/dyson/pipesore/badge.svg?branch=main)](https://coverallsio/github/dyson/pipesore?branch=main)
[![maintainability](https://api.codeclimate.com/v1/badges/a9de05463178f58c181f/maintainability)](https://codeclimate.com/github/dyson/pipesore/maintainability)
[![go report](https://goreportcard.com/badge/github.com/dyson/pipesore)](https://goreportcard.com/report/github.com/dyson/pipesore)
[![license](https://img.shields.io/github/license/dyson/pipesore.svg)](https://github.com/dyson/pipesore/blob/master/LICENSE)

*Pipe* because it's similar to unix like pipes and *sore* because the initial
hackathon version of this project was an eyesore.

Born from a proof of concept in using
[bitfield/script](https://github.com/bitfield/script) directly in the CLI
pipesore provides a number of text filters that you can pipe together to
process text. It takes input from stdin and writes the pipeline output to
stdout allowing it to be used alongside unix pipes.

## Motivation

Pipesore isn't intended to replace any of the well established cli text
processing tools (see
[https://tldp.org/LDP/abs/html/textproc.html](https://tldp.org/LDP/abs/html/textproc.html)
for a good list and examples). These tools do a single job well and have many
powerful features to accomplish anything you might want to do.

On the other hand there can be a bit of a learning curve to remember their
names , flags, and usage - even for basic tasks.

Pipesore is intended to be a single command that covers the most useful use
cases of these tools while being intuitive to even someone who has never seen
pipesore before.

## Installation

Download a stable [release](https://github.com/dyson/pipesore/releases) or install the latest development version with:

```bash
$ go install github.com/dyson/pipesore/cmd/pipesore@main
```

Optionally alias pipesore to something quicker to type, eg:

```bash
echo 'alias sp="pipesore"' >> ~/.bash_profile
```

## Basic Usage

A contrived example:

```bash
$ echo "cat cat cat dog bird bird bird bird" | pipesore 'Replace(" ", "\n") | Frequency() | First(1)'
4 bird
```

## Filters

All filters can be '|' (piped) together in any order, although not all ordering is logical.

All filter arguments are required. There a no assumptions about default values.

A filter prefixed with an "!" will return the opposite result of the non
prefixed filter of the same name. For example `First(1)` would return only the
first line of the input and `!First(1)` (read as not first) would skip the
first line of the input and return all other lines.

| Filter                                          |         |
| ------                                          | ------- |
| Columns(delimiter *string*, columns *string*)   | Returns the selected `columns` in order where `columns` is a 1-indexed comma separated list of column positions. Columns are defined by splitting with the 'delimiter'. |
| CountLines()                                    | Returns the line count. Lines are delimited by `\r?\n`. |
| CountRunes()                                    | Returns the rune (Unicode code points) count. Erroneous and short encodings are treated as single runes of width 1 byte. |
| CountWords()                                    | Returns the word count. Words are delimited by<br />`\t\|\n\|\v\|\f\|\r\| \|0x85\|0xA0`. |
| First(n int)                                    | Returns first `n` lines where `n` is a positive integer. If the input has less than `n` lines, all lines are returned. |
| !First(n int)                                   | Returns all but the the first `n` lines where `n` is a positive integer. If the input has less than `n` lines, no lines are returned. |
| Frequency()                                     | Ruturns a descending list containing frequency and unique line. Lines with equal frequency are sorted alphabetically. |
| Join(delimiter *string*)                        | Joins all lines together seperated by `delimiter`. |
| Last(n int)                                     | Returns last `n` lines where `n` is a positive integer. If the input has less than `n` lines, all lines are returned. |
| !Last(n int)                                    | Returns all but the last `n` lines where `n` is a positive integer. If the input has less than `n` lines, no lines are returned. |
| Match(substring *string*)                       | Returns all lines that contain `substring`. |
| !Match(substring *string*)                      | Returns all lines that don't contain `substring`. |
| MatchRegex(regex *string*)                      | Returns all lines that match the compiled regular expression 'regex'. Regex is in the form of [Re2](https://github.com/google/re2/wiki/Syntax). |
| !MatchRegex(regex *string*)                     | Returns all lines that don't match the compiled regular expression 'regex'. Regex is in the form of [Re2](https://github.com/google/re2/wiki/Syntax). |
| Replace(old *string*, replace *string*)         | Replaces all non-overlapping instances of `old` with `replace`. |
| ReplaceRegex(regex *string*, replace *string*)  | Replaces all matches of the compiled regular expression `regex` with `replace`. Inside `replace`, `$` signs represent submatches. For example `$1` represents the text of the first submatch. |

## License
See [LICENSE](https://github.com/dyson/pipesore/blob/master/LICENSE) file.
