# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest


def pytest_addoption(parser):
    parser.addoption('--product',
                     action='store',
                     dest='product',
                     metavar='str',
                     default='firefox-stub',
                     help='product under test')


@pytest.fixture
def product(request):
    return request.config.getoption('product')
