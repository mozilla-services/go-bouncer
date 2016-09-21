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
        # With Bug 1233779, WinXP bouncer is configured to redirect users
        # to Firefox version 43.0.1 if they visit firefox-latest and firefox-44.0
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)')

        self._verify_winxp_redirect_rules(base_url, product_alias, user_agent_ie6)

    @pytest.mark.parametrize(('product_alias'), _winxp_products)
    def test_ie6_winxp_useragent_5_2_redirects_to_correct_version(self, base_url, product_alias):
        # With Bug 1233779, WinXP bouncer is configured to redirect users
        # to Firefox version 43.0.1 if they visit firefox-latest and firefox-44.0
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1)')

        self._verify_winxp_redirect_rules(base_url, product_alias, user_agent_ie6)

    def _verify_winxp_redirect_rules(self, base_url, product_alias, user_agent,):
        param = {
            'product': 'firefox-' + product_alias,
            'lang': 'en-US',
            'os': 'win'
        }
        response = self._head_request(base_url, user_agent=user_agent, params=param)
        parsed_url = urlparse(response.url)
        if product_alias in ['stub', 'latest', '42.0', '43.0.1', '44.0']:
            assert '49.0.exe' in parsed_url.path
        elif 'esr' in product_alias:
            assert '38.5.1esr.exe' in parsed_url.path
        elif product_alias in ['49.0b8', '49.0b10', '49.0b37']:
            # beta versions are pinned to 49.0b10
            assert '49.0b10.exe' in parsed_url.path
        elif product_alias in ['firefox-sha1']:
            assert '49.0b8.exe' in parsed_url.path
        elif product_alias in ['beta-latest']:
            # beta-latest is a moving target, check that it is never
            # less than 49
            extracted_ver_num = parsed_url.path.split('Firefox%20Setup%20')[1].split('.')[0]
            assert 49 <= int(extracted_ver_num)
        elif product_alias in ['beta']:
            assert '49.0.exe' in parsed_url.path
        else:
            assert '49.0.exe' in parsed_url.path

    def test_ie6_vista_6_0_redirects_to_correct_version(self, base_url):
        user_agent = ('Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 6.0; SV1; .NET CLR 2.0.50727)')
        param = {
            'product': 'firefox-stub',
            'lang': 'en-US',
            'os': 'win'
        }
        response = self._head_request(base_url, user_agent=user_agent, params=param)
        parsed_url = urlparse(response.url)
        assert '49.0.exe' in parsed_url.path

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

    @pytest.mark.parametrize('product_alias', [
        {'product_name': 'firefox-beta-latest', 'lang': 'en-US'},
        {'product_name': 'firefox-latest-euballot', 'lang': 'en-GB'},
        {'product_name': 'firefox-latest', 'lang': 'en-US'},
        {'product_name': 'firefox-beta-stub', 'lang': 'en-US'},
        {'product_name': 'firefox-nightly-latest', 'lang': 'en-US'},
    ])
    def test_redirect_for_firefox_aliases(self, base_url, product_alias):
        param = {
            'product': product_alias['product_name'],
            'os': 'win',
            'lang': product_alias['lang']
        }

        response = self._head_request(base_url, params=param)

        parsed_url = urlparse(response.url)

        if not (
            product_alias['product_name'] == 'firefox-latest-euballot' and
            "download.allizom.org" in base_url
        ):
            url_scheme = 'http'
            if product_alias['product_name'] == 'firefox-beta-stub':
                url_scheme = 'https'
            elif product_alias['product_name'] in ['firefox-beta-latest', 'firefox-latest']:
                # Can be served over SSL or not
                # https://bugzilla.mozilla.org/show_bug.cgi?id=1299163#c1
                url_scheme = ['http', 'https']

            assert requests.codes.ok == response.status_code, \
                'Redirect failed with HTTP status. %s' % \
                self.response_info_failure_message(base_url, param, response)

            assert parsed_url.scheme in url_scheme, \
                'Failed to redirect to the correct scheme. %s' % \
                self.response_info_failure_message(base_url, param, response)

            assert parsed_url.netloc in ['download.cdn.mozilla.net',
                                         'edgecastcdn.net',
                                         'download-installer.cdn.mozilla.net',
                                         'cloudfront.net',
                                         'ftp.mozilla.org'], \
                'Failed, redirected to unknown host. %s' % \
                self.response_info_failure_message(base_url, param, response)

            if (
                product_alias['product_name'] != 'firefox-nightly-latest' and
                product_alias['product_name'] != 'firefox-aurora-latest' and
                product_alias['product_name'] != 'firefox-latest-euballot'
            ):
                assert '/win32/' in parsed_url.path, self.response_info(response)
