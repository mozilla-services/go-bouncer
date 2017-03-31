# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.

import pytest
import requests
from urllib import urlencode
from urlparse import urlparse

from base import Base
import releng_utils as utils


class TestRedirects(Base):

    _locales = utils.get_firefox_locales()
    _os = ('win', 'win64', 'linux', 'linux64', 'osx')
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
        response = self.request_with_headers(base_url, user_agent=user_agent_ie6, params=param)
        assert '52.0.2esr.exe' in response.url, param

    @pytest.mark.parametrize(('product_alias'), _winxp_products)
    def test_ie6_winxp_useragent_5_2_redirects_to_correct_version(self, base_url, product_alias):
        user_agent_ie6 = ('Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1)')
        param = {
            'product': 'firefox-' + product_alias,
            'lang': 'en-US',
            'os': 'win'
        }
        response = self.request_with_headers(base_url, user_agent=user_agent_ie6, params=param)
        assert '52.0.2esr.exe' in response.url, param

    def _extract_windows_version_num(self, path):
        return int(path.split('Firefox%20Setup%20')[1].split('.')[0])

    def test_that_checks_redirect_using_incorrect_query_values(self, base_url):
        param = {
            'product': 'firefox-47.0.1',
            'lang': 'kitty_language',
            'os': 'stella'
        }
        response = self.request_with_headers(base_url, params=param)

        assert requests.codes.not_found == response.status_code, \
            self.response_info_failure_message(base_url, param, response)

        parsed_url = urlparse(response.url)

        assert 'http' == parsed_url.scheme, 'Failed to redirect to the correct scheme. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert urlparse(base_url).netloc == parsed_url.netloc, \
            self.response_info_failure_message(base_url, param, response)

        assert urlencode(param) == parsed_url.query, \
            self.response_info_failure_message(base_url, param, response)

    @pytest.mark.parametrize('os', _os)
    @pytest.mark.parametrize('locale', _locales)
    def test_verify_locales_redirect_to_the_expected_product(self, base_url, locale, os):
        """Verifies the downloaded version of Firefox matches the expected version number
        and filename when Firefox is requested for a specific OS and locale.

        The test verifies the following aliases: firefox-latest, firefox-esr-latest,
        firefox-nightly-latest, firefox-beta-latest, firefox-aurora-latest.
        """
        lang = locale.lang
        # Ja locale has a macOS-specific locale
        if lang == 'ja' and os == 'osx':
            lang = 'ja-JP-mac'

        for version in locale.versions:
            get_params = {
                'product': 'firefox-' + version,
                'lang': lang,
                'os': os
            }

            fx_pkg_name = self.get_expected_fx_pkg_str(os, 'latest', version)
            self.verify_redirect_to_correct_product(base_url, fx_pkg_name, get_params)

    def test_stub_installer_redirect_for_en_us_and_win(self, base_url):
        param = {
            'product': 'firefox-stub',
            'lang': 'en-US',
            'os': 'win'
        }

        response = self.request_with_headers(base_url, params=param)

        parsed_url = urlparse(response.url)

        assert requests.codes.ok == response.status_code, \
            'Redirect failed with HTTP status. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert 'https' == parsed_url.scheme, \
            'Failed to redirect to the correct scheme. %s' % \
            self.response_info_failure_message(base_url, param, response)

        assert parsed_url.netloc in self.cdn_netloc_locations, \
            'Failed to redirect to the correct host. %s' % \
            self.response_info_failure_message(base_url, param, response)
