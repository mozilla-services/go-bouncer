[![Updates](https://pyup.io/repos/github/mozilla-services/go-bouncer/shield.svg)](https://pyup.io/repos/github/mozilla-services/go-bouncer/)
[![Python 3](https://pyup.io/repos/github/mozilla-services/go-bouncer/python-3-shield.svg)](https://pyup.io/repos/github/mozilla-services/go-bouncer/)

# Bouncer Tests for download.allizom.org

Thank you for checking out Mozilla's Bouncer test suite. Mozilla and [Web QA team](https://quality.mozilla.org/teams/web-qa/) are grateful for the help and hard work of many contributors [past](https://github.com/mozilla/bouncer-tests/graphs/contributors) and [present](https://github.com/mozilla-services/go-bouncer/graphs/contributors) like yourself.

## Getting involved

We love working with contributors to improve test coverage our projects, but it
does require a few skills. By contributing to our test suite you will have an
opportunity to learn and/or improve your skills with Python, Selenium
WebDriver, GitHub, virtual environments, the Page Object Model, and more.

Our [new contributor guide][guide] should help you to get started, and will
also point you in the right direction if you need to ask questions.

## How to run the tests

### Clone the repository

If you have cloned this project already, then you can skip this; otherwise
you'll need to clone this repo using Git. If you do not know how to clone a
GitHub repository, check out this [help page][git clone] from GitHub.

If you think you would like to contribute to the tests by writing or
maintaining them in the future, it would be a good idea to create a fork of
this repository first, and then clone that. GitHub also has great instructions
for [forking a repository][git fork].

### Running tests locally

Then you can run the tests using [Docker][]:

```bash
$ docker build -t bouncer-tests .
$ docker run -it bouncer-tests
```

Use `-k` to run a specific test. For example,

```bash
docker run -it bouncer-tests pytest -k test_that_checks_redirect_using_incorrect_query_values
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

[Docker]: https://www.docker.com
[guide]: http://firefox-test-engineering.readthedocs.io/en/latest/guide/index.html
[git clone]: https://help.github.com/articles/cloning-a-repository/
[git fork]: https://help.github.com/articles/fork-a-repo/
