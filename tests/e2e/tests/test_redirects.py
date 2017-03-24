# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest
import requests
from urlparse import urlparse
from urllib import urlencode

from base import Base

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
    'zh-TW'
)

OS = ('win', 'win64', 'linux', 'linux64', 'osx')


class TestRedirects(Base):

    _winxp_products = [
        '38.5.1esr',
        '38.5.2esr',
        '38.5.3esr',
        '38.6.3esr',
        '40.0.0esr',
        'stub',
        'latest',
        'sha1',
        '42.0',
        '43.0.1',
        '44.0',
        'beta',
        'beta-latest',
        '49.0b8',
        '49.0b10',
        '49.0b37'
    ]

    @pytest.mark.parametrize(('product_alias'), _winxp_products)
    def test_ie6_winxp_useragent_5_1_redirects_to_correct_version(self, base_url, product_alias):
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)')
        param = {
            'product': 'firefox-' + product_alias,
            'lang': 'en-US',
            'os': 'win'
        }
        response = self._head_request(base_url, user_agent=user_agent_ie6, params=param)
        assert '52.0.1esr.exe' in response.url, param

    @pytest.mark.parametrize(('product_alias'), _winxp_products)
    def test_ie6_winxp_useragent_5_2_redirects_to_correct_version(self, base_url, product_alias):
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1)')
        param = {
            'product': 'firefox-' + product_alias,
            'lang': 'en-US',
            'os': 'win'
        }
        response = self._head_request(base_url, user_agent=user_agent_ie6, params=param)
        assert '52.0.1esr.exe' in response.url, param

    def _extract_windows_version_num(self, path):
        return int(path.split('Firefox%20Setup%20')[1].split('.')[0])

    def test_that_checks_redirect_using_incorrect_query_values(self, base_url):
        param = {
            'product': 'firefox-47.0.1',
            'lang': 'kitty_language',
            'os': 'stella'
        }
        response = self._head_request(base_url, params=param)

        assert requests.codes.not_found == response.status_code, \
            self.response_info_failure_message(base_url, param, response)

        parsed_url = urlparse(response.url)

        assert 'http' == parsed_url.scheme, 'Failed to redirect to the correct scheme. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert urlparse(base_url).netloc == parsed_url.netloc, \
            self.response_info_failure_message(base_url, param, response)

        assert urlencode(param) == parsed_url.query, \
            self.response_info_failure_message(base_url, param, response)

    @pytest.mark.parametrize('os', OS)
    @pytest.mark.parametrize('lang', LOCALES)
    def test_that_checks_redirect_using_locales_and_os(
        self,
        base_url,
        lang,
        os
    ):
        # Ja locale has a special code for mac
        if lang == 'ja' and os == 'osx':
            lang = 'ja-JP-mac'

        param = {
            'product': 'firefox-47.0.1',
            'lang': lang,
            'os': os
        }

        response = self._head_request(base_url, params=param)

        parsed_url = urlparse(response.url)

        assert requests.codes.ok == response.status_code, \
            'Redirect failed with HTTP status. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert 'http' == parsed_url.scheme, 'Failed to redirect to the correct scheme. %s' % \
            self.response_info_failure_message(base_url, param, response)

    def test_stub_installer_redirect_for_en_us_and_win(self, base_url, product):
        param = {
            'product': product,
            'lang': 'en-US',
            'os': 'win'
        }

        response = self._head_request(base_url, params=param)

        parsed_url = urlparse(response.url)

        assert requests.codes.ok == response.status_code, \
            'Redirect failed with HTTP status. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert 'https' == parsed_url.scheme, \
            'Failed to redirect to the correct scheme. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert 'download-installer.cdn.mozilla.net' == parsed_url.netloc, \
            'Failed by redirected to incorrect host. %s' % \
            self.response_info_failure_message(base_url, param, response)
