# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest

LOCALES = (
    'ach',
    'af',
    'an',
    'ar',
    'as',
    'ast',
    'be',
    'bg',
    'bn-BD',
    'bn-IN',
    'br',
    'bs',
    'ca',
    'cs',
    'csb',
    'cy',
    'da',
    'de',
    'el',
    'en-GB',
    'en-ZA',
    'eo',
    'es-AR',
    'es-CL',
    'es-ES',
    'es-MX',
    'et',
    'eu',
    'fa',
    'ff',
    'fi',
    'fr',
    'fy-NL',
    'ga-IE',
    'gd',
    'gl',
    'gu-IN',
    'he',
    'hi-IN',
    'hr',
    'hsb',
    'hu',
    'hy-AM',
    'id',
    'is',
    'it',
    'ja',
    'kk',
    'km',
    'kn',
    'ko',
    'ku',
    'lij',
    'lt',
    'lv',
    'mai',
    'mk',
    'ml',
    'mr',
    'ms',
    'nb-NO',
    'nl',
    'nn-NO',
    'or',
    'pa-IN',
    'pl',
    'pt-BR',
    'pt-PT',
    'rm',
    'ro',
    'ru',
    'si',
    'sk',
    'sl',
    'son',
    'sq',
    'sr',
    'sv-SE',
    'ta',
    'te',
    'th',
    'tr',
    'uk',
    'vi',
    'xh',
    'zh-CN',
    'zh-TW',
    'zu'
)

OS = ('win', 'win64', 'linux', 'linux64', 'osx')


def pytest_generate_tests(metafunc):
    if 'lang' in metafunc.funcargnames:
        metafunc.parametrize('lang', LOCALES)

    if 'os' in metafunc.funcargnames:
        metafunc.parametrize('os', OS)


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
