# Bouncer Tests for download.allizom.org

Thank you for checking out Mozilla's Bouncer test suite. Mozilla and [Web QA team](https://quality.mozilla.org/teams/web-qa/) are grateful for the help and hard work of many contributors [past](https://github.com/mozilla/bouncer-tests/graphs/contributors) and [present](https://github.com/mozilla-services/go-bouncer/graphs/contributors) like yourself.

## Getting involved as a contributor

We love working with contributors to fill out the test coverage for Bouncer Tests, but it does require a few skills. You will need to know some Python and you will need some basic familiarity with [GitHub](https://guides.github.com/).

If you need to brush up on programming but are eager to start contributing immediately, please consider helping us [find bugs in Mozilla Firefox](https://oneanddone.mozilla.org/team/2/) or [find bugs in the Mozilla websites](https://oneanddone.mozilla.org/team/6/) tested by the Web QA team.

To brush up on Python skills before engaging with us, [Dive Into Python](http://www.diveintopython.net/toc/) is an excellent resource. MIT also has [lecture notes on Python](http://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-189-a-gentle-introduction-to-programming-using-python-january-iap-2011/) available through their open courseware. The programming concepts you will need to know include functions, working with classes, and some object oriented programming basics.

## Questions are always welcome

While we take pains to keep our documentation updated, the best source of information is those of us who work on the project. Don't be afraid to join us in [irc.mozilla.org](https://wiki.mozilla.org/IRC) [#mozwebqa](http://chat.mibbit.com/?server=irc.mozilla.org&channel=#mozwebqa) to ask questions about Bouncer Tests. Mozilla also hosts the [#mozillians](http://chat.mibbit.com/?server=irc.mozilla.org&channel=#mozillians) chat room to answer your general questions about contributing to Mozilla.

## How to set up and build Bouncer tests locally

This repository contains tests suite used to test Mozilla's Bouncer. Mozilla maintains a guide to run automated tests on our [QMO website](https://quality.mozilla.org/docs/webqa/running-webqa-automated-tests/).

You will need to install the following:

* **Git**: If you have cloned this project already then you can skip this! GitHub has excellent guides for [Windows](https://help.github.com/articles/set-up-git/#platform-windows), [OS X](https://help.github.com/articles/set-up-git/#platform-mac) and [Linux](https://help.github.com/articles/set-up-git/#platform-linux).
* **Python**: Before you will be able to run these tests you will need to have [Python 2.6](https://www.python.org/download/releases/2.6/) installed.
* **Tox**: You will need to [install Tox](https://testrun.org/tox/latest/install.html) to manage the virtual environments and run the tests.

### Running tests locally

To run these tests, use:

```bash
tox -e e2e -- --base-url=http://bouncer-bouncer.stage.mozaws.net
```

Use `-k` to run a specific test. For example,

```bash
tox -e e2e -- --base-url=http://bouncer-bouncer.stage.mozaws.net -k test_that_checks_redirect_using_incorrect_query_values
```

## Writing tests

If you want to get involved and add more tests then there's just a few things we'd like to ask you to do:

1. Use an existing file from this repository as a template for all new tests
2. Follow our simple [style guide](https://wiki.mozilla.org/QA/Execution/Web_Testing/Docs/Automation/StyleGuide)
3. Fork this project with your own GitHub account
4. Add your test into the `tests` folder
5. Make sure all tests are passing and submit a pull request with your changes

## License

This software is licensed under the MPL 2.0:

```
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
```
